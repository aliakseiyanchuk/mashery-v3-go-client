package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
)

func showServiceSecurityProfileData(ctx context.Context, cl v3client.Client, rawIds interface{}) int {
	ids, _ := rawIds.([]string)

	for _, id := range ids {
		if srv, err := cl.GetServiceOAuthSecurityProfile(ctx, id); err == nil {
			fmt.Printf("OAuth Security Profile for Service %s:", id)
			fmt.Println()

			_ = jsonEncoder.Encode(&srv)
		} else {
			fmt.Printf("ERROR: Failed to retrieve OAuth proile for service %s: %s", id, err)
		}
		fmt.Println()
	}

	return 0
}

func showServiceSecurityProfileDataArgParser() (bool, error) {
	if argAt(0) == "service" && argAt(1) == "sec-profile" && argAt(2) == "show" {
		if len(os.Args) > 2 {
			handler = showServiceSecurityProfileData
			handlerArgs = os.Args[4:]
			return true, nil
		} else {
			return true, errors.New("service sec-profile show requires at least once service Id parameter")
		}
	}

	return false, nil
}

func init() {
	argParsers = append(argParsers, showServiceSecurityProfileDataArgParser)
}
