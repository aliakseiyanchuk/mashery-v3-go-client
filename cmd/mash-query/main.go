package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
	"time"
)

const logHeader = "--------------------------------------------------------------------------------------"

const customTokenFileOpt = "token-file"
const customNetTTLOpt = "net-ttl"
const endpointOpt = "endpoint"
const tokenEnvironmentOpt = "token-env"
const qpsOps = "qps"

var cmdTokenFile string
var qps int64
var travelTimeComp string
var endpoint string
var tokenEnv string
var subCmd []string
var jsonEncoder *json.Encoder

var argParsers []func() (bool, error)

var handler func(context.Context, v3client.Client, interface{}) int = nil
var handlerArgs interface{}

func argAt(idx int) string {
	if len(subCmd) > idx {
		return subCmd[idx]
	} else {
		return ""
	}
}

func init() {
	jsonEncoder = json.NewEncoder(os.Stdout)
	jsonEncoder.SetIndent("", "  ")
}

func authorizer() (transport.Authorizer, error) {
	if len(tokenEnv) > 0 {
		if vaultToken := os.Getenv(tokenEnv); len(vaultToken) > 0 {
			return transport.NewVaultAuthorizer(vaultToken), nil
		}
	}

	return nil, errors.New("no suitable authorization supplied")
}

func main() {
	fmt.Println("----------------------------------")

	flag.StringVar(&cmdTokenFile, customTokenFileOpt, v3client.DefaultSavedAccessTokenFilePath(), "Use locally saved token file")
	flag.Int64Var(&qps, qpsOps, 2, "Observe specified queries-per-second while querying")
	flag.StringVar(&travelTimeComp, customNetTTLOpt, "173ms", "Consider specified network travel time")
	flag.StringVar(&endpoint, endpointOpt, "", "A non-standard endpoint to connect to")
	flag.StringVar(&tokenEnv, tokenEnvironmentOpt, "", "An environment ")
	flag.Parse()

	subCmd = flag.Args()
	fmt.Println(subCmd)

	fmt.Println(len(argParsers))
	for _, p := range argParsers {
		rec, err := p()
		if rec {
			if err != nil {
				fmt.Println("Error in command:")
				fmt.Println(err)
				os.Exit(1)
			} else {
				break
			}
		}
	}

	if handler == nil {
		fmt.Println("Unrecognized command")
		for _, p := range subCmd {
			fmt.Println(p)
		}
		os.Exit(1)
	}

	// Arguments have been parsed correctly.
	if tknProvider, err := authorizer(); err != nil {
		fmt.Printf("Access token provider is not ready: %s", err)
		fmt.Println()
		os.Exit(1)
	} else {
		ctx := context.TODO()

		dur, durErr := time.ParseDuration(travelTimeComp)
		if durErr != nil {
			dur = 173 * time.Millisecond
		}

		cl := v3client.NewHttpClient(v3client.Params{
			MashEndpoint:  endpoint,
			Authorizer:    tknProvider,
			QPS:           qps,
			AvgNetLatency: dur,
		})

		exitCode := handler(ctx, cl, handlerArgs)
		os.Exit(exitCode)
	}
}
