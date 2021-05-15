package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func showPublicDomains(ctx context.Context, cl v3client.Client, args interface{}) int {
	if domains, err := cl.GetPublicDomains(ctx); err != nil {
		fmt.Println(err)
		return 1
	} else {
		fmt.Printf("There are %d public domains", len(domains))
		fmt.Println()

		for idx, d := range domains {
			fmt.Printf("%d. %s", (idx + 1), d)
			fmt.Println()
		}

		return 0
	}
}

func showSystemDomains(ctx context.Context, cl v3client.Client, args interface{}) int {
	if domains, err := cl.GetSystemDomains(ctx); err != nil {
		fmt.Println(err)
		return 1
	} else {
		fmt.Printf("There are %d system domains", len(domains))
		fmt.Println()

		for idx, d := range domains {
			fmt.Printf("%d. %s", (idx + 1), d)
			fmt.Println()
		}

		return 0
	}
}

func showDomainsParser() (bool, error) {
	if argAt(0) == "domains" && argAt(1) == "list" {
		if argAt(2) == "public" {
			handler = showPublicDomains
			return true, nil
		} else if argAt(2) == "system" {
			handler = showSystemDomains
			return true, nil
		} else {
			return true, errors.New("domain list requires domain type to list, either public or system")
		}
	}

	return false, nil
}

func init() {
	argParsers = append(argParsers, showDomainsParser)
}
