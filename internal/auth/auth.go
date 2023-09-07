package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"

	"github.com/mamaart/go-learn/pkg/cookiejar"
	"github.com/mamaart/go-learn/pkg/functools"
)

type AuthManager struct {
	cookiepath, credentialspath string
	jar                         *cookiejar.CookieJar
	check                       func() bool
	login                       func(Credentials) error
}

type Options struct {
	Credentials     *Credentials
	CookiePath      string
	CredentialsPath string
	CheckAuthorized func(cli *http.Client) (*http.Response, error)
	Login           func(cli *http.Client, username, password string) (*http.Response, error)
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (opts Options) validate() error {
	if opts.CookiePath == "" {
		u, err := user.Current()
		if err != nil {
			return err
		}
		opts.CookiePath = fmt.Sprintf("%s/.dtu_cookies.json", u.HomeDir)
	}
	if opts.Login == nil {
		return fmt.Errorf("no login function provided")
	}
	if opts.CheckAuthorized == nil {
		return fmt.Errorf("no check function provided")
	}
	return nil
}

func New(opts Options) (*AuthManager, error) {
	if err := opts.validate(); err != nil {
		return nil, fmt.Errorf("failed to successfully validate options: %s", err)
	}
	jar := cookiejar.NewCookieJar(opts.CookiePath) //TODO what if cookiefile is not writable?
	cli := &http.Client{Jar: jar}
	return AuthManager{
		jar:             jar,
		credentialspath: opts.CredentialsPath,
		cookiepath:      opts.CookiePath,
		check: func() bool {
			_, err := opts.CheckAuthorized(cli)
			return err == nil
		},
		login: func(creds Credentials) error {
			if _, err := opts.Login(cli, creds.Username, creds.Password); err != nil {
				return err
			}
			if err := jar.Save(); err != nil {
				return fmt.Errorf("failed to persist cookies: %s", err)
			}
			return nil
		},
	}.init(opts.Credentials)
}

func (m AuthManager) init(creds *Credentials) (*AuthManager, error) {
	if creds != nil {
		if err := m.Login(*creds); err != nil {
			return nil, fmt.Errorf("failed to login: %s", err)
		}
	}
	if m.check() {
		return &m, nil
	}
	return &m, m.AutoLogin()
}

func (m *AuthManager) WithClient(fn func(*http.Client) (*http.Response, error)) (*http.Response, error) {
	resp, err := fn(&http.Client{Jar: m.jar})
	if err != nil {
		return nil, err
	}
	m.jar.Save()
	return resp, nil
}

func (m *AuthManager) Authorized() bool {
	return m.check()
}

func (m *AuthManager) Login(creds Credentials) error {
	if err := m.login(creds); err != nil {
		return err
	}
	if m.credentialspath != "" {
		data := functools.MustV(json.Marshal(creds))
		if err := os.WriteFile(m.credentialspath, data, 0644); err != nil {
			return fmt.Errorf("failed to persist credentials: %s", err)
		}
	}
	return nil
}

func (m *AuthManager) AutoLogin() error {
	if m.credentialspath == "" {
		return fmt.Errorf("no credentials path defined")
	}

	data, err := os.ReadFile(m.credentialspath)
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %s", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return fmt.Errorf("failed to parse credentials from file: %s", err)
	}
	return m.login(creds)
}
