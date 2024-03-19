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

func (ppsemi PackagePlanServiceEndpointMethodIdentifier) GetPackagePlanServiceEndpointIdentifier() PackagePlanServiceEndpointIdentifier {
	return PackagePlanServiceEndpointIdentifier{
		PackagePlanIdentifier: PackagePlanIdentifier{
			PackageIdentifier: PackageIdentifier{PackageId: ppsemi.PackageId},
			PlanId:            ppsemi.PlanId,
		},
		ServiceEndpointIdentifier: ServiceEndpointIdentifier{
			ServiceIdentifier: ServiceIdentifier{ServiceId: ppsemi.ServiceId},
			EndpointId:        ppsemi.EndpointId,
		},
	}
}

type PackagePlanServiceEndpointMethodIdentifier struct {
	PackagePlanIdentifier
	ServiceEndpointMethodIdentifier
}

type PackagePlanServiceEndpointMethodFilterIdentifier struct {
	PackagePlanIdentifier
	ServiceEndpointMethodFilterIdentifier
}

func (ppsemfi PackagePlanServiceEndpointMethodFilterIdentifier) AsPackagePlanServiceEndpointMethodIdentifier() PackagePlanServiceEndpointMethodIdentifier {
	return PackagePlanServiceEndpointMethodIdentifier{
		PackagePlanIdentifier: PackagePlanIdentifier{
			PackageIdentifier: PackageIdentifier{PackageId: ppsemfi.PackageId},
			PlanId:            ppsemfi.PlanId,
		},
		ServiceEndpointMethodIdentifier: ServiceEndpointMethodIdentifier{
			ServiceEndpointIdentifier: ServiceEndpointIdentifier{
				ServiceIdentifier: ServiceIdentifier{ServiceId: ppsemfi.ServiceId},
				EndpointId:        ppsemfi.EndpointId,
			},
			MethodId: ppsemfi.MethodId,
		},
	}
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

type ErrorSetMessageIdentifier struct {
	ErrorSetIdentifier
	ErrorSetMessageId string `json:"esmid"`
}
