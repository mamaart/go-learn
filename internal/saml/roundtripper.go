package saml

import (
	"fmt"
	"net/http"
)

type LoggingRoundTripper struct {
	Transport http.RoundTripper
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Perform the HTTP request using the original transport
	fmt.Println(req.Header["User-Agent"])
	req.Header.Add("User-Agent", "James Bond")
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Log the cookies in the response
	cookies := resp.Cookies()
	fmt.Printf("Cookies in Response for %s:\n", req.URL)
	for _, cookie := range cookies {
		fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
	}

	// You can perform additional logging or processing here if needed

	return resp, nil
}
