package main

import (
	"context"
	"errors"
	"fmt"
	mashery_v3_go_client "github.com/aliakseiyanchuk/mashery-v3-go-client"
	"os"
)

func showPackageData(ctx context.Context, cl *mashery_v3_go_client.Client, rawIds interface{}) int {
	ids, _ := rawIds.([]string)

	for _, id := range ids {
		if srv, err := cl.GetPackage(ctx, id); err == nil {
			fmt.Printf("Package %s:", id)
			fmt.Println()

			_ = jsonEncoder.Encode(&srv)
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
