package akamai

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/comcast/go-edgegrid/edgegrid"
)

// Reference: "https://developer.akamai.com/api/luna/diagnostic-tools/overview.html"
// Akamai EdgeGrid Auth: https://github.com/Comcast/go-edgegrid

type AkamaiCredentials struct {
	AccessToken  string
	ClientSecret string
	ClientToken  string
	Hostname     string
}

type AkamaiErrorTranslateRequest struct {
	ErrorCode        string
	TraceForwardLogs bool
}

type AkamaiErrorTranslateResponse struct {
	CreatedBy       string
	CreatedTime     string
	ExecutionStatus string
	Link            string
	Request         AkamaiErrorTranslateRequest
	RequestID       string
	RetryAfter      int
}

func TranslateErrorString(akamaiErrorString string, traceEnabled bool, akmCreds AkamaiCredentials) (string, error) {

	// Nasty but necessary for converting the bool type to a string for the stirngs.NewReader()
	convertedTraceEnabled := _convertBool(traceEnabled)

	akmRequestPayload := strings.NewReader("{\"traceForwardLogs\":" + convertedTraceEnabled + ",\"errorCode\":\"" + akamaiErrorString + "\"}")

	// This triggers the async generation of translating the error code. Akamai will return a payload body with an ID, status, and retryAfer value
	// We need to sleep for the retryAfter value and then do a GET request for the requestID to actually get the error translation
	akamaiErrorTranslateRequest, err := makeAkamaiErrorTranslateRequest(akmCreds, akmRequestPayload)

	// Print the output from the async request just for reference
	fmt.Println(string(akamaiErrorTranslateRequest))

	if err != nil {
		return "", err
	}

	// TODO: Break this out into a function to make testing easier
	var akmErrTransRes AkamaiErrorTranslateResponse
	json.Unmarshal(akamaiErrorTranslateRequest, &akmErrTransRes)

	// Sleep for the retryAfter value in the akamaiErrorTranslateRequest response
	// TODO: Add fancy loading animation
	fmt.Printf("Sleeping for %d minutes while the request processes....\n", akmErrTransRes.RetryAfter)
	_sleep(akmErrTransRes.RetryAfter)

	// After sleeping, attempt to get the error translation
	akmResponse, err := getAkamaiErrorTranslation(akmCreds, akmErrTransRes.Link)

	if err != nil {
		return "", err
	}

	return string(akmResponse), nil
}

func makeAkamaiErrorTranslateRequest(akmCreds AkamaiCredentials, payload io.Reader) ([]byte, error) {
	url := "https://" + akmCreds.Hostname + "/edge-diagnostics/v1/error-translator"

	req, _ := http.NewRequest("POST", url, payload)

	// Akamai EdgeGridAuth https://github.com/Comcast/go-edgegrid#alternative-usage
	params := edgegrid.NewAuthParams(req, akmCreds.AccessToken, akmCreds.ClientToken, akmCreds.ClientSecret)
	auth := edgegrid.Auth(params)

	req.Header.Add("Authorization", auth)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

// TODO: Lots of repeat code, optimize with helper function to reduce repeat code from makeAkamaiErrorTranslateRequest() func above
func getAkamaiErrorTranslation(akmCreds AkamaiCredentials, link string) ([]byte, error) {
	url := "https://" + akmCreds.Hostname + link

	req, _ := http.NewRequest("GET", url, nil)

	// Akamai EdgeGridAuth https://github.com/Comcast/go-edgegrid#alternative-usage
	params := edgegrid.NewAuthParams(req, akmCreds.AccessToken, akmCreds.ClientToken, akmCreds.ClientSecret)
	auth := edgegrid.Auth(params)

	req.Header.Add("Authorization", auth)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func _convertBool(b bool) string {
	return strconv.FormatBool(b)
}

func _sleep(sleepDurationMinutes int) {
	time.Sleep(time.Duration(sleepDurationMinutes) * time.Minute)
}
