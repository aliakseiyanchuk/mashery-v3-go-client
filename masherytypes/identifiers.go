package masherytypes

// MemberIdentifier Member identifier
type MemberIdentifier struct {
	MemberId string `json:"mid"`
	Username string `json:"u"`
}

type ApplicationIdentifier struct {
	ApplicationId string `json:"aid"`
}

type PackageIdentifier struct {
	PackageId string `json:"pid"`
}

type PackagePlanIdentifier struct {
	PackageIdentifier
	PlanId string `json:"plid"`
}

type PackagePlanServiceIdentifier struct {
	PackagePlanIdentifier
	ServiceIdentifier
}

type PackagePlanServiceEndpointIdentifier struct {
	PackagePlanIdentifier
	ServiceEndpointIdentifier
}

type PackagePlanServiceEndpointMethodIdentifier struct {
	PackagePlanIdentifier
	ServiceEndpointMethodIdentifier
}

type PackagePlanServiceEndpointMethodFilterIdentifier struct {
	PackagePlanServiceIdentifier
	ServiceEndpointMethodFilterIdentifier
}

type ServiceIdentifier struct {
	ServiceId string `json:"sid"`
}

type ServiceEndpointIdentifier struct {
	ServiceIdentifier
	EndpointId string `json:"eid"`
}

type ServiceEndpointMethodIdentifier struct {
	ServiceEndpointIdentifier
	MethodId string `json:"mthid"`
}

type ServiceEndpointMethodFilterIdentifier struct {
	ServiceEndpointMethodIdentifier
	FilterId string `json:"fid"`
}

type PackageKeyIdentifier struct {
	PackageKeyId string `json:"pkid"`
}

type ErrorSetIdentifier struct {
	ServiceIdentifier
	ErrorSetId string `json:"esid"`
}
