package main

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func listPackages(ctx context.Context, cl v3client.Client, rawArg interface{}) int {
	if dat, err := cl.ListPackages(ctx); err == nil {
		fmt.Printf("Found %d packages", len(dat))
		fmt.Println()

		for idx, srv := range dat {
			fmt.Printf("%d. Package %s, id=%s, created=%s", idx, srv.Name, srv.Id, srv.Created.ToString())
			fmt.Println()
		}
	} else {
		fmt.Printf("Could not have package listed: %s", err)
		fmt.Println()
		return 1
	}

	return 0
}

func listPackageArgParser() (bool, error) {
	if argAt(0) == "package" && argAt(1) == "list" {
		handler = listPackages
		return true, nil
	}

	return false, nil
}

func init() {
	argParsers = append(argParsers, listPackageArgParser)
}
