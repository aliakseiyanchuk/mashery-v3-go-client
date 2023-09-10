package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
)

func showPackageData(ctx context.Context, cl v3client.Client, rawIds interface{}) int {
	ids, _ := rawIds.([]string)

	for _, id := range ids {
		if srv, exists, err := cl.GetPackage(ctx, masherytypes.PackageIdentityFrom(id)); err == nil {
			fmt.Printf("Package %s:", id)
			fmt.Println()

			if exists {
				_ = jsonEncoder.Encode(&srv)
			}
		} else {
			fmt.Println(logHeader)
			fmt.Printf("ERROR: Failed to retrieve package %s: %s", id, err)
			fmt.Println(logHeader)
		}
		fmt.Println()
	}

	return 0
}

func showSPackageDataArgParser() (bool, error) {
	if argAt(0) == "package" && argAt(1) == "show" {
		if len(os.Args) > 2 {
			handler = showPackageData
			handlerArgs = os.Args[3:]
			return true, nil
		} else {
			return true, errors.New("package show requires at least once package Id parameter")
		}
	}

	return false, nil
}

func init() {
	argParsers = append(argParsers, showSPackageDataArgParser)
}
