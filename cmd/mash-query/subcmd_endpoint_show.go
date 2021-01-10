package main

import (
	"context"
	"flag"
	"fmt"
	mashery_v3_go_client "github.com/aliakseiyanchuk/mashery-v3-go-client"
	"os"
)

type ShowServiceEndpointData struct {
	serviceId  string
	endpointId string
}

func showServiceEndpoints(ctx context.Context, cl *mashery_v3_go_client.Client, args interface{}) int {
	p, _ := args.(ShowServiceEndpointData)

	if p.endpointId == "" {
		if srv, gerr := cl.ListEndpoints(ctx, p.serviceId); gerr == nil {
			fmt.Printf("Service %s defines %d endpoints:", p.serviceId, len(srv))
			fmt.Println()

			for idx, e := range srv {
				fmt.Printf("%d. %s (id=%s)", (idx + 1), e.Name, e.Id)
				fmt.Println()
			}

			return 0
		} else {
			fmt.Println(gerr)
			return 1
		}
	} else {
		if endp, gerr := cl.GetEndpoint(ctx, p.serviceId, p.endpointId); gerr == nil {
			_ = jsonEncoder.Encode(&endp)
		} else {
			fmt.Println(gerr)
			return 1
		}

	}
	return 0
}

func showServiceEndpointsArgParser() (bool, error) {
	if argAt(0) == "endpoint" && argAt(1) == "show" {
		fg := flag.NewFlagSet("endpoint show", flag.ExitOnError)
		var param ShowServiceEndpointData

		fg.StringVar(&param.serviceId, "s", "", "Service ID")
		fg.StringVar(&param.endpointId, "e", "", "Endpoint ID")

		err := fg.Parse(os.Args[3:])

		handler = showServiceEndpoints
		handlerArgs = param

		return true, err
	} else {
		return false, nil
	}
}

func init() {
	argParsers = append(argParsers, showServiceEndpointsArgParser)
}
