package auth

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/mamaart/go-learn/internal/saml"
)

func LoginToLearn(username, password string) (*http.Client, error) {
	login, err := saml.New().GetLogin(
		"https://learn.inside.dtu.dk/",
		"https://learn.inside.dtu.dk/d2l/lp/auth/login/samlLogin.d2l",
	)
	if err != nil {
		return nil, err
	}
	return login(username, password)
}

func LoginToInside(username, password string) (*http.Client, error) {
	login, err := saml.New().GetLogin(
		"https://cn.inside.dtu.dk/",
		"https://auth.dtu.dk/dtu/",
	)
	if err != nil {
		return nil, fmt.Errorf("saml flow failed: %s", err)
	}
	cli, err := login(username, password)
	if err != nil {
		return nil, fmt.Errorf("login failed: %s", err)
	}

	resp, err := cli.Get("https://cn.inside.dtu.dk/cnnet/element/682907")
	if err != nil {
		return nil, fmt.Errorf("failed request url that requires session cookie: %s", err)
	}
	_, _, err = follow(cli, resp)
	if err != nil {
		return nil, fmt.Errorf("following failed while trying to get session cookie: %s", err)
	}
	return cli, nil
}

// follow is used to get a ticket which is needed to get access to inside
// the ticket flow happens after the user is logged in
func follow(cli *http.Client, resp *http.Response) (*http.Response, []byte, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	re := regexp.MustCompile(`window\.location\.href='(https://[^']+)';`)

	match := re.FindSubmatch(data)
	if len(match) > 1 {
		resp, err := cli.Get(string(match[1]))
		if err != nil {
			return nil, nil, err
		}
		return follow(cli, resp)
	}
	return resp, data, nil
}
