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

var handler func(context.Context, *v3client.HttpTransport, interface{}) int = nil
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

func main() {
	fmt.Println("----------------------------------")

	flag.StringVar(&cmdTokenFile, customTokenFileOpt, v3client.SavedAccessTokenFile(), "Use locally saved token file")
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
	var exitCode = 0

	if tkn, err := v3client.ReadSavedV3TokenData(cmdTokenFile); err == nil && tkn != nil {
		ctx := context.TODO()

		dur, durErr := time.ParseDuration(travelTimeComp)
		if durErr != nil {
			dur = 173 * time.Millisecond
		}
		cl := v3client.NewHttpClient(v3client.NewFixedTokenProvider(tkn.AccessToken), qps, dur)

		exitCode = handler(ctx, &cl, handlerArgs)
	} else {
		fmt.Println("Could not load token file:")
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}
