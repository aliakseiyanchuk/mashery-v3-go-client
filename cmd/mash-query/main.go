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

const customNetTTLOpt = "net-ttl"
const endpointOpt = "endpoint"
const bearerEnvironmentOpt = "bearer-token-env"
const tokenEnvironmentOpt = "vault-token-env"
const tokenResourceEnvironmentOpt = "vault-token-resource-env"
const tokenResourceOpt = "vault-token-resource"
const qpsOps = "qps"
const outputJsonOps = "as-json"
const helpOpt = "help"

var qps int64
var travelTimeComp string
var endpoint string
var envBearerToken string
var envVaultToken string
var cliVaultTokenResource string
var envVaultTokenResource string
var globalOptOutputJson bool
var showHelp bool
var jsonEncoder *json.Encoder

type ExecutorFunc func(context.Context, v3client.Client, []string) int

var subCommandFinders []*SubcommandFinder

func init() {
	jsonEncoder = json.NewEncoder(os.Stdout)
	jsonEncoder.SetIndent("", "  ")
}

func getEffectiveVaultTokenResource() string {
	if len(cliVaultTokenResource) > 0 {
		return cliVaultTokenResource
	}
	return os.Getenv(envVaultTokenResource)
}

func authorizer() (transport.Authorizer, error) {
	vaultTokenResource := getEffectiveVaultTokenResource()
	vaultToken := transport.VaultToken(os.Getenv(envVaultToken))

	if len(vaultTokenResource) > 0 && len(vaultToken) > 0 && len(endpoint) == 0 {
		return transport.NewVaultTokenResourceAuthorizer(vaultTokenResource, vaultToken), nil
	}

	if len(vaultToken) > 0 && len(endpoint) > 0 {
		return transport.NewVaultAuthorizer(vaultToken), nil
	}

	if len(endpoint) == 0 && len(envBearerToken) > 0 {
		if bearerToken := os.Getenv(envBearerToken); len(bearerToken) > 0 {
			return transport.NewBearerAuthorizer(bearerToken), nil
		}
	}

	return nil, errors.New("no suitable authorization supplied")
}

func enableSubcommand(cmd *SubcommandFinder) {
	subCommandFinders = append(subCommandFinders, cmd)
}

func main() {
	flag.Int64Var(&qps, qpsOps, 2, "Observe specified queries-per-second while querying")
	flag.StringVar(&travelTimeComp, customNetTTLOpt, "173ms", "Consider specified network travel time")
	flag.StringVar(&endpoint, endpointOpt, "", "A non-standard endpoint to connect to")
	flag.StringVar(&envBearerToken, bearerEnvironmentOpt, "MASH_BEARER_TOKEN", "An environment variable containing bearer token")
	flag.StringVar(&envVaultToken, tokenEnvironmentOpt, "VAULT_TOKEN", "An environment variable containing HashiCorp vault access token")
	flag.StringVar(&envVaultTokenResource, tokenResourceEnvironmentOpt, "VAULT_TOKEN_RESOURCE", "An environment variable containing HashiCorp resource that will provide the V3 access token")
	flag.StringVar(&cliVaultTokenResource, tokenResourceOpt, "", "URL of the resource in the Vault to read an access token from. Requires specifying Vault credentials.")
	flag.BoolVar(&globalOptOutputJson, outputJsonOps, false, "Output JSON rather than a pretty-printed template")
	flag.BoolVar(&showHelp, helpOpt, false, "Show help options")
	flag.Parse()

	if showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	subCmd := flag.Args()
	if len(subCmd) == 0 {
		fmt.Println("Sub-command required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	execFunc := locateSubCommandExecutor(subCmd)

	if execFunc == nil {
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

		exitCode := execFunc(ctx, cl, subCmd)
		os.Exit(exitCode)
	}
}

func locateSubCommandExecutor(subCmd []string) ExecutorFunc {
	specificity := 0
	var execFunc ExecutorFunc

	for _, p := range subCommandFinders {
		if match := p.Matches(subCmd); match > specificity {
			execFunc = p.Executor
		}
	}
	return execFunc
}
