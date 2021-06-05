package main

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
)

func listRoles(ctx context.Context, cl v3client.Client, _ interface{}) int {
	if dat, err := cl.ListRoles(ctx); err == nil {
		fmt.Printf("Found %d roles", len(dat))
		fmt.Println()

		for idx, role := range dat {
			fmt.Printf("%d. Role %s, id=%s, created=%s, predefine=%t, org=%t, assignable=%t", idx, role.Name, role.Id, role.Created.ToString(), role.Predefined, role.OrgRole, role.Assignable)
			fmt.Println()
		}
	} else {
		fmt.Printf("Could not have roles listed: %s", err)
		fmt.Println()
		return 1
	}

	return 0
}

func listRolesArgParser() (bool, error) {
	if argAt(0) == "role" {
		if argAt(1) == "list" {
			handler = listRoles
			return true, nil
		}
	}

	return false, nil
}

func init() {
	argParsers = append(argParsers, listRolesArgParser)
}
