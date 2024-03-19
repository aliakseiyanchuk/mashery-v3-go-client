package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func validateApplicationShowArg(arg *masherytypes.ApplicationIdentifier) error {
	if len(arg.ApplicationId) == 0 {
		return errors.New("application identifier required")
	}

	return nil
}

func execApplicationShow(ctx context.Context, cl v3client.Client, id masherytypes.ApplicationIdentifier) (ObjectWithExists[masherytypes.ApplicationIdentifier, masherytypes.Application], error) {
	rv, appExists, err := cl.GetApplication(ctx, id)

	if appExists {
		if appkeys, appKeysErr := cl.GetApplicationPackageKeys(ctx, id); appKeysErr != nil {
			err = appKeysErr
		} else {
			rv.PackageKeys = &appkeys
		}
	}

	return ObjectWithExists[masherytypes.ApplicationIdentifier, masherytypes.Application]{
		Identifier: id,
		Object:     rv,
		Exists:     appExists,
	}, err
}

//go:embed templates/application_show.tmpl
var applicationShowTemplate string
var subCmdApplicationShow *SubcommandTemplate[masherytypes.ApplicationIdentifier, ObjectWithExists[masherytypes.ApplicationIdentifier, masherytypes.Application]]

func initApplicationShowFlagSet(arg *masherytypes.ApplicationIdentifier, fs *flag.FlagSet) {
	fs.StringVar(&arg.ApplicationId, "app-id", "", "Application identifier")
}

func init() {
	subCmdApplicationShow = &SubcommandTemplate[masherytypes.ApplicationIdentifier, ObjectWithExists[masherytypes.ApplicationIdentifier, masherytypes.Application]]{
		Command:     []string{"application", "show"},
		FlagSetInit: initApplicationShowFlagSet,
		Validator:   validateApplicationShowArg,
		Executor:    execApplicationShow,
		Template:    mustTemplate(applicationShowTemplate),
	}

	enableSubcommand(subCmdApplicationShow.Finder())
}
