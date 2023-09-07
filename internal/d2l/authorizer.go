package d2l

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mamaart/go-learn/internal/saml"
	"github.com/mamaart/go-learn/pkg/functools"
)

func CheckAuthState(cli *http.Client) (*http.Response, error) {
	u := "https://learn.inside.dtu.dk/d2l/api/lp/1.2/users/whoami"
	resp, err := cli.Get(u)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		fmt.Println(string(functools.MustV(io.ReadAll(resp.Body))))
		return resp, errors.New("unauthorized")
	}
	return resp, nil
}

func LoginToLearn(cli *http.Client, username, password string) (*http.Response, error) {
	resp, err := CheckAuthState(cli)
	if err != nil {
		login, err := saml.New(cli).GetLogin(
			"https://learn.inside.dtu.dk/",
			"https://learn.inside.dtu.dk/d2l/lp/auth/login/samlLogin.d2l",
		)
		if err != nil {
			return nil, err
		}
		return login(username, password)
	}
	return resp, errors.New("already logged in")
}
