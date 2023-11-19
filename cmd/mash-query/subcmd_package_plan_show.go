package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func validatePackagePLanShowArg(arg *masherytypes.PackagePlanIdentifier) error {
	if len(arg.PackageId) == 0 {
		return errors.New("package identifier required")
	} else if len(arg.PlanId) == 0 {
		return errors.New("package plan identifier required")
	}

	return nil
}

func execPackagePlanShow(ctx context.Context, cl v3client.Client, id masherytypes.PackagePlanIdentifier) (ObjectWithExists[masherytypes.PackagePlanIdentifier, masherytypes.Plan], error) {
	rv, serviceExists, err := cl.GetPlan(ctx, id)

	return ObjectWithExists[masherytypes.PackagePlanIdentifier, masherytypes.Plan]{
		Identifier: id,
		Object:     rv,
		Exists:     serviceExists,
	}, err
}

//go:embed templates/package_plan_show.tmpl
var packagePlanShowTemplate string
var subCmdPackagePlanShow *SubcommandTemplate[masherytypes.PackagePlanIdentifier, ObjectWithExists[masherytypes.PackagePlanIdentifier, masherytypes.Plan]]

func initPackagePlanShowFlagSet(arg *masherytypes.PackagePlanIdentifier, fs *flag.FlagSet) {
	initPackageShowFlagSet(&arg.PackageIdentifier, fs)
	fs.StringVar(&arg.PlanId, "plan-id", "", "Package plan identifier")
}

func initPackagePlanShowEnvFlagSet(arg *masherytypes.PackagePlanIdentifier) []EnvFlag {
	rv := initPackageShowEnvFlagSet(&arg.PackageIdentifier)
	rv = append(rv, EnvFlag{
		Dest:   &arg.PlanId,
		EnvVar: "MASH_PACKAGE_PLAN_ID",
		Option: "plan-id",
	})

	return rv
}

func init() {
	subCmdPackagePlanShow = &SubcommandTemplate[masherytypes.PackagePlanIdentifier, ObjectWithExists[masherytypes.PackagePlanIdentifier, masherytypes.Plan]]{
		Command:        []string{"package", "plan", "show"},
		FlagSetInit:    initPackagePlanShowFlagSet,
		EnvFlagSetInit: initPackagePlanShowEnvFlagSet,
		Validator:      validatePackagePLanShowArg,
		Executor:       execPackagePlanShow,
		Template:       mustTemplate(packagePlanShowTemplate),
	}

	enableSubcommand(subCmdPackagePlanShow.Finder())
}
