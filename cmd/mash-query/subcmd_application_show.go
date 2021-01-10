package main

import (
	"context"
	"errors"
	"fmt"
	mashery_v3_go_client "github.com/aliakseiyanchuk/mashery-v3-go-client"
	"os"
)

func showAppData(ctx context.Context, cl *mashery_v3_go_client.Client, rawIds interface{}) int {
	ids, _ := rawIds.([]string)

	for _, id := range ids {
		if srv, err := cl.GetFullApplication(ctx, id); err == nil {
			fmt.Printf("Application %s (id=%d):", srv.Name, id)
			fmt.Println()

			_ = jsonEncoder.Encode(&srv)
		} else {
			fmt.Printf("ERROR: Failed to retrieve service %s: %s", id, err)
		}
		fmt.Println()
	}

	return 0
}

func showAppDataArgParser() (bool, error) {
	if argAt(0) == "application" && argAt(1) == "show" {
		if len(os.Args) > 2 {
			handler = showAppData
			handlerArgs = os.Args[3:]
			return true, nil
		} else {
			return true, errors.New("application show requires at least once application Id parameter")
		}
	}

	return false, nil
}

func init() {
	argParsers = append(argParsers, showServiceDataArgParser)
}
