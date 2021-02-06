package main

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func listPackageKeys(ctx context.Context, cl *v3client.HttpTransport, rawArg interface{}) int {
	if dat, err := cl.ListPackageKeys(ctx); err == nil {
		fmt.Printf("Found %d packages keys", len(dat))
		fmt.Println()

		for idx, srv := range dat {
			fmt.Printf("%d. Package key %s, id=%s, created=%s", idx, srv.Apikey, srv.Id, srv.Created.ToString())
			fmt.Println()
		}
	} else {
		fmt.Printf("Could not have package listed: %s", err)
		fmt.Println()
		return 1
	}

	return 0
}

func listPackageKeysArgParser() (bool, error) {
	if argAt(0) == "packageKey" && argAt(1) == "list" {
		handler = listPackageKeys
		return true, nil
	}

	return false, nil
}

func init() {
	argParsers = append(argParsers, listPackageKeysArgParser)
}
