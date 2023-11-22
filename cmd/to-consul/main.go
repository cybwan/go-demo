package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/hashicorp/go-hclog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cybwan/go-demo/catalog"
	"github.com/cybwan/go-demo/controller"
)

// AppendSliceValue implements the flag.Value interface and allows multiple
// calls to the same variable to append a list.
type AppendSliceValue []string

func (s *AppendSliceValue) String() string {
	return strings.Join(*s, ",")
}

func (s *AppendSliceValue) Set(value string) error {
	if *s == nil {
		*s = make([]string, 0, 1)
	}

	*s = append(*s, value)
	return nil
}

// ToSet creates a set from s.
func ToSet(s []string) mapset.Set {
	set := mapset.NewSet()
	for _, allow := range s {
		set.Add(allow)
	}
	return set
}

var (
	flags *flag.FlagSet

	flagListen                string
	flagToConsul              bool
	flagToK8S                 bool
	flagConsulDomain          string
	flagConsulK8STag          string
	flagConsulNodeName        string
	flagK8SDefault            bool
	flagK8SServicePrefix      string
	flagConsulServicePrefix   string
	flagK8SSourceNamespace    string
	flagK8SWriteNamespace     string
	flagConsulWritePeriod     time.Duration
	flagSyncClusterIPServices bool
	flagSyncLBEndpoints       bool
	flagNodePortSyncType      string
	flagAddK8SNamespaceSuffix bool
	flagLogLevel              string
	flagLogJSON               bool

	// Flags to support namespaces
	flagEnableNamespaces           bool     // Use namespacing on all components
	flagConsulDestinationNamespace string   // Consul namespace to register everything if not mirroring
	flagAllowK8sNamespacesList     []string // K8s namespaces to explicitly inject
	flagDenyK8sNamespacesList      []string // K8s namespaces to deny injection (has precedence)
	flagEnableK8SNSMirroring       bool     // Enables mirroring of k8s namespaces into Consul
	flagK8SNSMirroringPrefix       string   // Prefix added to Consul namespaces created when mirroring
	flagCrossNamespaceACLPolicy    string   // The name of the ACL policy to add to every created namespace if ACLs are enabled

	// Flags to support Kubernetes Ingress resources
	flagEnableIngress   bool // Register services using the hostname from an ingress resource
	flagLoadBalancerIPs bool // Use the load balancer IP of an ingress resource instead of the hostname

	logger hclog.Logger

	sigCh chan os.Signal

	// ready indicates whether this controller is ready to sync services. This will be changed to true once the
	// consul-server-connection-manager has finished initial initialization.
	ready bool
)

func initA() {
	flags = flag.NewFlagSet("", flag.ContinueOnError)
	flags.StringVar(&flagListen, "listen", ":8080", "Address to bind listener to.")
	flags.BoolVar(&flagToConsul, "to-consul", true,
		"If true, K8S services will be synced to Consul.")
	flags.BoolVar(&flagToK8S, "to-k8s", true,
		"If true, Consul services will be synced to Kubernetes.")
	flags.BoolVar(&flagK8SDefault, "k8s-default-sync", true,
		"If true, all valid services in K8S are synced by default. If false, "+
			"the service must be annotated properly to sync. In either case "+
			"an annotation can override the default")
	flags.StringVar(&flagK8SServicePrefix, "k8s-service-prefix", "",
		"A prefix to prepend to all services written to Kubernetes from Consul. "+
			"If this is not set then services will have no prefix.")
	flags.StringVar(&flagConsulServicePrefix, "consul-service-prefix", "",
		"A prefix to prepend to all services written to Consul from Kubernetes. "+
			"If this is not set then services will have no prefix.")
	flags.StringVar(&flagK8SSourceNamespace, "k8s-source-namespace", metav1.NamespaceAll,
		"The Kubernetes namespace to watch for service changes and sync to Consul. "+
			"If this is not set then it will default to all namespaces.")
	flags.StringVar(&flagK8SWriteNamespace, "k8s-write-namespace", metav1.NamespaceDefault,
		"The Kubernetes namespace to write to for services from Consul. "+
			"If this is not set then it will default to the default namespace.")
	flags.StringVar(&flagConsulDomain, "consul-domain", "consul",
		"The domain for Consul services to use when writing services to "+
			"Kubernetes. Defaults to consul.")
	flags.StringVar(&flagConsulK8STag, "consul-k8s-tag", "k8s",
		"Tag value for K8S services registered in Consul")
	flags.StringVar(&flagConsulNodeName, "consul-node-name", "k8s-sync",
		"The Consul node name to register for catalog sync. Defaults to k8s-sync. To be discoverable "+
			"via DNS, the name should only contain alpha-numerics and dashes.")
	flags.DurationVar(&flagConsulWritePeriod, "consul-write-interval", 30*time.Second,
		"The interval to perform syncing operations creating Consul services, formatted "+
			"as a time.Duration. All changes are merged and write calls are only made "+
			"on this interval. Defaults to 30 seconds (30s).")
	flags.BoolVar(&flagSyncClusterIPServices, "sync-clusterip-services", true,
		"If true, all valid ClusterIP services in K8S are synced by default. If false, "+
			"ClusterIP services are not synced to Consul.")
	flags.BoolVar(&flagSyncLBEndpoints, "sync-lb-services-endpoints", false,
		"If true, LoadBalancer service endpoints instead of ingress addresses will be synced to Consul. If false, "+
			"LoadBalancer endpoints are not synced to Consul.")
	flags.StringVar(&flagNodePortSyncType, "node-port-sync-type", "ExternalOnly",
		"Defines the type of sync for NodePort services. Valid options are ExternalOnly, "+
			"InternalOnly and ExternalFirst.")
	flags.BoolVar(&flagAddK8SNamespaceSuffix, "add-k8s-namespace-suffix", false,
		"If true, Kubernetes namespace will be appended to service names synced to Consul separated by a dash. "+
			"If false, no suffix will be appended to the service names in Consul. "+
			"If the service name annotation is provided, the suffix is not appended.")
	flags.StringVar(&flagLogLevel, "log-level", "info",
		"Log verbosity level. Supported values (in order of detail) are \"trace\", "+
			"\"debug\", \"info\", \"warn\", and \"error\".")
	flags.BoolVar(&flagLogJSON, "log-json", false,
		"Enable or disable JSON output format for logging.")

	flags.Var((*AppendSliceValue)(&flagAllowK8sNamespacesList), "allow-k8s-namespace",
		"K8s namespaces to explicitly allow. May be specified multiple times.")
	flags.Var((*AppendSliceValue)(&flagDenyK8sNamespacesList), "deny-k8s-namespace",
		"K8s namespaces to explicitly deny. Takes precedence over allow. May be specified multiple times.")
	flags.BoolVar(&flagEnableNamespaces, "enable-namespaces", false,
		"[Enterprise Only] Enables namespaces, in either a single Consul namespace or mirrored.")
	flags.StringVar(&flagConsulDestinationNamespace, "consul-destination-namespace", "default",
		"[Enterprise Only] Defines which Consul namespace to register all synced services into. If '-enable-k8s-namespace-mirroring' "+
			"is true, this is not used.")
	flags.BoolVar(&flagEnableK8SNSMirroring, "enable-k8s-namespace-mirroring", false, "[Enterprise Only] Enables "+
		"namespace mirroring.")
	flags.StringVar(&flagK8SNSMirroringPrefix, "k8s-namespace-mirroring-prefix", "",
		"[Enterprise Only] Prefix that will be added to all k8s namespaces mirrored into Consul if mirroring is enabled.")
	flags.StringVar(&flagCrossNamespaceACLPolicy, "consul-cross-namespace-acl-policy", "",
		"[Enterprise Only] Name of the ACL policy to attach to all created Consul namespaces to allow service "+
			"discovery across Consul namespaces. Only necessary if ACLs are enabled.")

	flags.BoolVar(&flagEnableIngress, "enable-ingress", false,
		"[Enterprise Only] Enables namespaces, in either a single Consul namespace or mirrored.")
	flags.BoolVar(&flagLoadBalancerIPs, "loadBalancer-ips", false,
		"[Enterprise Only] Enables namespaces, in either a single Consul namespace or mirrored.")

	// Set up logging
	if logger == nil {
		parsedLevel := hclog.LevelFromString(flagLogLevel)

		logger = hclog.New(&hclog.LoggerOptions{
			JSONFormat: flagLogJSON,
			Level:      parsedLevel,
			Output:     os.Stderr,
		}).Named("")
	}

	sigCh = make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
}

func main() {
	initA()

	// Convert allow/deny lists to sets
	allowSet := ToSet(flagAllowK8sNamespacesList)
	denySet := ToSet(flagDenyK8sNamespacesList)
	if flagK8SSourceNamespace != "" {
		// For backwards compatibility, if `flagK8SSourceNamespace` is set,
		// it will be the only allowed namespace
		allowSet = mapset.NewSet(flagK8SSourceNamespace)
	}

	allowSet.Add("httpbin")
	ready = true

	// Create the context we'll use to cancel everything
	ctx, cancelF := context.WithCancel(context.Background())

	syncer := &catalog.ConsulSyncer{
		Log:                     logger.Named("to-consul/sink"),
		EnableNamespaces:        flagEnableNamespaces,
		CrossNamespaceACLPolicy: flagCrossNamespaceACLPolicy,
		SyncPeriod:              flagConsulWritePeriod,
		ServicePollPeriod:       flagConsulWritePeriod * 2,
		ConsulK8STag:            flagConsulK8STag,
		ConsulNodeName:          flagConsulNodeName,
	}
	go syncer.Run(ctx)

	clientset := catalog.GetKubeClient()

	flagSyncClusterIPServices = true
	flagK8SDefault = true

	serviceResource := catalog.ServiceResource{
		Log:                        logger.Named("to-consul/source"),
		Client:                     clientset,
		Syncer:                     syncer,
		Ctx:                        ctx,
		AllowK8sNamespacesSet:      allowSet,
		DenyK8sNamespacesSet:       denySet,
		ExplicitEnable:             !flagK8SDefault,
		ClusterIPSync:              flagSyncClusterIPServices,
		LoadBalancerEndpointsSync:  flagSyncLBEndpoints,
		NodePortSync:               catalog.NodePortSyncType(flagNodePortSyncType),
		ConsulK8STag:               flagConsulK8STag,
		ConsulServicePrefix:        flagConsulServicePrefix,
		AddK8SNamespaceSuffix:      flagAddK8SNamespaceSuffix,
		EnableNamespaces:           flagEnableNamespaces,
		ConsulDestinationNamespace: flagConsulDestinationNamespace,
		EnableK8SNSMirroring:       flagEnableK8SNSMirroring,
		K8SNSMirroringPrefix:       flagK8SNSMirroringPrefix,
		ConsulNodeName:             flagConsulNodeName,
		EnableIngress:              flagEnableIngress,
		SyncLoadBalancerIPs:        flagLoadBalancerIPs,
	}

	// Build the controller and start it
	ctl := &controller.Controller{
		Log:      logger.Named("to-consul/controller"),
		Resource: &serviceResource,
	}
	ctl.Log.SetLevel(hclog.Debug)

	// Start the K8S-to-Consul syncer
	var toConsulCh chan struct{}

	toConsulCh = make(chan struct{})
	go func() {
		defer close(toConsulCh)
		ctl.Run(ctx.Done())
	}()

	// Start healthcheck handler
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health/ready", handleReady)
		var handler http.Handler = mux

		if err := http.ListenAndServe(flagListen, handler); err != nil {
			fmt.Printf("Error listening: %s\n", err)
		}
	}()

	select {
	// Unexpected exit
	case <-toConsulCh:
		cancelF()
		return

	// Interrupted/terminated, gracefully exit
	case sig := <-sigCh:
		logger.Info(fmt.Sprintf("%s received, shutting down", sig))
		cancelF()
		if toConsulCh != nil {
			<-toConsulCh
		}
		return
	}
}

func handleReady(rw http.ResponseWriter, r *http.Request) {
	if !ready {
		rw.WriteHeader(500)
		return
	}
	rw.WriteHeader(204)
}
