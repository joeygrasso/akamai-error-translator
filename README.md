# Akamai Error Translator

This is a quick little tool that took a couple hours to write that takes an Akamai Reference ID and translates it for debugging and troubleshooting purposes. It accepts an Akamai Ref ID as input and an optional flag to enable more verbose trace logging. See usage for details...

**Your mileage will vary... there's little to no error handling and handling of edge cases. This was a tool that was created in a couple hours because I didn't want to click through the Akamai UI.**

## Requirements

To use this tool, you need an Akamai API key with proper permissions granted. These can be generated through the Akamai control panel.

[Akamai API Credentials Documentation](https://techdocs.akamai.com/developer/docs/set-up-authentication-credentials)

Also you'll need a properly setup golang environment.

## Usage

You will need to set the following env vars once you have generated an API key following the linked doc above.
```
export AKM_ACCESS_TOKEN="abcdefg1337"
export AKM_CLIENT_SECRET="$3c4etsRfun"
export AKM_CLIENT_TOKEN="jrrToken"
export AKM_HOSTNAME="uniqueURL.luna.akamaiapis.net"
```

### Quick Example:
Run without compiling:
```
go run cmd/akamai-error-translator/akamai-error-translator.go --error="00.51680117.1679419172.5545c6b5"
```

Run without compiling w/ additional tracing:
```
go run cmd/akamai-error-translator/akamai-error-translator.go --error="00.51680117.1679419172.5545c6b5" --tracing=true
```

Both of the above examples will dump out a JSON response. You can use IO redirection to put it into a file or you can pipe to JQ for a human readable format `| jq .`

## Contributing

Pull requests are always welcome.

## License

Use this freely for your own purposes.