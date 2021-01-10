package mashery_v3_go_client

import (
	"encoding/json"
	"errors"
	"fmt"
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

type Cache struct {
	ClientSurrogateControlEnabled bool     `json:"clientSurrogateControlEnabled"`
	ContentCacheKeyHeaders        []string `json:"contentCacheKeyHeaders"`
}

type Cors struct {
	AllDomainsEnabled bool `json:"boolean"`
	MaxAge            int  `json:"maxAge"`
}

type ScheduledMaintenanceEvent struct {
	Id            string                `json:"id"`
	Name          string                `json:"name"`
	StartDateTime MasheryJSONTime       `json:"startDateTime"`
	EndDateTime   MasheryJSONTime       `json:"endDateTime"`
	Endpoints     []AddressableV3Object `json:"endpoints"`
}

type Method struct {
	Id                 string          `json:"id"`
	Name               string          `json:"name"`
	Created            MasheryJSONTime `json:"created"`
	Updated            MasheryJSONTime `json:"updated"`
	SampleJsonResponse string          `json:"sampleJsonResponse"`
	SampleXmlResponse  string          `json:"sampleXmlResponse"`
}

type ResponseFilter struct {
	Id               string          `json:"id"`
	Name             string          `json:"name"`
	Created          MasheryJSONTime `json:"created"`
	Updated          MasheryJSONTime `json:"updated"`
	Notes            string          `json:"notes"`
	XmlFilterFields  string          `json:"xmlFilterFields"`
	JsonFilterFields string          `json:"jsonFilterFields"`
}

type SystemDomainAuthentication struct {
	Type        string `json:"type"`
	Username    string `json:"username"`
	Certificate string `json:"certificate"`
	Password    string `json:"password"`
}

type Processor struct {
	PreProcessEnabled  bool                `json:"preProcessEnabled"`
	PostProcessEnabled bool                `json:"postProcessEnabled"`
	PostInputs         map[string](string) `json:"postInputs"`
	PreInputs          map[string](string) `json:"preInputs"`
	Adapter            string              `json:"adapter"`
}

type Domain struct {
	Address string `json:"address"`
}

type MasheryEndpoint struct {
	AddressableV3Object

	AllowMissingApiKey                         bool                       `json:"allowMissingApiKey"`
	ApiKeyValueLocationKey                     string                     `json:"apiKeyValueLocationKey"`
	ApiKeyValueLocations                       []string                   `json:"apiKeyValueLocations"`
	ApiMethodDetectionKey                      string                     `json:"apiMethodDetectionKey"`
	ApiMethodDetectionLocations                []string                   `json:"apiMethodDetectionLocations"`
	Cache                                      *Cache                     `json:"cache,omitempty"`
	ConnectionTimeoutForSystemDomainRequest    int                        `json:"connectionTimeoutForSystemDomainRequest"`
	ConnectionTimeoutForSystemDomainResponse   int                        `json:"connectionTimeoutForSystemDomainResponse"`
	CookiesDuringHttpRedirectsEnabled          bool                       `json:"cookiesDuringHttpRedirectsEnabled"`
	Cors                                       *Cors                      `json:"cors,omitempty"`
	CustomRequestAuthenticationAdapter         string                     `json:"customRequestAuthenticationAdapter"`
	DropApiKeyFromIncomingCall                 bool                       `json:"dropApiKeyFromIncomingCall"`
	ForceGzipOfBackendCall                     bool                       `json:"forceGzipOfBackendCall"`
	GzipPassthroughSupportEnabled              bool                       `json:"gzipPassthroughSupportEnabled"`
	HeadersToExcludeFromIncomingCall           []string                   `json:"headersToExcludeFromIncomingCall"`
	HighSecurity                               bool                       `json:"highSecurity"`
	HostPassthroughIncludedInBackendCallHeader bool                       `json:"hostPassthroughIncludedInBackendCallHeader"`
	InboundSslRequired                         bool                       `json:"inboundSslRequired"`
	JsonpCallbackParameter                     string                     `json:"jsonpCallbackParameter"`
	JsonpCallbackParameterValue                string                     `json:"jsonpCallbackParameterValue"`
	ScheduledMaintenanceEvent                  ScheduledMaintenanceEvent  `json:"scheduledMaintenanceEvent"`
	ForwardedHeaders                           []string                   `json:"forwardedHeaders"`
	ReturnedHeaders                            []string                   `json:"returnedHeaders"`
	Methods                                    []Method                   `json:"methods"`
	NumberOfHttpRedirectsToFollow              int                        `json:"numberOfHttpRedirectsToFollow"`
	OutboundRequestTargetPath                  string                     `json:"outboundRequestTargetPath"`
	OutboundRequestTargetQueryParameters       string                     `json:"outboundRequestTargetQueryParameters"`
	OutboundTransportProtocol                  string                     `json:"outboundTransportProtocol"`
	Processor                                  Processor                  `json:"processor"`
	PublicDomains                              []Domain                   `json:"publicDomains"`
	RequestAuthenticationType                  string                     `json:"requestAuthenticationType"`
	RequestPathAlias                           string                     `json:"requestPathAlias"`
	RequestProtocol                            string                     `json:"requestProtocol"`
	OauthGrantTypes                            []string                   `json:"oauthGrantTypes"`
	StringsToTrimFromApiKey                    string                     `json:"stringsToTrimFromApiKey"`
	SupportedHttpMethods                       []string                   `json:"supportedHttpMethods"`
	SystemDomainAuthentication                 SystemDomainAuthentication `json:"systemDomainAuthentication"`
	SystemDomains                              []Domain                   `json:"systemDomains"`
	TrafficManagerDomain                       string                     `json:"trafficManagerDomain"`
	UseSystemDomainCredentials                 bool                       `json:"useSystemDomainCredentials"`
	SystemDomainCredentialKey                  string                     `json:"systemDomainCredentialKey"`
	SystemDomainCredentialSecret               string                     `json:"systemDomainCredentialSecret"`
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
	MACAlgorithm                string   `json:"macAlgorithm"`
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
	Status       int    `json:"status"`
	DetailHeader string `json:"detailHeader"`
	ResponseBody string `json:"responseBody"`
}

type MasheryErrorSet struct {
	Name          string                `json:"name"`
	Type          string                `json:"type"`
	JSONP         bool                  `json:"jsonp"`
	JSONPType     string                `json:"jsonpType"`
	ErrorMessages []MasheryErrorMessage `json:"errorMessages"`
}

type MasheryJSONTime time.Time

type AddressableV3Object struct {
	Id      string          `json:"id"`
	Name    string          `json:"name"`
	Created MasheryJSONTime `json:"created"`
	Updated MasheryJSONTime `json:"updated"`
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

	return errors.New(fmt.Sprint("unknown Mashery JSON date: '%s'", s))
}

func (t *MasheryJSONTime) MarshalJSON() ([]byte, error) {
	return time.Time(*t).MarshalJSON()
}

func (t *MasheryJSONTime) ToString() string {
	return time.Time(*t).Format("January 02, 2006 15:04")
}

type MasheryService struct {
	AddressableV3Object
	Endpoints         []MasheryEndpoint      `json:"endpoints"`
	EditorHandle      string                 `json:"editorHandle"`
	RevisionNumber    int                    `json:"revisionNumber"`
	RobotsPolicy      string                 `json:"robotsPolicy"`
	CrossdomainPolicy string                 `json:"crossdomainPolicy"`
	Description       string                 `json:"description"`
	ErrorSets         []MasheryErrorSet      `json:"errorSets"`
	QpsLimitOverall   int64                  `json:"qpsLimitOverall"`
	RFC3986Encode     bool                   `json:"rfc3986Encode"`
	SecurityProfile   MasherySecurityProfile `json:"securityProfile"`
	Version           string                 `json:"version"`
}

// -----------------------------------------------------------------------------
// Mashery package-related types

type EAV map[string]string

type MasheryPlan struct {
	AddressableV3Object
	Description                       string           `json:"description"`
	Eav                               EAV              `json:"eav"`
	SelfServiceKeyProvisioningEnabled bool             `json:"selfServiceKeyProvisioningEnabled"`
	AdminKeyProvisioningEnabled       bool             `json:"adminKeyProvisioningEnabled"`
	Notes                             string           `json:"notes"`
	MaxNumKeysAllowed                 int              `json:"maxNumKeysAllowed"`
	NumKeysBeforeReview               int              `json:"numKeysBeforeReview"`
	QpsLimitCeiling                   int64            `json:"qpsLimitCeiling"`
	QpsLimitExempt                    bool             `json:"qpsLimitExempt"`
	QpsLimitKeyOverrideAllowed        bool             `json:"qpsLimitKeyOverrideAllowed"`
	RateLimitCeiling                  int64            `json:"rateLimitCeiling"`
	RateLimitExempt                   bool             `json:"rateLimitExempt"`
	RateLimitKeyOverrideAllowed       bool             `json:"rateLimitKeyOverrideAllowed"`
	RateLimitPeriod                   string           `json:"rateLimitPeriod"`
	ResponseFilterOverrideAllowed     bool             `json:"responseFilterOverrideAllowed"`
	Status                            string           `json:"status"`
	EmailTemplateSetId                int64            `json:"emailTemplateSetId"`
	Services                          []MasheryService `json:"services"`

	// Identity of the contxt object
	ParentPackageId string
}

type PlanFilter struct {
	AddressableV3Object
	Notes            string `json:"notes"`
	XmlFilterFields  string `json:"xmlFilterFields"`
	JsonFilterFields string `json:"jsonFilterFields"`
}

type MasheryPackage struct {
	Id                          string          `json:"id"`
	Created                     MasheryJSONTime `json:"created"`
	Updated                     MasheryJSONTime `json:"updated"`
	Name                        string          `json:"name"`
	Description                 string          `json:"description"`
	NotifyDeveloperPeriod       string          `json:"notifyDeveloperPeriod"`
	NotifyDeveloperNearQuota    bool            `json:"notifyDeveloperNearQuota"`
	NotifyDeveloperOverQuota    bool            `json:"notifyDeveloperOverQuota"`
	NotifyDeveloperOverThrottle bool            `json:"notifyDeveloperOverThrottle"`
	NotifyAdminPeriod           string          `json:"notifyAdminPeriod"`
	NotifyAdminNearQuota        bool            `json:"notifyAdminNearQuota"`
	NotifyAdminOverQuota        bool            `json:"notifyAdminOverQuota"`
	NotifyAdminOverThrottle     bool            `json:"notifyAdminOverThrottle"`
	NotifyAdminEmails           string          `json:"notifyAdminEmails"`
	NearQuotaThreshold          int             `json:"nearQuotaThreshold"`
	Eav                         EAV             `json:"eav"`
	KeyAdapter                  string          `json:"keyAdapter"`
	KeyLength                   int             `json:"keyLength"`
	SharedSecretLength          int             `json:"sharedSecretLength"`
	Plans                       []MasheryPlan   `json:"plans"`
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
	Apikey           string          `json:"apikey"`
	Secret           string          `json:"secret"`
	RateLimitCeiling int64           `json:"rateLimitCeiling"`
	RateLimitExempt  bool            `json:"rateLimitExempt"`
	QpsLimitCeiling  int64           `json:"qpsLimitCeiling"`
	QpsLimitExempt   bool            `json:"qpsLimitExempt"`
	Status           string          `json:"status"`
	Limits           []Limit         `json:"limits"`
	Package          *MasheryPackage `json:"package"`
	Plan             *MasheryPlan    `json:"plan"`
}

type MasheryApplication struct {
	AddressableV3Object

	Username          string              `json:"username"`
	Description       string              `json:"description"`
	Type              string              `json:"type"`
	Commercial        bool                `json:"commercial"`
	Ads               bool                `json:"ads"`
	AdsSystem         string              `json:"adsSystem"`
	UsageModel        string              `json:"usageModel"`
	Tags              string              `json:"tags"`
	Notes             string              `json:"notes"`
	HowDidYouHear     string              `json:"howDidYouHear"`
	PreferredProtocol string              `json:"preferredProtocol"`
	PreferredOutput   string              `json:"preferredOutput"`
	ExternalId        string              `json:"externalId"`
	Uri               string              `json:"uri"`
	PauthRedirectUri  string              `json:"oauthRedirectUri"`
	PackageKeys       []MasheryPackageKey `json:"packageKeys"`
}

type MasheryRole struct {
	AddressableV3Object
	// TODO: Map Group object.
}

type MasheryMember struct {
	AddressableV3Object

	Username     string                `json:"username"`
	Email        string                `json:"email"`
	DisplayName  string                `json:"displayName"`
	Uri          string                `json:"uri"`
	Blog         string                `json:"blog"`
	Im           string                `json:"im"`
	Imsvc        string                `json:"imsvc"`
	Phone        string                `json:"phone"`
	Company      string                `json:"company"`
	Address1     string                `json:"address1"`
	Address2     string                `json:"address2"`
	Locality     string                `json:"locality"`
	Region       string                `json:"region"`
	PostalCode   string                `json:"postalCode"`
	CountryCode  string                `json:"countryCode"`
	FirstName    string                `json:"firstName"`
	LastName     string                `json:"lastName"`
	AreaStatus   string                `json:"areaStatus"`
	ExternalId   string                `json:"externalId"`
	PasswdNew    string                `json:"passwdNew"`
	Applications *[]MasheryApplication `json:"applications,omitempty"`
	PackageKeys  *[]MasheryPackageKey  `json:"packageKeys"`
	Roles        *[]MasheryRole        `json:"roles"`
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

// --------------------------------------------
// Factory routines

func MasheryServiceFactory() interface{} {
	return MasheryService{}
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

func ParseMasheryServiceArray(dat []byte) (interface{}, int, error) {
	var rv []MasheryService
	err := json.Unmarshal(dat, &rv)
	return rv, len(rv), err
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

func MasheryEndpointFactory() MasheryEndpoint {
	return MasheryEndpoint{}
}
