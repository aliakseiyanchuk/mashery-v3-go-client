package main

import (
	"context"
	_ "embed"
	"flag"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

//go:embed templates/application_list.tmpl
var applicationListTemplate string
var subCmdApplicationList *SubcommandTemplate[ApplicationFilter, []masherytypes.Application]

type ApplicationFilter struct {
	Name     string
	Username string
}

func (af ApplicationFilter) Params() map[string]string {
	rv := map[string]string{}
	if len(af.Name) > 0 {
		rv["name"] = af.Name
	}
	if len(af.Username) > 0 {
		rv["username"] = af.Username
	}

	return rv
}

func execApplicationList(ctx context.Context, cl v3client.Client, cfg ApplicationFilter) ([]masherytypes.Application, error) {
	filterParams := cfg.Params()
	if len(filterParams) > 0 {
		return cl.ListApplicationsFiltered(ctx, filterParams)
	} else {
		return cl.ListApplications(ctx)
	}
}

func initApplicationListFlagSet(arg *ApplicationFilter, fs *flag.FlagSet) {
	fs.StringVar(&arg.Name, "name", "", "Filter to application name")
	fs.StringVar(&arg.Username, "username", "", "Filter to owning user name")
}

func init() {
	subCmdApplicationList = &SubcommandTemplate[ApplicationFilter, []masherytypes.Application]{
		Command:     []string{"application", "list"},
		Arg:         ApplicationFilter{},
		FlagSetInit: initApplicationListFlagSet,
		Executor:    execApplicationList,
		Template:    mustTemplate(applicationListTemplate),
	}

	enableSubcommand(subCmdApplicationList.Finder())
}
