package masherytypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
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

type Endpoint struct {
	AddressableV3Object

	AllowMissingApiKey                         bool                        `json:"allowMissingApiKey"`
	ApiKeyValueLocationKey                     string                      `json:"apiKeyValueLocationKey"`
	ApiKeyValueLocations                       []string                    `json:"apiKeyValueLocations"`
	ApiMethodDetectionKey                      string                      `json:"apiMethodDetectionKey"`
	ApiMethodDetectionLocations                []string                    `json:"apiMethodDetectionLocations,omitempty"`
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
	InboundMutualSslRequired                   bool                        `json:"inboundMsslRequired"`
	JsonpCallbackParameter                     string                      `json:"jsonpCallbackParameter,omitempty"`
	JsonpCallbackParameterValue                string                      `json:"jsonpCallbackParameterValue,omitempty"`
	ScheduledMaintenanceEvent                  *ScheduledMaintenanceEvent  `json:"scheduledMaintenanceEvent,omitempty"`
	ForwardedHeaders                           []string                    `json:"forwardedHeaders,omitempty"`
	ReturnedHeaders                            []string                    `json:"returnedHeaders,omitempty"`
	Methods                                    *[]ServiceEndpointMethod    `json:"methods,omitempty"`
	NumberOfHttpRedirectsToFollow              int                         `json:"numberOfHttpRedirectsToFollow,omitempty"`
	OutboundRequestTargetPath                  string                      `json:"outboundRequestTargetPath,omitempty"`
	OutboundRequestTargetQueryParameters       string                      `json:"outboundRequestTargetQueryParameters,omitempty"`
	OutboundTransportProtocol                  string                      `json:"outboundTransportProtocol,omitempty"`
	Processor                                  *Processor                  `json:"processor,omitempty"`
	PublicDomains                              []Domain                    `json:"publicDomains,omitempty"`
	RequestAuthenticationType                  string                      `json:"requestAuthenticationType,omitempty"`
	RequestPathAlias                           string                      `json:"requestPathAlias,omitempty"`
	RequestProtocol                            string                      `json:"requestProtocol,omitempty"`
	OAuthGrantTypes                            []string                    `json:"oauthGrantTypes,omitempty"`
	StringsToTrimFromApiKey                    string                      `json:"stringsToTrimFromApiKey,omitempty"`
	SupportedHttpMethods                       []string                    `json:"supportedHttpMethods,omitempty"`
	SystemDomainAuthentication                 *SystemDomainAuthentication `json:"systemDomainAuthentication,omitempty"`
	SystemDomains                              []Domain                    `json:"systemDomains,omitempty"`
	TrafficManagerDomain                       string                      `json:"trafficManagerDomain"`
	UseSystemDomainCredentials                 bool                        `json:"useSystemDomainCredentials"`
	SystemDomainCredentialKey                  *string                     `json:"systemDomainCredentialKey,omitempty"`
	SystemDomainCredentialSecret               *string                     `json:"systemDomainCredentialSecret,omitempty"`

	ParentServiceId ServiceIdentifier
}

func (e *Endpoint) Identifier() ServiceEndpointIdentifier {
	return ServiceEndpointIdentifier{
		ServiceIdentifier: e.ParentServiceId,
		EndpointId:        e.Id,
	}
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

	ParentService ServiceIdentifier `json:"-"`
}

type MasherySecurityProfile struct {
	OAuth *MasheryOAuth `json:"oauth,omitempty"`
}

type MasheryErrorMessage struct {
	Id           string `json:"id"`
	Code         int    `json:"code"`
	Status       string `json:"status"`
	DetailHeader string `json:"detailHeader"`
	ResponseBody string `json:"responseBody"`
}

type ErrorSet struct {
	AddressableV3Object
	Type          string                 `json:"type,omitempty"`
	JSONP         bool                   `json:"jsonp"`
	JSONPType     string                 `json:"jsonpType,omitempty"`
	ErrorMessages *[]MasheryErrorMessage `json:"errorMessages,omitempty"`

	ParentServiceId ServiceIdentifier
}

func (es *ErrorSet) Identifier() ErrorSetIdentifier {
	return ErrorSetIdentifier{
		ServiceIdentifier: es.ParentServiceId,
		ErrorSetId:        es.Id,
	}
}

type MasheryJSONTime time.Time

// IdReferenced Id referenced data structure, used internally.
type IdReferenced struct {
	IdRef string `json:"id"`
}

// AddressableV3Object base properties of any object that is read through Mashery V3 API
type AddressableV3Object struct {
	Id        string           `json:"id,omitempty"`
	Name      string           `json:"name,omitempty"`
	Created   *MasheryJSONTime `json:"created,omitempty"`
	Updated   *MasheryJSONTime `json:"updated,omitempty"`
	Retrieved time.Time        `json:"-"`
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

type Service struct {
	AddressableV3Object
	Cache             *ServiceCache           `json:"cache,omitempty"`
	Endpoints         []Endpoint              `json:"endpoints,omitempty"`
	EditorHandle      string                  `json:"editorHandle,omitempty"`
	RevisionNumber    int                     `json:"revisionNumber,omitempty"`
	RobotsPolicy      string                  `json:"robotsPolicy,omitempty"`
	CrossdomainPolicy string                  `json:"crossdomainPolicy,omitempty"`
	Description       string                  `json:"description,omitempty"`
	ErrorSets         *[]ErrorSet             `json:"errorSets,omitempty"`
	QpsLimitOverall   *int64                  `json:"qpsLimitOverall,omitempty"`
	RFC3986Encode     bool                    `json:"rfc3986Encode,omitempty"`
	SecurityProfile   *MasherySecurityProfile `json:"securityProfile,omitempty"`
	Version           string                  `json:"version,omitempty"`
	Roles             *[]RolePermission       `json:"roles,omitempty"`
}

func (s *Service) Identifier() ServiceIdentifier {
	return ServiceIdentifier{
		ServiceId: s.Id,
	}
}

type ServiceCache struct {
	CacheTtl int `json:"cacheTtl"`
}

// -----------------------------------------------------------------------------
// Mashery package-related types

type EAV map[string]string

type Plan struct {
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
	EmailTemplateSetId string     `json:"emailTemplateSetId,omitempty"`
	Services           *[]Service `json:"services,omitempty"`

	// Identity of the context object
	ParentPackageId PackageIdentifier
}

func (p *Plan) Identifier() PackagePlanIdentifier {
	return PackagePlanIdentifier{
		PackageIdentifier: p.ParentPackageId,
		PlanId:            p.Id,
	}
}

type PlanFilter struct {
	AddressableV3Object
	Notes            string `json:"notes"`
	XmlFilterFields  string `json:"xmlFilterFields"`
	JsonFilterFields string `json:"jsonFilterFields"`
}

type Package struct {
	AddressableV3Object
	Description                 string `json:"description,omitempty"`
	NotifyDeveloperPeriod       string `json:"notifyDeveloperPeriod,omitempty"`
	NotifyDeveloperNearQuota    bool   `json:"notifyDeveloperNearQuota"`
	NotifyDeveloperOverQuota    bool   `json:"notifyDeveloperOverQuota"`
	NotifyDeveloperOverThrottle bool   `json:"notifyDeveloperOverThrottle"`
	NotifyAdminPeriod           string `json:"notifyAdminPeriod,omitempty"`
	NotifyAdminNearQuota        bool   `json:"notifyAdminNearQuota"`
	NotifyAdminOverQuota        bool   `json:"notifyAdminOverQuota"`
	NotifyAdminOverThrottle     bool   `json:"notifyAdminOverThrottle"`
	NotifyAdminEmails           string `json:"notifyAdminEmails,omitempty"`
	NearQuotaThreshold          *int   `json:"nearQuotaThreshold,omitempty"`
	Eav                         EAV    `json:"eav,omitempty"`
	KeyAdapter                  string `json:"keyAdapter,omitempty"`
	KeyLength                   *int   `json:"keyLength,omitempty"`
	SharedSecretLength          *int   `json:"sharedSecretLength,omitempty"`
	Plans                       []Plan `json:"plans,omitempty"`
}

func (p *Package) Identifier() PackageIdentifier {
	return PackageIdentifier{PackageId: p.Id}
}

// ---------------------------------------------------------------------------------
// Members, Applications, and Package keys

type Limit struct {
	Period  string `json:"period"`
	Source  string `json:"source"`
	Ceiling int64  `json:"ceiling"`
}

type PackageKey struct {
	AddressableV3Object
	Apikey           *string  `json:"apikey"`
	Secret           *string  `json:"secret"`
	RateLimitCeiling *int64   `json:"rateLimitCeiling"`
	RateLimitExempt  bool     `json:"rateLimitExempt"`
	QpsLimitCeiling  *int64   `json:"qpsLimitCeiling"`
	QpsLimitExempt   bool     `json:"qpsLimitExempt"`
	Status           string   `json:"status"`
	Limits           *[]Limit `json:"limits"`
	Package          *Package `json:"package"`
	Plan             *Plan    `json:"plan"`
}

func (mpk *PackageKey) Identifier() PackageKeyIdentifier {
	return PackageKeyIdentifier{PackageKeyId: mpk.Id}
}

func (mpk *PackageKey) LinksPackageAndPlan() bool {
	return mpk.Package != nil && mpk.Package.Id != "" &&
		mpk.Plan != nil && mpk.Plan.Id != ""
}

type Application struct {
	AddressableV3Object

	Username          string        `json:"username"`
	Description       string        `json:"description,omitempty"`
	Type              string        `json:"type,omitempty"`
	Commercial        bool          `json:"commercial"`
	Ads               bool          `json:"ads"`
	AdsSystem         string        `json:"adsSystem,omitempty"`
	UsageModel        string        `json:"usageModel,omitempty"`
	Tags              string        `json:"tags,omitempty"`
	Notes             string        `json:"notes,omitempty"`
	HowDidYouHear     string        `json:"howDidYouHear,omitempty"`
	PreferredProtocol string        `json:"preferredProtocol,omitempty"`
	PreferredOutput   string        `json:"preferredOutput,omitempty"`
	ExternalId        string        `json:"externalId,omitempty"`
	Uri               string        `json:"uri,omitempty"`
	OAuthRedirectUri  string        `json:"oauthRedirectUri,omitempty"`
	PackageKeys       *[]PackageKey `json:"packageKeys,omitempty"`
	Eav               *EAV          `json:"eav,omitempty"`
}

func (a *Application) Identifier() ApplicationIdentifier {
	return ApplicationIdentifier{
		ApplicationId: a.Id,
	}
}

type Role struct {
	AddressableV3Object
	Description string `json:"description,omitempty"`
	Predefined  bool   `json:"isPredefined,omitempty"`
	OrgRole     bool   `json:"isOrgrole,omitempty"`
	Assignable  bool   `json:"isAssignable,omitempty"`
}

type RolePermission struct {
	Role
	Action string `json:"action"`
}

type Member struct {
	AddressableV3Object

	Username     string         `json:"username"`
	Email        string         `json:"email"`
	DisplayName  string         `json:"displayName,omitempty"`
	Uri          string         `json:"uri,omitempty"`
	Blog         string         `json:"blog,omitempty"`
	Im           string         `json:"im,omitempty"`
	Imsvc        string         `json:"imsvc,omitempty"`
	Phone        string         `json:"phone,omitempty"`
	Company      string         `json:"company,omitempty"`
	Address1     string         `json:"address1,omitempty"`
	Address2     string         `json:"address2,omitempty"`
	Locality     string         `json:"locality,omitempty"`
	Region       string         `json:"region,omitempty"`
	PostalCode   string         `json:"postalCode,omitempty"`
	CountryCode  string         `json:"countryCode,omitempty"`
	FirstName    string         `json:"firstName,omitempty"`
	LastName     string         `json:"lastName,omitempty"`
	AreaStatus   string         `json:"areaStatus,omitempty"`
	ExternalId   string         `json:"externalId,omitempty"`
	PasswdNew    *string        `json:"passwdNew,omitempty"`
	Applications *[]Application `json:"applications,omitempty"`
	PackageKeys  *[]PackageKey  `json:"packageKeys,omitempty"`
	Roles        *[]Role        `json:"roles,omitempty"`
}

func (m *Member) Identifier() MemberIdentifier {
	return MemberIdentifier{
		MemberId: m.Id,
		Username: m.Username,
	}
}

// -------------------------------------------------------------------
// Synthetic path identifier
//
// Service-related

type ServiceEndpoint struct {
	ServiceId  string
	EndpointId string
}
type ServiceMethod struct {
	ServiceEndpoint
	MethodId string
}
type ServiceMethodFilter struct {
	ServiceMethod
	FilterId string
}

// ----------------------------------------
//

// Package-plan related

// EmailTemplateSet email template set
type EmailTemplateSet struct {
	AddressableV3Object
	Type           string           `json:"type,omitempty"`
	EmailTemplates *[]EmailTemplate `json:"emailTemplates,omitempty"`
}

// EmailTemplate email template
type EmailTemplate struct {
	AddressableV3Object
	Type    string `json:"type"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// -----------------------------------------------------------------------------
// Methods

type BaseMethod struct {
	AddressableV3Object
	SampleJsonResponse string `json:"sampleJsonResponse,omitempty"`
	SampleXmlResponse  string `json:"sampleXmlResponse,omitempty"`
}

type ServiceEndpointMethod struct {
	BaseMethod
	ParentEndpointId ServiceEndpointIdentifier
}

func (m *ServiceEndpointMethod) Identifier() ServiceEndpointMethodIdentifier {
	return ServiceEndpointMethodIdentifier{
		MethodId:                  m.Id,
		ServiceEndpointIdentifier: m.ParentEndpointId,
	}
}

type PackagePlanServiceEndpointMethod struct {
	BaseMethod

	PackagePlanServiceEndpoint PackagePlanServiceEndpointIdentifier
}

func (ppsem *PackagePlanServiceEndpointMethod) Identifier() PackagePlanServiceEndpointMethodIdentifier {
	return PackagePlanServiceEndpointMethodIdentifier{
		PackagePlanIdentifier: ppsem.PackagePlanServiceEndpoint.PackagePlanIdentifier,
		ServiceEndpointMethodIdentifier: ServiceEndpointMethodIdentifier{
			MethodId: ppsem.Id,
			ServiceEndpointIdentifier: ServiceEndpointIdentifier{
				EndpointId: ppsem.PackagePlanServiceEndpoint.EndpointId,
				ServiceIdentifier: ServiceIdentifier{
					ServiceId: ppsem.PackagePlanServiceEndpoint.ServiceId,
				},
			},
		},
	}
}

func ParseServiceEndpointMethod(inp []byte) (interface{}, int, error) {
	var rv ServiceEndpointMethod
	err := json.Unmarshal(inp, &rv)
	return rv, 1, err
}

func ParsePacakgePlanServiceEndpointMethod(inp []byte) (interface{}, int, error) {
	var rv PackagePlanServiceEndpointMethod
	err := json.Unmarshal(inp, &rv)
	return rv, 1, err
}

func ParseServiceEndpointMethodArray(inp []byte) (interface{}, int, error) {
	var rv []ServiceEndpointMethod
	err := json.Unmarshal(inp, &rv)
	return rv, len(rv), err
}

type ResponseFilter struct {
	AddressableV3Object
	Notes            string `json:"notes"`
	XmlFilterFields  string `json:"xmlFilterFields,omitempty"`
	JsonFilterFields string `json:"jsonFilterFields,omitempty"`
}

type ServiceEndpointMethodFilter struct {
	ResponseFilter

	ServiceEndpointMethod ServiceEndpointMethodIdentifier
}

type PackagePlanServiceEndpointMethodFilter struct {
	ResponseFilter

	PackagePlanServiceEndpointMethod PackagePlanServiceEndpointMethodIdentifier
}

func (ppsemf *PackagePlanServiceEndpointMethodFilter) Identifier() PackagePlanServiceEndpointMethodFilterIdentifier {
	return PackagePlanServiceEndpointMethodFilterIdentifier{
		PackagePlanServiceIdentifier: PackagePlanServiceIdentifier{
			PackagePlanIdentifier: ppsemf.PackagePlanServiceEndpointMethod.PackagePlanIdentifier,
			ServiceIdentifier: ServiceIdentifier{
				ServiceId: ppsemf.PackagePlanServiceEndpointMethod.ServiceId,
			},
		},
		ServiceEndpointMethodFilterIdentifier: ServiceEndpointMethodFilterIdentifier{
			FilterId: ppsemf.Id,
			ServiceEndpointMethodIdentifier: ServiceEndpointMethodIdentifier{
				MethodId: ppsemf.PackagePlanServiceEndpointMethod.MethodId,
				ServiceEndpointIdentifier: ServiceEndpointIdentifier{
					EndpointId: ppsemf.PackagePlanServiceEndpointMethod.EndpointId,
					ServiceIdentifier: ServiceIdentifier{
						ServiceId: ppsemf.PackagePlanServiceEndpointMethod.ServiceId,
					},
				},
			},
		},
	}
}

func (emrf *ServiceEndpointMethodFilter) Identifier() ServiceEndpointMethodFilterIdentifier {
	return ServiceEndpointMethodFilterIdentifier{
		FilterId:                        emrf.Id,
		ServiceEndpointMethodIdentifier: emrf.ServiceEndpointMethod,
	}
}

func ParsePackagePlanServiceEndpointMethodFilterArray(inp []byte) (interface{}, int, error) {
	var rv []PackagePlanServiceEndpointMethodFilter
	err := json.Unmarshal(inp, &rv)
	return rv, len(rv), err
}

func ParseServiceEndpointMethodFilterArray(inp []byte) (interface{}, int, error) {
	var rv []ServiceEndpointMethodFilter
	err := json.Unmarshal(inp, &rv)
	return rv, len(rv), err
}

func ParsePackagePlanServiceEndpointMethodFilter(inp []byte) (interface{}, int, error) {
	var rv PackagePlanServiceEndpointMethodFilter
	err := json.Unmarshal(inp, &rv)
	return rv, 1, err
}

func ParseServiceEndpointMethodFilter(inp []byte) (interface{}, int, error) {
	var rv ServiceEndpointMethodFilter
	err := json.Unmarshal(inp, &rv)
	return rv, 1, err
}

// -----------------------------------------------------------------------------
// Type conversion functions

func (m *Service) AddressableEndpoints() []AddressableV3Object {
	return AddressableEndpoints(m.Endpoints)
}

func AddressableEndpoints(inp []Endpoint) []AddressableV3Object {
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

func (e *V3GenericErrorResponse) HasData() bool {
	return len(e.ErrorCode) > 0 || len(e.ErrorMessage) > 0
}

func (e *V3GenericErrorResponse) Error() string {
	return fmt.Sprintf("%s (mashery V3 code %s)", e.ErrorMessage, e.ErrorCode)
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

// TimedAccessTokenResponse Timed access token response, suitable for storing in a log file.
type TimedAccessTokenResponse struct {
	Obtained   time.Time `json:"obtained"`
	ServerTime time.Time `json:"-"`
	QPS        int       `json:"qps"`
	AccessTokenResponse
}

func (a AccessTokenResponse) ObtainedNow() *TimedAccessTokenResponse {
	rv := TimedAccessTokenResponse{
		AccessTokenResponse: a,
		Obtained:            time.Now(),
	}
	return &rv
}

func (t *TimedAccessTokenResponse) Expired() bool {
	now := time.Now()
	secondsDiff := now.Unix() - t.Obtained.Unix()

	if secondsDiff > int64(math.Round(float64(t.ExpiresIn)*0.95)) {
		return true
	}

	return false
}

func (t *TimedAccessTokenResponse) ExpiryTime() time.Time {
	return time.Unix(t.Obtained.Unix()+int64(t.ExpiresIn), 0)
}

// TimeLeft Returns number of seconds that are still left in this access tokens.
func (t *TimedAccessTokenResponse) TimeLeft() int {
	diff := t.Obtained.Unix() + int64(t.ExpiresIn) - time.Now().Unix()
	if diff > 0 {
		return int(diff)
	} else {
		return 0
	}
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
	var rv Endpoint
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryEndpointArray(dat []byte) (interface{}, int, error) {
	var rv []Endpoint
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func NilParser(dat []byte) (interface{}, int, error) {
	return nil, 0, nil
}

func ParseService(dat []byte) (interface{}, int, error) {
	var rv Service
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseServiceCache(dat []byte) (interface{}, int, error) {
	var rv ServiceCache
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryServiceSecurityProfileOAuth(dat []byte) (interface{}, int, error) {
	var rv MasheryOAuth
	err := json.Unmarshal(dat, &rv)

	return rv, 1, err
}

func ParseMasheryServiceArray(dat []byte) (interface{}, int, error) {
	var rv []Service
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseServiceErrorSetArray(dat []byte) (interface{}, int, error) {
	var rv []ErrorSet
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseErrorSet(dat []byte) (interface{}, int, error) {
	var rv ErrorSet
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseErrorSetMessage(dat []byte) (interface{}, int, error) {
	var rv MasheryErrorMessage
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryPackage(dat []byte) (interface{}, int, error) {
	var rv Package
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryPackageArray(dat []byte) (interface{}, int, error) {
	var rv []Package
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryPlan(dat []byte) (interface{}, int, error) {
	var rv Plan
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryPlanArray(dat []byte) (interface{}, int, error) {
	var rv []Plan
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryPackageKey(dat []byte) (interface{}, int, error) {
	var rv PackageKey
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryPackageKeyArray(dat []byte) (interface{}, int, error) {
	var rv []PackageKey
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryApplication(dat []byte) (interface{}, int, error) {
	var rv Application
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryApplicationArray(dat []byte) (interface{}, int, error) {
	var rv []Application
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryMember(dat []byte) (interface{}, int, error) {
	var rv Member
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryMemberArray(dat []byte) (interface{}, int, error) {
	var rv []Member
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryRoleArray(dat []byte) (interface{}, int, error) {
	var rv []Role
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryRole(dat []byte) (interface{}, int, error) {
	var rv Role
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseRolePermissionArray(dat []byte) (interface{}, int, error) {
	var rv []RolePermission
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryRolePermission(dat []byte) (interface{}, int, error) {
	var rv RolePermission
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryEmailTemplateSetArray(dat []byte) (interface{}, int, error) {
	var rv []EmailTemplateSet
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}

func ParseMasheryEmailTemplateSet(dat []byte) (interface{}, int, error) {
	var rv EmailTemplateSet
	err := json.Unmarshal(dat, &rv)
	return rv, 1, err
}

func ParseMasheryDomainAddressArray(dat []byte) (interface{}, int, error) {
	var rv []DomainAddress
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
}
