package main

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func listApplications(ctx context.Context, cl v3client.Client, rawArg interface{}) int {
	if dat, err := cl.ListApplications(ctx); err == nil {
		fmt.Printf("Found %d applications", len(dat))
		fmt.Println()

		for idx, app := range dat {
			fmt.Printf("%d. Application %s, id=%s, created=%s", idx, app.Name, app.Id, app.Created.ToString())
			fmt.Println()
		}
	} else {
		fmt.Printf("Could not have applications listed: %s", err)
		fmt.Println()
		return 1
	}

	return 0
}

func listApplicationsArgParser() (bool, error) {
	if argAt(0) == "application" && argAt(1) == "list" {
		handler = listApplications
		return true, nil
	}

	return false, nil
}

func init() {
	argParsers = append(argParsers, listApplicationsArgParser)
}
