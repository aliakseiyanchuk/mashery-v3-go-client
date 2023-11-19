package main

import (
	"context"
	_ "embed"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func execServiceEndpointList(ctx context.Context, cl v3client.Client, id masherytypes.ServiceIdentifier) ([]masherytypes.AddressableV3Object, error) {
	return cl.ListEndpoints(ctx, id)
}

//go:embed templates/service_endpoint_list.tmpl
var serviceEndpointListTemplate string
var subCmdServiceEndpointList *SubcommandTemplate[masherytypes.ServiceIdentifier, []masherytypes.AddressableV3Object]

func init() {
	subCmdServiceEndpointList = &SubcommandTemplate[masherytypes.ServiceIdentifier, []masherytypes.AddressableV3Object]{
		Command:        []string{"service", "endpoint", "list"},
		FlagSetInit:    initServiceShowFlagSet,
		EnvFlagSetInit: initServiceShowEnvFlagSet,
		Validator:      validateServiceShowArg,
		Executor:       execServiceEndpointList,
		Template:       mustTemplate(serviceEndpointListTemplate),
	}

	enableSubcommand(subCmdServiceEndpointList.Finder())
}
