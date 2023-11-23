package client

import (
	"context"
	"time"
)

// Client provides a client to the Consul API
type Client struct {
}

// Catalog can be used to query the Catalog endpoints
type Catalog struct {
	c *Client
}

// QueryOptions are used to parameterize a query
type QueryOptions struct {
	// Namespace overrides the `default` namespace
	// Note: Namespaces are available only in Consul Enterprise
	Namespace string

	// Partition overrides the `default` partition
	// Note: Partitions are available only in Consul Enterprise
	Partition string

	// Providing a datacenter overwrites the DC provided
	// by the Config
	Datacenter string

	// Providing a peer name in the query option
	Peer string

	// AllowStale allows any Consul server (non-leader) to service
	// a read. This allows for lower latency and higher throughput
	AllowStale bool

	// RequireConsistent forces the read to be fully consistent.
	// This is more expensive but prevents ever performing a stale
	// read.
	RequireConsistent bool

	// UseCache requests that the agent cache results locally. See
	// https://www.consul.io/api/features/caching.html for more details on the
	// semantics.
	UseCache bool

	// MaxAge limits how old a cached value will be returned if UseCache is true.
	// If there is a cached response that is older than the MaxAge, it is treated
	// as a cache miss and a new fetch invoked. If the fetch fails, the error is
	// returned. Clients that wish to allow for stale results on error can set
	// StaleIfError to a longer duration to change this behavior. It is ignored
	// if the endpoint supports background refresh caching. See
	// https://www.consul.io/api/features/caching.html for more details.
	MaxAge time.Duration

	// StaleIfError specifies how stale the client will accept a cached response
	// if the servers are unavailable to fetch a fresh one. Only makes sense when
	// UseCache is true and MaxAge is set to a lower, non-zero value. It is
	// ignored if the endpoint supports background refresh caching. See
	// https://www.consul.io/api/features/caching.html for more details.
	StaleIfError time.Duration

	// WaitIndex is used to enable a blocking query. Waits
	// until the timeout or the next index is reached
	WaitIndex uint64

	// WaitHash is used by some endpoints instead of WaitIndex to perform blocking
	// on state based on a hash of the response rather than a monotonic index.
	// This is required when the state being blocked on is not stored in Raft, for
	// example agent-local proxy configuration.
	WaitHash string

	// WaitTime is used to bound the duration of a wait.
	// Defaults to that of the Config, but can be overridden.
	WaitTime time.Duration

	// Token is used to provide a per-request ACL token
	// which overrides the agent's default token.
	Token string

	// Near is used to provide a node name that will sort the results
	// in ascending order based on the estimated round trip time from
	// that node. Setting this to "_agent" will use the agent's node
	// for the sort.
	Near string

	// NodeMeta is used to filter results by nodes with the given
	// metadata key/value pairs. Currently, only one key/value pair can
	// be provided for filtering.
	NodeMeta map[string]string

	// RelayFactor is used in keyring operations to cause responses to be
	// relayed back to the sender through N other random nodes. Must be
	// a value from 0 to 5 (inclusive).
	RelayFactor uint8

	// LocalOnly is used in keyring list operation to force the keyring
	// query to only hit local servers (no WAN traffic).
	LocalOnly bool

	// Connect filters prepared query execution to only include Connect-capable
	// services. This currently affects prepared query execution.
	Connect bool

	// ctx is an optional context pass through to the underlying HTTP
	// request layer. Use Context() and WithContext() to manage this.
	ctx context.Context

	// Filter requests filtering data prior to it being returned. The string
	// is a go-bexpr compatible expression.
	Filter string

	// MergeCentralConfig returns a service definition merged with the
	// proxy-defaults/global and service-defaults/:service config entries.
	// This can be used to ensure a full service definition is returned in the response
	// especially when the service might not be written into the catalog that way.
	MergeCentralConfig bool

	// Global is used to request information from all datacenters. Currently only
	// used for operator usage requests.
	Global bool
}

// QueryMeta is used to return meta data about a query
type QueryMeta struct {
	// LastIndex. This can be used as a WaitIndex to perform
	// a blocking query
	LastIndex uint64

	// LastContentHash. This can be used as a WaitHash to perform a blocking query
	// for endpoints that support hash-based blocking. Endpoints that do not
	// support it will return an empty hash.
	LastContentHash string

	// Time of last contact from the leader for the
	// server servicing the request
	LastContact time.Duration

	// Is there a known leader
	KnownLeader bool

	// How long did the request take
	RequestTime time.Duration

	// Is address translation enabled for HTTP responses on this agent
	AddressTranslationEnabled bool

	// CacheHit is true if the result was served from agent-local cache.
	CacheHit bool

	// CacheAge is set if request was ?cached and indicates how stale the cached
	// response is.
	CacheAge time.Duration

	// QueryBackend represent which backend served the request.
	QueryBackend string

	// DefaultACLPolicy is used to control the ACL interaction when there is no
	// defined policy. This can be "allow" which means ACLs are used to
	// deny-list, or "deny" which means ACLs are allow-lists.
	DefaultACLPolicy string

	// ResultsFilteredByACLs is true when some of the query's results were
	// filtered out by enforcing ACLs. It may be false because nothing was
	// removed, or because the endpoint does not yet support this flag.
	ResultsFilteredByACLs bool
}
