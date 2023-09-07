package saml

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

type SAML struct {
	cli *http.Client
}

func New(cli *http.Client) *SAML {
	return &SAML{cli}
}

func (s *SAML) GetLogin(initUrl, authUrl string) (func(username, password string) (*http.Response, error), error) {
	u, err := s.getSAMLRequest(initUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get SAML request: %s", err)
	}
	return func(username, password string) (*http.Response, error) {
		samlResponse, err := s.getSAMLResponse(u, username, password)
		if err != nil {
			return nil, fmt.Errorf("failed to get SAML response: %s", err)
		}

		r, err := s.cli.PostForm(authUrl, url.Values{
			"SAMLResponse": []string{samlResponse},
		})
		if err != nil {
			return nil, fmt.Errorf("failed at login request: %s", err)
		}

		return r, nil
	}, nil
}

func (s *SAML) getSAMLRequest(u string) (string, error) {
	resp, err := s.cli.Get(u)
	if err != nil {
		return "", fmt.Errorf("failed at first request: %s", err)
	}

	saml := resp.Request.URL.Query().Get("SAMLRequest")
	if saml == "" {
		fmt.Println(resp.Request.URL.String())
		//fmt.Println(string(functools.MustV(io.ReadAll(resp.Body))))
		// Reason might be that we are already authorized
		// TODO figure out: what to do?
		return "", fmt.Errorf("no saml request returned")
	}

	resp, err = s.cli.PostForm(resp.Request.URL.String(), url.Values{
		"HomeRealmSelection": []string{"AD+AUTHORITY"},
		"Email":              []string{""},
	})
	if err != nil {
		return "", fmt.Errorf("failed at second request: %s", err)
	}
	return resp.Request.URL.String(), nil
}

func (s *SAML) getSAMLResponse(u, username, password string) (string, error) {
	resp, err := s.cli.PostForm(u, url.Values{
		"Username":   []string{fmt.Sprintf("win\\%s", username)},
		"Password":   []string{password},
		"AuthMethod": []string{"FormsAuthenticate"},
	})
	if err != nil {
		return "", fmt.Errorf("failed at third request: %s", err)
	}
	x, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read data: %s", err)
	}
	match := regexp.MustCompile(`name="SAMLResponse" value="([^"]+)"`).FindSubmatch(x)
	if len(match) != 2 {
		return "", errors.New("didn't find SAMLResponse")
	}
	return string(match[1]), nil
}
