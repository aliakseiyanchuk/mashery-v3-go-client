package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
)

func showMemberData(ctx context.Context, cl v3client.Client, rawIds interface{}) int {
	ids, _ := rawIds.([]string)

	for _, id := range ids {
		if srv, err := cl.GetMember(ctx, id); err == nil {
			fmt.Printf("Member %s:", id)
			fmt.Println()

			_ = jsonEncoder.Encode(&srv)
		} else {
			fmt.Printf("ERROR: Failed to retrieve member %s: %s", id, err)
		}
		fmt.Println()
	}

	return 0
}

func showMemberDataArgParser() (bool, error) {
	if argAt(0) == "member" && argAt(1) == "show" {
		if len(os.Args) > 2 {
			handler = showMemberData
			handlerArgs = os.Args[3:]
			return true, nil
		} else {
			return true, errors.New("member show requires at least one member Id parameter")
		}
	}

	return false, nil
}

func init() {
	argParsers = append(argParsers, showMemberDataArgParser)
}
