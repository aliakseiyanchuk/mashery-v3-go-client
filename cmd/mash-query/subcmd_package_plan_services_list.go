package main

import (
	"context"
	_ "embed"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func asyncFetchPlanServiceEndpoint(ctx context.Context, cl v3client.Client, id masherytypes.PackagePlanIdentifier, service *masherytypes.Service, ch chan error) {
	ident := masherytypes.PackagePlanServiceIdentifier{
		PackagePlanIdentifier: id,
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: service.Id,
		},
	}

	addr, err := cl.ListPlanEndpoints(ctx, ident)
	if err == nil {
		rv := make([]masherytypes.Endpoint, len(addr))
		for i, a := range addr {
			rv[i] = masherytypes.Endpoint{
				AddressableV3Object: a,
				ParentServiceId:     ident.ServiceIdentifier,
			}
		}
		service.Endpoints = rv
	}

	ch <- err
}

func execPackagePlanListShow(ctx context.Context, cl v3client.Client, id masherytypes.PackagePlanIdentifier) ([]masherytypes.Service, error) {
	services, err := cl.ListPlanServices(ctx, id)

	// Async retrieve service endpoints
	if err == nil {
		errChan := make(chan error)
		defer close(errChan)

		for i := range services {
			go asyncFetchPlanServiceEndpoint(ctx, cl, id, &services[i], errChan)
		}

		for _ = range services {
			if fetchErr := <-errChan; fetchErr != nil {
				err = fetchErr
			}
		}
	}
	return services, err
}

//go:embed templates/package_plan_services_list.tmpl
var packagePlanServicesListTemplate string
var subCmdPackagePlanServiceList *SubcommandTemplate[masherytypes.PackagePlanIdentifier, []masherytypes.Service]

func init() {
	subCmdPackagePlanServiceList = &SubcommandTemplate[masherytypes.PackagePlanIdentifier, []masherytypes.Service]{
		Command:        []string{"package", "plan", "services", "list"},
		FlagSetInit:    initPackagePlanShowFlagSet,
		EnvFlagSetInit: initPackagePlanShowEnvFlagSet,
		Validator:      validatePackagePLanShowArg,
		Executor:       execPackagePlanListShow,
		Template:       mustTemplate(packagePlanServicesListTemplate),
	}

	enableSubcommand(subCmdPackagePlanServiceList.Finder())
}
