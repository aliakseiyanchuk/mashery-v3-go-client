package v3client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Definition of the Mashery V3 Types used by Mashery.

type AccessTokenResponse struct {
	TokenType    string `json:"token_type"`
	ApiKey       string `json:"mapi"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type DomainAddress struct {
	Address string `json:"address"`
}

type Cache struct {
	ClientSurrogateControlEnabled bool     `json:"clientSurrogateControlEnabled"`
	ContentCacheKeyHeaders        []string `json:"contentCacheKeyHeaders"`
}

func (c *Cache) IsEmpty() bool {
	if c != nil {
		return len(c.ContentCacheKeyHeaders) == 0
	} else {
		return true
	}
}

type Cors struct {
	AllDomainsEnabled bool `json:"allDomainsEnabled"`
	MaxAge            int  `json:"maxAge"`
}

type ScheduledMaintenanceEvent struct {
	Id            string                `json:"id"`
	Name          string                `json:"name"`
	StartDateTime MasheryJSONTime       `json:"startDateTime"`
	EndDateTime   MasheryJSONTime       `json:"endDateTime"`
	Endpoints     []AddressableV3Object `json:"endpoints"`
}

type SystemDomainAuthentication struct {
	Type        string  `json:"type"`
	Username    *string `json:"username,omitempty"`
	Certificate *string `json:"certificate,omitempty"`
	Password    *string `json:"password,omitempty"`
}

type Processor struct {
	PreProcessEnabled  bool              `json:"preProcessEnabled"`
	PostProcessEnabled bool              `json:"postProcessEnabled"`
	PostInputs         map[string]string `json:"postInputs"`
	PreInputs          map[string]string `json:"preInputs"`
	Adapter            string            `json:"adapter"`
}

// Checks if the pre-processor structure is empty, i.e. doesn't convey any adapter information.
func (p *Processor) IsEmpty() bool {
	if p != nil {
		return len(p.Adapter) == 0 && !p.PreProcessEnabled && !p.PostProcessEnabled &&
			len(p.PreInputs) == 0 && len(p.PostInputs) == 0
	} else {
		return true
	}
}

type Domain struct {
	Address string `json:"address"`
}

type MasheryEndpoint struct {
	AddressableV3Object

	AllowMissingApiKey                         bool                        `json:"allowMissingApiKey"`
	ApiKeyValueLocationKey                     string                      `json:"apiKeyValueLocationKey"`
	ApiKeyValueLocations                       []string                    `json:"apiKeyValueLocations"`
	ApiMethodDetectionKey                      string                      `json:"apiMethodDetectionKey"`
	ApiMethodDetectionLocations                []string                    `json:"apiMethodDetectionLocations"`
	Cache                                      *Cache                      `json:"cache,omitempty"`
	ConnectionTimeoutForSystemDomainRequest    int                         `json:"connectionTimeoutForSystemDomainRequest"`
	ConnectionTimeoutForSystemDomainResponse   int                         `json:"connectionTimeoutForSystemDomainResponse"`
	CookiesDuringHttpRedirectsEnabled          bool                        `json:"cookiesDuringHttpRedirectsEnabled"`
	Cors                                       *Cors                       `json:"cors,omitempty"`
	CustomRequestAuthenticationAdapter         *string                     `json:"customRequestAuthenticationAdapter,omitempty"`
	DropApiKeyFromIncomingCall                 bool                        `json:"dropApiKeyFromIncomingCall"`
	ForceGzipOfBackendCall                     bool                        `json:"forceGzipOfBackendCall"`
	GzipPassthroughSupportEnabled              bool                        `json:"gzipPassthroughSupportEnabled"`
	HeadersToExcludeFromIncomingCall           []string                    `json:"headersToExcludeFromIncomingCall,omitempty"`
	HighSecurity                               bool                        `json:"highSecurity"`
	HostPassthroughIncludedInBackendCallHeader bool                        `json:"hostPassthroughIncludedInBackendCallHeader"`
	InboundSslRequired                         bool                        `json:"inboundSslRequired"`
	JsonpCallbackParameter                     string                      `json:"jsonpCallbackParameter"`
	JsonpCallbackParameterValue                string                      `json:"jsonpCallbackParameterValue"`
	ScheduledMaintenanceEvent                  *ScheduledMaintenanceEvent  `json:"scheduledMaintenanceEvent"`
	ForwardedHeaders                           []string                    `json:"forwardedHeaders,omitempty"`
	ReturnedHeaders                            []string                    `json:"returnedHeaders,omitempty"`
	Methods                                    *[]MasheryMethod            `json:"methods,omitempty"`
	NumberOfHttpRedirectsToFollow              int                         `json:"numberOfHttpRedirectsToFollow"`
	OutboundRequestTargetPath                  string                      `json:"outboundRequestTargetPath"`
	OutboundRequestTargetQueryParameters       string                      `json:"outboundRequestTargetQueryParameters"`
	OutboundTransportProtocol                  string                      `json:"outboundTransportProtocol"`
	Processor                                  *Processor                  `json:"processor"`
	PublicDomains                              []Domain                    `json:"publicDomains"`
	RequestAuthenticationType                  string                      `json:"requestAuthenticationType"`
	RequestPathAlias                           string                      `json:"requestPathAlias"`
	RequestProtocol                            string                      `json:"requestProtocol"`
	OAuthGrantTypes                            []string                    `json:"oauthGrantTypes"`
	StringsToTrimFromApiKey                    string                      `json:"stringsToTrimFromApiKey"`
	SupportedHttpMethods                       []string                    `json:"supportedHttpMethods"`
	SystemDomainAuthentication                 *SystemDomainAuthentication `json:"systemDomainAuthentication"`
	SystemDomains                              []Domain                    `json:"systemDomains"`
	TrafficManagerDomain                       string                      `json:"trafficManagerDomain"`
	UseSystemDomainCredentials                 bool                        `json:"useSystemDomainCredentials"`
	SystemDomainCredentialKey                  *string                     `json:"systemDomainCredentialKey"`
	SystemDomainCredentialSecret               *string                     `json:"systemDomainCredentialSecret"`
}

type MasheryOAuth struct {
	AccessTokenTtlEnabled       bool     `json:"accessTokenTtlEnabled"`
	AccessTokenTtl              int      `json:"accessTokenTtl"`
	AccessTokenType             string   `json:"accessTokenType"`
	AllowMultipleToken          bool     `json:"allowMultipleToken"`
	AuthorizationCodeTtl        int      `json:"authorizationCodeTtl"`
	ForwardedHeaders            []string `json:"forwardedHeaders"`
	MasheryTokenApiEnabled      bool     `json:"masheryTokenApiEnabled"`
	RefreshTokenEnabled         bool     `json:"refreshTokenEnabled"`
	EnableRefreshTokenTtl       bool     `json:"enableRefreshTokenTtl"`
	TokenBasedRateLimitsEnabled bool     `json:"tokenBasedRateLimitsEnabled"`
	ForceOauthRedirectUrl       bool     `json:"forceOauthRedirectUrl"`
	ForceSslRedirectUrlEnabled  bool     `json:"forceSslRedirectUrlEnabled"`
	GrantTypes                  []string `json:"grantTypes"`
	MACAlgorithm                string   `json:"macAlgorithm,omitempty"`
	QPSLimitCeiling             int64    `json:"qpsLimitCeiling"`
	RateLimitCeiling            int64    `json:"rateLimitCeiling"`
	RefreshTokenTtl             int64    `json:"refreshTokenTtl"`
	SecureTokensEnabled         bool     `json:"secureTokensEnabled"`
}

type MasherySecurityProfile struct {
	OAuth MasheryOAuth `json:"oauth"`
}

type MasheryErrorMessage struct {
	Id           string `json:"id"`
	Code         int    `json:"code"`
	Status       string `json:"status"`
	DetailHeader string `json:"detailHeader"`
	ResponseBody string `json:"responseBody"`
}

type MasheryErrorSet struct {
	AddressableV3Object
	Type          string                 `json:"type,omitempty"`
	JSONP         bool                   `json:"jsonp"`
	JSONPType     string                 `json:"jsonpType,omitempty"`
	ErrorMessages *[]MasheryErrorMessage `json:"errorMessages,omitempty"`
}

type MasheryJSONTime time.Time

// Id referenced data structure, used internally.
type IdReferenced struct {
	IdRef string `json:"id"`
}

type AddressableV3Object struct {
	Id      string           `json:"id,omitempty"`
	Name    string           `json:"name,omitempty"`
	Created *MasheryJSONTime `json:"created,omitempty"`
	Updated *MasheryJSONTime `json:"updated,omitempty"`
}

func (t *MasheryJSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if rv, err := time.Parse("2006-01-02T15:04:05Z0700", s); err == nil {
		*t = MasheryJSONTime(rv)
		return nil
	}

	if rv, err := time.Parse("2006-01-02T15:04:05Z07:00", s); err == nil {
		*t = MasheryJSONTime(rv)
		return nil
	}

	return errors.New(fmt.Sprintf("unknown Mashery JSON date: '%s'", s))
}

func (t *MasheryJSONTime) MarshalJSON() ([]byte, error) {
	return time.Time(*t).MarshalJSON()
}

func (t *MasheryJSONTime) ToString() string {
	if t == nil {
		return "null-time"
	}

	return time.Time(*t).Format(time.RFC3339)
}

type MasheryService struct {
	AddressableV3Object
	Cache             *MasheryServiceCache    `json:"cache,omitempty"`
	Endpoints         []MasheryEndpoint       `json:"endpoints,omitempty"`
	EditorHandle      string                  `json:"editorHandle,omitempty"`
	RevisionNumber    int                     `json:"revisionNumber,omitempty"`
	RobotsPolicy      string                  `json:"robotsPolicy,omitempty"`
	CrossdomainPolicy string                  `json:"crossdomainPolicy,omitempty"`
	Description       string                  `json:"description,omitempty"`
	ErrorSets         *[]MasheryErrorSet      `json:"errorSets,omitempty"`
	QpsLimitOverall   *int64                  `json:"qpsLimitOverall,omitempty"`
	RFC3986Encode     bool                    `json:"rfc3986Encode,omitempty"`
	SecurityProfile   *MasherySecurityProfile `json:"securityProfile,omitempty"`
	Version           string                  `json:"version,omitempty"`
	Roles             *[]MasheryRole          `json:"roles,omitempty"`
}

type MasheryServiceCache struct {
	CacheTtl int `json:"cacheTtl"`
}

// -----------------------------------------------------------------------------
// Mashery package-related types

type EAV map[string]string

type MasheryPlan struct {
	AddressableV3Object
	Description                       string `json:"description,omitempty"`
	Eav                               *EAV   `json:"eav,omitempty"`
	SelfServiceKeyProvisioningEnabled bool   `json:"selfServiceKeyProvisioningEnabled"`
	AdminKeyProvisioningEnabled       bool   `json:"adminKeyProvisioningEnabled"`
	Notes                             string `json:"notes,omitempty"`
	MaxNumKeysAllowed                 int    `json:"maxNumKeysAllowed"`
	NumKeysBeforeReview               int    `json:"numKeysBeforeReview"`
	QpsLimitCeiling                   *int64 `json:"qpsLimitCeiling,omitempty"`
	QpsLimitExempt                    bool   `json:"qpsLimitExempt"`
	QpsLimitKeyOverrideAllowed        bool   `json:"qpsLimitKeyOverrideAllowed"`
	RateLimitCeiling                  *int64 `json:"rateLimitCeiling,omitempty"`
	RateLimitExempt                   bool   `json:"rateLimitExempt"`
	RateLimitKeyOverrideAllowed       bool   `json:"rateLimitKeyOverrideAllowed"`
	RateLimitPeriod                   string `json:"rateLimitPeriod,omitempty"`
	ResponseFilterOverrideAllowed     bool   `json:"responseFilterOverrideAllowed"`
	Status                            string `json:"status,omitempty"`
	// Mashery's documentation erroneously quotes this field
	// as int. It is actually UUID of the email template set Id.
	EmailTemplateSetId string            `json:"emailTemplateSetId,omitempty"`
	Services           *[]MasheryService `json:"services,omitempty"`

	// Identity of the context object
	ParentPackageId string
}

type PlanFilter struct {
	AddressableV3Object
	Notes            string `json:"notes"`
	XmlFilterFields  string `json:"xmlFilterFields"`
	JsonFilterFields string `json:"jsonFilterFields"`
}

type MasheryPackage struct {
	AddressableV3Object
	Description                 string        `json:"description,omitempty"`
	NotifyDeveloperPeriod       string        `json:"notifyDeveloperPeriod,omitempty"`
	NotifyDeveloperNearQuota    bool          `json:"notifyDeveloperNearQuota"`
	NotifyDeveloperOverQuota    bool          `json:"notifyDeveloperOverQuota"`
	NotifyDeveloperOverThrottle bool          `json:"notifyDeveloperOverThrottle"`
	NotifyAdminPeriod           string        `json:"notifyAdminPeriod,omitempty"`
	NotifyAdminNearQuota        bool          `json:"notifyAdminNearQuota"`
	NotifyAdminOverQuota        bool          `json:"notifyAdminOverQuota"`
	NotifyAdminOverThrottle     bool          `json:"notifyAdminOverThrottle"`
	NotifyAdminEmails           string        `json:"notifyAdminEmails,omitempty"`
	NearQuotaThreshold          *int          `json:"nearQuotaThreshold,omitempty"`
	Eav                         EAV           `json:"eav,omitempty"`
	KeyAdapter                  string        `json:"keyAdapter,omitempty"`
	KeyLength                   *int          `json:"keyLength,omitempty"`
	SharedSecretLength          *int          `json:"sharedSecretLength,omitempty"`
	Plans                       []MasheryPlan `json:"plans,omitempty"`
}

// ---------------------------------------------------------------------------------
// Members, Applications, and Package keys

type Limit struct {
	Period  string `json:"period"`
	Source  string `json:"source"`
	Ceiling int64  `json:"ceiling"`
}

type MasheryPackageKey struct {
	AddressableV3Object
	Apikey           *string         `json:"apikey"`
	Secret           *string         `json:"secret"`
	RateLimitCeiling *int64          `json:"rateLimitCeiling"`
	RateLimitExempt  bool            `json:"rateLimitExempt"`
	QpsLimitCeiling  *int64          `json:"qpsLimitCeiling"`
	QpsLimitExempt   bool            `json:"qpsLimitExempt"`
	Status           string          `json:"status"`
	Limits           *[]Limit        `json:"limits"`
	Package          *MasheryPackage `json:"package"`
	Plan             *MasheryPlan    `json:"plan"`
}

func (mpk *MasheryPackageKey) LinksPackageAndPlan() bool {
	return mpk.Package != nil && mpk.Package.Id != "" &&
		mpk.Plan != nil && mpk.Plan.Id != ""
}

type MasheryApplication struct {
	AddressableV3Object

	Username          string               `json:"username"`
	Description       string               `json:"description,omitempty"`
	Type              string               `json:"type,omitempty"`
	Commercial        bool                 `json:"commercial"`
	Ads               bool                 `json:"ads"`
	AdsSystem         string               `json:"adsSystem,omitempty"`
	UsageModel        string               `json:"usageModel,omitempty"`
	Tags              string               `json:"tags,omitempty"`
	Notes             string               `json:"notes,omitempty"`
	HowDidYouHear     string               `json:"howDidYouHear,omitempty"`
	PreferredProtocol string               `json:"preferredProtocol,omitempty"`
	PreferredOutput   string               `json:"preferredOutput,omitempty"`
	ExternalId        string               `json:"externalId,omitempty"`
	Uri               string               `json:"uri,omitempty"`
	OAuthRedirectUri  string               `json:"oauthRedirectUri,omitempty"`
	PackageKeys       *[]MasheryPackageKey `json:"packageKeys,omitempty"`
	Eav               *EAV                 `json:"eav,omitempty"`
}

type MasheryRole struct {
	AddressableV3Object
	Description string `json:"description,omitempty"`
	Predefined  bool   `json:"isPredefined"`
	OrgRole     bool   `json:"isOrgrole"`
	Assignable  bool   `json:"isAssignable"`
}

type MasheryMember struct {
	AddressableV3Object

	Username     string                `json:"username"`
	Email        string                `json:"email"`
	DisplayName  string                `json:"displayName,omitempty"`
	Uri          string                `json:"uri,omitempty"`
	Blog         string                `json:"blog,omitempty"`
	Im           string                `json:"im,omitempty"`
	Imsvc        string                `json:"imsvc,omitempty"`
	Phone        string                `json:"phone,omitempty"`
	Company      string                `json:"company,omitempty"`
	Address1     string                `json:"address1,omitempty"`
	Address2     string                `json:"address2,omitempty"`
	Locality     string                `json:"locality,omitempty"`
	Region       string                `json:"region,omitempty"`
	PostalCode   string                `json:"postalCode,omitempty"`
	CountryCode  string                `json:"countryCode,omitempty"`
	FirstName    string                `json:"firstName,omitempty"`
	LastName     string                `json:"lastName,omitempty"`
	AreaStatus   string                `json:"areaStatus,omitempty"`
	ExternalId   string                `json:"externalId,omitempty"`
	PasswdNew    *string               `json:"passwdNew,omitempty"`
	Applications *[]MasheryApplication `json:"applications,omitempty"`
	PackageKeys  *[]MasheryPackageKey  `json:"packageKeys,omitempty"`
	Roles        *[]MasheryRole        `json:"roles,omitempty"`
}

// -------------------------------------------------------------------
// Synthetic path identifier
//
// Service-related

type MasheryServiceEndpoint struct {
	ServiceId  string
	EndpointId string
}
type MasheryServiceMethod struct {
	MasheryServiceEndpoint
	MethodId string
}
type MasheryServiceMethodFilter struct {
	MasheryServiceMethod
	FilterId string
}

// Package-plan related
type MasheryPlanService struct {
	PackageId string
	PlanId    string
	ServiceId string
}

type MasheryPlanServiceEndpoint struct {
	MasheryPlanService
	EndpointId string
}

type MasheryPlanServiceEndpointMethod struct {
	MasheryPlanServiceEndpoint
	MethodId string
}

type MasheryPlanServiceEndpointMethodFilter struct {
	MasheryPlanServiceEndpointMethod
	FilterId string
}

// MasheryEmailTemplateSet Mashery email template set
type MasheryEmailTemplateSet struct {
	AddressableV3Object
	Type           string                  `json:"type,omitempty"`
	EmailTemplates *[]MasheryEmailTemplate `json:"emailTemplates,omitempty"`
}

// MasheryEmailTemplate Mashery email template
type MasheryEmailTemplate struct {
	AddressableV3Object
	Type    string `json:"type"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// -----------------------------------------------------------------------------
// Mashery Methods

type MasheryMethod struct {
	AddressableV3Object
	SampleJsonResponse string `json:"sampleJsonResponse,omitempty"`
	SampleXmlResponse  string `json:"sampleXmlResponse,omitempty"`
}

func ParseMasheryMethod(inp []byte) (interface{}, int, error) {
	var rv MasheryMethod
	err := json.Unmarshal(inp, &rv)
	return rv, 1, err
}

func ParseMasheryMethodArray(inp []byte) (interface{}, int, error) {
	var rv []MasheryMethod
	err := json.Unmarshal(inp, &rv)
	return rv, len(rv), err
}

type MasheryResponseFilter struct {
	AddressableV3Object
	Notes            string `json:"notes"`
	XmlFilterFields  string `json:"xmlFilterFields,omitempty"`
	JsonFilterFields string `json:"jsonFilterFields,omitempty"`
}

func ParseMasheryResponseFilterArray(inp []byte) (interface{}, int, error) {
	var rv []MasheryResponseFilter
	err := json.Unmarshal(inp, &rv)
	return rv, len(rv), err
}

func ParseMasheryResponseFilter(inp []byte) (interface{}, int, error) {
	var rv MasheryResponseFilter
	err := json.Unmarshal(inp, &rv)
	return rv, 1, err
}

// -----------------------------------------------------------------------------
// Type conversion functions

func (m *MasheryService) AddressableEndpoints() []AddressableV3Object {
	return AddressableEndpoints(m.Endpoints)
}

func AddressableEndpoints(inp []MasheryEndpoint) []AddressableV3Object {
	rv := make([]AddressableV3Object, len(inp))
	for i, v := range inp {
		rv[i] = AddressableV3Object{
			Id:      v.Id,
			Name:    v.Name,
			Created: v.Created,
			Updated: v.Updated,
		}
	}

	return rv
}

// ----------------------------------------------
// Errors

type V3GenericErrorResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (e *V3GenericErrorResponse) hasData() bool {
	return len(e.ErrorCode) > 0 || len(e.ErrorMessage) > 0
}

func (e *V3GenericErrorResponse) Error() string {
	return fmt.Sprintf("%s (Mashery V3 code %s)", e.ErrorMessage, e.ErrorCode)
}

type V3PropertyErrorMessage struct {
	Property string `json:"property"`
	Message  string `json:"message"`
}

type V3PropertyErrorMessages struct {
	Errors []V3PropertyErrorMessage `json:"errors"`
}

func (e *V3PropertyErrorMessages) Error() string {
	if len(e.Errors) == 0 {
		return "no errors returned in the error response"
	}

	sb := strings.Builder{}
	for _, e := range e.Errors {
		sb.WriteString("error in property ")
		sb.WriteString(e.Property)
		sb.WriteString(": ")
		sb.WriteString(e.Message)
		sb.WriteString(";")
	}

	return sb.String()
}

type V3UndeterminedError struct {
	Code   int
	Header http.Header
	Body   []byte
}

func (e *V3UndeterminedError) Error() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("server responded with code %d; ", e.Code))

	for k, v := range e.Header {
		sb.WriteString(fmt.Sprintf("%s::%s", k, strings.Join(v, ",")))
	}

	if len(e.Body) > 0 {
		sb.WriteString("Body:")
		sb.WriteString(string(e.Body))
	} else {
		sb.WriteString("(Empty body)")
	}

	return sb.String()
}

// --------------------------------------------
// Factory routines

func ParseMasheryAddressableObject(dat []byte) (interface{}, int, error) {
	var rv AddressableV3Object
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryAddressableObjectsArray(dat []byte) (interface{}, int, error) {
	var rv []AddressableV3Object
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryEndpoint(dat []byte) (interface{}, int, error) {
	var rv MasheryEndpoint
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryEndpointArray(dat []byte) (interface{}, int, error) {
	var rv []MasheryEndpoint
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryService(dat []byte) (interface{}, int, error) {
	var rv MasheryService
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryServiceCache(dat []byte) (interface{}, int, error) {
	var rv MasheryServiceCache
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryServiceSecurityProfileOAuth(dat []byte) (interface{}, int, error) {
	var rv MasheryOAuth
	err := json.Unmarshal(dat, &rv)

	return rv, 1, err
}

func ParseMasheryServiceArray(dat []byte) (interface{}, int, error) {
	var rv []MasheryService
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseServiceErrorSetArray(dat []byte) (interface{}, int, error) {
	var rv []MasheryErrorSet
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseErrorSet(dat []byte) (interface{}, int, error) {
	var rv MasheryErrorSet
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseErrorSetMessage(dat []byte) (interface{}, int, error) {
	var rv MasheryErrorMessage
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryPackage(dat []byte) (interface{}, int, error) {
	var rv MasheryPackage
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryPackageArray(dat []byte) (interface{}, int, error) {
	var rv []MasheryPackage
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryPlan(dat []byte) (interface{}, int, error) {
	var rv MasheryPlan
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryPlanArray(dat []byte) (interface{}, int, error) {
	var rv []MasheryPlan
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryPackageKey(dat []byte) (interface{}, int, error) {
	var rv MasheryPackageKey
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryPackageKeyArray(dat []byte) (interface{}, int, error) {
	var rv []MasheryPackageKey
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryApplication(dat []byte) (interface{}, int, error) {
	var rv MasheryApplication
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryApplicationArray(dat []byte) (interface{}, int, error) {
	var rv []MasheryApplication
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryMember(dat []byte) (interface{}, int, error) {
	var rv MasheryMember
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryMemberArray(dat []byte) (interface{}, int, error) {
	var rv []MasheryMember
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryRoleArray(dat []byte) (interface{}, int, error) {
	var rv []MasheryRole
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryRole(dat []byte) (interface{}, int, error) {
	var rv MasheryRole
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryEmailTemplateSetArray(dat []byte) (interface{}, int, error) {
	var rv []MasheryEmailTemplateSet
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryEmailTemplateSet(dat []byte) (interface{}, int, error) {
	var rv MasheryEmailTemplateSet
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryDomainAddressArray(dat []byte) (interface{}, int, error) {
	var rv []DomainAddress
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}
