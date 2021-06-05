package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

type Client interface {
	GetPublicDomains(ctx context.Context) ([]string, error)
	GetSystemDomains(ctx context.Context) ([]string, error)

	// GetApplication Retrieve the details of the application
	GetApplication(ctx context.Context, appId string) (*MasheryApplication, error)
	GetApplicationPackageKeys(ctx context.Context, appId string) ([]MasheryPackageKey, error)
	CountApplicationPackageKeys(ctx context.Context, appId string) (int64, error)
	GetFullApplication(ctx context.Context, id string) (*MasheryApplication, error)
	CreateApplication(ctx context.Context, memberId string, member MasheryApplication) (*MasheryApplication, error)
	UpdateApplication(ctx context.Context, app MasheryApplication) (*MasheryApplication, error)
	DeleteApplication(ctx context.Context, appId string) error
	CountApplicationsOfMember(ctx context.Context, memberId string) (int64, error)
	ListApplications(ctx context.Context) ([]MasheryApplication, error)

	// Email template sets
	GetEmailTemplateSet(ctx context.Context, id string) (*MasheryEmailTemplateSet, error)
	ListEmailTemplateSets(ctx context.Context) ([]MasheryEmailTemplateSet, error)
	ListEmailTemplateSetsFiltered(ctx context.Context, params map[string]string, fields []string) ([]MasheryEmailTemplateSet, error)

	// Endpoints
	ListEndpoints(ctx context.Context, serviceId string) ([]AddressableV3Object, error)
	ListEndpointsWithFullInfo(ctx context.Context, serviceId string) ([]MasheryEndpoint, error)
	CreateEndpoint(ctx context.Context, serviceId string, endp MasheryEndpoint) (*MasheryEndpoint, error)
	UpdateEndpoint(ctx context.Context, serviceId string, endp MasheryEndpoint) (*MasheryEndpoint, error)
	GetEndpoint(ctx context.Context, serviceId string, endpointId string) (*MasheryEndpoint, error)
	DeleteEndpoint(ctx context.Context, serviceId, endpointId string) error
	CountEndpointsOf(ctx context.Context, serviceId string) (int64, error)

	// Endpoint methods
	ListEndpointMethods(ctx context.Context, serviceId, endpointId string) ([]MasheryMethod, error)
	ListEndpointMethodsWithFullInfo(ctx context.Context, serviceId, endpointId string) ([]MasheryMethod, error)
	CreateEndpointMethod(ctx context.Context, serviceId, endpointId string, methoUpsert MasheryMethod) (*MasheryMethod, error)
	UpdateEndpointMethod(ctx context.Context, serviceId, endpointId string, methUpsert MasheryMethod) (*MasheryMethod, error)
	GetEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string) (*MasheryMethod, error)
	DeleteEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string) error
	CountEndpointsMethodsOf(ctx context.Context, serviceId, endpointId string) (int64, error)

	// Endpoint method filters
	ListEndpointMethodFilters(ctx context.Context, serviceId, endpointId, methodId string) ([]MasheryResponseFilter, error)
	ListEndpointMethodFiltersWithFullInfo(ctx context.Context, serviceId, endpointId, methodId string) ([]MasheryResponseFilter, error)
	CreateEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId string, filterUpsert MasheryResponseFilter) (*MasheryResponseFilter, error)
	UpdateEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId string, methUpsert MasheryResponseFilter) (*MasheryResponseFilter, error)
	GetEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId, filterId string) (*MasheryResponseFilter, error)
	DeleteEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId, filterId string) error
	CountEndpointsMethodsFiltersOf(ctx context.Context, serviceId, endpointId, methodId string) (int64, error)

	// Member
	GetMember(ctx context.Context, id string) (*MasheryMember, error)
	GetFullMember(ctx context.Context, id string) (*MasheryMember, error)
	CreateMember(ctx context.Context, member MasheryMember) (*MasheryMember, error)
	UpdateMember(ctx context.Context, member MasheryMember) (*MasheryMember, error)
	DeleteMember(ctx context.Context, memberId string) error
	ListMembers(ctx context.Context) ([]MasheryMember, error)

	// Packages
	GetPackage(ctx context.Context, id string) (*MasheryPackage, error)
	CreatePackage(ctx context.Context, pack MasheryPackage) (*MasheryPackage, error)
	UpdatePackage(ctx context.Context, pack MasheryPackage) (*MasheryPackage, error)
	DeletePackage(ctx context.Context, packId string) error
	ListPackages(ctx context.Context) ([]MasheryPackage, error)

	// Package plans
	CreatePlanService(ctx context.Context, planService MasheryPlanService) (*AddressableV3Object, error)
	DeletePlanService(ctx context.Context, planService MasheryPlanService) error
	CreatePlanEndpoint(ctx context.Context, planEndp MasheryPlanServiceEndpoint) (*AddressableV3Object, error)
	DeletePlanEndpoint(ctx context.Context, planEndp MasheryPlanServiceEndpoint) error
	ListPlanEndpoints(ctx context.Context, planService MasheryPlanService) ([]AddressableV3Object, error)
	CountPlanService(ctx context.Context, packageId, planId string) (int64, error)
	GetPlan(ctx context.Context, packageId string, planId string) (*MasheryPlan, error)
	CreatePlan(ctx context.Context, packageId string, plan MasheryPlan) (*MasheryPlan, error)
	UpdatePlan(ctx context.Context, plan MasheryPlan) (*MasheryPlan, error)
	DeletePlan(ctx context.Context, packageId, planId string) error
	CountPlans(ctx context.Context, packageId string) (int64, error)
	ListPlans(ctx context.Context, packageId string) ([]MasheryPlan, error)
	ListPlanServices(ctx context.Context, packageId string, planId string) ([]MasheryService, error)

	CountPlanEndpoints(ctx context.Context, planService MasheryPlanService) (int64, error)

	// Plan methods
	ListPackagePlanMethods(ctx context.Context, id MasheryPlanServiceEndpoint) ([]MasheryMethod, error)
	GetPackagePlanMethod(ctx context.Context, id MasheryPlanServiceEndpointMethod) (*MasheryMethod, error)
	CreatePackagePlanMethod(ctx context.Context, id MasheryPlanServiceEndpoint, upsert MasheryMethod) (*MasheryMethod, error)
	DeletePackagePlanMethod(ctx context.Context, id MasheryPlanServiceEndpointMethod) error

	// PLan method filter
	GetPackagePlanMethodFilter(ctx context.Context, id MasheryPlanServiceEndpointMethod) (*MasheryResponseFilter, error)
	CreatePackagePlanMethodFilter(ctx context.Context, id MasheryPlanServiceEndpointMethod, ref MasheryServiceMethodFilter) (*MasheryResponseFilter, error)
	DeletePackagePlanMethodFilter(ctx context.Context, id MasheryPlanServiceEndpointMethod) error

	// Package key
	GetPackageKey(ctx context.Context, id string) (*MasheryPackageKey, error)
	CreatePackageKey(ctx context.Context, appId string, packageKey MasheryPackageKey) (*MasheryPackageKey, error)
	UpdatePackageKey(ctx context.Context, packageKey MasheryPackageKey) (*MasheryPackageKey, error)
	DeletePackageKey(ctx context.Context, keyId string) error
	ListPackageKeysFiltered(ctx context.Context, params map[string]string, fields []string) ([]MasheryPackageKey, error)
	ListPackageKeys(ctx context.Context) ([]MasheryPackageKey, error)

	// Roles
	GetRole(ctx context.Context, id string) (*MasheryRole, error)
	ListRoles(ctx context.Context) ([]MasheryRole, error)
	ListRolesFiltered(ctx context.Context, params map[string]string, fields []string) ([]MasheryRole, error)

	// GetService retrieves service based on the service identifier
	GetService(ctx context.Context, id string) (*MasheryService, error)
	CreateService(ctx context.Context, service MasheryService) (*MasheryService, error)
	UpdateService(ctx context.Context, service MasheryService) (*MasheryService, error)
	DeleteService(ctx context.Context, serviceId string) error
	ListServicesFiltered(ctx context.Context, params map[string]string, fields []string) ([]MasheryService, error)
	ListServices(ctx context.Context) ([]MasheryService, error)
	CountServices(ctx context.Context, params map[string]string) (int64, error)

	ListErrorSets(ctx context.Context, serviceId string, qs url.Values) ([]MasheryErrorSet, error)
	GetErrorSet(ctx context.Context, serviceId, setId string) (*MasheryErrorSet, error)
	CreateErrorSet(ctx context.Context, serviceId string, set MasheryErrorSet) (*MasheryErrorSet, error)
	UpdateErrorSet(ctx context.Context, serviceId string, setData MasheryErrorSet) (*MasheryErrorSet, error)
	DeleteErrorSet(ctx context.Context, serviceId, setId string) error
	UpdateErrorSetMessage(ctx context.Context, serviceId string, setId string, msg MasheryErrorMessage) (*MasheryErrorMessage, error)

	GetServiceRoles(ctx context.Context, serviceId string) ([]MasheryRolePermission, error)
	SetServiceRoles(ctx context.Context, id string, roles []MasheryRolePermission) error

	// Service cache
	GetServiceCache(ctx context.Context, id string) (*MasheryServiceCache, error)
	CreateServiceCache(ctx context.Context, id string, service MasheryServiceCache) (*MasheryServiceCache, error)
	UpdateServiceCache(ctx context.Context, id string, service MasheryServiceCache) (*MasheryServiceCache, error)
	DeleteServiceCache(ctx context.Context, id string) error

	// Service OAuth
	GetServiceOAuthSecurityProfile(ctx context.Context, id string) (*MasheryOAuth, error)
	CreateServiceOAuthSecurityProfile(ctx context.Context, id string, service MasheryOAuth) (*MasheryOAuth, error)
	UpdateServiceOAuthSecurityProfile(ctx context.Context, id string, service MasheryOAuth) (*MasheryOAuth, error)
	DeleteServiceOAuthSecurityProfile(ctx context.Context, id string) error
}

type PluggableClient struct {
	schema    *ClientMethodSchema
	transport *HttpTransport
}

// FixedSchemeClient Fixed method scheme client that will not allow changing schema after it was created.
type FixedSchemeClient struct {
	PluggableClient
}

func (fsc *FixedSchemeClient) AssumeSchema(_ *ClientMethodSchema) {
	panic("This client cannot change method schema after it was created")
}

// Allows a pluggable client to assume a schema. This method should be mainly used in the test contexts.
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
	GetPublicDomains func(ctx context.Context, transport *HttpTransport) ([]string, error)
	GetSystemDomains func(ctx context.Context, transport *HttpTransport) ([]string, error)

	// Applications
	GetApplicationContext       func(ctx context.Context, appId string, transport *HttpTransport) (*MasheryApplication, error)
	GetApplicationPackageKeys   func(ctx context.Context, appId string, transport *HttpTransport) ([]MasheryPackageKey, error)
	CountApplicationPackageKeys func(ctx context.Context, appId string, c *HttpTransport) (int64, error)
	GetFullApplication          func(ctx context.Context, id string, c *HttpTransport) (*MasheryApplication, error)
	CreateApplication           func(ctx context.Context, memberId string, member MasheryApplication, c *HttpTransport) (*MasheryApplication, error)
	UpdateApplication           func(ctx context.Context, app MasheryApplication, c *HttpTransport) (*MasheryApplication, error)
	DeleteApplication           func(ctx context.Context, appId string, c *HttpTransport) error
	CountApplicationsOfMember   func(ctx context.Context, memberId string, c *HttpTransport) (int64, error)
	ListApplications            func(ctx context.Context, c *HttpTransport) ([]MasheryApplication, error)

	// Email template set
	GetEmailTemplateSet           func(ctx context.Context, id string, c *HttpTransport) (*MasheryEmailTemplateSet, error)
	ListEmailTemplateSets         func(ctx context.Context, c *HttpTransport) ([]MasheryEmailTemplateSet, error)
	ListEmailTemplateSetsFiltered func(ctx context.Context, params map[string]string, fields []string, c *HttpTransport) ([]MasheryEmailTemplateSet, error)

	// Endpoints
	ListEndpoints             func(ctx context.Context, serviceId string, c *HttpTransport) ([]AddressableV3Object, error)
	ListEndpointsWithFullInfo func(ctx context.Context, serviceId string, c *HttpTransport) ([]MasheryEndpoint, error)
	CreateEndpoint            func(ctx context.Context, serviceId string, endp MasheryEndpoint, c *HttpTransport) (*MasheryEndpoint, error)
	UpdateEndpoint            func(ctx context.Context, serviceId string, endp MasheryEndpoint, c *HttpTransport) (*MasheryEndpoint, error)
	GetEndpoint               func(ctx context.Context, serviceId string, endpointId string, c *HttpTransport) (*MasheryEndpoint, error)
	DeleteEndpoint            func(ctx context.Context, serviceId, endpointId string, c *HttpTransport) error
	CountEndpointsOf          func(ctx context.Context, serviceId string, c *HttpTransport) (int64, error)

	// Endpoint Methods
	ListEndpointMethods             func(ctx context.Context, serviceId, endpointId string, c *HttpTransport) ([]MasheryMethod, error)
	ListEndpointMethodsWithFullInfo func(ctx context.Context, serviceId, endpointId string, c *HttpTransport) ([]MasheryMethod, error)
	CreateEndpointMethod            func(ctx context.Context, serviceId, endpointId string, methoUpsert MasheryMethod, c *HttpTransport) (*MasheryMethod, error)
	UpdateEndpointMethod            func(ctx context.Context, serviceId, endpointId string, methUpsert MasheryMethod, c *HttpTransport) (*MasheryMethod, error)
	GetEndpointMethod               func(ctx context.Context, serviceId, endpointId, methodId string, c *HttpTransport) (*MasheryMethod, error)
	DeleteEndpointMethod            func(ctx context.Context, serviceId, endpointId, methodId string, c *HttpTransport) error
	CountEndpointsMethodsOf         func(ctx context.Context, serviceId, endpointId string, c *HttpTransport) (int64, error)

	// Endpoint method filters
	ListEndpointMethodFilters             func(ctx context.Context, serviceId, endpointId, methodId string, c *HttpTransport) ([]MasheryResponseFilter, error)
	ListEndpointMethodFiltersWithFullInfo func(ctx context.Context, serviceId, endpointId, methodId string, c *HttpTransport) ([]MasheryResponseFilter, error)
	CreateEndpointMethodFilter            func(ctx context.Context, serviceId, endpointId, methodId string, filterUpsert MasheryResponseFilter, c *HttpTransport) (*MasheryResponseFilter, error)
	UpdateEndpointMethodFilter            func(ctx context.Context, serviceId, endpointId, methodId string, methUpsert MasheryResponseFilter, c *HttpTransport) (*MasheryResponseFilter, error)
	GetEndpointMethodFilter               func(ctx context.Context, serviceId, endpointId, methodId, filterId string, c *HttpTransport) (*MasheryResponseFilter, error)
	DeleteEndpointMethodFilter            func(ctx context.Context, serviceId, endpointId, methodId, filterId string, c *HttpTransport) error
	CountEndpointsMethodsFiltersOf        func(ctx context.Context, serviceId, endpointId, methodId string, c *HttpTransport) (int64, error)

	// Member
	GetMember     func(ctx context.Context, id string, c *HttpTransport) (*MasheryMember, error)
	GetFullMember func(ctx context.Context, id string, c *HttpTransport) (*MasheryMember, error)
	CreateMember  func(ctx context.Context, member MasheryMember, c *HttpTransport) (*MasheryMember, error)
	UpdateMember  func(ctx context.Context, member MasheryMember, c *HttpTransport) (*MasheryMember, error)
	DeleteMember  func(ctx context.Context, memberId string, c *HttpTransport) error
	ListMembers   func(ctx context.Context, c *HttpTransport) ([]MasheryMember, error)

	// Packages
	GetPackage    func(ctx context.Context, id string, c *HttpTransport) (*MasheryPackage, error)
	CreatePackage func(ctx context.Context, pack MasheryPackage, c *HttpTransport) (*MasheryPackage, error)
	UpdatePackage func(ctx context.Context, pack MasheryPackage, c *HttpTransport) (*MasheryPackage, error)
	DeletePackage func(ctx context.Context, packId string, c *HttpTransport) error
	ListPackages  func(ctx context.Context, c *HttpTransport) ([]MasheryPackage, error)

	// Package plans
	CreatePlanService  func(ctx context.Context, planService MasheryPlanService, c *HttpTransport) (*AddressableV3Object, error)
	DeletePlanService  func(ctx context.Context, planService MasheryPlanService, c *HttpTransport) error
	CreatePlanEndpoint func(ctx context.Context, planEndp MasheryPlanServiceEndpoint, c *HttpTransport) (*AddressableV3Object, error)
	DeletePlanEndpoint func(ctx context.Context, planEndp MasheryPlanServiceEndpoint, c *HttpTransport) error
	ListPlanEndpoints  func(ctx context.Context, planService MasheryPlanService, c *HttpTransport) ([]AddressableV3Object, error)

	CountPlanEndpoints func(ctx context.Context, planService MasheryPlanService, c *HttpTransport) (int64, error)
	CountPlanService   func(ctx context.Context, packageId, planId string, c *HttpTransport) (int64, error)
	GetPlan            func(ctx context.Context, packageId string, planId string, c *HttpTransport) (*MasheryPlan, error)
	CreatePlan         func(ctx context.Context, packageId string, plan MasheryPlan, c *HttpTransport) (*MasheryPlan, error)
	UpdatePlan         func(ctx context.Context, plan MasheryPlan, c *HttpTransport) (*MasheryPlan, error)
	DeletePlan         func(ctx context.Context, packageId, planId string, c *HttpTransport) error
	CountPlans         func(ctx context.Context, packageId string, c *HttpTransport) (int64, error)
	ListPlans          func(ctx context.Context, packageId string, c *HttpTransport) ([]MasheryPlan, error)
	ListPlanServices   func(ctx context.Context, packageId string, planId string, c *HttpTransport) ([]MasheryService, error)

	// Plan methods
	ListPackagePlanMethods  func(ctx context.Context, id MasheryPlanServiceEndpoint, c *HttpTransport) ([]MasheryMethod, error)
	GetPackagePlanMethod    func(ctx context.Context, id MasheryPlanServiceEndpointMethod, c *HttpTransport) (*MasheryMethod, error)
	CreatePackagePlanMethod func(ctx context.Context, id MasheryPlanServiceEndpoint, upsert MasheryMethod, c *HttpTransport) (*MasheryMethod, error)
	DeletePackagePlanMethod func(ctx context.Context, id MasheryPlanServiceEndpointMethod, c *HttpTransport) error

	// Plan method filter
	GetPackagePlanMethodFilter    func(ctx context.Context, id MasheryPlanServiceEndpointMethod, c *HttpTransport) (*MasheryResponseFilter, error)
	CreatePackagePlanMethodFilter func(ctx context.Context, id MasheryPlanServiceEndpointMethod, ref MasheryServiceMethodFilter, c *HttpTransport) (*MasheryResponseFilter, error)
	DeletePackagePlanMethodFilter func(ctx context.Context, id MasheryPlanServiceEndpointMethod, c *HttpTransport) error

	// Packge key
	GetPackageKey           func(ctx context.Context, id string, c *HttpTransport) (*MasheryPackageKey, error)
	CreatePackageKey        func(ctx context.Context, appId string, packageKey MasheryPackageKey, c *HttpTransport) (*MasheryPackageKey, error)
	UpdatePackageKey        func(ctx context.Context, packageKey MasheryPackageKey, c *HttpTransport) (*MasheryPackageKey, error)
	DeletePackageKey        func(ctx context.Context, keyId string, c *HttpTransport) error
	ListPackageKeysFiltered func(ctx context.Context, params map[string]string, fields []string, c *HttpTransport) ([]MasheryPackageKey, error)
	ListPackageKeys         func(ctx context.Context, c *HttpTransport) ([]MasheryPackageKey, error)

	// Roles
	GetRole           func(ctx context.Context, id string, c *HttpTransport) (*MasheryRole, error)
	ListRoles         func(ctx context.Context, c *HttpTransport) ([]MasheryRole, error)
	ListRolesFiltered func(ctx context.Context, params map[string]string, fields []string, c *HttpTransport) ([]MasheryRole, error)

	// Services
	GetService           func(ctx context.Context, id string, c *HttpTransport) (*MasheryService, error)
	CreateService        func(ctx context.Context, service MasheryService, c *HttpTransport) (*MasheryService, error)
	UpdateService        func(ctx context.Context, service MasheryService, c *HttpTransport) (*MasheryService, error)
	DeleteService        func(ctx context.Context, serviceId string, c *HttpTransport) error
	ListServicesFiltered func(ctx context.Context, params map[string]string, fields []string, c *HttpTransport) ([]MasheryService, error)
	ListServices         func(ctx context.Context, c *HttpTransport) ([]MasheryService, error)
	CountServices        func(ctx context.Context, params map[string]string, c *HttpTransport) (int64, error)

	ListErrorSets         func(ctx context.Context, serviceId string, qs url.Values, c *HttpTransport) ([]MasheryErrorSet, error)
	GetErrorSet           func(ctx context.Context, serviceId, setId string, c *HttpTransport) (*MasheryErrorSet, error)
	CreateErrorSet        func(ctx context.Context, serviceId string, set MasheryErrorSet, c *HttpTransport) (*MasheryErrorSet, error)
	UpdateErrorSet        func(ctx context.Context, serviceId string, setData MasheryErrorSet, c *HttpTransport) (*MasheryErrorSet, error)
	DeleteErrorSet        func(ctx context.Context, serviceId, setId string, c *HttpTransport) error
	UpdateErrorSetMessage func(ctx context.Context, serviceId string, setId string, msg MasheryErrorMessage, c *HttpTransport) (*MasheryErrorMessage, error)

	GetServiceRoles func(ctx context.Context, serviceId string, c *HttpTransport) ([]MasheryRolePermission, error)
	SetServiceRoles func(ctx context.Context, id string, roles []MasheryRolePermission, c *HttpTransport) error

	// Service cache
	GetServiceCache    func(ctx context.Context, id string, c *HttpTransport) (*MasheryServiceCache, error)
	CreateServiceCache func(ctx context.Context, id string, service MasheryServiceCache, c *HttpTransport) (*MasheryServiceCache, error)
	UpdateServiceCache func(ctx context.Context, id string, service MasheryServiceCache, c *HttpTransport) (*MasheryServiceCache, error)
	DeleteServiceCache func(ctx context.Context, id string, c *HttpTransport) error

	// Service OAuth
	GetServiceOAuthSecurityProfile    func(ctx context.Context, id string, c *HttpTransport) (*MasheryOAuth, error)
	CreateServiceOAuthSecurityProfile func(ctx context.Context, id string, service MasheryOAuth, c *HttpTransport) (*MasheryOAuth, error)
	UpdateServiceOAuthSecurityProfile func(ctx context.Context, id string, service MasheryOAuth, c *HttpTransport) (*MasheryOAuth, error)
	DeleteServiceOAuthSecurityProfile func(ctx context.Context, id string, c *HttpTransport) error
}

func (c *PluggableClient) ListErrorSets(ctx context.Context, serviceId string, qs url.Values) ([]MasheryErrorSet, error) {
	if c.schema.ListErrorSets != nil {
		return c.schema.ListErrorSets(ctx, serviceId, qs, c.transport)
	} else {
		return []MasheryErrorSet{}, c.notImplemented("ListErrorSets")
	}
}

func (c *PluggableClient) GetErrorSet(ctx context.Context, serviceId, setId string) (*MasheryErrorSet, error) {
	if c.schema.GetErrorSet != nil {
		return c.schema.GetErrorSet(ctx, serviceId, setId, c.transport)
	} else {
		return nil, c.notImplemented("GetErrorSet")
	}
}

func (c *PluggableClient) CreateErrorSet(ctx context.Context, serviceId string, set MasheryErrorSet) (*MasheryErrorSet, error) {
	if c.schema.CreateErrorSet != nil {
		return c.schema.CreateErrorSet(ctx, serviceId, set, c.transport)
	} else {
		return nil, c.notImplemented("CreateErrorSet")
	}
}

func (c *PluggableClient) UpdateErrorSet(ctx context.Context, serviceId string, setData MasheryErrorSet) (*MasheryErrorSet, error) {
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

func (c *PluggableClient) UpdateErrorSetMessage(ctx context.Context, serviceId string, setId string, msg MasheryErrorMessage) (*MasheryErrorMessage, error) {
	if c.schema.UpdateErrorSetMessage != nil {
		return c.schema.UpdateErrorSetMessage(ctx, serviceId, setId, msg, c.transport)
	} else {
		return nil, c.notImplemented("UpdateErrorSetMessage")
	}
}

func (c *PluggableClient) GetServiceRoles(ctx context.Context, serviceId string) ([]MasheryRolePermission, error) {
	if c.schema.GetServiceRoles != nil {
		return c.schema.GetServiceRoles(ctx, serviceId, c.transport)
	} else {
		return []MasheryRolePermission{}, c.notImplemented("GetServiceRoles")
	}
}

func (c *PluggableClient) SetServiceRoles(ctx context.Context, serviceId string, perms []MasheryRolePermission) error {
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

func (c *PluggableClient) GetApplication(ctx context.Context, appId string) (*MasheryApplication, error) {
	if c.schema.GetApplicationContext != nil {
		return c.schema.GetApplicationContext(ctx, appId, c.transport)
	} else {
		return nil, c.notImplemented("GetApplication")
	}
}

func (c *PluggableClient) GetApplicationPackageKeys(ctx context.Context, appId string) ([]MasheryPackageKey, error) {
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

func (c *PluggableClient) GetFullApplication(ctx context.Context, id string) (*MasheryApplication, error) {
	if c.schema.GetFullApplication != nil {
		return c.schema.GetFullApplication(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetFullApplication")
	}
}

func (c *PluggableClient) CreateApplication(ctx context.Context, memberId string, member MasheryApplication) (*MasheryApplication, error) {
	if c.schema.CreateApplication != nil {
		return c.schema.CreateApplication(ctx, memberId, member, c.transport)
	} else {
		return nil, c.notImplemented("CreateApplication")
	}
}

func (c *PluggableClient) UpdateApplication(ctx context.Context, app MasheryApplication) (*MasheryApplication, error) {
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

func (c *PluggableClient) ListApplications(ctx context.Context) ([]MasheryApplication, error) {
	if c.schema.ListApplications != nil {
		return c.schema.ListApplications(ctx, c.transport)
	} else {
		return []MasheryApplication{}, c.notImplemented("ListApplications")
	}
}

// -----------------------------------------------------------------------------------------------------------------
// Email template set
// -----------------------------------------------------------------------------------------------------------------

func (c *PluggableClient) GetEmailTemplateSet(ctx context.Context, id string) (*MasheryEmailTemplateSet, error) {
	if c.schema.GetEmailTemplateSet != nil {
		return c.schema.GetEmailTemplateSet(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetEmailTemplateSet")
	}
}

func (c *PluggableClient) ListEmailTemplateSets(ctx context.Context) ([]MasheryEmailTemplateSet, error) {
	if c.schema.ListEmailTemplateSets != nil {
		return c.schema.ListEmailTemplateSets(ctx, c.transport)
	} else {
		return nil, c.notImplemented("ListEmailTemplateSets")
	}
}

func (c *PluggableClient) ListEmailTemplateSetsFiltered(ctx context.Context, params map[string]string, fields []string) ([]MasheryEmailTemplateSet, error) {
	if c.schema.ListEmailTemplateSetsFiltered != nil {
		return c.schema.ListEmailTemplateSetsFiltered(ctx, params, fields, c.transport)
	} else {
		return nil, c.notImplemented("ListEmailTemplateSetsFiltered")
	}
}

// -----------------------------------------------------------------------------------------------------------------
// Endpoints
// -----------------------------------------------------------------------------------------------------------------

func (c *PluggableClient) ListEndpoints(ctx context.Context, serviceId string) ([]AddressableV3Object, error) {
	if c.schema.ListEndpoints != nil {
		return c.schema.ListEndpoints(ctx, serviceId, c.transport)
	} else {
		return nil, c.notImplemented("ListEndpoints")
	}
}

func (c *PluggableClient) ListEndpointsWithFullInfo(ctx context.Context, serviceId string) ([]MasheryEndpoint, error) {
	if c.schema.ListEndpointsWithFullInfo != nil {
		return c.schema.ListEndpointsWithFullInfo(ctx, serviceId, c.transport)
	} else {
		return nil, c.notImplemented("ListEndpointsWithFullInfo")
	}
}

func (c *PluggableClient) CreateEndpoint(ctx context.Context, serviceId string, endp MasheryEndpoint) (*MasheryEndpoint, error) {
	if c.schema.CreateEndpoint != nil {
		return c.schema.CreateEndpoint(ctx, serviceId, endp, c.transport)
	} else {
		return nil, c.notImplemented("CreateEndpoint")
	}
}

func (c *PluggableClient) UpdateEndpoint(ctx context.Context, serviceId string, endp MasheryEndpoint) (*MasheryEndpoint, error) {
	if c.schema.UpdateEndpoint != nil {
		return c.schema.UpdateEndpoint(ctx, serviceId, endp, c.transport)
	} else {
		return nil, c.notImplemented("UpdateEndpoint")
	}
}

func (c *PluggableClient) GetEndpoint(ctx context.Context, serviceId string, endpointId string) (*MasheryEndpoint, error) {
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

func (c *PluggableClient) ListEndpointMethods(ctx context.Context, serviceId, endpointId string) ([]MasheryMethod, error) {
	if c.schema.ListEndpointMethods != nil {
		return c.schema.ListEndpointMethods(ctx, serviceId, endpointId, c.transport)
	} else {
		return []MasheryMethod{}, c.notImplemented("ListEndpointMethods")
	}
}

func (c *PluggableClient) ListEndpointMethodsWithFullInfo(ctx context.Context, serviceId, endpointId string) ([]MasheryMethod, error) {
	if c.schema.ListEndpointMethodsWithFullInfo != nil {
		return c.schema.ListEndpointMethodsWithFullInfo(ctx, serviceId, endpointId, c.transport)
	} else {
		return []MasheryMethod{}, c.notImplemented("ListEndpointMethodsWithFullInfo")
	}
}

func (c *PluggableClient) CreateEndpointMethod(ctx context.Context, serviceId, endpointId string, methoUpsert MasheryMethod) (*MasheryMethod, error) {
	if c.schema.CreateEndpointMethod != nil {
		return c.schema.CreateEndpointMethod(ctx, serviceId, endpointId, methoUpsert, c.transport)
	} else {
		return nil, c.notImplemented("CreateEndpointMethod")
	}
}

func (c *PluggableClient) UpdateEndpointMethod(ctx context.Context, serviceId, endpointId string, methUpsert MasheryMethod) (*MasheryMethod, error) {
	if c.schema.UpdateEndpointMethod != nil {
		return c.schema.UpdateEndpointMethod(ctx, serviceId, endpointId, methUpsert, c.transport)
	} else {
		return nil, c.notImplemented("UpdateEndpointMethod")
	}
}

func (c *PluggableClient) GetEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string) (*MasheryMethod, error) {
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

func (c *PluggableClient) ListEndpointMethodFilters(ctx context.Context, serviceId, endpointId, methodId string) ([]MasheryResponseFilter, error) {
	if c.schema.ListEndpointMethodFilters != nil {
		return c.schema.ListEndpointMethodFilters(ctx, serviceId, endpointId, methodId, c.transport)
	} else {
		return nil, c.notImplemented("ListEndpointMethodFilters")
	}
}

func (c *PluggableClient) ListEndpointMethodFiltersWithFullInfo(ctx context.Context, serviceId, endpointId, methodId string) ([]MasheryResponseFilter, error) {
	if c.schema.ListEndpointMethodFiltersWithFullInfo != nil {
		return c.schema.ListEndpointMethodFiltersWithFullInfo(ctx, serviceId, endpointId, methodId, c.transport)
	} else {
		return nil, c.notImplemented("ListEndpointMethodFiltersWithFullInfo")
	}
}

func (c *PluggableClient) CreateEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId string, filterUpsert MasheryResponseFilter) (*MasheryResponseFilter, error) {
	if c.schema.CreateEndpointMethodFilter != nil {
		return c.schema.CreateEndpointMethodFilter(ctx, serviceId, endpointId, methodId, filterUpsert, c.transport)
	} else {
		return nil, c.notImplemented("CreateEndpointMethodFilter")
	}
}

func (c *PluggableClient) UpdateEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId string, filterUpsert MasheryResponseFilter) (*MasheryResponseFilter, error) {
	if c.schema.UpdateEndpointMethodFilter != nil {
		return c.schema.UpdateEndpointMethodFilter(ctx, serviceId, endpointId, methodId, filterUpsert, c.transport)
	} else {
		return nil, c.notImplemented("UpdateEndpointMethodFilter")
	}
}

func (c *PluggableClient) GetEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId, filterId string) (*MasheryResponseFilter, error) {
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

func (c *PluggableClient) GetMember(ctx context.Context, id string) (*MasheryMember, error) {
	if c.schema.GetMember != nil {
		return c.schema.GetMember(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetMember")
	}
}

func (c *PluggableClient) GetFullMember(ctx context.Context, id string) (*MasheryMember, error) {
	if c.schema.GetFullMember != nil {
		return c.schema.GetFullMember(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetFullMember")
	}
}

func (c *PluggableClient) CreateMember(ctx context.Context, member MasheryMember) (*MasheryMember, error) {
	if c.schema.CreateMember != nil {
		return c.schema.CreateMember(ctx, member, c.transport)
	} else {
		return nil, c.notImplemented("CreateMember")
	}
}

func (c *PluggableClient) UpdateMember(ctx context.Context, member MasheryMember) (*MasheryMember, error) {
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

func (c *PluggableClient) ListMembers(ctx context.Context) ([]MasheryMember, error) {
	if c.schema.ListMembers != nil {
		return c.schema.ListMembers(ctx, c.transport)
	} else {
		return []MasheryMember{}, c.notImplemented("ListMembers")
	}
}

// ---------------------------------------------
// Packages

func (c *PluggableClient) GetPackage(ctx context.Context, id string) (*MasheryPackage, error) {
	if c.schema.GetPackage != nil {
		return c.schema.GetPackage(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetPackage")
	}
}

func (c *PluggableClient) CreatePackage(ctx context.Context, pack MasheryPackage) (*MasheryPackage, error) {
	if c.schema.CreatePackage != nil {
		return c.schema.CreatePackage(ctx, pack, c.transport)
	} else {
		return nil, c.notImplemented("CreatePackage")
	}
}

func (c *PluggableClient) UpdatePackage(ctx context.Context, pack MasheryPackage) (*MasheryPackage, error) {
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

func (c *PluggableClient) ListPackages(ctx context.Context) ([]MasheryPackage, error) {
	if c.schema.ListPackages != nil {
		return c.schema.ListPackages(ctx, c.transport)
	} else {
		return []MasheryPackage{}, c.notImplemented("ListPackages")
	}
}

// --------------------------------------
// Package plans

func (c *PluggableClient) CreatePlanService(ctx context.Context, planService MasheryPlanService) (*AddressableV3Object, error) {
	if c.schema.CreatePlanService != nil {
		return c.schema.CreatePlanService(ctx, planService, c.transport)
	} else {
		return nil, c.notImplemented("CreatePlanService")
	}
}

func (c *PluggableClient) DeletePlanService(ctx context.Context, planService MasheryPlanService) error {
	if c.schema.DeletePlanService != nil {
		return c.schema.DeletePlanService(ctx, planService, c.transport)
	} else {
		return c.notImplemented("CreatePlanService")
	}
}

func (c *PluggableClient) CreatePlanEndpoint(ctx context.Context, planEndp MasheryPlanServiceEndpoint) (*AddressableV3Object, error) {
	if c.schema.CreatePlanEndpoint != nil {
		return c.schema.CreatePlanEndpoint(ctx, planEndp, c.transport)
	} else {
		return nil, c.notImplemented("CreatePlanEndpoint")
	}
}

func (c *PluggableClient) DeletePlanEndpoint(ctx context.Context, planEndp MasheryPlanServiceEndpoint) error {
	if c.schema.DeletePlanEndpoint != nil {
		return c.schema.DeletePlanEndpoint(ctx, planEndp, c.transport)
	} else {
		return c.notImplemented("DeletePlanEndpoint")
	}
}

func (c *PluggableClient) ListPlanEndpoints(ctx context.Context, planService MasheryPlanService) ([]AddressableV3Object, error) {
	if c.schema.ListPlanEndpoints != nil {
		return c.schema.ListPlanEndpoints(ctx, planService, c.transport)
	} else {
		return []AddressableV3Object{}, c.notImplemented("ListPlanEndpoints")
	}
}

func (c *PluggableClient) ListPlanServices(ctx context.Context, packageId string, planId string) ([]MasheryService, error) {
	if c.schema.ListPlanServices != nil {
		return c.schema.ListPlanServices(ctx, packageId, planId, c.transport)
	} else {
		return []MasheryService{}, c.notImplemented("ListPlanServices")
	}
}

func (c *PluggableClient) ListPlans(ctx context.Context, packageId string) ([]MasheryPlan, error) {
	if c.schema.ListPlans != nil {
		return c.schema.ListPlans(ctx, packageId, c.transport)
	} else {
		return []MasheryPlan{}, c.notImplemented("ListPlans")
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

func (c *PluggableClient) UpdatePlan(ctx context.Context, plan MasheryPlan) (*MasheryPlan, error) {
	if c.schema.UpdatePlan != nil {
		return c.schema.UpdatePlan(ctx, plan, c.transport)
	} else {
		return nil, c.notImplemented("UpdatePlan")
	}
}

func (c *PluggableClient) CreatePlan(ctx context.Context, packageId string, plan MasheryPlan) (*MasheryPlan, error) {
	if c.schema.CreatePlan != nil {
		return c.schema.CreatePlan(ctx, packageId, plan, c.transport)
	} else {
		return nil, c.notImplemented("CreatePlan")
	}
}

func (c *PluggableClient) GetPlan(ctx context.Context, packageId string, planId string) (*MasheryPlan, error) {
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

func (c *PluggableClient) CountPlanEndpoints(ctx context.Context, planService MasheryPlanService) (int64, error) {
	if c.schema.CountPlanEndpoints != nil {
		return c.schema.CountPlanEndpoints(ctx, planService, c.transport)
	} else {
		return 0, c.notImplemented("CountPlanEndpoints")
	}
}

//---------------------------------------------
// Package plan methods

func (c *PluggableClient) ListPackagePlanMethods(ctx context.Context, id MasheryPlanServiceEndpoint) ([]MasheryMethod, error) {
	if c.schema.ListPackagePlanMethods != nil {
		return c.schema.ListPackagePlanMethods(ctx, id, c.transport)
	} else {
		return []MasheryMethod{}, c.notImplemented("ListPackagePlanMethods")
	}
}

func (c *PluggableClient) GetPackagePlanMethod(ctx context.Context, id MasheryPlanServiceEndpointMethod) (*MasheryMethod, error) {
	if c.schema.GetPackagePlanMethod != nil {
		return c.schema.GetPackagePlanMethod(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetPackagePlanMethod")
	}
}

func (c *PluggableClient) CreatePackagePlanMethod(ctx context.Context, id MasheryPlanServiceEndpoint, upsert MasheryMethod) (*MasheryMethod, error) {
	if c.schema.CreatePackagePlanMethod != nil {
		return c.schema.CreatePackagePlanMethod(ctx, id, upsert, c.transport)
	} else {
		return nil, c.notImplemented("CreatePackagePlanMethod")
	}
}

func (c *PluggableClient) DeletePackagePlanMethod(ctx context.Context, id MasheryPlanServiceEndpointMethod) error {
	if c.schema.DeletePackagePlanMethod != nil {
		return c.schema.DeletePackagePlanMethod(ctx, id, c.transport)
	} else {
		return c.notImplemented("CreatePackagePlanMethod")
	}
}

// ----------------------------------------
// Plan method filter

func (c *PluggableClient) GetPackagePlanMethodFilter(ctx context.Context, id MasheryPlanServiceEndpointMethod) (*MasheryResponseFilter, error) {
	if c.schema.GetPackagePlanMethodFilter != nil {
		return c.schema.GetPackagePlanMethodFilter(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetPackagePlanMethodFilter")
	}
}

func (c *PluggableClient) CreatePackagePlanMethodFilter(ctx context.Context, id MasheryPlanServiceEndpointMethod, ref MasheryServiceMethodFilter) (*MasheryResponseFilter, error) {
	if c.schema.CreatePackagePlanMethodFilter != nil {
		return c.schema.CreatePackagePlanMethodFilter(ctx, id, ref, c.transport)
	} else {
		return nil, c.notImplemented("CreatePackagePlanMethodFilter")
	}
}

func (c *PluggableClient) DeletePackagePlanMethodFilter(ctx context.Context, id MasheryPlanServiceEndpointMethod) error {
	if c.schema.DeletePackagePlanMethodFilter != nil {
		return c.schema.DeletePackagePlanMethodFilter(ctx, id, c.transport)
	} else {
		return c.notImplemented("DeletePackagePlanMethodFilter")
	}
}

// ------------------------------------------------------------
// Package key

func (c *PluggableClient) GetPackageKey(ctx context.Context, id string) (*MasheryPackageKey, error) {
	if c.schema.GetPackageKey != nil {
		return c.schema.GetPackageKey(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetPackageKey")
	}
}

func (c *PluggableClient) CreatePackageKey(ctx context.Context, appId string, packageKey MasheryPackageKey) (*MasheryPackageKey, error) {
	if c.schema.CreatePackageKey != nil {
		return c.schema.CreatePackageKey(ctx, appId, packageKey, c.transport)
	} else {
		return nil, c.notImplemented("CreatePackageKey")
	}
}

func (c *PluggableClient) UpdatePackageKey(ctx context.Context, packageKey MasheryPackageKey) (*MasheryPackageKey, error) {
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

func (c *PluggableClient) ListPackageKeysFiltered(ctx context.Context, params map[string]string, fields []string) ([]MasheryPackageKey, error) {
	if c.schema.ListPackageKeysFiltered != nil {
		return c.schema.ListPackageKeysFiltered(ctx, params, fields, c.transport)
	} else {
		return []MasheryPackageKey{}, c.notImplemented("ListPackageKeysFiltered")
	}
}

func (c *PluggableClient) ListPackageKeys(ctx context.Context) ([]MasheryPackageKey, error) {
	if c.schema.ListPackageKeys != nil {
		return c.schema.ListPackageKeys(ctx, c.transport)
	} else {
		return []MasheryPackageKey{}, c.notImplemented("ListPackageKeys")
	}
}

// ---------------------
// Roles

func (c *PluggableClient) GetRole(ctx context.Context, id string) (*MasheryRole, error) {
	if c.schema.GetRole != nil {
		return c.schema.GetRole(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetRole")
	}
}

func (c *PluggableClient) ListRoles(ctx context.Context) ([]MasheryRole, error) {
	if c.schema.ListRoles != nil {
		return c.schema.ListRoles(ctx, c.transport)
	} else {
		return []MasheryRole{}, c.notImplemented("ListRoles")
	}
}

func (c *PluggableClient) ListRolesFiltered(ctx context.Context, params map[string]string, fields []string) ([]MasheryRole, error) {
	if c.schema.ListRolesFiltered != nil {
		return c.schema.ListRolesFiltered(ctx, params, fields, c.transport)
	} else {
		return []MasheryRole{}, c.notImplemented("ListRolesFiltered")
	}
}

// ------------------------------
// Service

func (c *PluggableClient) GetService(ctx context.Context, id string) (*MasheryService, error) {
	if c.schema.GetService != nil {
		return c.schema.GetService(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetService")
	}
}

func (c *PluggableClient) CreateService(ctx context.Context, service MasheryService) (*MasheryService, error) {
	if c.schema.CreateService != nil {
		return c.schema.CreateService(ctx, service, c.transport)
	} else {
		return nil, c.notImplemented("CreateService")
	}
}

func (c *PluggableClient) UpdateService(ctx context.Context, service MasheryService) (*MasheryService, error) {
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

func (c *PluggableClient) ListServicesFiltered(ctx context.Context, params map[string]string, fields []string) ([]MasheryService, error) {
	if c.schema.ListServicesFiltered != nil {
		return c.schema.ListServicesFiltered(ctx, params, fields, c.transport)
	} else {
		return []MasheryService{}, c.notImplemented("ListServicesFiltered")
	}
}

func (c *PluggableClient) ListServices(ctx context.Context) ([]MasheryService, error) {
	if c.schema.ListServices != nil {
		return c.schema.ListServices(ctx, c.transport)
	} else {
		return []MasheryService{}, c.notImplemented("ListServices")
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

func (c *PluggableClient) GetServiceCache(ctx context.Context, id string) (*MasheryServiceCache, error) {
	if c.schema.GetServiceCache != nil {
		return c.schema.GetServiceCache(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetServiceCache")
	}
}

func (c *PluggableClient) CreateServiceCache(ctx context.Context, id string, service MasheryServiceCache) (*MasheryServiceCache, error) {
	if c.schema.CreateServiceCache != nil {
		return c.schema.CreateServiceCache(ctx, id, service, c.transport)
	} else {
		return nil, c.notImplemented("CreateServiceCache")
	}
}

func (c *PluggableClient) UpdateServiceCache(ctx context.Context, id string, service MasheryServiceCache) (*MasheryServiceCache, error) {
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
func (c *PluggableClient) GetServiceOAuthSecurityProfile(ctx context.Context, id string) (*MasheryOAuth, error) {
	if c.schema.GetServiceOAuthSecurityProfile != nil {
		return c.schema.GetServiceOAuthSecurityProfile(ctx, id, c.transport)
	} else {
		return nil, c.notImplemented("GetServiceOAuthSecurityProfile")
	}
}

func (c *PluggableClient) CreateServiceOAuthSecurityProfile(ctx context.Context, id string, service MasheryOAuth) (*MasheryOAuth, error) {
	if c.schema.CreateServiceOAuthSecurityProfile != nil {
		return c.schema.CreateServiceOAuthSecurityProfile(ctx, id, service, c.transport)
	} else {
		return nil, c.notImplemented("CreateServiceOAuthSecurityProfile")
	}
}

func (c *PluggableClient) UpdateServiceOAuthSecurityProfile(ctx context.Context, id string, service MasheryOAuth) (*MasheryOAuth, error) {
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
