package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	mashery_v3_go_client "github.com/aliakseiyanchuk/mashery-v3-go-client"
	"os"
)

func showPackagePlansServices(ctx context.Context, cl *mashery_v3_go_client.Client, args interface{}) int {
	p, _ := args.(ShowPlanData)

	if srv, gerr := cl.ListPlanServices(ctx, p.packageId, p.planId); gerr == nil {
		fmt.Printf("Package %s, plan %s includes %d services:", p.packageId, p.planId, len(srv))
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

	return 0
}

func showPackagePlansServicesArgParser() (bool, error) {
	if argAt(0) == "plan" && argAt(1) == "services" {
		fg := flag.NewFlagSet("plan services", flag.ExitOnError)
		var param ShowPlanData

		fg.StringVar(&param.packageId, "p", "", "Package ID")
		fg.StringVar(&param.planId, "pln", "", "Plan ID")

		err := fg.Parse(os.Args[3:])

		if param.packageId == "" || param.planId == "" {
			return true, errors.New("plan services required both package and plan Ids to be specified")
		}

		handler = showPackagePlansServices
		handlerArgs = param

		return true, err
	} else {
		return false, nil
	}
}

func init() {
	argParsers = append(argParsers, showPackagePlansServicesArgParser)
}
