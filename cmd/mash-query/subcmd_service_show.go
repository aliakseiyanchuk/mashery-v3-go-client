package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
)

func showServiceData(ctx context.Context, cl v3client.Client, rawIds interface{}) int {
	ids, _ := rawIds.([]string)

	for _, id := range ids {
		if srv, err := cl.GetService(ctx, id); err == nil {
			fmt.Printf("Service %s:", id)
			fmt.Println()

			_ = jsonEncoder.Encode(&srv)
		} else {
			fmt.Printf("ERROR: Failed to retrieve service %s: %s", id, err)
		}
		fmt.Println()
	}

	return 0
}

func showServiceDataArgParser() (bool, error) {
	if argAt(0) == "service" && argAt(1) == "show" {
		if len(os.Args) > 2 {
			handler = showServiceData
			handlerArgs = os.Args[3:]
			return true, nil
		} else {
			return true, errors.New("service show requires at least once service Id parameter")
		}
	}

	return false, nil
}

func init() {
	argParsers = append(argParsers, showServiceDataArgParser)
}
