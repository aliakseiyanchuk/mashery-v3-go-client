package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
)

func showAppKeysData(ctx context.Context, cl *v3client.HttpTransport, rawIds interface{}) int {
	ids, _ := rawIds.([]string)

	for _, id := range ids {
		if srv, err := cl.GetApplicationPackageKeys(ctx, id); err == nil {
			fmt.Printf("Application %s has %d package keys:", id, len(srv))
			fmt.Println()

			for idx, srv := range srv {
				fmt.Printf("%d. Package key %s, id=%s, created=%s for package %s", (idx + 1), srv.Apikey, srv.Id, srv.Created.ToString(), srv.Package.Name)
				fmt.Println()
			}
		} else {
			fmt.Printf("ERROR: Failed to retrieve service %s: %s", id, err)
		}
	}

	return 0
}

func showAppKeysDataArgParser() (bool, error) {
	if argAt(0) == "application" && argAt(1) == "keys" {
		if len(os.Args) > 2 {
			handler = showAppKeysData
			handlerArgs = os.Args[3:]
			return true, nil
		} else {
			return true, errors.New("application show requires at least once application Id parameter")
		}
	}

	return false, nil
}

func init() {
	argParsers = append(argParsers, showAppKeysDataArgParser)
}
