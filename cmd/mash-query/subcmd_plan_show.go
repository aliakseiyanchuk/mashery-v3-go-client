package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
)

type ShowPlanData struct {
	packageId string
	planId    string
}

func showPackagePlans(ctx context.Context, cl v3client.Client, args interface{}) int {
	p, _ := args.(ShowPlanData)

	if p.planId == "" {
		if srv, gerr := cl.ListPlans(ctx, masherytypes.PackageIdentityFrom(p.packageId)); gerr == nil {
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
		if endp, exists, gerr := cl.GetPlan(ctx, masherytypes.PackagePlanIdentityFrom(p.packageId, p.planId)); gerr == nil {
			if exists {
				_ = jsonEncoder.Encode(&endp)
			}
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
