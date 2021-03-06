package go_mauth_client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

/*
Much of this was heavily informed by:
https://medium.com/@matryer/the-http-handlerfunc-wrapper-technique-in-golang-c60bf76e6124#.9yl4dj1gd
and
https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81#.xj15k9f5k
*/

//go:generate go run gen.go

// Get the Version for this Client
func GetVersion() string {
	return VersionString
}

// isJSON tries to work out if the content is JSON, so it can add the correct Content-Type to the Headers
// taken from http://stackoverflow.com/a/22129435/1638744
func isJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}

// makeRequest formulates the message, including the MAuth Headers and returns a http.Request, ready to send
func (mauthApp *MAuthApp) makeRequest(method string, rawurl string, body string) (req *http.Request, err error) {
	// Use the url.URL to assist with path management
	url2, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	// this needs to persist
	secondsSinceEpoch := time.Now().Unix()
	// build the MWS string
	stringToSign := MakeSignatureString(mauthApp, method, url2.Path, body, secondsSinceEpoch)
	// Sign the string
	signedString, err := SignString(mauthApp, stringToSign)
	if err != nil {
		return nil, err
	}
	// create a new request object
	req, err = http.NewRequest(method, rawurl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	// take everything and build the structure of the MAuth Headers
	made_headers := MakeAuthenticationHeaders(mauthApp, signedString, secondsSinceEpoch)
	for header, value := range made_headers {
		req.Header.Set(header, value)
	}
	// Detect JSON, send appropriate Content-Type if detected
	if isJSON(body) {
		req.Header.Set("Content-Type", "application/json")
	}
	// Add the User-Agent using the Client Version
	req.Header.Set("User-Agent", "go-mauth-client/"+GetVersion())
	return req, nil
}
