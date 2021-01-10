package main

import (
	"context"
	"flag"
	"fmt"
	mashery_v3_go_client "github.com/aliakseiyanchuk/mashery-v3-go-client"
	"os"
)

type ShowPlanData struct {
	packageId string
	planId    string
}

func showPackagePlans(ctx context.Context, cl *mashery_v3_go_client.Client, args interface{}) int {
	p, _ := args.(ShowPlanData)

	if p.planId == "" {
		if srv, gerr := cl.ListPlans(ctx, p.packageId); gerr == nil {
			fmt.Printf("Package %s defines %d plans:", p.packageId, len(srv))
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
		if endp, gerr := cl.GetPlan(ctx, p.packageId, p.planId); gerr == nil {
			_ = jsonEncoder.Encode(&endp)
		} else {
			fmt.Println(gerr)
			return 1
		}

	}
	return 0
}

func showPackagePansArgParser() (bool, error) {
	if argAt(0) == "plan" && argAt(1) == "show" {
		fg := flag.NewFlagSet("plan show", flag.ExitOnError)
		var param ShowPlanData

		fg.StringVar(&param.packageId, "p", "", "Package ID")
		fg.StringVar(&param.planId, "pln", "", "Plan ID")

		err := fg.Parse(os.Args[3:])

		handler = showPackagePlans
		handlerArgs = param

		return true, err
	} else {
		return false, nil
	}
}

func init() {
	argParsers = append(argParsers, showPackagePansArgParser)
}
