package main

import (
	"context"
	_ "embed"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

//go:embed templates/member_list.tmpl
var membersListTemplate string
var subCmdMembersList *SubcommandTemplate[int, []masherytypes.Member]

func execMembersList(ctx context.Context, cl v3client.Client, _ int, params []string) ([]masherytypes.Member, error) {
	if len(params) > 0 {
		return cl.ListMembersFiltered(ctx, kvArrayToMap(params))
	} else {
		return cl.ListMembers(ctx)
	}
}

func init() {
	subCmdMembersList = &SubcommandTemplate[int, []masherytypes.Member]{
		Command:               []string{"member", "list"},
		ParameterizedExecutor: execMembersList,
		Template:              mustTemplate(membersListTemplate),
	}

	enableSubcommand(subCmdMembersList.Finder())
}
