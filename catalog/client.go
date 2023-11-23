package catalog

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/consul-k8s/control-plane/namespaces"
	"github.com/hashicorp/consul/api"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	discoveryClient *DiscoveryClient
	consulclient    *api.Client
	kubeclient      *kubernetes.Clientset
)

type DiscoveryClient struct {
	consulClient *api.Client
}

// Service is used to query catalog entries for a given service
func (dc *DiscoveryClient) Service(service, tag string, q *QueryOptions) ([]*CatalogService, error) {
	// Set up query options
	opts := api.QueryOptions{}
	opts.AllowStale = q.AllowStale
	opts.Namespace = q.Namespace

	// Only consider services that are tagged from k8s
	services, _, err := dc.consulClient.Catalog().Service(service, tag, &opts)
	if err != nil {
		return nil, err
	}

	catalogServices := make([]*CatalogService, len(services))
	for idx, svc := range services {
		catalogServices[idx] = new(CatalogService)
		catalogServices[idx].fromConsul(svc)
	}
	return catalogServices, nil
}

func (dc *DiscoveryClient) NodeServiceList(node string, q *QueryOptions) (*CatalogNodeServiceList, error) {
	nodeServices, meta, err := dc.consulClient.Catalog().NodeServiceList(node, q.toConsul())
	if err != nil {
		return nil, err
	}

	// Update our blocking index
	q.WaitIndex = meta.LastIndex

	nodeServiceList := new(CatalogNodeServiceList)
	nodeServiceList.fromConsul(nodeServices)
	return nodeServiceList, nil
}

func (dc *DiscoveryClient) Deregister(dereg *CatalogDeregistration) error {
	_, err := dc.consulClient.Catalog().Deregister(dereg.toConsul(), nil)
	return err
}

func (dc *DiscoveryClient) Register(reg *CatalogRegistration) error {
	_, err := dc.consulClient.Catalog().Register(reg.toConsul(), nil)
	return err
}

// EnsureExists ensures a Consul namespace with name ns exists. If it doesn't,
// it will create it and set crossNSACLPolicy as a policy default.
// Boolean return value indicates if the namespace was created by this call.
func (dc *DiscoveryClient) EnsureExists(ns string, crossNSAClPolicy string) (bool, error) {
	return namespaces.EnsureExists(dc.consulClient, ns, crossNSAClPolicy)
}

func getDiscoveryClient() *DiscoveryClient {
	if discoveryClient == nil {
		discoveryClient = new(DiscoveryClient)
		c := api.DefaultConfig()
		c.Address = "127.0.0.1:8500"
		discoveryClient.consulClient, _ = api.NewClient(c)
	}
	return discoveryClient
}

func GetKubeClient() *kubernetes.Clientset {
	if kubeclient == nil {
		// Initialize kube config and client
		kubeConfigFile := "/Users/baili/.kube/config"
		kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		kubeclient = kubernetes.NewForConfigOrDie(kubeConfig)
	}
	return kubeclient
}

const (
	// HealthAny is special, and is used as a wild card,
	// not as a specific state.
	HealthAny      = "any"
	HealthPassing  = "passing"
	HealthWarning  = "warning"
	HealthCritical = "critical"
	HealthMaint    = "maintenance"
)

// ConsulNamespace returns the consul namespace that a service should be
// registered in based on the namespace options. It returns an
// empty string if namespaces aren't enabled.
func ConsulNamespace(kubeNS string, enableConsulNamespaces bool, consulDestNS string, enableMirroring bool, mirroringPrefix string) string {
	if !enableConsulNamespaces {
		return ""
	}

	// Mirroring takes precedence.
	if enableMirroring {
		return fmt.Sprintf("%s%s", mirroringPrefix, kubeNS)
	}

	return consulDestNS
}

// ParseTags parses the tags annotation into a slice of tags.
// Tags are split on commas (except for escaped commas "\,").
func ParseTags(tagsAnno string) []string {

	// This algorithm parses the tagsAnno string into a slice of strings.
	// Ideally we'd just split on commas but since Consul tags support commas,
	// we allow users to escape commas so they're included in the tag, e.g.
	// the annotation "tag\,with\,commas,tag2" will become the tags:
	// ["tag,with,commas", "tag2"].

	var tags []string
	// nextTag is built up char by char until we see a comma. Then we
	// append it to tags.
	var nextTag string

	for _, runeChar := range tagsAnno {
		runeStr := fmt.Sprintf("%c", runeChar)

		// Not a comma, just append to nextTag.
		if runeStr != "," {
			nextTag += runeStr
			continue
		}

		// Reached a comma but there's nothing in nextTag,
		// skip. (e.g. "a,,b" => ["a", "b"])
		if len(nextTag) == 0 {
			continue
		}

		// Check if the comma was escaped comma, e.g. "a\,b".
		if string(nextTag[len(nextTag)-1]) == `\` {
			// Replace the backslash with a comma.
			nextTag = nextTag[0:len(nextTag)-1] + ","
			continue
		}

		// Non-escaped comma. We're ready to push nextTag onto tags and reset nextTag.
		tags = append(tags, strings.TrimSpace(nextTag))
		nextTag = ""
	}

	// We're done the loop but nextTag still contains the last tag.
	if len(nextTag) > 0 {
		tags = append(tags, strings.TrimSpace(nextTag))
	}

	return tags
}
