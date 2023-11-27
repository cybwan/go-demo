package catalog

import (
	"github.com/hashicorp/consul-k8s/control-plane/namespaces"
	"github.com/hashicorp/consul/api"
)

var (
	consulDiscoveryClient *ConsulDiscoveryClient
)

type ConsulDiscoveryClient struct {
	consulClient *api.Client
}

// CatalogService is used to query catalog entries for a given service
func (dc *ConsulDiscoveryClient) CatalogService(service, tag string, q *QueryOptions) ([]*CatalogService, error) {
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

func (dc *ConsulDiscoveryClient) NodeServiceList(node string, q *QueryOptions) (*CatalogNodeServiceList, error) {
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

func (dc *ConsulDiscoveryClient) Deregister(dereg *CatalogDeregistration) error {
	_, err := dc.consulClient.Catalog().Deregister(dereg.toConsul(), nil)
	return err
}

func (dc *ConsulDiscoveryClient) Register(reg *CatalogRegistration) error {
	_, err := dc.consulClient.Catalog().Register(reg.toConsul(), nil)
	return err
}

// EnsureNamespaceExists ensures a Consul namespace with name ns exists. If it doesn't,
// it will create it and set crossNSACLPolicy as a policy default.
// Boolean return value indicates if the namespace was created by this call.
func (dc *ConsulDiscoveryClient) EnsureNamespaceExists(ns string, crossNSAClPolicy string) (bool, error) {
	return namespaces.EnsureExists(dc.consulClient, ns, crossNSAClPolicy)
}

func GetConsulDiscoveryClient() *ConsulDiscoveryClient {
	if consulDiscoveryClient == nil {
		consulDiscoveryClient = new(ConsulDiscoveryClient)
		c := api.DefaultConfig()
		c.Address = "127.0.0.1:8500"
		consulDiscoveryClient.consulClient, _ = api.NewClient(c)
	}
	return consulDiscoveryClient
}
