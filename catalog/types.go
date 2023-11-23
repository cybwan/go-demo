package catalog

import (
	"time"

	"github.com/hashicorp/consul/api"
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
	Namespace string
}

func (cdr *CatalogDeregistration) toConsul() *api.CatalogDeregistration {
	r := new(api.CatalogDeregistration)
	r.Node = cdr.Node
	r.ServiceID = cdr.ServiceID
	r.Namespace = cdr.Namespace
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
