package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func validatePackageShowArg(arg *masherytypes.PackageIdentifier) error {
	if len(arg.PackageId) == 0 {
		return errors.New("package identifier required")
	}

	return nil
}

func execPackageShow(ctx context.Context, cl v3client.Client, id masherytypes.PackageIdentifier) (ObjectWithExists[masherytypes.PackageIdentifier, masherytypes.Package], error) {
	rv, packageExists, err := cl.GetPackage(ctx, id)

	if packageExists {
		if plans, plansErr := cl.ListPlans(ctx, id); plansErr != nil {
			err = plansErr
		} else {
			rv.Plans = plans
		}
	}

	return ObjectWithExists[masherytypes.PackageIdentifier, masherytypes.Package]{
		Identifier: id,
		Object:     rv,
		Exists:     packageExists,
	}, err
}

//go:embed templates/package_show.tmpl
var packageShowTemplate string
var subCmdPackageShow *SubcommandTemplate[masherytypes.PackageIdentifier, ObjectWithExists[masherytypes.PackageIdentifier, masherytypes.Package]]

func initPackageShowFlagSet(arg *masherytypes.PackageIdentifier, fs *flag.FlagSet) {
	fs.StringVar(&arg.PackageId, "package-id", "", "Package identifier")
}

func initPackageShowEnvFlagSet(arg *masherytypes.PackageIdentifier) []EnvFlag {
	return []EnvFlag{
		{
			Dest:   &arg.PackageId,
			EnvVar: "MASH_PACKAGE_ID",
			Option: "package-id",
		},
	}
}

func init() {
	subCmdPackageShow = &SubcommandTemplate[masherytypes.PackageIdentifier, ObjectWithExists[masherytypes.PackageIdentifier, masherytypes.Package]]{
		Command:        []string{"package", "show"},
		FlagSetInit:    initPackageShowFlagSet,
		EnvFlagSetInit: initPackageShowEnvFlagSet,
		Validator:      validatePackageShowArg,
		Executor:       execPackageShow,
		Template:       mustTemplate(packageShowTemplate),
	}

	enableSubcommand(subCmdPackageShow.Finder())
}
