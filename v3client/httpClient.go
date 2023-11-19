package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
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
	AccessToken(context context.Context) (string, error)
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

	Pipeline []transport.ChainedMiddlewareFunc
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

	if len(p.Pipeline) == 0 {
		p.Pipeline = []transport.ChainedMiddlewareFunc{
			transport.ThrottleFunc,
			transport.BreakOnDeveloperOverRateFunc,
			transport.BackOffOnDeveloperOverQPSFunc,
			transport.ErrorOn404Func,
			transport.RetryOn400Func,
			transport.EnsureBodyWasRead,
			transport.UnmarshalServerError,
		}
	}
}

func NewHttpClient(p Params) Client {
	return newHttpClientWithSchema(p, StandardClientMethodSchema())
}

func NewHttpClientWithBadRequestAutoRetries(p Params) Client {
	return newHttpClientWithSchema(p, AutoRetryOnBadRequestMethodSchema())
}

func newHttpClientWithSchema(p Params, s *ClientMethodSchema) Client {
	p.FillDefaults()

	impl := createHTTPTransport(p)

	// TODO Refactor this to allow accepting the mocks on the HTTP transport.
	rv := FixedSchemeClient{
		PluggableClient{
			schema:    s,
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

		HttpExecutor: p.CreateHttpExecutor(),

		ExchangeListener: p.ExchangeListener,
		Pipeline:         p.Pipeline,
	}
}

func AutoRetryOnBadRequestMethodSchema() *ClientMethodSchema {
	rv := StandardClientMethodSchema()

	// Application methods that needs supporting
	rv.GetApplicationContext = autoRetryBadGetRequest(rv.GetApplicationContext)
	rv.GetApplicationPackageKeys = autoRetryBadRequest(rv.GetApplicationPackageKeys)
	rv.CountApplicationPackageKeys = autoRetryBadRequest(rv.CountApplicationPackageKeys)
	rv.GetFullApplication = autoRetryBadGetRequest(rv.GetFullApplication)
	rv.UpdateApplication = autoRetryBadRequest(rv.UpdateApplication)

	rv.CreateEndpoint = autoRetryBadBiRequest(rv.CreateEndpoint)
	rv.UpdateEndpoint = autoRetryBadRequest(rv.UpdateEndpoint)
	rv.GetEndpoint = autoRetryBadGetRequest(rv.GetEndpoint)

	// Endpoint methods
	rv.ListEndpointMethods = autoRetryBadRequest(rv.ListEndpointMethods)
	rv.CreateEndpointMethod = autoRetryBiBadRequest(rv.CreateEndpointMethod)
	rv.UpdateEndpointMethod = autoRetryBadRequest(rv.UpdateEndpointMethod)
	rv.GetEndpointMethod = autoRetryBadGetRequest(rv.GetEndpointMethod)

	// Endpoint method filters
	rv.ListEndpointMethodFilters = autoRetryBadRequest(rv.ListEndpointMethodFilters)
	rv.CreateEndpointMethodFilter = autoRetryBiBadRequest(rv.CreateEndpointMethodFilter)
	rv.UpdateEndpointMethodFilter = autoRetryBadRequest(rv.UpdateEndpointMethodFilter)
	rv.GetEndpointMethodFilter = autoRetryBadGetRequest(rv.GetEndpointMethodFilter)
	rv.CountEndpointsMethodsFiltersOf = autoRetryBadRequest(rv.CountEndpointsMethodsFiltersOf)

	// Package plans
	rv.CreatePlanService = autoRetryBadRequest(rv.CreatePlanService)
	rv.CheckPlanServiceExists = autoRetryBadRequest(rv.CheckPlanServiceExists)
	rv.CreatePlanEndpoint = autoRetryBadRequest(rv.CreatePlanEndpoint)
	rv.CheckPlanEndpointExists = autoRetryBadRequest(rv.CheckPlanEndpointExists)
	rv.ListPlanEndpoints = autoRetryBadRequest(rv.ListPlanEndpoints)

	rv.CountPlanEndpoints = autoRetryBadRequest(rv.CountPlanEndpoints)
	rv.CountPlanService = autoRetryBadRequest(rv.CountPlanService)
	rv.GetPlan = autoRetryBadGetRequest(rv.GetPlan)
	rv.CreatePlan = autoRetryBiBadRequest(rv.CreatePlan)
	rv.UpdatePlan = autoRetryBadRequest(rv.UpdatePlan)
	rv.CountPlans = autoRetryBadRequest(rv.CountPlans)
	rv.ListPlans = autoRetryBadRequest(rv.ListPlans)
	rv.ListPlanServices = autoRetryBadRequest(rv.ListPlanServices)

	// Plan methods
	rv.ListPackagePlanMethods = autoRetryBadRequest(rv.ListPackagePlanMethods)
	rv.GetPackagePlanMethod = autoRetryBadGetRequest(rv.GetPackagePlanMethod)
	rv.CreatePackagePlanMethod = autoRetryBadRequest(rv.CreatePackagePlanMethod)

	// Plan method filter
	rv.GetPackagePlanMethodFilter = autoRetryBadGetRequest(rv.GetPackagePlanMethodFilter)
	rv.CreatePackagePlanMethodFilter = autoRetryBadRequest(rv.CreatePackagePlanMethodFilter)

	return rv
}

func StandardClientMethodSchema() *ClientMethodSchema {
	return &ClientMethodSchema{
		GetPublicDomains: RootFetcher[int, masherytypes.DomainAddress](publicDomainsCRUD.FetchAll, 0),
		GetSystemDomains: RootFetcher[int, masherytypes.DomainAddress](systemDomainsCRUD.FetchAll, 0),

		// Application method schema
		GetApplicationContext:       applicationCRUD.Get,
		GetApplicationPackageKeys:   applicationPackageKeyCRUD.FetchAll,
		CountApplicationPackageKeys: applicationPackageKeyCRUD.Count,
		GetFullApplication:          GetWithFields[masherytypes.ApplicationIdentifier, masherytypes.Application](applicationDeepFields, applicationCRUD.Get),
		CreateApplication:           applicationCRUD.Create,
		UpdateApplication:           applicationCRUD.Update,
		DeleteApplication:           applicationCRUD.Delete,
		CountApplicationsOfMember:   applicationCRUD.Count,
		ListApplications: func(ctx context.Context, c *transport.HttpTransport) ([]masherytypes.Application, error) {
			return applicationCRUD.FetchAll(ctx, masherytypes.MemberIdentifier{}, c)
		},
		ListApplicationsFiltered: func(ctx context.Context, params map[string]string, c *transport.HttpTransport) ([]masherytypes.Application, error) {
			return applicationCRUD.FetchFiltered(ctx, masherytypes.MemberIdentifier{}, params, c)
		},

		// Email sets
		GetEmailTemplateSet:           emailTemplateSetCRUD.Get,
		ListEmailTemplateSets:         RootFetcher[int, masherytypes.EmailTemplateSet](emailTemplateSetCRUD.FetchAll, 0),
		ListEmailTemplateSetsFiltered: RootFilteredFetcher[int, masherytypes.EmailTemplateSet](emailTemplateSetCRUD.FetchFiltered, 0),

		// Endpoints
		ListEndpoints: endpointCRUD.FetchAllAsAddressable,

		ListEndpointsWithFullInfo: endpointCRUD.FetchAll,
		CreateEndpoint:            endpointCRUD.Create,
		UpdateEndpoint:            endpointCRUD.Update,
		GetEndpoint:               endpointCRUD.Get,
		DeleteEndpoint:            endpointCRUD.Delete,
		CountEndpointsOf:          endpointCRUD.Count,

		// Endpoint methods
		ListEndpointMethods: endpointMethodCRUD.FetchAllAsAddressable,

		ListEndpointMethodsWithFullInfo: endpointMethodCRUD.FetchAll,
		CreateEndpointMethod:            endpointMethodCRUD.Create,
		UpdateEndpointMethod:            endpointMethodCRUD.Update,
		GetEndpointMethod:               endpointMethodCRUD.Get,
		DeleteEndpointMethod:            endpointMethodCRUD.Delete,
		CountEndpointsMethodsOf:         endpointMethodCRUD.Count,

		// Endpoint method filters
		ListEndpointMethodFilters:             endpointMethodFilterCRUD.FetchAllAsAddressable,
		ListEndpointMethodFiltersWithFullInfo: endpointMethodFilterCRUD.FetchAll,
		CreateEndpointMethodFilter:            endpointMethodFilterCRUD.Create,
		UpdateEndpointMethodFilter:            endpointMethodFilterCRUD.Update,
		GetEndpointMethodFilter:               endpointMethodFilterCRUD.Get,
		DeleteEndpointMethodFilter:            endpointMethodFilterCRUD.Delete,
		CountEndpointsMethodsFiltersOf:        endpointMethodFilterCRUD.Count,

		// Member
		GetMember:     memberCRUD.Get,
		GetFullMember: GetWithFields(memberDeepFields, memberCRUD.Get),
		CreateMember:  RootCreator(memberCRUD.Create, 0),
		UpdateMember:  memberCRUD.Update,
		DeleteMember:  memberCRUD.Delete,
		ListMembers:   RootFetcher(memberCRUD.FetchAll, 0),

		// Packages
		GetPackage:            packageCRUD.Get,
		CreatePackage:         RootCreator(packageCRUD.Create, 0),
		UpdatePackage:         packageCRUD.Update,
		DeletePackage:         packageCRUD.Delete,
		ListPackages:          RootFetcher(packageCRUD.FetchAll, 0),
		ResetPackageOwnership: ResetPackageOwnership,

		// Package plans
		CreatePlanService:       CreatePlanService,
		CheckPlanServiceExists:  CheckPlanServiceExists,
		DeletePlanService:       DeletePlanService,
		CreatePlanEndpoint:      CreatePlanEndpoint,
		CheckPlanEndpointExists: CheckPlanEndpointExists,
		DeletePlanEndpoint:      DeletePlanEndpoint,
		ListPlanEndpoints:       ListPlanEndpoints,

		GetPlan:    packagePlanCRDU.Get,
		CreatePlan: packagePlanCRDU.Create,
		UpdatePlan: packagePlanCRDU.Update,
		DeletePlan: packagePlanCRDU.Delete,
		CountPlans: packagePlanCRDU.Count,
		ListPlans:  packagePlanCRDU.FetchAll,

		CountPlanService:   CountPlanService,
		CountPlanEndpoints: CountPlanEndpoints,
		ListPlanServices:   ListPlanServices,

		// Plan methods
		ListPackagePlanMethods:  ListPackagePlanMethods,
		GetPackagePlanMethod:    GetPackagePlanMethod,
		CreatePackagePlanMethod: CreatePackagePlanServiceEndpointMethod,
		DeletePackagePlanMethod: DeletePackagePlanMethod,

		// Plan method filter
		GetPackagePlanMethodFilter:    packagePlanServiceEndpointMethodFilterCRUD.Get,
		CreatePackagePlanMethodFilter: CreatePackagePlanMethodFilter,
		DeletePackagePlanMethodFilter: packagePlanServiceEndpointMethodFilterCRUD.Delete,

		// Package key
		GetPackageKey:    packageKeyCRUD.Get,
		UpdatePackageKey: packageKeyCRUD.Update,

		CreatePackageKey: CreatePackageKey,

		DeletePackageKey:        packageKeyCRUD.Delete,
		ListPackageKeysFiltered: RootFilteredFetcher(packageKeyCRUD.FetchFiltered, 0),
		ListPackageKeys:         RootFetcher(packageKeyCRUD.FetchAll, 0),

		// Roles
		GetRole:           roleCRUD.Get,
		ListRoles:         RootFetcher(roleCRUD.FetchAll, 0),
		ListRolesFiltered: RootFilteredFetcher(roleCRUD.FetchFiltered, 0),

		// Service
		GetService:           serviceCRUD.Get,
		CreateService:        RootCreator(serviceCRUD.Create, 0),
		UpdateService:        serviceCRUD.Update,
		DeleteService:        serviceCRUD.Delete,
		ListServicesFiltered: RootFilteredFetcher(serviceCRUD.FetchFiltered, 0),
		ListServices:         RootFetcher(serviceCRUD.FetchAll, 0),
		CountServices:        RootFilteredCounter[int, masherytypes.Service](serviceCRUD.CountFiltered, 0),

		ListErrorSets:         errorSetCRUD.FetchFiltered,
		GetErrorSet:           errorSetCRUD.Get,
		CreateErrorSet:        errorSetCRUD.Create,
		UpdateErrorSet:        errorSetCRUD.Update,
		DeleteErrorSet:        errorSetCRUD.Delete,
		UpdateErrorSetMessage: errorSetMessageCRUD.Update,

		GetServiceRoles:    GetServiceRoles,
		SetServiceRoles:    SetServiceRoles,
		DeleteServiceRoles: DeleteServiceRoles,

		// Service cache,
		GetServiceCache:    serviceCacheCRUD.Get,
		CreateServiceCache: serviceCacheCRUD.Create,
		UpdateServiceCache: serviceCacheCRUD.Update,
		DeleteServiceCache: serviceCacheCRUD.Delete,

		// Service OAuth
		GetServiceOAuthSecurityProfile:    serviceOAuthCRUD.Get,
		CreateServiceOAuthSecurityProfile: serviceOAuthCRUD.Create,
		UpdateServiceOAuthSecurityProfile: serviceOAuthCRUD.Update,
		DeleteServiceOAuthSecurityProfile: serviceOAuthCRUD.Delete,

		// List organizations
		ListOrganizations:         RootFetcher(organizationCRUD.FetchAll, 0),
		ListOrganizationsFiltered: RootFilteredFetcher(organizationCRUD.FetchFiltered, 0),
	}
}

type AsyncFetchResult struct {
	Data *http.Response
	Err  error
}
