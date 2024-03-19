package main

import (
	"context"
	_ "embed"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

//go:embed templates/roles_list.tmpl
var rolesListTemplate string
var subCmdRoleList *SubcommandTemplate[int, []masherytypes.Role]

func execRoleList(ctx context.Context, cl v3client.Client, _ int) ([]masherytypes.Role, error) {
	return cl.ListRoles(ctx)
}

func init() {
	subCmdRoleList = &SubcommandTemplate[int, []masherytypes.Role]{
		Command:  []string{"role", "list"},
		Executor: execRoleList,
		Template: mustTemplate(rolesListTemplate),
	}

	enableSubcommand(subCmdRoleList.Finder())
}
