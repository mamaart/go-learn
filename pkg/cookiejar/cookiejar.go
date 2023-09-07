package cookiejar

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sync"

	"github.com/mamaart/go-learn/pkg/functools"
)

type CookieJar struct {
	sync.RWMutex
	jar        *cookiejar.Jar
	allCookies map[string][]*http.Cookie
	path       string
}

func NewCookieJar(path string) *CookieJar {
	realJar, _ := cookiejar.New(nil)

	from := functools.MustV(load(path))
	for k, v := range from {
		url := functools.MustV(url.Parse(k))
		realJar.SetCookies(url, v)
	}

	return &CookieJar{
		jar:        realJar,
		allCookies: from,
		path:       path,
	}
}

func (jar *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.Lock()
	defer jar.Unlock()
	jar.allCookies[u.String()] = cookies
	jar.jar.SetCookies(u, cookies)
}

func (jar *CookieJar) Cookies(u *url.URL) (out []*http.Cookie) {
	return jar.jar.Cookies(u)
}

func (jar *CookieJar) Save() error {
	jar.RLock()
	defer jar.RUnlock()
	bs, err := json.MarshalIndent(jar.allCookies, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(jar.path, bs, 0644)
}

func load(path string) (map[string][]*http.Cookie, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	out := make(map[string][]*http.Cookie)
	if err := json.NewDecoder(file).Decode(&out); err != nil {
		if err == io.EOF {
			return out, nil
		}
		return nil, err
	}
	return out, nil
}
