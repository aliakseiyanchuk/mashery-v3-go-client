package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
)

func showServiceData(ctx context.Context, cl v3client.Client, rawIds interface{}) int {
	ids, _ := rawIds.([]string)

	for _, id := range ids {
		if srv, exists, err := cl.GetService(ctx, masherytypes.ServiceIdentityFrom(id)); exists && err == nil {
			fmt.Printf("Service %s:", id)
			fmt.Println()

			_ = jsonEncoder.Encode(&srv)

			errorSets, _ := cl.ListErrorSets(ctx, srv.Identifier(), v3client.EmptyQuery)
			fmt.Printf("Service contains %d error sets\n", len(errorSets))
			for v, idx := range errorSets {
				fmt.Printf("%d. %s type=%s, jsonp: %t, jsonpType=%s (id=%s)", v+1, idx.Name, idx.Type, idx.JSONP, idx.JSONPType, idx.Id)

				if es, errorSetExists, errorSetFetchErr := cl.GetErrorSet(ctx, idx.Identifier()); errorSetFetchErr != nil {
					fmt.Printf("Can't retrieve error set: %s", errorSetFetchErr)
				} else {
					if errorSetExists {
						for idx, v := range *es.ErrorMessages {
							fmt.Printf("%d. %s\n", idx+1, v.Id)
						}
					}
				}
			}
		} else if !exists && err == nil {
			fmt.Println("Service with this identifier is not found")
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
