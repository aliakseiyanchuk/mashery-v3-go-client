package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
)

type ShowEndpointMethodData struct {
	ShowServiceEndpointData
	methodId string
	filterId string
}

func showMethodFilter(ctx context.Context, cl v3client.Client, args interface{}) int {
	p, _ := args.(ShowEndpointMethodData)

	if srv, exists, gerr := cl.GetEndpointMethodFilter(ctx, masherytypes.ServiceEndpointMethodFilterIdentityFrom(p.serviceId, p.endpointId, p.methodId, p.filterId)); gerr == nil {
		fmt.Println()
		if exists {
			_ = jsonEncoder.Encode(&srv)
		}

		return 0
	} else {
		fmt.Println(gerr)
		return 1
	}
}

func listMethodFilters(ctx context.Context, cl v3client.Client, args interface{}) int {
	p, _ := args.(ShowEndpointMethodData)

	if srv, gerr := cl.ListEndpointMethodFiltersWithFullInfo(ctx, masherytypes.ServiceEndpointMethodIdentityFrom(p.serviceId, p.endpointId, p.methodId)); gerr == nil {
		fmt.Printf("There are %d filters", len(srv))
		fmt.Println()

		_ = jsonEncoder.Encode(&srv)

		return 0
	} else {
		fmt.Println(gerr)
		return 1
	}
}

func showMethodFilterArgParser() (bool, error) {
	if argAt(0) == "endpoint" && argAt(1) == "method" && argAt(2) == "filter" {
		fg := flag.NewFlagSet("endpoint show", flag.ExitOnError)
		var param ShowEndpointMethodData

		fg.StringVar(&param.serviceId, "s", "nil", "Service ID")
		fg.StringVar(&param.endpointId, "e", "nil", "Endpoint ID")
		fg.StringVar(&param.methodId, "m", "nil", "Method ID")
		fg.StringVar(&param.filterId, "f", "nil", "Filter ID")

		err := fg.Parse(os.Args[5:])

		handlerArgs = param
		if argAt(3) == "show" {
			handler = showMethodFilter
		} else if (argAt(3)) == "list" {
			handler = listMethodFilters
		} else {
			return false, nil
		}

		return true, err
	} else {
		return false, nil
	}
}

func init() {
	argParsers = append(argParsers, showMethodFilterArgParser)
}
