package v3client

import (
	"net/url"
	"strings"
)

const (
	MasheryTokenSystemProperty = "MasheryV3AccessToken"
	MasheryTokenEndpoint       = "https://api.mashery.com/v3/token"
)

var EmptyQuery url.Values = map[string][]string{}

// MasheryServiceFields Usual mashery service fields without endpoints
var MasheryServiceFields = []string{
	"id", "name", "created", "updated", "editorHandle", "revisionNumber", "robotsPolicy",
	"crossdomainPolicy", "description", "errorSets", "qpsLimitOverall", "rfc3986Encode",
	"securityProfile", "version", "cache",
}

// MasheryErrorSetFields all fields of Mashery error set
var MasheryErrorSetFields = []string{
	"id", "name", "type", "jsonp", "jsonpType", "errorMessages",
}

var MasheryEndpointFields = []string{
	"id", "allowMissingApiKey", "apiKeyValueLocationKey", "apiKeyValueLocations",
	"apiMethodDetectionKey", "apiMethodDetectionLocations", "cache", "connectionTimeoutForSystemDomainRequest",
	"connectionTimeoutForSystemDomainResponse", "cookiesDuringHttpRedirectsEnabled", "cors",
	"created", "customRequestAuthenticationAdapter", "dropApiKeyFromIncomingCall", "forceGzipOfBackendCall",
	"gzipPassthroughSupportEnabled", "headersToExcludeFromIncomingCall", "highSecurity",
	"hostPassthroughIncludedInBackendCallHeader", "inboundSslRequired", "inboundMsslRequired", "jsonpCallbackParameter",
	"jsonpCallbackParameterValue", "scheduledMaintenanceEvent", "forwardedHeaders", "returnedHeaders",
	"methods", "name", "numberOfHttpRedirectsToFollow", "outboundRequestTargetPath", "outboundRequestTargetQueryParameters",
	"outboundTransportProtocol", "processor", "publicDomains", "requestAuthenticationType",
	"requestPathAlias", "requestProtocol", "oauthGrantTypes", "stringsToTrimFromApiKey", "supportedHttpMethods",
	"systemDomainAuthentication", "systemDomains", "trafficManagerDomain", "updated", "useSystemDomainCredentials",
	"systemDomainCredentialKey", "systemDomainCredentialSecret",
}

var MasheryPlanFields = []string{
	"id", "created", "updated", "name", "description", "eav", "selfServiceKeyProvisioningEnabled",
	"adminKeyProvisioningEnabled", "notes", "maxNumKeysAllowed", "numKeysBeforeReview", "qpsLimitCeiling", "qpsLimitExempt",
	"qpsLimitKeyOverrideAllowed", "rateLimitCeiling", "rateLimitExempt", "rateLimitKeyOverrideAllowed", "rateLimitPeriod",
	"responseFilterOverrideAllowed", "status", "emailTemplateSetId", "services.id", "services.endpoints",
}

var MasheryServiceFieldsWithEndpoinds = append(MasheryServiceFields, "endpoints")

var MasheryPackageFields = []string{
	"id", "name", "created", "updated", "description", "notifyDeveloperPeriod",
	"notifyDeveloperNearQuota", "notifyDeveloperOverQuota", "notifyDeveloperOverThrottle", "notifyAdminNearQuota",
	"notifyAdminPeriod", "notifyAdminOverQuota", "notifyAdminOverThrottle", "notifyAdminEmails",
	"nearQuotaThreshold", "eav", "keyAdapter", "keyLength",
	"sharedSecretLength", "plans",
}
var MasheryPackageFieldsStr = strings.Join(MasheryPackageFields, ",")

// Joined strings, representing all fields of Mashery service. Used (frequently) in
// {@link HttpTransport#GetService}
var MasheryServiceFullFieldsStr = strings.Join(MasheryServiceFieldsWithEndpoinds, ",")
var MasheryEndpointFieldsStr = strings.Join(MasheryEndpointFields, ",")
var MasheryPlanFieldsStr = strings.Join(MasheryPlanFields, ",")

var memberFields = []string{
	"id", "username", "created", "updated", "email", "displayName", "uri", "blog", "im", "imsvc", "phone",
	"company", "address1", "address2", "locality", "region", "postalCode", "countryCode", "firstName",
	"lastName", "areaStatus", "externalId",
}

var memberFieldsStr = strings.Join(memberFields, ",")

var memberDeepFields = append(memberFields,
	"applications", "packageKeys", "roles")

var MasheryPackageKeyFields = []string{
	"id", "apikey", "secret", "created", "updated", "rateLimitCeiling", "rateLimitExempt", "qpsLimitCeiling",
	"qpsLimitExempt", "status", "limits",
}

var MasheryPackageKeyFullFields = append(MasheryPackageKeyFields, "package", "plan")

var MasheryPackageKeyFieldsStr = strings.Join(MasheryPackageKeyFields, ",")
var MasheryPackageKeyFullFieldsStr = strings.Join(MasheryPackageKeyFullFields, ",")

var MasheryEmailTemplateSetFields = []string{"id", "created", "updated", "name", "type", "emailTemplates"}
var MasheryEmailTemplateSetFieldsStr = strings.Join(MasheryEmailTemplateSetFields, ",")

var MasheryMethodsFields = []string{"id", "name", "created", "updated", "sampleJsonResponse", "sampleXmlResponse"}
var MasheryMethodsFieldsStr = strings.Join(MasheryMethodsFields, ",")

var MasheryResponseFilterFields = []string{"id", "name", "created", "updated", "notes", "notes", "xmlFilterFields", "jsonFilterFields"}
var MasheryResponseFilterFieldsStr = strings.Join(MasheryResponseFilterFields, ",")
