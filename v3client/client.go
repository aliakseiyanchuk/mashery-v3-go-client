package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/http"
	"net/url"
)

type WildcardClient interface {
	// FetchAny fetch an arbitrary resource from Mashery V3
	FetchAny(ctx context.Context, resource string, qs *url.Values) (*http.Response, error)
	// DeleteAny Delete an arbitrary resource from Mashery V3
	DeleteAny(ctx context.Context, resource string) (*http.Response, error)
	// PostAny post any value to an arbitrary resource
	PostAny(ctx context.Context, resource string, body interface{}) (*http.Response, error)
	// PutAny put any value to an arbitrary resource
	PutAny(ctx context.Context, resource string, body interface{}) (*http.Response, error)

	Close(ctx context.Context)
}

type Client interface {
	// GetPublicDomains Get public domains in this area
	GetPublicDomains(ctx context.Context) ([]string, error)
	GetSystemDomains(ctx context.Context) ([]string, error)

	// GetApplication Retrieve the details of the application
	GetApplication(ctx context.Context, appId string) (*masherytypes.MasheryApplication, error)
	GetApplicationPackageKeys(ctx context.Context, appId string) ([]masherytypes.MasheryPackageKey, error)
	CountApplicationPackageKeys(ctx context.Context, appId string) (int64, error)
	GetFullApplication(ctx context.Context, id string) (*masherytypes.MasheryApplication, error)
	CreateApplication(ctx context.Context, memberId string, member masherytypes.MasheryApplication) (*masherytypes.MasheryApplication, error)
	UpdateApplication(ctx context.Context, app masherytypes.MasheryApplication) (*masherytypes.MasheryApplication, error)
	DeleteApplication(ctx context.Context, appId string) error
	CountApplicationsOfMember(ctx context.Context, memberId string) (int64, error)
	ListApplications(ctx context.Context) ([]masherytypes.MasheryApplication, error)

	// Email template sets
	GetEmailTemplateSet(ctx context.Context, id string) (*masherytypes.MasheryEmailTemplateSet, error)
	ListEmailTemplateSets(ctx context.Context) ([]masherytypes.MasheryEmailTemplateSet, error)
	ListEmailTemplateSetsFiltered(ctx context.Context, params map[string]string, fields []string) ([]masherytypes.MasheryEmailTemplateSet, error)

	// Endpoints
	ListEndpoints(ctx context.Context, serviceId string) ([]masherytypes.AddressableV3Object, error)
	ListEndpointsWithFullInfo(ctx context.Context, serviceId string) ([]masherytypes.MasheryEndpoint, error)
	CreateEndpoint(ctx context.Context, serviceId string, endp masherytypes.MasheryEndpoint) (*masherytypes.MasheryEndpoint, error)
	UpdateEndpoint(ctx context.Context, serviceId string, endp masherytypes.MasheryEndpoint) (*masherytypes.MasheryEndpoint, error)
	GetEndpoint(ctx context.Context, serviceId string, endpointId string) (*masherytypes.MasheryEndpoint, error)
	DeleteEndpoint(ctx context.Context, serviceId, endpointId string) error
	CountEndpointsOf(ctx context.Context, serviceId string) (int64, error)

	// Endpoint methods
	ListEndpointMethods(ctx context.Context, serviceId, endpointId string) ([]masherytypes.MasheryMethod, error)
	ListEndpointMethodsWithFullInfo(ctx context.Context, serviceId, endpointId string) ([]masherytypes.MasheryMethod, error)
	CreateEndpointMethod(ctx context.Context, serviceId, endpointId string, methoUpsert masherytypes.MasheryMethod) (*masherytypes.MasheryMethod, error)
	UpdateEndpointMethod(ctx context.Context, serviceId, endpointId string, methUpsert masherytypes.MasheryMethod) (*masherytypes.MasheryMethod, error)
	GetEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string) (*masherytypes.MasheryMethod, error)
	DeleteEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string) error
	CountEndpointsMethodsOf(ctx context.Context, serviceId, endpointId string) (int64, error)

	// Endpoint method filters
	ListEndpointMethodFilters(ctx context.Context, serviceId, endpointId, methodId string) ([]masherytypes.MasheryResponseFilter, error)
	ListEndpointMethodFiltersWithFullInfo(ctx context.Context, serviceId, endpointId, methodId string) ([]masherytypes.MasheryResponseFilter, error)
	CreateEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId string, filterUpsert masherytypes.MasheryResponseFilter) (*masherytypes.MasheryResponseFilter, error)
	UpdateEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId string, methUpsert masherytypes.MasheryResponseFilter) (*masherytypes.MasheryResponseFilter, error)
	GetEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId, filterId string) (*masherytypes.MasheryResponseFilter, error)
	DeleteEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId, filterId string) error
	CountEndpointsMethodsFiltersOf(ctx context.Context, serviceId, endpointId, methodId string) (int64, error)

	// Member
	GetMember(ctx context.Context, id string) (*masherytypes.MasheryMember, error)
	GetFullMember(ctx context.Context, id string) (*masherytypes.MasheryMember, error)
	CreateMember(ctx context.Context, member masherytypes.MasheryMember) (*masherytypes.MasheryMember, error)
	UpdateMember(ctx context.Context, member masherytypes.MasheryMember) (*masherytypes.MasheryMember, error)
	DeleteMember(ctx context.Context, memberId string) error
	ListMembers(ctx context.Context) ([]masherytypes.MasheryMember, error)

	// Packages
	GetPackage(ctx context.Context, id string) (*masherytypes.MasheryPackage, error)
	CreatePackage(ctx context.Context, pack masherytypes.MasheryPackage) (*masherytypes.MasheryPackage, error)
	UpdatePackage(ctx context.Context, pack masherytypes.MasheryPackage) (*masherytypes.MasheryPackage, error)
	DeletePackage(ctx context.Context, packId string) error
	ListPackages(ctx context.Context) ([]masherytypes.MasheryPackage, error)

	// Package plans
	CreatePlanService(ctx context.Context, planService masherytypes.MasheryPlanService) (*masherytypes.AddressableV3Object, error)
	DeletePlanService(ctx context.Context, planService masherytypes.MasheryPlanService) error
	CreatePlanEndpoint(ctx context.Context, planEndp masherytypes.MasheryPlanServiceEndpoint) (*masherytypes.AddressableV3Object, error)
	DeletePlanEndpoint(ctx context.Context, planEndp masherytypes.MasheryPlanServiceEndpoint) error
	ListPlanEndpoints(ctx context.Context, planService masherytypes.MasheryPlanService) ([]masherytypes.AddressableV3Object, error)
	CountPlanService(ctx context.Context, packageId, planId string) (int64, error)
	GetPlan(ctx context.Context, packageId string, planId string) (*masherytypes.MasheryPlan, error)
	CreatePlan(ctx context.Context, packageId string, plan masherytypes.MasheryPlan) (*masherytypes.MasheryPlan, error)
	UpdatePlan(ctx context.Context, plan masherytypes.MasheryPlan) (*masherytypes.MasheryPlan, error)
	DeletePlan(ctx context.Context, packageId, planId string) error
	CountPlans(ctx context.Context, packageId string) (int64, error)
	ListPlans(ctx context.Context, packageId string) ([]masherytypes.MasheryPlan, error)
	ListPlanServices(ctx context.Context, packageId string, planId string) ([]masherytypes.MasheryService, error)

	CountPlanEndpoints(ctx context.Context, planService masherytypes.MasheryPlanService) (int64, error)

	// Plan methods
	ListPackagePlanMethods(ctx context.Context, id masherytypes.MasheryPlanServiceEndpoint) ([]masherytypes.MasheryMethod, error)
	GetPackagePlanMethod(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod) (*masherytypes.MasheryMethod, error)
	CreatePackagePlanMethod(ctx context.Context, id masherytypes.MasheryPlanServiceEndpoint, upsert masherytypes.MasheryMethod) (*masherytypes.MasheryMethod, error)
	DeletePackagePlanMethod(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod) error

	// PLan method filter
	GetPackagePlanMethodFilter(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod) (*masherytypes.MasheryResponseFilter, error)
	CreatePackagePlanMethodFilter(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod, ref masherytypes.MasheryServiceMethodFilter) (*masherytypes.MasheryResponseFilter, error)
	DeletePackagePlanMethodFilter(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod) error

	// Package key
	GetPackageKey(ctx context.Context, id string) (*masherytypes.MasheryPackageKey, error)
	CreatePackageKey(ctx context.Context, appId string, packageKey masherytypes.MasheryPackageKey) (*masherytypes.MasheryPackageKey, error)
	UpdatePackageKey(ctx context.Context, packageKey masherytypes.MasheryPackageKey) (*masherytypes.MasheryPackageKey, error)
	DeletePackageKey(ctx context.Context, keyId string) error
	ListPackageKeysFiltered(ctx context.Context, params map[string]string, fields []string) ([]masherytypes.MasheryPackageKey, error)
	ListPackageKeys(ctx context.Context) ([]masherytypes.MasheryPackageKey, error)

	// Roles
	GetRole(ctx context.Context, id string) (*masherytypes.MasheryRole, error)
	ListRoles(ctx context.Context) ([]masherytypes.MasheryRole, error)
	ListRolesFiltered(ctx context.Context, params map[string]string, fields []string) ([]masherytypes.MasheryRole, error)

	// GetService retrieves service based on the service identifier
	GetService(ctx context.Context, id string) (*masherytypes.MasheryService, error)
	CreateService(ctx context.Context, service masherytypes.MasheryService) (*masherytypes.MasheryService, error)
	UpdateService(ctx context.Context, service masherytypes.MasheryService) (*masherytypes.MasheryService, error)
	DeleteService(ctx context.Context, serviceId string) error
	ListServicesFiltered(ctx context.Context, params map[string]string, fields []string) ([]masherytypes.MasheryService, error)
	ListServices(ctx context.Context) ([]masherytypes.MasheryService, error)
	CountServices(ctx context.Context, params map[string]string) (int64, error)

	ListErrorSets(ctx context.Context, serviceId string, qs url.Values) ([]masherytypes.MasheryErrorSet, error)
	GetErrorSet(ctx context.Context, serviceId, setId string) (*masherytypes.MasheryErrorSet, error)
	CreateErrorSet(ctx context.Context, serviceId string, set masherytypes.MasheryErrorSet) (*masherytypes.MasheryErrorSet, error)
	UpdateErrorSet(ctx context.Context, serviceId string, setData masherytypes.MasheryErrorSet) (*masherytypes.MasheryErrorSet, error)
	DeleteErrorSet(ctx context.Context, serviceId, setId string) error
	UpdateErrorSetMessage(ctx context.Context, serviceId string, setId string, msg masherytypes.MasheryErrorMessage) (*masherytypes.MasheryErrorMessage, error)

	GetServiceRoles(ctx context.Context, serviceId string) ([]masherytypes.MasheryRolePermission, error)
	SetServiceRoles(ctx context.Context, id string, roles []masherytypes.MasheryRolePermission) error

	// Service cache
	GetServiceCache(ctx context.Context, id string) (*masherytypes.MasheryServiceCache, error)
	CreateServiceCache(ctx context.Context, id string, service masherytypes.MasheryServiceCache) (*masherytypes.MasheryServiceCache, error)
	UpdateServiceCache(ctx context.Context, id string, service masherytypes.MasheryServiceCache) (*masherytypes.MasheryServiceCache, error)
	DeleteServiceCache(ctx context.Context, id string) error

	// GetServiceOAuthSecurityProfile Service OAuth
	GetServiceOAuthSecurityProfile(ctx context.Context, id string) (*masherytypes.MasheryOAuth, error)
	CreateServiceOAuthSecurityProfile(ctx context.Context, id string, service masherytypes.MasheryOAuth) (*masherytypes.MasheryOAuth, error)
	UpdateServiceOAuthSecurityProfile(ctx context.Context, id string, service masherytypes.MasheryOAuth) (*masherytypes.MasheryOAuth, error)
	DeleteServiceOAuthSecurityProfile(ctx context.Context, id string) error
}

type PluggableClient struct {
	schema    *ClientMethodSchema
	transport *transport.V3Transport
}

// FixedSchemeClient Fixed method scheme client that will not allow changing schema after it was created.
type FixedSchemeClient struct {
	PluggableClient
}

func (fsc *FixedSchemeClient) AssumeSchema(_ *ClientMethodSchema) {
	panic("This client cannot change method schema after it was created")
}

// AssumeSchema Allows a pluggable client to assume a schema. This method should be mainly used in the test contexts.
func (pc *PluggableClient) AssumeSchema(sch *ClientMethodSchema) {
	pc.schema = sch
}

func (c *PluggableClient) notImplemented(meth string) error {
	return errors.New(fmt.Sprintf("No implementation method was supplied for method %s", meth))
}

// -----------------------------------------------------------------------------------------------------------------
// Application methods
// -----------------------------------------------------------------------------------------------------------------

type ClientMethodSchema struct {
	// Public and System Domains
	GetPublicDomains func(ctx context.Context, transport *transport.V3Transport) ([]string, error)
	GetSystemDomains func(ctx context.Context, transport *transport.V3Transport) ([]string, error)

	// Applications
	GetApplicationContext       func(ctx context.Context, appId string, transport *transport.V3Transport) (*masherytypes.MasheryApplication, error)
	GetApplicationPackageKeys   func(ctx context.Context, appId string, transport *transport.V3Transport) ([]masherytypes.MasheryPackageKey, error)
	CountApplicationPackageKeys func(ctx context.Context, appId string, c *transport.V3Transport) (int64, error)
	GetFullApplication          func(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryApplication, error)
	CreateApplication           func(ctx context.Context, memberId string, member masherytypes.MasheryApplication, c *transport.V3Transport) (*masherytypes.MasheryApplication, error)
	UpdateApplication           func(ctx context.Context, app masherytypes.MasheryApplication, c *transport.V3Transport) (*masherytypes.MasheryApplication, error)
	DeleteApplication           func(ctx context.Context, appId string, c *transport.V3Transport) error
	CountApplicationsOfMember   func(ctx context.Context, memberId string, c *transport.V3Transport) (int64, error)
	ListApplications            func(ctx context.Context, c *transport.V3Transport) ([]masherytypes.MasheryApplication, error)

	// Email template set
	GetEmailTemplateSet           func(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryEmailTemplateSet, error)
	ListEmailTemplateSets         func(ctx context.Context, c *transport.V3Transport) ([]masherytypes.MasheryEmailTemplateSet, error)
	ListEmailTemplateSetsFiltered func(ctx context.Context, params map[string]string, fields []string, c *transport.V3Transport) ([]masherytypes.MasheryEmailTemplateSet, error)

	// Endpoints
	ListEndpoints             func(ctx context.Context, serviceId string, c *transport.V3Transport) ([]masherytypes.AddressableV3Object, error)
	ListEndpointsWithFullInfo func(ctx context.Context, serviceId string, c *transport.V3Transport) ([]masherytypes.MasheryEndpoint, error)
	CreateEndpoint            func(ctx context.Context, serviceId string, endp masherytypes.MasheryEndpoint, c *transport.V3Transport) (*masherytypes.MasheryEndpoint, error)
	UpdateEndpoint            func(ctx context.Context, serviceId string, endp masherytypes.MasheryEndpoint, c *transport.V3Transport) (*masherytypes.MasheryEndpoint, error)
	GetEndpoint               func(ctx context.Context, serviceId string, endpointId string, c *transport.V3Transport) (*masherytypes.MasheryEndpoint, error)
	DeleteEndpoint            func(ctx context.Context, serviceId, endpointId string, c *transport.V3Transport) error
	CountEndpointsOf          func(ctx context.Context, serviceId string, c *transport.V3Transport) (int64, error)

	// Endpoint Methods
	ListEndpointMethods             func(ctx context.Context, serviceId, endpointId string, c *transport.V3Transport) ([]masherytypes.MasheryMethod, error)
	ListEndpointMethodsWithFullInfo func(ctx context.Context, serviceId, endpointId string, c *transport.V3Transport) ([]masherytypes.MasheryMethod, error)
	CreateEndpointMethod            func(ctx context.Context, serviceId, endpointId string, methoUpsert masherytypes.MasheryMethod, c *transport.V3Transport) (*masherytypes.MasheryMethod, error)
	UpdateEndpointMethod            func(ctx context.Context, serviceId, endpointId string, methUpsert masherytypes.MasheryMethod, c *transport.V3Transport) (*masherytypes.MasheryMethod, error)
	GetEndpointMethod               func(ctx context.Context, serviceId, endpointId, methodId string, c *transport.V3Transport) (*masherytypes.MasheryMethod, error)
	DeleteEndpointMethod            func(ctx context.Context, serviceId, endpointId, methodId string, c *transport.V3Transport) error
	CountEndpointsMethodsOf         func(ctx context.Context, serviceId, endpointId string, c *transport.V3Transport) (int64, error)

	// Endpoint method filters
	ListEndpointMethodFilters             func(ctx context.Context, serviceId, endpointId, methodId string, c *transport.V3Transport) ([]masherytypes.MasheryResponseFilter, error)
	ListEndpointMethodFiltersWithFullInfo func(ctx context.Context, serviceId, endpointId, methodId string, c *transport.V3Transport) ([]masherytypes.MasheryResponseFilter, error)
	CreateEndpointMethodFilter            func(ctx context.Context, serviceId, endpointId, methodId string, filterUpsert masherytypes.MasheryResponseFilter, c *transport.V3Transport) (*masherytypes.MasheryResponseFilter, error)
	UpdateEndpointMethodFilter            func(ctx context.Context, serviceId, endpointId, methodId string, methUpsert masherytypes.MasheryResponseFilter, c *transport.V3Transport) (*masherytypes.MasheryResponseFilter, error)
	GetEndpointMethodFilter               func(ctx context.Context, serviceId, endpointId, methodId, filterId string, c *transport.V3Transport) (*masherytypes.MasheryResponseFilter, error)
	DeleteEndpointMethodFilter            func(ctx context.Context, serviceId, endpointId, methodId, filterId string, c *transport.V3Transport) error
	CountEndpointsMethodsFiltersOf        func(ctx context.Context, serviceId, endpointId, methodId string, c *transport.V3Transport) (int64, error)

	// Member
	GetMember     func(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryMember, error)
	GetFullMember func(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryMember, error)
	CreateMember  func(ctx context.Context, member masherytypes.MasheryMember, c *transport.V3Transport) (*masherytypes.MasheryMember, error)
	UpdateMember  func(ctx context.Context, member masherytypes.MasheryMember, c *transport.V3Transport) (*masherytypes.MasheryMember, error)
	DeleteMember  func(ctx context.Context, memberId string, c *transport.V3Transport) error
	ListMembers   func(ctx context.Context, c *transport.V3Transport) ([]masherytypes.MasheryMember, error)

	// Packages
	GetPackage    func(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryPackage, error)
	CreatePackage func(ctx context.Context, pack masherytypes.MasheryPackage, c *transport.V3Transport) (*masherytypes.MasheryPackage, error)
	UpdatePackage func(ctx context.Context, pack masherytypes.MasheryPackage, c *transport.V3Transport) (*masherytypes.MasheryPackage, error)
	DeletePackage func(ctx context.Context, packId string, c *transport.V3Transport) error
	ListPackages  func(ctx context.Context, c *transport.V3Transport) ([]masherytypes.MasheryPackage, error)

	// Package plans
	CreatePlanService  func(ctx context.Context, planService masherytypes.MasheryPlanService, c *transport.V3Transport) (*masherytypes.AddressableV3Object, error)
	DeletePlanService  func(ctx context.Context, planService masherytypes.MasheryPlanService, c *transport.V3Transport) error
	CreatePlanEndpoint func(ctx context.Context, planEndp masherytypes.MasheryPlanServiceEndpoint, c *transport.V3Transport) (*masherytypes.AddressableV3Object, error)
	DeletePlanEndpoint func(ctx context.Context, planEndp masherytypes.MasheryPlanServiceEndpoint, c *transport.V3Transport) error
	ListPlanEndpoints  func(ctx context.Context, planService masherytypes.MasheryPlanService, c *transport.V3Transport) ([]masherytypes.AddressableV3Object, error)

	CountPlanEndpoints func(ctx context.Context, planService masherytypes.MasheryPlanService, c *transport.V3Transport) (int64, error)
	CountPlanService   func(ctx context.Context, packageId, planId string, c *transport.V3Transport) (int64, error)
	GetPlan            func(ctx context.Context, packageId string, planId string, c *transport.V3Transport) (*masherytypes.MasheryPlan, error)
	CreatePlan         func(ctx context.Context, packageId string, plan masherytypes.MasheryPlan, c *transport.V3Transport) (*masherytypes.MasheryPlan, error)
	UpdatePlan         func(ctx context.Context, plan masherytypes.MasheryPlan, c *transport.V3Transport) (*masherytypes.MasheryPlan, error)
	DeletePlan         func(ctx context.Context, packageId, planId string, c *transport.V3Transport) error
	CountPlans         func(ctx context.Context, packageId string, c *transport.V3Transport) (int64, error)
	ListPlans          func(ctx context.Context, packageId string, c *transport.V3Transport) ([]masherytypes.MasheryPlan, error)
	ListPlanServices   func(ctx context.Context, packageId string, planId string, c *transport.V3Transport) ([]masherytypes.MasheryService, error)

	// Plan methods
	ListPackagePlanMethods  func(ctx context.Context, id masherytypes.MasheryPlanServiceEndpoint, c *transport.V3Transport) ([]masherytypes.MasheryMethod, error)
	GetPackagePlanMethod    func(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod, c *transport.V3Transport) (*masherytypes.MasheryMethod, error)
	CreatePackagePlanMethod func(ctx context.Context, id masherytypes.MasheryPlanServiceEndpoint, upsert masherytypes.MasheryMethod, c *transport.V3Transport) (*masherytypes.MasheryMethod, error)
	DeletePackagePlanMethod func(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod, c *transport.V3Transport) error

	// Plan method filter
	GetPackagePlanMethodFilter    func(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod, c *transport.V3Transport) (*masherytypes.MasheryResponseFilter, error)
	CreatePackagePlanMethodFilter func(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod, ref masherytypes.MasheryServiceMethodFilter, c *transport.V3Transport) (*masherytypes.MasheryResponseFilter, error)
	DeletePackagePlanMethodFilter func(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod, c *transport.V3Transport) error

	// Packge key
	GetPackageKey           func(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryPackageKey, error)
	CreatePackageKey        func(ctx context.Context, appId string, packageKey masherytypes.MasheryPackageKey, c *transport.V3Transport) (*masherytypes.MasheryPackageKey, error)
	UpdatePackageKey        func(ctx context.Context, packageKey masherytypes.MasheryPackageKey, c *transport.V3Transport) (*masherytypes.MasheryPackageKey, error)
	DeletePackageKey        func(ctx context.Context, keyId string, c *transport.V3Transport) error
	ListPackageKeysFiltered func(ctx context.Context, params map[string]string, fields []string, c *transport.V3Transport) ([]masherytypes.MasheryPackageKey, error)
	ListPackageKeys         func(ctx context.Context, c *transport.V3Transport) ([]masherytypes.MasheryPackageKey, error)

	// Roles
	GetRole           func(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryRole, error)
	ListRoles         func(ctx context.Context, c *transport.V3Transport) ([]masherytypes.MasheryRole, error)
	ListRolesFiltered func(ctx context.Context, params map[string]string, fields []string, c *transport.V3Transport) ([]masherytypes.MasheryRole, error)

	// Services
	GetService           func(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryService, error)
	CreateService        func(ctx context.Context, service masherytypes.MasheryService, c *transport.V3Transport) (*masherytypes.MasheryService, error)
	UpdateService        func(ctx context.Context, service masherytypes.MasheryService, c *transport.V3Transport) (*masherytypes.MasheryService, error)
	DeleteService        func(ctx context.Context, serviceId string, c *transport.V3Transport) error
	ListServicesFiltered func(ctx context.Context, params map[string]string, fields []string, c *transport.V3Transport) ([]masherytypes.MasheryService, error)
	ListServices         func(ctx context.Context, c *transport.V3Transport) ([]masherytypes.MasheryService, error)
	CountServices        func(ctx context.Context, params map[string]string, c *transport.V3Transport) (int64, error)

	ListErrorSets         func(ctx context.Context, serviceId string, qs url.Values, c *transport.V3Transport) ([]masherytypes.MasheryErrorSet, error)
	GetErrorSet           func(ctx context.Context, serviceId, setId string, c *transport.V3Transport) (*masherytypes.MasheryErrorSet, error)
	CreateErrorSet        func(ctx context.Context, serviceId string, set masherytypes.MasheryErrorSet, c *transport.V3Transport) (*masherytypes.MasheryErrorSet, error)
	UpdateErrorSet        func(ctx context.Context, serviceId string, setData masherytypes.MasheryErrorSet, c *transport.V3Transport) (*masherytypes.MasheryErrorSet, error)
	DeleteErrorSet        func(ctx context.Context, serviceId, setId string, c *transport.V3Transport) error
	UpdateErrorSetMessage func(ctx context.Context, serviceId string, setId string, msg masherytypes.MasheryErrorMessage, c *transport.V3Transport) (*masherytypes.MasheryErrorMessage, error)

	GetServiceRoles func(ctx context.Context, serviceId string, c *transport.V3Transport) ([]masherytypes.MasheryRolePermission, error)
	SetServiceRoles func(ctx context.Context, id string, roles []masherytypes.MasheryRolePermission, c *transport.V3Transport) error

	// Service cache
	GetServiceCache    func(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryServiceCache, error)
	CreateServiceCache func(ctx context.Context, id string, service masherytypes.MasheryServiceCache, c *transport.V3Transport) (*masherytypes.MasheryServiceCache, error)
	UpdateServiceCache func(ctx context.Context, id string, service masherytypes.MasheryServiceCache, c *transport.V3Transport) (*masherytypes.MasheryServiceCache, error)
	DeleteServiceCache func(ctx context.Context, id string, c *transport.V3Transport) error

	// Service OAuth
	GetServiceOAuthSecurityProfile    func(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryOAuth, error)
	CreateServiceOAuthSecurityProfile func(ctx context.Context, id string, service masherytypes.MasheryOAuth, c *transport.V3Transport) (*masherytypes.MasheryOAuth, error)
	UpdateServiceOAuthSecurityProfile func(ctx context.Context, id string, service masherytypes.MasheryOAuth, c *transport.V3Transport) (*masherytypes.MasheryOAuth, error)
	DeleteServiceOAuthSecurityProfile func(ctx context.Context, id string, c *transport.V3Transport) error
}

func (c *PluggableClient) ListErrorSets(ctx context.Context, serviceId string, qs url.Values) ([]masherytypes.MasheryErrorSet, error) {
	if c.schema.ListErrorSets != nil {
		return c.schema.ListErrorSets(ctx, serviceId, qs, c.transport)
	} else {
		return []masherytypes.MasheryErrorSet{}, c.notImplemented("ListErrorSets")
	}
}

func (c *PluggableClient) GetErrorSet(ctx context.Context, serviceId, setId string) (*masherytypes.MasheryErrorSet, error) {
	if c.schema.GetErrorSet != nil {
		return c.schema.GetErrorSet(ctx, serviceId, setId, c.transport)
	} else {
		return nil, c.notImplemented("GetErrorSet")
	}
}

func (c *PluggableClient) CreateErrorSet(ctx context.Context, serviceId string, set masherytypes.MasheryErrorSet) (*masherytypes.MasheryErrorSet, error) {
	if c.schema.CreateErrorSet != nil {
		return c.schema.CreateErrorSet(ctx, serviceId, set, c.transport)
	} else {
		return nil, c.notImplemented("CreateErrorSet")
	}
}

func (c *PluggableClient) UpdateErrorSet(ctx context.Context, serviceId string, setData masherytypes.MasheryErrorSet) (*masherytypes.MasheryErrorSet, error) {
	if c.schema.UpdateErrorSet != nil {
		return c.schema.UpdateErrorSet(ctx, serviceId, setData, c.transport)
	} else {
		return nil, c.notImplemented("UpdateErrorSet")
	}
}

func (c *PluggableClient) DeleteErrorSet(ctx context.Context, serviceId, setId string) error {
	if c.schema.DeleteErrorSet != nil {
		return c.schema.DeleteErrorSet(ctx, serviceId, setId, c.transport)
	} else {
		return c.notImplemented("DeleteErrorSet")
	}
}

func (c *PluggableClient) UpdateErrorSetMessage(ctx context.Context, serviceId string, setId string, msg masherytypes.MasheryErrorMessage) (*masherytypes.MasheryErrorMessage, error) {
	if c.schema.UpdateErrorSetMessage != nil {
		return c.schema.UpdateErrorSetMessage(ctx, serviceId, setId, msg, c.transport)
	} else {
		return nil, c.notImplemented("UpdateErrorSetMessage")
	}
}

func (c *PluggableClient) GetServiceRoles(ctx context.Context, serviceId string) ([]masherytypes.MasheryRolePermission, error) {
	if c.schema.GetServiceRoles != nil {
		return c.schema.GetServiceRoles(ctx, serviceId, c.transport)
	} else {
		return []masherytypes.MasheryRolePermission{}, c.notImplemented("GetServiceRoles")
	}
}

func (c *PluggableClient) SetServiceRoles(ctx context.Context, serviceId string, perms []masherytypes.MasheryRolePermission) error {
	if c.schema.SetServiceRoles != nil {
		return c.schema.SetServiceRoles(ctx, serviceId, perms, c.transport)
	} else {
		return c.notImplemented("SetServiceRoles")
	}
}

func (c *PluggableClient) GetPublicDomains(ctx context.Context) ([]string, error) {
	if c.schema.GetPublicDomains != nil {
		return c.schema.GetPublicDomains(ctx, c.transport)
	} else {
		return []string{}, c.notImplemented("GetPublicDomains")
	}
}

func (c *PluggableClient) GetSystemDomains(ctx context.Context) ([]string, error) {
	if c.schema.GetSystemDomains != nil {
		return c.schema.GetSystemDomains(ctx, c.transport)
	} else {
		return []string{}, c.notImplemented("GetSystemDomains")
	}
}

func (c *PluggableClient) GetApplication(ctx context.Context, appId string) (*masherytypes.MasheryApplication, error) {
	if c.schema.GetApplicationContext != nil {
		return c.schema.GetApplicationContext(ctx, appId, c.transport)
	} else {
		return nil, c.notImplemented("GetApplication")
	}
}

func (c *PluggableClient) GetApplicationPackageKeys(ctx context.Context, appId string) ([]masherytypes.MasheryPackageKey, error) {
	if c.schema.GetApplicationContext != nil {
		return c.schema.GetApplicationPackageKeys(ctx, appId, c.transport)
	} else {
		return nil, c.notImplemented("GetApplicationPackageKeys")
	}
}

func (c *PluggableClient) CountApplicationPackageKeys(ctx context.Context, appId string) (int64, error) {
	if c.schema.CountApplicationPackageKeys != nil {
		return c.schema.CountApplicationPackageKeys(ctx, appId, c.transport)
	} else {
		return 0, c.notImplemented("CountApplicationPackageKeys")
	}
}

func (c *PluggableClient) GetFullApplication(ctx context.Context, id string) (*masherytypes.MasheryApplication, error) {
	if c.schema.GetFullApplication != nil {
		return c.schema.GetFullApplication(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetFullApplication")
	}
}

func (c *PluggableClient) CreateApplication(ctx context.Context, memberId string, member masherytypes.MasheryApplication) (*masherytypes.MasheryApplication, error) {
	if c.schema.CreateApplication != nil {
		return c.schema.CreateApplication(ctx, memberId, member, c.transport)
	} else {
		return nil, c.notImplemented("CreateApplication")
	}
}

func (c *PluggableClient) UpdateApplication(ctx context.Context, app masherytypes.MasheryApplication) (*masherytypes.MasheryApplication, error) {
	if c.schema.UpdateApplication != nil {
		return c.schema.UpdateApplication(ctx, app, c.transport)
	} else {
		return nil, c.notImplemented("UpdateApplication")
	}
}

func (c *PluggableClient) DeleteApplication(ctx context.Context, appId string) error {
	if c.schema.DeleteApplication != nil {
		return c.schema.DeleteApplication(ctx, appId, c.transport)
	} else {
		return c.notImplemented("DeleteApplication")
	}
}

func (c *PluggableClient) CountApplicationsOfMember(ctx context.Context, memberId string) (int64, error) {
	if c.schema.CountApplicationsOfMember != nil {
		return c.schema.CountApplicationsOfMember(ctx, memberId, c.transport)
	} else {
		return 0, c.notImplemented("DeleteApplication")
	}
}

func (c *PluggableClient) ListApplications(ctx context.Context) ([]masherytypes.MasheryApplication, error) {
	if c.schema.ListApplications != nil {
		return c.schema.ListApplications(ctx, c.transport)
	} else {
		return []masherytypes.MasheryApplication{}, c.notImplemented("ListApplications")
	}
}

// -----------------------------------------------------------------------------------------------------------------
// Email template set
// -----------------------------------------------------------------------------------------------------------------

func (c *PluggableClient) GetEmailTemplateSet(ctx context.Context, id string) (*masherytypes.MasheryEmailTemplateSet, error) {
	if c.schema.GetEmailTemplateSet != nil {
		return c.schema.GetEmailTemplateSet(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetEmailTemplateSet")
	}
}

func (c *PluggableClient) ListEmailTemplateSets(ctx context.Context) ([]masherytypes.MasheryEmailTemplateSet, error) {
	if c.schema.ListEmailTemplateSets != nil {
		return c.schema.ListEmailTemplateSets(ctx, c.transport)
	} else {
		return nil, c.notImplemented("ListEmailTemplateSets")
	}
}

func (c *PluggableClient) ListEmailTemplateSetsFiltered(ctx context.Context, params map[string]string, fields []string) ([]masherytypes.MasheryEmailTemplateSet, error) {
	if c.schema.ListEmailTemplateSetsFiltered != nil {
		return c.schema.ListEmailTemplateSetsFiltered(ctx, params, fields, c.transport)
	} else {
		return nil, c.notImplemented("ListEmailTemplateSetsFiltered")
	}
}

// -----------------------------------------------------------------------------------------------------------------
// Endpoints
// -----------------------------------------------------------------------------------------------------------------

func (c *PluggableClient) ListEndpoints(ctx context.Context, serviceId string) ([]masherytypes.AddressableV3Object, error) {
	if c.schema.ListEndpoints != nil {
		return c.schema.ListEndpoints(ctx, serviceId, c.transport)
	} else {
		return nil, c.notImplemented("ListEndpoints")
	}
}

func (c *PluggableClient) ListEndpointsWithFullInfo(ctx context.Context, serviceId string) ([]masherytypes.MasheryEndpoint, error) {
	if c.schema.ListEndpointsWithFullInfo != nil {
		return c.schema.ListEndpointsWithFullInfo(ctx, serviceId, c.transport)
	} else {
		return nil, c.notImplemented("ListEndpointsWithFullInfo")
	}
}

func (c *PluggableClient) CreateEndpoint(ctx context.Context, serviceId string, endp masherytypes.MasheryEndpoint) (*masherytypes.MasheryEndpoint, error) {
	if c.schema.CreateEndpoint != nil {
		return c.schema.CreateEndpoint(ctx, serviceId, endp, c.transport)
	} else {
		return nil, c.notImplemented("CreateEndpoint")
	}
}

func (c *PluggableClient) UpdateEndpoint(ctx context.Context, serviceId string, endp masherytypes.MasheryEndpoint) (*masherytypes.MasheryEndpoint, error) {
	if c.schema.UpdateEndpoint != nil {
		return c.schema.UpdateEndpoint(ctx, serviceId, endp, c.transport)
	} else {
		return nil, c.notImplemented("UpdateEndpoint")
	}
}

func (c *PluggableClient) GetEndpoint(ctx context.Context, serviceId string, endpointId string) (*masherytypes.MasheryEndpoint, error) {
	if c.schema.GetEndpoint != nil {
		return c.schema.GetEndpoint(ctx, serviceId, endpointId, c.transport)
	} else {
		return nil, c.notImplemented("GetEndpoint")
	}
}

func (c *PluggableClient) DeleteEndpoint(ctx context.Context, serviceId, endpointId string) error {
	if c.schema.DeleteEndpoint != nil {
		return c.schema.DeleteEndpoint(ctx, serviceId, endpointId, c.transport)
	} else {
		return c.notImplemented("DeleteEndpoint")
	}
}

func (c *PluggableClient) CountEndpointsOf(ctx context.Context, serviceId string) (int64, error) {
	if c.schema.CountEndpointsOf != nil {
		return c.schema.CountEndpointsOf(ctx, serviceId, c.transport)
	} else {
		return 0, c.notImplemented("CountEndpointsOf")
	}
}

// -------------------------------------
// Endpoint methods

func (c *PluggableClient) ListEndpointMethods(ctx context.Context, serviceId, endpointId string) ([]masherytypes.MasheryMethod, error) {
	if c.schema.ListEndpointMethods != nil {
		return c.schema.ListEndpointMethods(ctx, serviceId, endpointId, c.transport)
	} else {
		return []masherytypes.MasheryMethod{}, c.notImplemented("ListEndpointMethods")
	}
}

func (c *PluggableClient) ListEndpointMethodsWithFullInfo(ctx context.Context, serviceId, endpointId string) ([]masherytypes.MasheryMethod, error) {
	if c.schema.ListEndpointMethodsWithFullInfo != nil {
		return c.schema.ListEndpointMethodsWithFullInfo(ctx, serviceId, endpointId, c.transport)
	} else {
		return []masherytypes.MasheryMethod{}, c.notImplemented("ListEndpointMethodsWithFullInfo")
	}
}

func (c *PluggableClient) CreateEndpointMethod(ctx context.Context, serviceId, endpointId string, methoUpsert masherytypes.MasheryMethod) (*masherytypes.MasheryMethod, error) {
	if c.schema.CreateEndpointMethod != nil {
		return c.schema.CreateEndpointMethod(ctx, serviceId, endpointId, methoUpsert, c.transport)
	} else {
		return nil, c.notImplemented("CreateEndpointMethod")
	}
}

func (c *PluggableClient) UpdateEndpointMethod(ctx context.Context, serviceId, endpointId string, methUpsert masherytypes.MasheryMethod) (*masherytypes.MasheryMethod, error) {
	if c.schema.UpdateEndpointMethod != nil {
		return c.schema.UpdateEndpointMethod(ctx, serviceId, endpointId, methUpsert, c.transport)
	} else {
		return nil, c.notImplemented("UpdateEndpointMethod")
	}
}

func (c *PluggableClient) GetEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string) (*masherytypes.MasheryMethod, error) {
	if c.schema.GetEndpointMethod != nil {
		return c.schema.GetEndpointMethod(ctx, serviceId, endpointId, methodId, c.transport)
	} else {
		return nil, c.notImplemented("GetEndpointMethod")
	}
}

func (c *PluggableClient) DeleteEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string) error {
	if c.schema.DeleteEndpointMethod != nil {
		return c.schema.DeleteEndpointMethod(ctx, serviceId, endpointId, methodId, c.transport)
	} else {
		return c.notImplemented("DeleteEndpointMethod")
	}
}

func (c *PluggableClient) CountEndpointsMethodsOf(ctx context.Context, serviceId, endpointId string) (int64, error) {
	if c.schema.CountEndpointsMethodsOf != nil {
		return c.schema.CountEndpointsMethodsOf(ctx, serviceId, endpointId, c.transport)
	} else {
		return 0, c.notImplemented("CountEndpointsMethodsOf")
	}
}

// -----------------------------------------------------------------------------------
// Endpoint method filters.

func (c *PluggableClient) ListEndpointMethodFilters(ctx context.Context, serviceId, endpointId, methodId string) ([]masherytypes.MasheryResponseFilter, error) {
	if c.schema.ListEndpointMethodFilters != nil {
		return c.schema.ListEndpointMethodFilters(ctx, serviceId, endpointId, methodId, c.transport)
	} else {
		return nil, c.notImplemented("ListEndpointMethodFilters")
	}
}

func (c *PluggableClient) ListEndpointMethodFiltersWithFullInfo(ctx context.Context, serviceId, endpointId, methodId string) ([]masherytypes.MasheryResponseFilter, error) {
	if c.schema.ListEndpointMethodFiltersWithFullInfo != nil {
		return c.schema.ListEndpointMethodFiltersWithFullInfo(ctx, serviceId, endpointId, methodId, c.transport)
	} else {
		return nil, c.notImplemented("ListEndpointMethodFiltersWithFullInfo")
	}
}

func (c *PluggableClient) CreateEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId string, filterUpsert masherytypes.MasheryResponseFilter) (*masherytypes.MasheryResponseFilter, error) {
	if c.schema.CreateEndpointMethodFilter != nil {
		return c.schema.CreateEndpointMethodFilter(ctx, serviceId, endpointId, methodId, filterUpsert, c.transport)
	} else {
		return nil, c.notImplemented("CreateEndpointMethodFilter")
	}
}

func (c *PluggableClient) UpdateEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId string, filterUpsert masherytypes.MasheryResponseFilter) (*masherytypes.MasheryResponseFilter, error) {
	if c.schema.UpdateEndpointMethodFilter != nil {
		return c.schema.UpdateEndpointMethodFilter(ctx, serviceId, endpointId, methodId, filterUpsert, c.transport)
	} else {
		return nil, c.notImplemented("UpdateEndpointMethodFilter")
	}
}

func (c *PluggableClient) GetEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId, filterId string) (*masherytypes.MasheryResponseFilter, error) {
	if c.schema.GetEndpointMethodFilter != nil {
		return c.schema.GetEndpointMethodFilter(ctx, serviceId, endpointId, methodId, filterId, c.transport)
	} else {
		return nil, c.notImplemented("GetEndpointMethodFilter")
	}
}

func (c *PluggableClient) DeleteEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId, filterId string) error {
	if c.schema.DeleteEndpointMethodFilter != nil {
		return c.schema.DeleteEndpointMethodFilter(ctx, serviceId, endpointId, methodId, filterId, c.transport)
	} else {
		return c.notImplemented("DeleteEndpointMethodFilter")
	}
}

func (c *PluggableClient) CountEndpointsMethodsFiltersOf(ctx context.Context, serviceId, endpointId, methodId string) (int64, error) {
	if c.schema.CountEndpointsMethodsFiltersOf != nil {
		return c.schema.CountEndpointsMethodsFiltersOf(ctx, serviceId, endpointId, methodId, c.transport)
	} else {
		return 0, c.notImplemented("CountEndpointsMethodsFiltersOf")
	}
}

func (c *PluggableClient) GetMember(ctx context.Context, id string) (*masherytypes.MasheryMember, error) {
	if c.schema.GetMember != nil {
		return c.schema.GetMember(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetMember")
	}
}

func (c *PluggableClient) GetFullMember(ctx context.Context, id string) (*masherytypes.MasheryMember, error) {
	if c.schema.GetFullMember != nil {
		return c.schema.GetFullMember(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetFullMember")
	}
}

func (c *PluggableClient) CreateMember(ctx context.Context, member masherytypes.MasheryMember) (*masherytypes.MasheryMember, error) {
	if c.schema.CreateMember != nil {
		return c.schema.CreateMember(ctx, member, c.transport)
	} else {
		return nil, c.notImplemented("CreateMember")
	}
}

func (c *PluggableClient) UpdateMember(ctx context.Context, member masherytypes.MasheryMember) (*masherytypes.MasheryMember, error) {
	if c.schema.UpdateMember != nil {
		return c.schema.UpdateMember(ctx, member, c.transport)
	} else {
		return nil, c.notImplemented("UpdateMember")
	}
}

func (c *PluggableClient) DeleteMember(ctx context.Context, memberId string) error {
	if c.schema.DeleteMember != nil {
		return c.schema.DeleteMember(ctx, memberId, c.transport)
	} else {
		return c.notImplemented("DeleteMember")
	}
}

func (c *PluggableClient) ListMembers(ctx context.Context) ([]masherytypes.MasheryMember, error) {
	if c.schema.ListMembers != nil {
		return c.schema.ListMembers(ctx, c.transport)
	} else {
		return []masherytypes.MasheryMember{}, c.notImplemented("ListMembers")
	}
}

// ---------------------------------------------
// Packages

func (c *PluggableClient) GetPackage(ctx context.Context, id string) (*masherytypes.MasheryPackage, error) {
	if c.schema.GetPackage != nil {
		return c.schema.GetPackage(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetPackage")
	}
}

func (c *PluggableClient) CreatePackage(ctx context.Context, pack masherytypes.MasheryPackage) (*masherytypes.MasheryPackage, error) {
	if c.schema.CreatePackage != nil {
		return c.schema.CreatePackage(ctx, pack, c.transport)
	} else {
		return nil, c.notImplemented("CreatePackage")
	}
}

func (c *PluggableClient) UpdatePackage(ctx context.Context, pack masherytypes.MasheryPackage) (*masherytypes.MasheryPackage, error) {
	if c.schema.UpdatePackage != nil {
		return c.schema.UpdatePackage(ctx, pack, c.transport)
	} else {
		return nil, c.notImplemented("UpdatePackage")
	}
}

func (c *PluggableClient) DeletePackage(ctx context.Context, packId string) error {
	if c.schema.DeletePackage != nil {
		return c.schema.DeletePackage(ctx, packId, c.transport)
	} else {
		return c.notImplemented("DeletePackage")
	}
}

func (c *PluggableClient) ListPackages(ctx context.Context) ([]masherytypes.MasheryPackage, error) {
	if c.schema.ListPackages != nil {
		return c.schema.ListPackages(ctx, c.transport)
	} else {
		return []masherytypes.MasheryPackage{}, c.notImplemented("ListPackages")
	}
}

// --------------------------------------
// Package plans

func (c *PluggableClient) CreatePlanService(ctx context.Context, planService masherytypes.MasheryPlanService) (*masherytypes.AddressableV3Object, error) {
	if c.schema.CreatePlanService != nil {
		return c.schema.CreatePlanService(ctx, planService, c.transport)
	} else {
		return nil, c.notImplemented("CreatePlanService")
	}
}

func (c *PluggableClient) DeletePlanService(ctx context.Context, planService masherytypes.MasheryPlanService) error {
	if c.schema.DeletePlanService != nil {
		return c.schema.DeletePlanService(ctx, planService, c.transport)
	} else {
		return c.notImplemented("CreatePlanService")
	}
}

func (c *PluggableClient) CreatePlanEndpoint(ctx context.Context, planEndp masherytypes.MasheryPlanServiceEndpoint) (*masherytypes.AddressableV3Object, error) {
	if c.schema.CreatePlanEndpoint != nil {
		return c.schema.CreatePlanEndpoint(ctx, planEndp, c.transport)
	} else {
		return nil, c.notImplemented("CreatePlanEndpoint")
	}
}

func (c *PluggableClient) DeletePlanEndpoint(ctx context.Context, planEndp masherytypes.MasheryPlanServiceEndpoint) error {
	if c.schema.DeletePlanEndpoint != nil {
		return c.schema.DeletePlanEndpoint(ctx, planEndp, c.transport)
	} else {
		return c.notImplemented("DeletePlanEndpoint")
	}
}

func (c *PluggableClient) ListPlanEndpoints(ctx context.Context, planService masherytypes.MasheryPlanService) ([]masherytypes.AddressableV3Object, error) {
	if c.schema.ListPlanEndpoints != nil {
		return c.schema.ListPlanEndpoints(ctx, planService, c.transport)
	} else {
		return []masherytypes.AddressableV3Object{}, c.notImplemented("ListPlanEndpoints")
	}
}

func (c *PluggableClient) ListPlanServices(ctx context.Context, packageId string, planId string) ([]masherytypes.MasheryService, error) {
	if c.schema.ListPlanServices != nil {
		return c.schema.ListPlanServices(ctx, packageId, planId, c.transport)
	} else {
		return []masherytypes.MasheryService{}, c.notImplemented("ListPlanServices")
	}
}

func (c *PluggableClient) ListPlans(ctx context.Context, packageId string) ([]masherytypes.MasheryPlan, error) {
	if c.schema.ListPlans != nil {
		return c.schema.ListPlans(ctx, packageId, c.transport)
	} else {
		return []masherytypes.MasheryPlan{}, c.notImplemented("ListPlans")
	}
}

func (c *PluggableClient) CountPlans(ctx context.Context, packageId string) (int64, error) {
	if c.schema.CountPlans != nil {
		return c.schema.CountPlans(ctx, packageId, c.transport)
	} else {
		return 0, c.notImplemented("CountPlans")
	}
}

func (c *PluggableClient) DeletePlan(ctx context.Context, packageId, planId string) error {
	if c.schema.DeletePlan != nil {
		return c.schema.DeletePlan(ctx, packageId, planId, c.transport)
	} else {
		return c.notImplemented("DeletePlan")
	}
}

func (c *PluggableClient) UpdatePlan(ctx context.Context, plan masherytypes.MasheryPlan) (*masherytypes.MasheryPlan, error) {
	if c.schema.UpdatePlan != nil {
		return c.schema.UpdatePlan(ctx, plan, c.transport)
	} else {
		return nil, c.notImplemented("UpdatePlan")
	}
}

func (c *PluggableClient) CreatePlan(ctx context.Context, packageId string, plan masherytypes.MasheryPlan) (*masherytypes.MasheryPlan, error) {
	if c.schema.CreatePlan != nil {
		return c.schema.CreatePlan(ctx, packageId, plan, c.transport)
	} else {
		return nil, c.notImplemented("CreatePlan")
	}
}

func (c *PluggableClient) GetPlan(ctx context.Context, packageId string, planId string) (*masherytypes.MasheryPlan, error) {
	if c.schema.GetPlan != nil {
		return c.schema.GetPlan(ctx, packageId, planId, c.transport)
	} else {
		return nil, c.notImplemented("GetPlan")
	}
}

func (c *PluggableClient) CountPlanService(ctx context.Context, packageId, planId string) (int64, error) {
	if c.schema.CountPlanService != nil {
		return c.schema.CountPlanService(ctx, packageId, planId, c.transport)
	} else {
		return 0, c.notImplemented("ListPlanEndpoints")
	}
}

func (c *PluggableClient) CountPlanEndpoints(ctx context.Context, planService masherytypes.MasheryPlanService) (int64, error) {
	if c.schema.CountPlanEndpoints != nil {
		return c.schema.CountPlanEndpoints(ctx, planService, c.transport)
	} else {
		return 0, c.notImplemented("CountPlanEndpoints")
	}
}

//---------------------------------------------
// Package plan methods

func (c *PluggableClient) ListPackagePlanMethods(ctx context.Context, id masherytypes.MasheryPlanServiceEndpoint) ([]masherytypes.MasheryMethod, error) {
	if c.schema.ListPackagePlanMethods != nil {
		return c.schema.ListPackagePlanMethods(ctx, id, c.transport)
	} else {
		return []masherytypes.MasheryMethod{}, c.notImplemented("ListPackagePlanMethods")
	}
}

func (c *PluggableClient) GetPackagePlanMethod(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod) (*masherytypes.MasheryMethod, error) {
	if c.schema.GetPackagePlanMethod != nil {
		return c.schema.GetPackagePlanMethod(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetPackagePlanMethod")
	}
}

func (c *PluggableClient) CreatePackagePlanMethod(ctx context.Context, id masherytypes.MasheryPlanServiceEndpoint, upsert masherytypes.MasheryMethod) (*masherytypes.MasheryMethod, error) {
	if c.schema.CreatePackagePlanMethod != nil {
		return c.schema.CreatePackagePlanMethod(ctx, id, upsert, c.transport)
	} else {
		return nil, c.notImplemented("CreatePackagePlanMethod")
	}
}

func (c *PluggableClient) DeletePackagePlanMethod(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod) error {
	if c.schema.DeletePackagePlanMethod != nil {
		return c.schema.DeletePackagePlanMethod(ctx, id, c.transport)
	} else {
		return c.notImplemented("CreatePackagePlanMethod")
	}
}

// ----------------------------------------
// Plan method filter

func (c *PluggableClient) GetPackagePlanMethodFilter(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod) (*masherytypes.MasheryResponseFilter, error) {
	if c.schema.GetPackagePlanMethodFilter != nil {
		return c.schema.GetPackagePlanMethodFilter(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetPackagePlanMethodFilter")
	}
}

func (c *PluggableClient) CreatePackagePlanMethodFilter(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod, ref masherytypes.MasheryServiceMethodFilter) (*masherytypes.MasheryResponseFilter, error) {
	if c.schema.CreatePackagePlanMethodFilter != nil {
		return c.schema.CreatePackagePlanMethodFilter(ctx, id, ref, c.transport)
	} else {
		return nil, c.notImplemented("CreatePackagePlanMethodFilter")
	}
}

func (c *PluggableClient) DeletePackagePlanMethodFilter(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod) error {
	if c.schema.DeletePackagePlanMethodFilter != nil {
		return c.schema.DeletePackagePlanMethodFilter(ctx, id, c.transport)
	} else {
		return c.notImplemented("DeletePackagePlanMethodFilter")
	}
}

// ------------------------------------------------------------
// Package key

func (c *PluggableClient) GetPackageKey(ctx context.Context, id string) (*masherytypes.MasheryPackageKey, error) {
	if c.schema.GetPackageKey != nil {
		return c.schema.GetPackageKey(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetPackageKey")
	}
}

func (c *PluggableClient) CreatePackageKey(ctx context.Context, appId string, packageKey masherytypes.MasheryPackageKey) (*masherytypes.MasheryPackageKey, error) {
	if c.schema.CreatePackageKey != nil {
		return c.schema.CreatePackageKey(ctx, appId, packageKey, c.transport)
	} else {
		return nil, c.notImplemented("CreatePackageKey")
	}
}

func (c *PluggableClient) UpdatePackageKey(ctx context.Context, packageKey masherytypes.MasheryPackageKey) (*masherytypes.MasheryPackageKey, error) {
	if c.schema.UpdatePackageKey != nil {
		return c.schema.UpdatePackageKey(ctx, packageKey, c.transport)
	} else {
		return nil, c.notImplemented("UpdatePackageKey")
	}
}

func (c *PluggableClient) DeletePackageKey(ctx context.Context, keyId string) error {
	if c.schema.DeletePackageKey != nil {
		return c.schema.DeletePackageKey(ctx, keyId, c.transport)
	} else {
		return c.notImplemented("DeletePackageKey")
	}
}

func (c *PluggableClient) ListPackageKeysFiltered(ctx context.Context, params map[string]string, fields []string) ([]masherytypes.MasheryPackageKey, error) {
	if c.schema.ListPackageKeysFiltered != nil {
		return c.schema.ListPackageKeysFiltered(ctx, params, fields, c.transport)
	} else {
		return []masherytypes.MasheryPackageKey{}, c.notImplemented("ListPackageKeysFiltered")
	}
}

func (c *PluggableClient) ListPackageKeys(ctx context.Context) ([]masherytypes.MasheryPackageKey, error) {
	if c.schema.ListPackageKeys != nil {
		return c.schema.ListPackageKeys(ctx, c.transport)
	} else {
		return []masherytypes.MasheryPackageKey{}, c.notImplemented("ListPackageKeys")
	}
}

// ---------------------
// Roles

func (c *PluggableClient) GetRole(ctx context.Context, id string) (*masherytypes.MasheryRole, error) {
	if c.schema.GetRole != nil {
		return c.schema.GetRole(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetRole")
	}
}

func (c *PluggableClient) ListRoles(ctx context.Context) ([]masherytypes.MasheryRole, error) {
	if c.schema.ListRoles != nil {
		return c.schema.ListRoles(ctx, c.transport)
	} else {
		return []masherytypes.MasheryRole{}, c.notImplemented("ListRoles")
	}
}

func (c *PluggableClient) ListRolesFiltered(ctx context.Context, params map[string]string, fields []string) ([]masherytypes.MasheryRole, error) {
	if c.schema.ListRolesFiltered != nil {
		return c.schema.ListRolesFiltered(ctx, params, fields, c.transport)
	} else {
		return []masherytypes.MasheryRole{}, c.notImplemented("ListRolesFiltered")
	}
}

// ------------------------------
// Service

func (c *PluggableClient) GetService(ctx context.Context, id string) (*masherytypes.MasheryService, error) {
	if c.schema.GetService != nil {
		return c.schema.GetService(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetService")
	}
}

func (c *PluggableClient) CreateService(ctx context.Context, service masherytypes.MasheryService) (*masherytypes.MasheryService, error) {
	if c.schema.CreateService != nil {
		return c.schema.CreateService(ctx, service, c.transport)
	} else {
		return nil, c.notImplemented("CreateService")
	}
}

func (c *PluggableClient) UpdateService(ctx context.Context, service masherytypes.MasheryService) (*masherytypes.MasheryService, error) {
	if c.schema.UpdateService != nil {
		return c.schema.UpdateService(ctx, service, c.transport)
	} else {
		return nil, c.notImplemented("UpdateService")
	}
}

func (c *PluggableClient) DeleteService(ctx context.Context, serviceId string) error {
	if c.schema.DeleteService != nil {
		return c.schema.DeleteService(ctx, serviceId, c.transport)
	} else {
		return c.notImplemented("UpdateService")
	}
}

func (c *PluggableClient) ListServicesFiltered(ctx context.Context, params map[string]string, fields []string) ([]masherytypes.MasheryService, error) {
	if c.schema.ListServicesFiltered != nil {
		return c.schema.ListServicesFiltered(ctx, params, fields, c.transport)
	} else {
		return []masherytypes.MasheryService{}, c.notImplemented("ListServicesFiltered")
	}
}

func (c *PluggableClient) ListServices(ctx context.Context) ([]masherytypes.MasheryService, error) {
	if c.schema.ListServices != nil {
		return c.schema.ListServices(ctx, c.transport)
	} else {
		return []masherytypes.MasheryService{}, c.notImplemented("ListServices")
	}
}

func (c *PluggableClient) CountServices(ctx context.Context, params map[string]string) (int64, error) {
	if c.schema.CountServices != nil {
		return c.schema.CountServices(ctx, params, c.transport)
	} else {
		return 0, c.notImplemented("CountServices")
	}
}

// --------------------------------
// Service cache

func (c *PluggableClient) GetServiceCache(ctx context.Context, id string) (*masherytypes.MasheryServiceCache, error) {
	if c.schema.GetServiceCache != nil {
		return c.schema.GetServiceCache(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetServiceCache")
	}
}

func (c *PluggableClient) CreateServiceCache(ctx context.Context, id string, service masherytypes.MasheryServiceCache) (*masherytypes.MasheryServiceCache, error) {
	if c.schema.CreateServiceCache != nil {
		return c.schema.CreateServiceCache(ctx, id, service, c.transport)
	} else {
		return nil, c.notImplemented("CreateServiceCache")
	}
}

func (c *PluggableClient) UpdateServiceCache(ctx context.Context, id string, service masherytypes.MasheryServiceCache) (*masherytypes.MasheryServiceCache, error) {
	if c.schema.UpdateServiceCache != nil {
		return c.schema.UpdateServiceCache(ctx, id, service, c.transport)
	} else {
		return nil, c.notImplemented("UpdateServiceCache")
	}
}

func (c *PluggableClient) DeleteServiceCache(ctx context.Context, id string) error {
	if c.schema.DeleteServiceCache != nil {
		return c.schema.DeleteServiceCache(ctx, id, c.transport)
	} else {
		return c.notImplemented("DeleteServiceCache")
	}
}

// Service OAtuh
func (c *PluggableClient) GetServiceOAuthSecurityProfile(ctx context.Context, id string) (*masherytypes.MasheryOAuth, error) {
	if c.schema.GetServiceOAuthSecurityProfile != nil {
		return c.schema.GetServiceOAuthSecurityProfile(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetServiceOAuthSecurityProfile")
	}
}

func (c *PluggableClient) CreateServiceOAuthSecurityProfile(ctx context.Context, id string, service masherytypes.MasheryOAuth) (*masherytypes.MasheryOAuth, error) {
	if c.schema.CreateServiceOAuthSecurityProfile != nil {
		return c.schema.CreateServiceOAuthSecurityProfile(ctx, id, service, c.transport)
	} else {
		return nil, c.notImplemented("CreateServiceOAuthSecurityProfile")
	}
}

func (c *PluggableClient) UpdateServiceOAuthSecurityProfile(ctx context.Context, id string, service masherytypes.MasheryOAuth) (*masherytypes.MasheryOAuth, error) {
	if c.schema.UpdateServiceOAuthSecurityProfile != nil {
		return c.schema.UpdateServiceOAuthSecurityProfile(ctx, id, service, c.transport)
	} else {
		return nil, c.notImplemented("UpdateServiceOAuthSecurityProfile")
	}
}

func (c *PluggableClient) DeleteServiceOAuthSecurityProfile(ctx context.Context, id string) error {
	if c.schema.DeleteServiceOAuthSecurityProfile != nil {
		return c.schema.DeleteServiceOAuthSecurityProfile(ctx, id, c.transport)
	} else {
		return c.notImplemented("DeleteServiceOAuthSecurityProfile")
	}
}
