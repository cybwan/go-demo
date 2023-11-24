package catalog

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hudl/fargo"
)

// AgentCheck represents a check known to the agent
type AgentCheck struct {
	CheckID   string
	ServiceID string
	Name      string
	Namespace string
	Type      string
	Status    string
	Output    string
}

func (ac *AgentCheck) toConsul() *api.AgentCheck {
	check := new(api.AgentCheck)
	check.CheckID = ac.CheckID
	check.ServiceID = ac.ServiceID
	check.Name = ac.Name
	check.Namespace = ac.Namespace
	check.Type = ac.Type
	check.Status = ac.Status
	check.Output = ac.Output
	return check
}

type AgentWeights struct {
	Passing int
	Warning int
}

func (aw *AgentWeights) toConsul() api.AgentWeights {
	return api.AgentWeights{
		Passing: aw.Passing,
		Warning: aw.Warning,
	}
}

func (aw *AgentWeights) fromConsul(w api.AgentWeights) {
	aw.Passing = w.Passing
	aw.Warning = w.Warning
}

// AgentService represents a service known to the agent
type AgentService struct {
	ID        string
	Service   string
	Namespace string
	Address   string
	Port      int
	Weights   AgentWeights
	Tags      []string
	Meta      map[string]string
}

func (as *AgentService) toConsul() *api.AgentService {
	agentService := new(api.AgentService)
	agentService.ID = as.ID
	agentService.Service = as.Service
	agentService.Namespace = as.Namespace
	agentService.Address = as.Address
	agentService.Port = as.Port
	agentService.Weights = as.Weights.toConsul()
	if len(as.Tags) > 0 {
		agentService.Tags = append(agentService.Tags, as.Tags...)
	}
	if len(as.Meta) > 0 {
		agentService.Meta = make(map[string]string)
		for k, v := range as.Meta {
			agentService.Meta[k] = v
		}
	}
	return agentService
}

func (as *AgentService) fromConsul(agentService *api.AgentService) {
	as.ID = agentService.ID
	as.Service = agentService.Service
	as.Namespace = agentService.Namespace
	as.Address = agentService.Address
	as.Port = agentService.Port
	as.Weights.fromConsul(agentService.Weights)
	if len(agentService.Tags) > 0 {
		as.Tags = append(as.Tags, agentService.Tags...)
	}
	if len(agentService.Meta) > 0 {
		as.Meta = make(map[string]string)
		for k, v := range agentService.Meta {
			as.Meta[k] = v
		}
	}
}

type CatalogDeregistration struct {
	Node      string
	ServiceID string
	Service   string
	Namespace string
}

func (cdr *CatalogDeregistration) toConsul() *api.CatalogDeregistration {
	r := new(api.CatalogDeregistration)
	r.Node = cdr.Node
	r.ServiceID = cdr.ServiceID
	r.Namespace = cdr.Namespace
	return r
}

func (cdr *CatalogDeregistration) toEureka() *fargo.Instance {
	r := new(fargo.Instance)
	r.InstanceId = cdr.ServiceID
	r.App = cdr.Service
	return r
}

type CatalogRegistration struct {
	Node           string
	Address        string
	NodeMeta       map[string]string
	Service        *AgentService
	Check          *AgentCheck
	SkipNodeUpdate bool
}

func (cr *CatalogRegistration) toConsul() *api.CatalogRegistration {
	r := new(api.CatalogRegistration)
	r.Node = cr.Node
	r.Address = cr.Address
	if len(cr.NodeMeta) > 0 {
		r.NodeMeta = make(map[string]string)
		for k, v := range cr.NodeMeta {
			r.NodeMeta[k] = v
		}
	}
	if cr.Service != nil {
		r.Service = cr.Service.toConsul()
	}
	if cr.Check != nil {
		r.Check = cr.Check.toConsul()
	}
	r.SkipNodeUpdate = cr.SkipNodeUpdate
	return r
}

func (cr *CatalogRegistration) toEureka() *fargo.Instance {
	r := new(fargo.Instance)
	if len(cr.NodeMeta) > 0 {
		for k, v := range cr.NodeMeta {
			r.SetMetadataString(k, v)
		}
	}
	if cr.Service != nil {
		r.UniqueID = func(i fargo.Instance) string {
			return cr.Service.ID
		}
		r.InstanceId = cr.Service.ID
		r.HostName = cr.Service.Address
		r.IPAddr = cr.Service.Address
		r.App = cr.Service.Service
		r.VipAddress = strings.ToLower(cr.Service.Service)
		r.SecureVipAddress = strings.ToLower(cr.Service.Service)
		r.Port = cr.Service.Port
		r.Status = fargo.UP
		r.DataCenterInfo = fargo.DataCenterInfo{Name: fargo.MyOwn}

		r.HomePageUrl = fmt.Sprintf("http://%s:%d/", cr.Service.Address, cr.Service.Port)
		r.StatusPageUrl = fmt.Sprintf("http://%s:%d/actuator/info", cr.Service.Address, cr.Service.Port)
		r.HealthCheckUrl = fmt.Sprintf("http://%s:%d/actuator/health", cr.Service.Address, cr.Service.Port)
	}
	return r
}

type CatalogService struct {
	Node        string
	ServiceID   string
	ServiceName string
}

func (cs *CatalogService) fromConsul(svc *api.CatalogService) {
	if svc == nil {
		return
	}
	cs.Node = svc.Node
	cs.ServiceID = svc.ServiceID
	cs.ServiceName = svc.ServiceName
}

func (cs *CatalogService) fromEureka(svc *fargo.Instance) {
	if svc == nil {
		return
	}
	cs.Node = svc.DataCenterInfo.Name
	cs.ServiceID = svc.Id()
	cs.ServiceName = svc.App
}

type CatalogNodeServiceList struct {
	Services []*AgentService
}

func (cnsl *CatalogNodeServiceList) fromConsul(svcList *api.CatalogNodeServiceList) {
	if svcList == nil || len(svcList.Services) == 0 {
		return
	}
	for _, svc := range svcList.Services {
		agentService := new(AgentService)
		agentService.fromConsul(svc)
		cnsl.Services = append(cnsl.Services, agentService)
	}
}

// QueryOptions are used to parameterize a query
type QueryOptions struct {
	// AllowStale allows any Consul server (non-leader) to service
	// a read. This allows for lower latency and higher throughput
	AllowStale bool

	// Namespace overrides the `default` namespace
	// Note: Namespaces are available only in Consul Enterprise
	Namespace string

	// WaitIndex is used to enable a blocking query. Waits
	// until the timeout or the next index is reached
	WaitIndex uint64

	// WaitTime is used to bound the duration of a wait.
	// Defaults to that of the Config, but can be overridden.
	WaitTime time.Duration

	// Providing a peer name in the query option
	Peer string

	// Filter requests filtering data prior to it being returned. The string
	// is a go-bexpr compatible expression.
	Filter string
}

func (qopts *QueryOptions) toConsul() *api.QueryOptions {
	opts := new(api.QueryOptions)
	opts.AllowStale = qopts.AllowStale
	opts.Namespace = qopts.Namespace
	opts.WaitIndex = qopts.WaitIndex
	opts.WaitTime = qopts.WaitTime
	opts.Peer = qopts.Peer
	opts.Filter = qopts.Filter
	return opts
}

type ServiceDiscoveryClient interface {
	NodeServiceList(node string, q *QueryOptions) (*CatalogNodeServiceList, error)
	Service(service, tag string, q *QueryOptions) ([]*CatalogService, error)
	Register(reg *CatalogRegistration) error
	Deregister(dereg *CatalogDeregistration) error
	EnsureNamespaceExists(ns string, crossNSAClPolicy string) (bool, error)
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
