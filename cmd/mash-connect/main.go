package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
)

const customTokenFileOpt = "token-file"
const windowsOpt = "win"

// Command line utility to log into the Mashery.
var credentialsFile string
var cmdTokenFile string

func refresh() int {

	_ = tokenCommandLine("refresh").Parse(os.Args[:2])

	if tkn, err := v3client.ReadSavedV3TokenData(cmdTokenFile); err == nil && tkn != nil {
		runtimeCredentials := v3client.DeriveAccessCredentials(credentialsFile, credentialsPassword(), nil)
		if !runtimeCredentials.FullySpecified() {
			fmt.Println("Insufficient input runtimeCredentials")
			return 1
		}

		ccProvider := v3client.NewOAuthHelper(v3client.OAuthHelperParams{})

		fmt.Println("Trying to refresh Mashery V3 access token...")
		if resp, err := ccProvider.ExchangeRefreshToken(&runtimeCredentials, tkn.RefreshToken); err == nil {
			if err = v3client.PersistV3TokenResponse(resp, cmdTokenFile); err != nil {
				fmt.Printf("Token was not refreshed: %s", err)
				fmt.Println()
				return 3
			} else {
				fmt.Println("Token has been successfully refreshed")
				return 0
			}
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
	var asWin32 = false

	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	exportCmd.StringVar(&cmdTokenFile, customTokenFileOpt, v3client.DefaultSavedAccessTokenFilePath(), "Use specified custom file")
	exportCmd.BoolVar(&asWin32, windowsOpt, false, "Export settings for batch file")

	_ = exportCmd.Parse(os.Args[2:])

	exportVar := "V3_ACCESS_TOKEN"
	if len(exportCmd.Args()) > 0 {
		exportVar = exportCmd.Args()[0]
	}

	if tkn, err := v3client.ReadSavedV3TokenData(cmdTokenFile); err == nil && tkn != nil {
		if asWin32 {
			fmt.Println("@echo off")
			fmt.Printf("SET %s=%s", exportVar, tkn.AccessToken)
		} else {
			fmt.Printf("export %s='%s'", exportVar, tkn.AccessToken)
		}
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
		f = v3client.DefaultSavedAccessTokenFilePath()
	}

	if tkn, err := v3client.ReadSavedV3TokenData(f); err == nil && tkn != nil {
		if tkn.Expired() {
			fmt.Printf("Your access token has already expired (on %s)", tkn.ExpiryTime())
		} else {
			minutesLeft := tkn.TimeLeft() / 60
			if minutesLeft <= 1 {
				fmt.Printf("Your token has 1 minute or less validatity time. Refresh before next operation")
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

func keepTokenAlive() int {
	if err := tokenCommandLine("keep-alive").Parse(os.Args[2:]); err == nil {
		runtimeCredentials := v3client.DeriveAccessCredentials(credentialsFile, credentialsPassword(), nil)
		if !runtimeCredentials.FullySpecified() {
			fmt.Println("Insufficient input credentials")
			return 1
		}

		provider := v3client.NewClientCredentialsProvider(runtimeCredentials, nil)

		if initTkn, err := provider.TokenData(); err != nil || initTkn == nil {
			if err != nil {
				fmt.Printf("Could not retrive initial access token: %s", err)
			} else {
				fmt.Println("Nil token and no error were returned.")
			}
			fmt.Println()
			return 2
		} else {
			if err = v3client.PersistV3TokenResponse(initTkn, cmdTokenFile); err != nil {
				fmt.Printf("Could not save initial token: %s", err)
				fmt.Println()
				return 3
			}

			exitChan := make(chan int)

			provider.OnPostRefresh(func() {
				if provider.Response == nil || provider.Response.Expired() {
					fmt.Printf("Refresh token failed, and current token has expired: %s", err)
					fmt.Println()
					exitChan <- 1
				} else if err := v3client.PersistV3TokenResponse(provider.Response, cmdTokenFile); err != nil {
					fmt.Printf("Could not save updated refresh token: %s", err)
					fmt.Println()
					exitChan <- 1
				}
			})

			provider.EnsureRefresh()
			fmt.Println("Initial token retrieved; will keep token alive until interrupted.")
			fmt.Println("------------------------------------------------------------------")

			// In case refresh has failed, let's make a provision to
			exitCode := <-exitChan

			fmt.Println("--------------------------------------------------")
			fmt.Println("Keep-alive stopped; perhaps due to the error above")

			return exitCode
		}
	} else {
		fmt.Printf("Could not parse command line arguments: %s", err)
		return 1
	}
}

func tokenCommandLine(name string) *flag.FlagSet {
	initCmd := flag.NewFlagSet(name, flag.ExitOnError)

	initCmd.StringVar(&credentialsFile, "credentials", v3client.DefaultCredentialsFile(), "Path to the credentials file")
	initCmd.StringVar(&cmdTokenFile, customTokenFileOpt, v3client.DefaultSavedAccessTokenFilePath(), "Use this file to read/write access token data")

	return initCmd
}

// credentialsPassword derive the credentials' password.
func credentialsPassword() string {
	if env := os.Getenv("MASH_CREDS_AES"); len(env) > 0 {
		return env
	}

	// If the password is being pipes via the standard input, the standard input
	// will be preferred over human input.
	fi, err := os.Stdin.Stat()

	if err == nil && fi.Mode()&os.ModeNamedPipe != 0 {
		rawPwd := make([]byte, 128)
		if n, err := os.Stdin.Read(rawPwd); err == nil && n > 0 {
			return string(rawPwd[:n])
		}
	}

	fmt.Println("No password found in standard input")
	fmt.Printf("Please enter credentials decryption password: ")
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')

	return line
}

func encryptCredentials() int {
	initCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	initCmd.StringVar(&credentialsFile, "credentials", v3client.DefaultCredentialsFile(), "Path to the credentials file")

	if err := initCmd.Parse(os.Args[2:]); err != nil {
		initCmd.PrintDefaults()
		return 1
	}

	if err := v3client.EncryptInPlace(credentialsFile, credentialsPassword()); err != nil {
		fmt.Printf("Could not encrypt file %s: %s", credentialsFile, err)
		return 1
	} else {
		return 0
	}
}

func initToken() int {
	if err := tokenCommandLine("init").Parse(os.Args[2:]); err == nil {

		runtimeCredentials := v3client.DeriveAccessCredentials(credentialsFile, credentialsPassword(), nil)
		if !runtimeCredentials.FullySpecified() {
			fmt.Println("Insufficient input credentials")
			return 1
		}

		provider := v3client.NewOAuthHelper(v3client.OAuthHelperParams{})

		if dat, err := provider.RetrieveAccessTokenFor(&runtimeCredentials); err == nil {
			if err = v3client.PersistV3TokenResponse(dat, cmdTokenFile); err != nil {
				fmt.Printf("Failed to save response: %s", err)
				return 2
			} else {
				return 0
			}
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
	if len(os.Args) == 0 {
		fmt.Println("expected at least a subcommand")
		os.Exit(1)
	}

	var exitCode int

	switch os.Args[1] {

	case "encrypt":
		exitCode = encryptCredentials()
		break
	case "init":
		exitCode = initToken()
		break
	case "keep-alive":
		exitCode = keepTokenAlive()
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
