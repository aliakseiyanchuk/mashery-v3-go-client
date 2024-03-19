package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

type ApplicationTagArg struct {
	masherytypes.ApplicationIdentifier
	Tag string
}

func validateApplicationShowTag(arg *ApplicationTagArg) error {
	if appIdErr := validateApplicationShowArg(&arg.ApplicationIdentifier); appIdErr != nil {
		return appIdErr
	}

	if len(arg.Tag) == 0 {
		return errors.New("application tag required")
	}

	return nil
}

func execApplicationTag(ctx context.Context, cl v3client.Client, arg ApplicationTagArg) (ObjectWithExists[masherytypes.ApplicationIdentifier, masherytypes.Application], error) {
	rv, appExists, err := cl.GetApplication(ctx, arg.ApplicationIdentifier)

	if appExists {
		// Attempt updating the application with the passed tag value
		rv.Tags = arg.Tag
		rv, err = cl.UpdateApplication(ctx, rv)

		if err == nil {
			if appkeys, appKeysErr := cl.GetApplicationPackageKeys(ctx, arg.ApplicationIdentifier); appKeysErr != nil {
				err = appKeysErr
			} else {
				rv.PackageKeys = &appkeys
			}
		}
	}

	return ObjectWithExists[masherytypes.ApplicationIdentifier, masherytypes.Application]{
		Identifier: arg.ApplicationIdentifier,
		Object:     rv,
		Exists:     appExists,
	}, err
}

var subCmdApplicationTag *SubcommandTemplate[ApplicationTagArg, ObjectWithExists[masherytypes.ApplicationIdentifier, masherytypes.Application]]

func initApplicationTagFlagSet(arg *ApplicationTagArg, fs *flag.FlagSet) {
	initApplicationShowFlagSet(&arg.ApplicationIdentifier, fs)
	fs.StringVar(&arg.Tag, "tag", "", "Tag to be supplied")
}

func init() {
	subCmdApplicationTag = &SubcommandTemplate[ApplicationTagArg, ObjectWithExists[masherytypes.ApplicationIdentifier, masherytypes.Application]]{
		Command:     []string{"application", "tag"},
		FlagSetInit: initApplicationTagFlagSet,
		Validator:   validateApplicationShowTag,
		Executor:    execApplicationTag,
		Template:    mustTemplate(applicationShowTemplate),
	}

	enableSubcommand(subCmdApplicationTag.Finder())
}
