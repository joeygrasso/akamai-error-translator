package main

import (
	"akamai-error-translator/pkg/akamai"
	"flag"
	"fmt"
	"os"
)

func main() {

	// Get Credentials from ENV Vars
	akmAccessToken := os.Getenv("AKM_ACCESS_TOKEN")
	akmClientSecret := os.Getenv("AKM_CLIENT_SECRET")
	akmClientToken := os.Getenv("AKM_CLIENT_TOKEN")
	akmHostname := os.Getenv("AKM_HOSTNAME")

	// Set Credentials
	akmCreds := akamai.AkamaiCredentials{
		AccessToken:  akmAccessToken,
		ClientSecret: akmClientSecret,
		ClientToken:  akmClientToken,
		Hostname:     akmHostname,
	}

	// Get Args
	var errorString = flag.String("error", "", "Akamai Error Code Example: 00.51680117.1679326682.4ab87915")
	var traceLoggingBool = flag.Bool("tracing", false, "Enable trace logging for additional data true/false")

	flag.Parse()

	// Call error translate func
	akmResponse, err := akamai.TranslateErrorString(*errorString, *traceLoggingBool, akmCreds)

	if err != nil {
		fmt.Println(err)
	}

	// Simply print the raw JSON response for the user
	fmt.Println(string(akmResponse))
}
