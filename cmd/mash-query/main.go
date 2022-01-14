package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
	"time"
)

const logHeader = "--------------------------------------------------------------------------------------"

const customTokenFileOpt = "token-file"
const customNetTTLOpt = "net-ttl"
const qpsOps = "qps"

var cmdTokenFile string
var qps int64
var travelTimeComp string
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

func tokenProvider() (v3client.V3AccessTokenProvider, error) {
	if envProp := os.Getenv(v3client.MasheryTokenSystemProperty); len(envProp) > 0 {
		return v3client.NewFixedTokenProvider(envProp), nil
	}

	fsProvider := v3client.NewFileSystemTokenProviderFrom(cmdTokenFile)
	_, err := fsProvider.AccessToken()

	return fsProvider, err
}

func main() {
	fmt.Println("----------------------------------")

	flag.StringVar(&cmdTokenFile, customTokenFileOpt, v3client.DefaultSavedAccessTokenFilePath(), "Use locally saved token file")
	flag.Int64Var(&qps, qpsOps, 2, "Observe specified queries-per-second while querying")
	flag.StringVar(&travelTimeComp, customNetTTLOpt, "173ms", "Consider specified network travel time")
	flag.Parse()
	subCmd = flag.Args()

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
		os.Exit(1)
	}

	// Arguments have been parsed correctly.
	if tknProvider, err := tokenProvider(); err != nil {
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
			Authorizer:    tknProvider,
			QPS:           qps,
			AvgNetLatency: dur,
		})

		exitCode := handler(ctx, cl, handlerArgs)
		os.Exit(exitCode)
	}
}
