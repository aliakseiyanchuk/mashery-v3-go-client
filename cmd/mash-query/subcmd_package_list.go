package main

import (
	"context"
	_ "embed"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

//go:embed templates/package_list.tmpl
var packageListTemplate string
var subCmdPackageList *SubcommandTemplate[int, []masherytypes.Package]

func execPackageList(ctx context.Context, cl v3client.Client, _ int) ([]masherytypes.Package, error) {
	return cl.ListPackages(ctx)
}

func init() {
	subCmdPackageList = &SubcommandTemplate[int, []masherytypes.Package]{
		Command:  []string{"package", "list"},
		Executor: execPackageList,
		Template: mustTemplate(packageListTemplate),
	}

	enableSubcommand(subCmdPackageList.Finder())
}
