package main

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func listServices(ctx context.Context, cl *v3client.Client, rawArg interface{}) int {
	if dat, err := cl.ListServices(ctx); err == nil {
		fmt.Printf("Found %d services", len(dat))
		fmt.Println()

		for idx, srv := range dat {
			fmt.Printf("%d. Service %s, id=%s, created=%s", idx, srv.Name, srv.Id, srv.Created.ToString())
			fmt.Println()
		}
	} else {
		fmt.Printf("Could not have services listed: %s", err)
		fmt.Println()
		return 1
	}

	return 0
}

func listServiceArgParser() (bool, error) {
	if argAt(0) == "service" && argAt(1) == "list" {
		handler = listServices
		return true, nil
	}

	return false, nil
}

func init() {
	argParsers = append(argParsers, listServiceArgParser)
}
