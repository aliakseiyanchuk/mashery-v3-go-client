package v3client

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/http"
	"sync"
	"time"
)

const tokenFile string = ".mashery-logon"

// V3AccessTokenProvider Access token provider that supplies the access token, depending on the strategy.
// There are three strategies:
// - FixedTokenProvider yields a fixed token. This method is useful for short deployments where an access
// token is obtained by an outside process and would be stored e.g. in-memory.
// - FileSystemTokenProvider yields a token that that was previously saved in the file system, e.g. using the `mash-login`
// command
// - Both these methods have limited applicability time-span of 1 hour, since Mashery V3 token would expire after 1
// hour, and repeated logon would be necessary.
// - ClientCredentialsProvider can support operations of exceeding 1 hour by using Mashery V3 API to retrieve and refresh
// the access token.
//
// The calling code has to pick an appropriate provider depending on the context.
type V3AccessTokenProvider interface {
	transport.Authorizer

	// AccessToken Yields an access token to be used in the next API call to Mashery
	AccessToken() (string, error)
}

func NewCustomClient(schema *ClientMethodSchema) Client {
	rv := FixedSchemeClient{
		PluggableClient{
			schema:    schema,
			transport: nil,
		},
	}
	return &rv
}

type Params struct {
	transport.HTTPClientParams

	MashEndpoint  string
	Authorizer    transport.Authorizer
	QPS           int64
	AvgNetLatency time.Duration
}

func (p *Params) FillDefaults() {
	if len(p.MashEndpoint) == 0 {
		p.MashEndpoint = "https://api.mashery.com/v3/rest"
	}

	if p.QPS <= 0 {
		p.QPS = 2
	}
	if p.Timeout <= 0 {
		p.Timeout = time.Minute * 2
	}
	if p.AvgNetLatency <= 0 {
		p.AvgNetLatency = time.Millisecond * 147
	}

	// Default is set if no TLS configuration is supplied, and where no explicit flag is set
	// to delegate the trust to the system settings.
	if p.TLSConfig == nil && !p.TLSConfigDelegateSystem {
		p.TLSConfig = transport.DefaultTLSConfig()
	}
}

func NewHttpClient(p Params) Client {
	p.FillDefaults()

	impl := transport.V3Transport{
		HttpTransport: createHTTPTransport(p),
	}

	rv := FixedSchemeClient{
		PluggableClient{
			schema:    StandardClientMethodSchema(),
			transport: &impl,
		},
	}

	return &rv
}

func createHTTPTransport(p Params) transport.HttpTransport {
	return transport.HttpTransport{
		MashEndpoint:  p.MashEndpoint,
		Authorizer:    p.Authorizer,
		AvgNetLatency: p.AvgNetLatency,
		Mutex:         &sync.Mutex{},
		MaxQPS:        p.QPS,

		HttpClient: p.CreateClient(),

		ExchangeListener: p.ExchangeListener,
	}
}

func StandardClientMethodSchema() *ClientMethodSchema {
	return &ClientMethodSchema{
		GetPublicDomains: ListPublicDomains,
		GetSystemDomains: ListSystemDomains,

		// Application method schema
		GetApplicationContext:       GetApplication,
		GetApplicationPackageKeys:   GetApplicationPackageKeys,
		CountApplicationPackageKeys: CountApplicationPackageKeys,
		GetFullApplication:          GetFullApplication,
		CreateApplication:           CreateApplication,
		UpdateApplication:           UpdateApplication,
		DeleteApplication:           DeleteApplication,
		CountApplicationsOfMember:   CountApplicationsOfMember,
		ListApplications:            ListApplications,

		// Email sets
		GetEmailTemplateSet:           GetEmailTemplateSet,
		ListEmailTemplateSets:         ListEmailTemplateSets,
		ListEmailTemplateSetsFiltered: ListEmailTemplateSetsFiltered,

		// Endpoints
		ListEndpoints:             ListEndpoints,
		ListEndpointsWithFullInfo: ListEndpointsWithFullInfo,
		CreateEndpoint:            CreateEndpoint,
		UpdateEndpoint:            UpdateEndpoint,
		GetEndpoint:               GetEndpoint,
		DeleteEndpoint:            DeleteEndpoint,
		CountEndpointsOf:          CountEndpointsOf,

		// Endpoint methods
		ListEndpointMethods:             ListEndpointMethods,
		ListEndpointMethodsWithFullInfo: ListEndpointMethodsWithFullInfo,
		CreateEndpointMethod:            CreateEndpointMethod,
		UpdateEndpointMethod:            UpdateEndpointMethod,
		GetEndpointMethod:               GetEndpointMethod,
		DeleteEndpointMethod:            DeleteEndpointMethod,
		CountEndpointsMethodsOf:         CountEndpointsMethodsOf,

		// Endpoint method filters
		ListEndpointMethodFilters:             ListEndpointMethodFilters,
		ListEndpointMethodFiltersWithFullInfo: ListEndpointMethodFiltersWithFullInfo,
		CreateEndpointMethodFilter:            CreateEndpointMethodFilter,
		UpdateEndpointMethodFilter:            UpdateEndpointMethodFilter,
		GetEndpointMethodFilter:               GetEndpointMethodFilter,
		DeleteEndpointMethodFilter:            DeleteEndpointMethodFilter,
		CountEndpointsMethodsFiltersOf:        CountEndpointsMethodsFiltersOf,

		// Member
		GetMember:     GetMember,
		GetFullMember: GetFullMember,
		CreateMember:  CreateMember,
		UpdateMember:  UpdateMember,
		DeleteMember:  DeleteMember,
		ListMembers:   ListMembers,

		// Packages
		GetPackage:    GetPackage,
		CreatePackage: CreatePackage,
		UpdatePackage: UpdatePackage,
		DeletePackage: DeletePackage,
		ListPackages:  ListPackages,

		// Package plans
		CreatePlanService:  CreatePlanService,
		DeletePlanService:  DeletePlanService,
		CreatePlanEndpoint: CreatePlanEndpoint,
		DeletePlanEndpoint: DeletePlanEndpoint,
		ListPlanEndpoints:  ListPlanEndpoints,

		CountPlanEndpoints: CountPlanEndpoints,
		CountPlanService:   CountPlanService,
		GetPlan:            GetPlan,
		CreatePlan:         CreatePlan,
		UpdatePlan:         UpdatePlan,
		DeletePlan:         DeletePlan,
		CountPlans:         CountPlans,
		ListPlans:          ListPlans,
		ListPlanServices:   ListPlanServices,

		// Plan methods
		ListPackagePlanMethods:  ListPackagePlanMethods,
		GetPackagePlanMethod:    GetPackagePlanMethod,
		CreatePackagePlanMethod: CreatePackagePlanServiceEndpointMethod,
		DeletePackagePlanMethod: DeletePackagePlanMethod,

		// Plan method filter
		GetPackagePlanMethodFilter:    GetPackagePlanMethodFilter,
		CreatePackagePlanMethodFilter: CreatePackagePlanMethodFilter,
		DeletePackagePlanMethodFilter: DeletePackagePlanMethodFilter,

		// Package key
		GetPackageKey:           GetPackageKey,
		CreatePackageKey:        CreatePackageKey,
		UpdatePackageKey:        UpdatePackageKey,
		DeletePackageKey:        DeletePackageKey,
		ListPackageKeysFiltered: ListPackageKeysFiltered,
		ListPackageKeys:         ListPackageKeys,

		// Roles
		GetRole:           GetRole,
		ListRoles:         ListRoles,
		ListRolesFiltered: ListRolesFiltered,

		// Service
		GetService:           GetService,
		CreateService:        CreateService,
		UpdateService:        UpdateService,
		DeleteService:        DeleteService,
		ListServicesFiltered: ListServicesFiltered,
		ListServices:         ListServices,
		CountServices:        CountServices,

		ListErrorSets:         ListErrorSets,
		GetErrorSet:           GetErrorSet,
		CreateErrorSet:        CreateErrorSet,
		UpdateErrorSet:        UpdateErrorSet,
		DeleteErrorSet:        DeleteErrorSet,
		UpdateErrorSetMessage: UpdateErrorSetMessage,

		GetServiceRoles:    GetServiceRoles,
		SetServiceRoles:    SetServiceRoles,
		DeleteServiceRoles: DeleteServiceRoles,

		// Service cache,
		GetServiceCache:    GetServiceCache,
		CreateServiceCache: CreateServiceCache,
		UpdateServiceCache: UpdateServiceCache,
		DeleteServiceCache: DeleteServiceCache,

		// Service OAuth
		GetServiceOAuthSecurityProfile:    GetServiceOAuthSecurityProfile,
		CreateServiceOAuthSecurityProfile: CreateServiceOAuthSecurityProfile,
		UpdateServiceOAuthSecurityProfile: UpdateServiceOAuthSecurityProfile,
		DeleteServiceOAuthSecurityProfile: DeleteServiceOAuthSecurityProfile,
	}
}

type AsyncFetchResult struct {
	Data *http.Response
	Err  error
}
