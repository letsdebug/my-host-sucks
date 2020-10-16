package cpanel

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// NewRemoteAPI provides an API client for a remote cPanel environment
func NewRemoteAPI(cpURL, username, password string, cl *http.Client) (API, error) {
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
		cl:       cl,
	}, nil
}

type remoteCpanel struct {
	URL      *url.URL
	Username string
	Password string

	cl *http.Client
}

func (c *remoteCpanel) UAPI(module, function string, args Args, out interface{}) error {
	return c.api("uapi", module, function, args, out)
}

func (c *remoteCpanel) API2(module, function string, args Args, out interface{}) error {
	return c.api("api2", module, function, args, out)
}

func (c *remoteCpanel) api(apiVersion, module, function string, args Args, out interface{}) error {
	reqArgs := url.Values{}
	for k, v := range args {
		reqArgs.Add(k, fmt.Sprintf("%v", v))
	}

	reqURL := *c.URL

	switch apiVersion {
	case "uapi":
		reqURL.Path = fmt.Sprintf("/execute/%s/%s", module, function)
		reqURL.RawQuery = reqArgs.Encode()
	case "api2":
		reqArgs.Add("cpanel_jsonapi_user", c.Username)
		reqArgs.Add("cpanel_jsonapi_apiversion", "2")
		reqArgs.Add("cpanel_jsonapi_module", module)
		reqArgs.Add("cpanel_jsonapi_func", function)
		reqURL.Path = "/json-api/cpanel"
		reqURL.RawQuery = reqArgs.Encode()
	default:
		return fmt.Errorf("unsupported API version: %s", apiVersion)
	}

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("user-agent", "my-host-sucks/dev (https://github.com/letsdebug/my-host-sucks)")

	if c.cl == nil {
		c.cl = http.DefaultClient
	}

	resp, err := c.cl.Do(req)
	if err != nil {
		return fmt.Errorf("API request %s:%s failed: %w", module, function, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request %s:%s failed: HTTP %s",
			module, function, resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(&out)
}
