package auth

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

func Login(username, password string) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	cli := &http.Client{
		Jar: jar,
	}

	return cli, login(cli, username, password)
}

func login(cli *http.Client, username, password string) error {
	u, err := getSAMLRequest(cli)
	if err != nil {
		return fmt.Errorf("failed to get SAML request: %s", err)
	}

	samlResponse, err := getSAMLResponse(cli, u, username, password)
	if err != nil {
		return fmt.Errorf("failed to get SAML response: %s", err)
	}

	authUrl := "https://learn.inside.dtu.dk/d2l/lp/auth/login/samlLogin.d2l"
	_, err = cli.PostForm(authUrl, url.Values{
		"SAMLResponse": []string{samlResponse},
	})
	if err != nil {
		return fmt.Errorf("failed at login request: %s", err)
	}
	return nil
}

func getSAMLRequest(cli *http.Client) (string, error) {
	resp, err := cli.Get("https://learn.inside.dtu.dk/")
	if err != nil {
		return "", fmt.Errorf("failed at first request: %s", err)
	}
	resp, err = cli.PostForm(resp.Request.URL.String(), url.Values{
		"HomeRealmSelection": []string{"AD+AUTHORITY"},
		"Email":              []string{""},
	})
	if err != nil {
		return "", fmt.Errorf("failed at second request: %s", err)
	}
	return resp.Request.URL.String(), nil
}

func getSAMLResponse(cli *http.Client, u, username, password string) (string, error) {
	resp, err := cli.PostForm(u, url.Values{
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
