package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

type PackageKeySecretRevealArg struct {
	masherytypes.PackageKeyIdentifier
	Full bool
}

type PackageKeyWithDisplay struct {
	masherytypes.PackageKey
	DisplayFull bool
}

func validatePackageKeySecretRevealArg(arg *PackageKeySecretRevealArg) error {
	if len(arg.PackageKeyId) == 0 {
		return errors.New("package key identifier required")
	}

	return nil
}

func execPackageKeyReveal(ctx context.Context, cl v3client.Client, id PackageKeySecretRevealArg) (ObjectWithExists[masherytypes.PackageKeyIdentifier, PackageKeyWithDisplay], error) {
	rv, keyExists, err := cl.GetPackageKey(ctx, id.PackageKeyIdentifier)

	return ObjectWithExists[masherytypes.PackageKeyIdentifier, PackageKeyWithDisplay]{
		Identifier: id.PackageKeyIdentifier,
		Object: PackageKeyWithDisplay{
			PackageKey:  rv,
			DisplayFull: id.Full,
		},
		Exists: keyExists,
	}, err
}

//go:embed templates/package_key_secret_reveal.tmpl
var packageKeySecretRevealTemplate string
var subCmdPackageKeyReveal *SubcommandTemplate[PackageKeySecretRevealArg, ObjectWithExists[masherytypes.PackageKeyIdentifier, PackageKeyWithDisplay]]

func initPackageKeySecretRevealFlagSet(arg *PackageKeySecretRevealArg, fs *flag.FlagSet) {
	fs.StringVar(&arg.PackageKeyId, "key-id", "", "package key identifier")
	fs.BoolVar(&arg.Full, "full", false, "If full, the complete display ")
}

func initPackageKeySecretRevealEnvFlagSet(arg *PackageKeySecretRevealArg) []EnvFlag {
	rv := []EnvFlag{}

	rv = append(rv, EnvFlag{
		Dest:   &arg.PackageKeyId,
		EnvVar: "MASH_PACKAGE_KEY_ID",
		Option: "key-id",
	})

	return rv
}

func init() {
	subCmdPackageKeyReveal = &SubcommandTemplate[PackageKeySecretRevealArg, ObjectWithExists[masherytypes.PackageKeyIdentifier, PackageKeyWithDisplay]]{
		Command:        []string{"package", "key", "secret", "reveal"},
		FlagSetInit:    initPackageKeySecretRevealFlagSet,
		EnvFlagSetInit: initPackageKeySecretRevealEnvFlagSet,
		Validator:      validatePackageKeySecretRevealArg,
		Executor:       execPackageKeyReveal,
		Template:       mustTemplate(packageKeySecretRevealTemplate),
	}

	enableSubcommand(subCmdPackageKeyReveal.Finder())
}
