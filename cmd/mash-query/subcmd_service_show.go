package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func validateServiceShowArg(arg *masherytypes.ServiceIdentifier) error {
	if len(arg.ServiceId) == 0 {
		return errors.New("service identifier required")
	}

	return nil
}

func execServiceShow(ctx context.Context, cl v3client.Client, id masherytypes.ServiceIdentifier) (ObjectWithExists[masherytypes.ServiceIdentifier, masherytypes.Service], error) {
	rv, serviceExists, err := cl.GetService(ctx, id)

	// If the service exists, retrieve the security profile for this service
	if err == nil && serviceExists {
		if secProfile, secProfileExists, secProfileErr := cl.GetServiceOAuthSecurityProfile(ctx, id); secProfileErr != nil {
			return ObjectWithExists[masherytypes.ServiceIdentifier, masherytypes.Service]{}, secProfileErr
		} else if secProfileExists {
			rv.SecurityProfile = &masherytypes.MasherySecurityProfile{
				OAuth: &secProfile,
			}
		}
	}
	return ObjectWithExists[masherytypes.ServiceIdentifier, masherytypes.Service]{
		Identifier: id,
		Object:     rv,
		Exists:     serviceExists,
	}, err
}

//go:embed templates/service_show.tmpl
var serviceShowTemplate string
var subCmdShowService *SubcommandTemplate[masherytypes.ServiceIdentifier, ObjectWithExists[masherytypes.ServiceIdentifier, masherytypes.Service]]

func initServiceShowFlagSet(arg *masherytypes.ServiceIdentifier, fs *flag.FlagSet) {
	fs.StringVar(&arg.ServiceId, "service-id", "", "Service identifier")
}

func initServiceShowEnvFlagSet(arg *masherytypes.ServiceIdentifier) []EnvFlag {
	return []EnvFlag{
		{
			Dest:   &arg.ServiceId,
			EnvVar: "MASH_SERVICE_ID",
			Option: "service-id",
		},
	}
}

func init() {
	subCmdShowService = &SubcommandTemplate[masherytypes.ServiceIdentifier, ObjectWithExists[masherytypes.ServiceIdentifier, masherytypes.Service]]{
		Command:        []string{"service", "show"},
		FlagSetInit:    initServiceShowFlagSet,
		EnvFlagSetInit: initServiceShowEnvFlagSet,
		Validator:      validateServiceShowArg,
		Executor:       execServiceShow,
		Template:       mustTemplate(serviceShowTemplate),
	}

	enableSubcommand(subCmdShowService.Finder())
}
