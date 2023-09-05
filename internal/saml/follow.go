package saml

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func (s *SAML) checkNext(resp *http.Response) (*http.Response, []byte, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	re := regexp.MustCompile(`window\.location\.href='(https://[^']+)';`)

	match := re.FindSubmatch(data)
	if len(match) > 1 {
		resp, err := s.cli.Get(string(match[1]))
		if err != nil {
			return nil, nil, err
		}
		return s.checkNext(resp)
	}
	return resp, data, nil

}

func follow(cli *http.Client, url string) (*http.Response, error) {
	resp, err := cli.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	fmt.Println(resp.Request.URL)
	fmt.Println(resp.Header.Get("location"))
	fmt.Println(resp.Status)
	if resp.StatusCode == http.StatusFound {
		location := resp.Header.Get("Location")
		if location == "" {
			return nil, errors.New("no location included in redirect")
		}
		return follow(cli, location)
	}
	return resp, nil
}
