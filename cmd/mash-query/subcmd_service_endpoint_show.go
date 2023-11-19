package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func validateServiceEndpointShowArg(arg *masherytypes.ServiceEndpointIdentifier) error {
	if len(arg.ServiceId) == 0 {
		return errors.New("service identifier required")
	} else if len(arg.EndpointId) == 0 {
		return errors.New("service endpoint identifier required")
	}

	return nil
}

func execServiceEndpointShow(ctx context.Context, cl v3client.Client, id masherytypes.ServiceEndpointIdentifier) (ObjectWithExists[masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint], error) {
	rv, serviceExists, err := cl.GetEndpoint(ctx, id)

	return ObjectWithExists[masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]{
		Identifier: id,
		Object:     rv,
		Exists:     serviceExists,
	}, err
}

//go:embed templates/service_endpoint_show.tmpl
var serviceEndpointShowTemplate string
var subCmdServiceEndpointShow *SubcommandTemplate[masherytypes.ServiceEndpointIdentifier, ObjectWithExists[masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]]

func initServiceEndpointShowFlagSet(arg *masherytypes.ServiceEndpointIdentifier, fs *flag.FlagSet) {
	initServiceShowFlagSet(&arg.ServiceIdentifier, fs)
	fs.StringVar(&arg.ServiceId, "endpoint-id", "", "Service endpoint identifier")
}

func initServiceEndpointShowEnvFlagSet(arg *masherytypes.ServiceEndpointIdentifier) []EnvFlag {
	rv := initServiceShowEnvFlagSet(&arg.ServiceIdentifier)
	rv = append(rv, EnvFlag{
		Dest:   &arg.EndpointId,
		EnvVar: "MASH_SERVICE_ENDPOINT_ID",
		Option: "endpoint-id",
	})

	return rv
}

func init() {
	subCmdServiceEndpointShow = &SubcommandTemplate[masherytypes.ServiceEndpointIdentifier, ObjectWithExists[masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]]{
		Command:        []string{"service", "endpoint", "show"},
		FlagSetInit:    initServiceEndpointShowFlagSet,
		EnvFlagSetInit: initServiceEndpointShowEnvFlagSet,
		Validator:      validateServiceEndpointShowArg,
		Executor:       execServiceEndpointShow,
		Template:       mustTemplate(serviceEndpointShowTemplate),
	}

	enableSubcommand(subCmdServiceEndpointShow.Finder())
}
