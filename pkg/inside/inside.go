package inside

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/user"
	"regexp"

	"github.com/mamaart/go-learn/internal/auth"
	"github.com/mamaart/go-learn/internal/inside"
	"github.com/mamaart/go-learn/internal/inside/models"
	"github.com/mamaart/go-learn/internal/inside/parsing"
	"github.com/mamaart/go-learn/pkg/functools"
)

type Inside struct {
	manager *auth.AuthManager
}

type Options struct {
	Credentials *auth.Credentials
}

func DefaultOptions() Options {
	return Options{nil}
}

func New(opts Options) (*Inside, error) {
	home := functools.MustV(user.Current()).HomeDir // TODO maybe get this from options?
	manager, err := auth.New(auth.Options{
		Credentials:     opts.Credentials,
		Login:           inside.LoginToInside,
		CheckAuthorized: inside.CheckAuthState,
		CookiePath:      fmt.Sprintf("%s/.dtu_inside_cookies", home),
		CredentialsPath: fmt.Sprintf("%s/.dtu_credentials", home),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create auth manager: %s", err)
	}
	return &Inside{manager}, nil
}

func (i *Inside) Get(url string) (*http.Response, error) {
	return i.manager.WithClient(func(c *http.Client) (*http.Response, error) {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15")
		resp, err := c.Do(req)
		if err != nil {
			return nil, err
		}

		// This check only happens if some 304 stuff has happened
		// it is probably not gonna solve every case of asking for
		// a redirect but the alternative would be to read the data
		// and check for either SAMLResponse payload or something
		// wiht js location href modification.
		if resp.Request.URL.String() != url {
			fmt.Println("refreshing")
			_, err := refreshSaml(c, resp)
			if err != nil {
				return nil, err
			}
			return i.Get(url) //Retry
		}
		return resp, nil
	})

}

func (i *Inside) Post(url string, data []byte) (*http.Response, error) {
	return i.manager.WithClient(func(c *http.Client) (*http.Response, error) {
		req, _ := http.NewRequest("POST", url, bytes.NewReader(data))
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/116.0")
		return c.Do(req)
	})
}

func (i *Inside) Whoami() (string, error) {
	resp, err := i.Get("https://www.inside.dtu.dk/dtuapi/dtucontext/getuserorcurrent?cwis=undefined")
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (i *Inside) Grades() ([]models.Grade, error) {
	resp, err := i.Get("https://cn.inside.dtu.dk/cnnet/Grades/Grades.aspx")
	if err != nil {
		return nil, fmt.Errorf("data parse failed: %s", err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parsing.ParseGradesHtml(data)
}

// TODO Check for division by zero, and validate that all grades are valid
// meaning they have both valid ects and grade, maybe ignore the ones that
// do not have valid info.
func (i *Inside) GPA() (int, error) {
	grades, err := i.Grades()
	if err != nil {
		return 0, fmt.Errorf("failed to get grades: %s", err)
	}
	fold := func(fn func(models.Grade) int) int {
		return functools.Foldr(grades, 0, func(g models.Grade, acc int) int {
			if g.Grade <= 0 {
				return 0 + acc
			}
			return fn(g) + acc
		})
	}
	gs := fold(func(g models.Grade) int { return g.Grade * g.Ects })
	cs := fold(func(g models.Grade) int { return g.Ects })

	if cs == 0 {
		return 0, fmt.Errorf("can't calculate because completed ectc is 0")
	}

	return gs / cs, nil
}

// TODO check how to put this inside saml package.
func refreshSaml(cli *http.Client, resp *http.Response) ([]byte, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	match := regexp.MustCompile(`name="SAMLResponse" value="([^"]+)"`).FindSubmatch(data)
	if len(match) != 2 {
		fmt.Println("SAMLResponse NOT FOUNT!")
		resp, data, err = follow(cli, resp)
		if err != nil {
			return nil, fmt.Errorf("follow response failed: %s", err)
		}
		return data, nil //errors.New("didn't find SAMLResponse")
	}
	fmt.Println(string(match[1]))
	// URL used here is from SAML authurl
	resp, err = cli.PostForm("https://auth.dtu.dk/dtu/", url.Values{
		"SAMLResponse": []string{string(match[1])},
	})
	if err != nil {
		return nil, fmt.Errorf("failed refreshing with saml: %s", err)
	}
	resp, data, err = follow(cli, resp)
	if err != nil {
		return nil, fmt.Errorf("follow response failed: %s", err)
	}
	return data, nil
}

// This function is just copied from inside.authorizer in this case it is used
// to redirect the client to previously requested page after a saml flow.
func follow(cli *http.Client, resp *http.Response) (*http.Response, []byte, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	re := regexp.MustCompile(`window\.location\.href='(https://[^']+)';`)
	match := re.FindSubmatch(data)
	if len(match) > 1 {
		fmt.Println(string(match[1]))
		resp, err := cli.Get(string(match[1]))
		if err != nil {
			return nil, nil, err
		}
		return resp, data, nil
	}
	return resp, data, nil
}
