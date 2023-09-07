package inside

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/mamaart/go-learn/internal/saml"
	"github.com/mamaart/go-learn/pkg/functools"
)

func CheckAuthState(cli *http.Client) (*http.Response, error) {
	req, _ := http.NewRequest("GET", "https://www.inside.dtu.dk/dtuapi/DtuContext/GetUserOrCurrent?cwis=undefined", nil)
	return functools.Retry(5, func() (*http.Response, bool) {
		resp, err := cli.Do(req)
		if err != nil {
			return nil, false
		}
		if resp.StatusCode != http.StatusOK {
			//Getting ASPXAUTH cookie used for whoami
			if err := getToken(cli, "https://www.inside.dtu.dk/da/undervisning"); err != nil {
				return nil, false
			}
			return nil, false
		}
		return resp, true
	}, errors.New("unauthorized"))
}

func LoginToInside(cli *http.Client, username, password string) (*http.Response, error) {
	//if resp, err := CheckAuthState(cli); err == nil {
	//	return resp, nil
	//}
	login, err := saml.New(cli).GetLogin(
		"https://cn.inside.dtu.dk/",
		"https://auth.dtu.dk/dtu/",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to do saml: %s", err)
	}
	resp, err := login(username, password)
	if err != nil {
		return nil, fmt.Errorf("login failed: %s", err)
	}

	//TODO analyze where the specific cookies arrive
	resp, _, err = follow(cli, resp)
	if err != nil {
		return nil, fmt.Errorf("follow response failed: %s", err)
	}

	//Getting  SessionID_v1 token used to get grades
	if err := getToken(cli, "https://cn.inside.dtu.dk/cnnet/element/682907"); err != nil {
		return resp, err
	}

	//Getting ASPXAUTH cookie used for whoami
	if err := getToken(cli, "https://www.inside.dtu.dk/da/undervisning"); err != nil {
		return resp, err
	}

	return resp, nil
}

func getToken(cli *http.Client, url string) error {
	resp, err := cli.Get(url)
	if err != nil {
		return fmt.Errorf("failed request url that requires session cookie: %s", err)
	}
	_, _, err = follow(cli, resp)
	if err != nil {
		return fmt.Errorf("following failed while trying to get session cookie: %s", err)
	}
	return nil
}

// follow is used to get a ticket which is needed to get access to inside
// the ticket flow happens after the user is logged in
func follow(cli *http.Client, resp *http.Response) (*http.Response, []byte, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(resp.Status, resp.Request.URL.String())
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
