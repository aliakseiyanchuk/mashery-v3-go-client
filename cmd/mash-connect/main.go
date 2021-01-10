package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client"
	"io/ioutil"
	"os"
)

const customTokenFileOpt = "token-file"
const accessTokenVariableOpt = "env-var"

// Command line utility to log into the Mashery.
var credsFile string
var cmdTokenFile string
var cmdCreds mashery_v3_go_client.MasheryV3Credentials

func refresh() int {
	refreshCmd := flag.NewFlagSet("refresh", flag.ExitOnError)
	refreshCmd.StringVar(&cmdTokenFile, customTokenFileOpt, mashery_v3_go_client.SavedAccessTokenFile(), "Use this file to read/write access token data")
	refreshCmd.StringVar(&credsFile, "creds", "", "Path to the credentials file")
	refreshCmd.StringVar(&(cmdCreds.ApiKey), "k", "", "Mashery V3 API key")
	refreshCmd.StringVar(&(cmdCreds.ApiKey), "apiKey", "", "Mashery V3 API key")
	refreshCmd.StringVar(&(cmdCreds.Secret), "s", "", "Mashery V3 API key secret")
	refreshCmd.StringVar(&(cmdCreds.Secret), "secret", "", "Mashery V3 API key secret")

	_ = refreshCmd.Parse(os.Args[:2])

	if tkn, err := mashery_v3_go_client.ReadSavedV3TokenData(cmdTokenFile); err == nil && tkn != nil {
		creds := mashery_v3_go_client.DeriveAccessCredentials(credsFile)
		creds.Inherit(&cmdCreds)

		ccProvider := mashery_v3_go_client.NewClientCredentialsProvider(creds)
		ccProvider.Response = tkn

		fmt.Println("Trying to refresh Mashery V3 access token...")
		if err := ccProvider.Refresh(); err == nil {
			return saveTokenResponse(ccProvider.Response)
		} else {
			fmt.Printf("Refresh failed: %s", err)
			return 2
		}
	} else {
		fmt.Println("Could not read saved V3 token data")

		if err != nil {
			fmt.Println(err)
		}

		return 1
	}
}

func export() int {
	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	exportCmd.StringVar(&cmdTokenFile, customTokenFileOpt, mashery_v3_go_client.SavedAccessTokenFile(), "Use specified custom file")

	_ = exportCmd.Parse(os.Args[:2])

	if tkn, err := mashery_v3_go_client.ReadSavedV3TokenData(cmdTokenFile); err == nil && tkn != nil {
		fmt.Print(tkn.AccessToken)
		return 0
	} else {
		fmt.Println("Could not read saved V3 token data")

		if err != nil {
			fmt.Println(err)
		}

		return 1
	}
}

func show() int {
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	showCmd.StringVar(&cmdTokenFile, customTokenFileOpt, "", "Use specified custom file")

	_ = showCmd.Parse(os.Args[2:])

	var f string
	if len(cmdTokenFile) > 0 {
		f = cmdTokenFile
	} else {
		f = mashery_v3_go_client.SavedAccessTokenFile()
	}

	if tkn, err := mashery_v3_go_client.ReadSavedV3TokenData(f); err == nil && tkn != nil {
		if tkn.Expired() {
			fmt.Printf("You access token has already expired (on %s)", tkn.ExpiryTime())
		} else {
			minutesLeft := tkn.TimeLeft() / 60
			if minutesLeft <= 1 {
				fmt.Printf("You token has 1 minute or less validatity time. Refresh before next operation")
			} else if minutesLeft < 3 {
				fmt.Printf("There are %d minutes left in your token. Refreshing is advised.", minutesLeft)
			} else {
				fmt.Printf("Your token is still valid for %d minutes", tkn.TimeLeft()/60)
			}
		}
		fmt.Println()
		return 0
	} else {
		fmt.Printf("Current Mashery V3 token information is not available.")
		fmt.Println()

		if err != nil {
			fmt.Printf("%s", err)
			fmt.Println()
		}
		return 1
	}
}

func saveTokenResponse(dat *mashery_v3_go_client.TimedAccessTokenResponse) int {
	if b, err := json.Marshal(dat); err == nil {
		if err = ioutil.WriteFile(cmdTokenFile, b, 0644); err == nil {
			fmt.Printf("Mashery V3 API access token has been successfuly initialized.")
			fmt.Println()

			fmt.Printf("Access token will expire in %d minutes (on %s)", dat.ExpiresIn/60, dat.ExpiryTime())
			return 0
		} else {
			fmt.Printf("ERROR: Could not save credentials to %s (%s)", cmdTokenFile, err)
			return 2
		}
	} else {
		fmt.Printf("Failed to unmarshal response: %s (%s)", dat, err)
		return 2
	}
}

func initToken() int {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	initCmd.StringVar(&credsFile, "creds", "", "Path to the credentials file")
	initCmd.StringVar(&cmdTokenFile, customTokenFileOpt, mashery_v3_go_client.SavedAccessTokenFile(), "Use this file to read/write access token data")
	initCmd.StringVar(&(cmdCreds.AreaId), "a", "", "Mashery V3 Area ID")
	initCmd.StringVar(&(cmdCreds.AreaId), "areaId", "", "Mashery V3 Area ID")
	initCmd.StringVar(&(cmdCreds.ApiKey), "k", "", "Mashery V3 API key")
	initCmd.StringVar(&(cmdCreds.ApiKey), "apiKey", "", "Mashery V3 API key")
	initCmd.StringVar(&(cmdCreds.Secret), "s", "", "Mashery V3 API key secret")
	initCmd.StringVar(&(cmdCreds.Secret), "secret", "", "Mashery V3 API key secret")
	initCmd.StringVar(&(cmdCreds.Username), "u", "", "Mashery V3 user name")
	initCmd.StringVar(&(cmdCreds.Username), "username", "", "Mashery V3 user name")
	initCmd.StringVar(&(cmdCreds.Password), "p", "", "Mashery V3 password")
	initCmd.StringVar(&(cmdCreds.Password), "password", "", "Mashery V3 password")

	if err := initCmd.Parse(os.Args[2:]); err == nil {
		// Derive credentials from all applicable sources, including the command line
		// The sequence of derivation is:
		// - Environment variables, overridden by
		// - User settings file, overridden by
		// - Credentials file in the working directory, overriden by
		// - Command line arguments
		creds := mashery_v3_go_client.DeriveAccessCredentials(credsFile)
		creds.Inherit(&cmdCreds)

		provider := mashery_v3_go_client.NewClientCredentialsProvider(creds)

		if dat, err := provider.TokenData(); err == nil {
			return saveTokenResponse(dat)
		} else {
			fmt.Printf("Error: Token data was not retrieved: %s", err)
			return 1
		}
	} else {
		fmt.Printf("Could not parse command line arguments: %s", err)
		return 1
	}
}

func main() {
	var exitCode int

	switch os.Args[1] {

	case "init":
		exitCode = initToken()
		break
	case "show":
		exitCode = show()
		break
	case "export":
		exitCode = export()
		break
	case "refresh":
		exitCode = refresh()
		break
	default:
		_ = fmt.Sprintf("Sub-command required")
		exitCode = 10
	}

	os.Exit(exitCode)
}
