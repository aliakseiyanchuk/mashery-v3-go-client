package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

//go:embed templates/domain_list.tmpl
var domainListTemplate string
var subCmdDomainList *SubcommandTemplate[DomainType, []masherytypes.DomainAddress]

type DomainType struct {
	SysDomain bool
	PubDomain bool
}

func execDomainListList(ctx context.Context, cl v3client.Client, dt DomainType) ([]masherytypes.DomainAddress, error) {
	if dt.PubDomain {
		return cl.GetPublicDomains(ctx)
	} else {
		return cl.GetSystemDomains(ctx)
	}
}

func initDomainListFlagSet(arg *DomainType, fs *flag.FlagSet) {
	fs.BoolVar(&arg.PubDomain, "public", false, "Show public domains")
	fs.BoolVar(&arg.SysDomain, "system", false, "Show system domains")
}

func validateDomainList(arg *DomainType) error {
	if arg.SysDomain && arg.PubDomain {
		return errors.New("system and public domains are mutually exclusive")
	} else if !arg.SysDomain && !arg.PubDomain {
		return errors.New("either -system or -public is required for this command")
	} else {
		return nil
	}
}

func init() {
	subCmdDomainList = &SubcommandTemplate[DomainType, []masherytypes.DomainAddress]{
		Command:     []string{"domain", "list"},
		Arg:         DomainType{},
		FlagSetInit: initDomainListFlagSet,
		Validator:   validateDomainList,
		Executor:    execDomainListList,
		Template:    mustTemplate(domainListTemplate),
	}

	enableSubcommand(subCmdDomainList.Finder())
}
