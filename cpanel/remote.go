package cpanel

import (
	"errors"
	"fmt"
	"net/url"
)

// NewRemoteAPI provides an API client for a remote cPanel environment
func NewRemoteAPI(cpURL, username, password string, insecure bool) (API, error) {
	if cpURL == "" || username == "" || password == "" {
		return nil, errors.New("not all required details (URL, username, password) were provided")
	}
	u, err := url.Parse(cpURL)
	if err != nil || u.Hostname() == "" {
		return nil, fmt.Errorf("'%s' is not a URL", cpURL)
	}

	return &remoteCpanel{
		URL:      u,
		Username: username,
		Password: password,
		Insecure: insecure,
	}, nil
}

type remoteCpanel struct {
	URL      *url.URL
	Username string
	Password string
	Insecure bool
}

func (c *remoteCpanel) UAPI(module, function string, args map[string]string, out interface{}) error {
	return errors.New("NYI")
}

func (c *remoteCpanel) API2(module, function string, args map[string]string, out interface{}) error {
	return errors.New("NYI")
}
