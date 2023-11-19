package main

import (
	"context"
	_ "embed"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

//go:embed templates/service_list.tmpl
var serviceListTemplate string
var subCmdServiceList *SubcommandTemplate[int, []masherytypes.Service]

func execServiceList(ctx context.Context, cl v3client.Client, _ int) ([]masherytypes.Service, error) {
	return cl.ListServices(ctx)
}

func init() {
	subCmdServiceList = &SubcommandTemplate[int, []masherytypes.Service]{
		Command:  []string{"service", "list"},
		Executor: execServiceList,
		Template: mustTemplate(serviceListTemplate),
	}

	enableSubcommand(subCmdServiceList.Finder())
}
