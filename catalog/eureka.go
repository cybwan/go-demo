package catalog

import (
	"github.com/hashicorp/consul/api"
	"github.com/hudl/fargo"
)

var (
	eurekaDiscoveryClient *EurekalDiscoveryClient
)

type EurekalDiscoveryClient struct {
	eurekaClient fargo.EurekaConnection
}

// Service is used to query catalog entries for a given service
func (dc *EurekalDiscoveryClient) Service(service, tag string, q *QueryOptions) ([]*CatalogService, error) {
	// Set up query options
	opts := api.QueryOptions{}
	opts.AllowStale = q.AllowStale
	opts.Namespace = q.Namespace

	// Only consider services that are tagged from k8s
	services, err := dc.eurekaClient.GetApp(service)
	//services, _, err := dc.eurekaClient.Catalog().Service(service, tag, &opts)
	if err != nil {
		return nil, err
	}
	catalogServices := make([]*CatalogService, len(services.Instances))
	for idx, ins := range services.Instances {
		catalogServices[idx] = new(CatalogService)
		catalogServices[idx].fromEureka(ins)
	}
	return catalogServices, nil
}

func (dc *EurekalDiscoveryClient) NodeServiceList(node string, q *QueryOptions) (*CatalogNodeServiceList, error) {
	return nil, nil
}

func (dc *EurekalDiscoveryClient) Deregister(dereg *CatalogDeregistration) error {
	err := dc.eurekaClient.DeregisterInstance(dereg.toEureka())
	return err
}

func (dc *EurekalDiscoveryClient) Register(reg *CatalogRegistration) error {
	err := dc.eurekaClient.RegisterInstance(reg.toEureka())
	return err
}

// EnsureNamespaceExists ensures a Consul namespace with name ns exists. If it doesn't,
// it will create it and set crossNSACLPolicy as a policy default.
// Boolean return value indicates if the namespace was created by this call.
func (dc *EurekalDiscoveryClient) EnsureNamespaceExists(ns string, crossNSAClPolicy string) (bool, error) {
	return false, nil
}

func GetEurekalDiscoveryClient() *EurekalDiscoveryClient {
	if eurekaDiscoveryClient == nil {
		eurekaDiscoveryClient = new(EurekalDiscoveryClient)
		httpAddr := "http://127.0.0.1:8761/eureka"
		eurekaDiscoveryClient.eurekaClient = fargo.NewConn(httpAddr)
	}
	return eurekaDiscoveryClient
}
