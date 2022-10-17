package masherytypes

func ServiceIdentityFrom(id string) ServiceIdentifier {
	return ServiceIdentifier{ServiceId: id}
}

func ServiceEndpointIdentityFrom(svcId, endpointId string) ServiceEndpointIdentifier {
	return ServiceEndpointIdentifier{
		ServiceIdentifier: ServiceIdentifier{ServiceId: svcId},
		EndpointId:        endpointId,
	}

}

func ServiceEndpointMethodIdentityFrom(svcId, endpId, methId string) ServiceEndpointMethodIdentifier {
	return ServiceEndpointMethodIdentifier{
		ServiceEndpointIdentifier: ServiceEndpointIdentifier{
			ServiceIdentifier: ServiceIdentifier{ServiceId: svcId},
			EndpointId:        endpId,
		},
		MethodId: methId,
	}
}

func ServiceEndpointMethodFilterIdentityFrom(svcId, endpId, methId, filterId string) ServiceEndpointMethodFilterIdentifier {
	return ServiceEndpointMethodFilterIdentifier{
		ServiceEndpointMethodIdentifier: ServiceEndpointMethodIdentifier{
			ServiceEndpointIdentifier: ServiceEndpointIdentifier{
				ServiceIdentifier: ServiceIdentifier{ServiceId: svcId},
				EndpointId:        endpId,
			},
			MethodId: methId,
		},
		FilterId: filterId,
	}
}

func PackageIdentityFrom(id string) PackageIdentifier {
	return PackageIdentifier{PackageId: id}
}

func PackagePlanIdentityFrom(packageId, planId string) PackagePlanIdentifier {
	return PackagePlanIdentifier{
		PackageIdentifier: PackageIdentifier{PackageId: packageId},
		PlanId:            planId,
	}
}
