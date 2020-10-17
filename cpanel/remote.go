package cpanel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// NewRemoteClient provides an API client for a remote cPanel environment
func NewRemoteClient(cpURL, username, password string, cl *http.Client) (*Client, error) {
	if cpURL == "" || username == "" || password == "" {
		return nil, errors.New("not all required details (URL, username, password) were provided")
	}
	u, err := url.Parse(cpURL)
	if err != nil || u.Hostname() == "" {
		return nil, fmt.Errorf("'%s' is not a URL", cpURL)
	}

	return &Client{a: &remoteCpanel{
		URL:      u,
		Username: username,
		Password: password,
		cl:       cl,
	}}, nil
}

type remoteCpanel struct {
	URL      *url.URL
	Username string
	Password string

	cl *http.Client
}

func (c *remoteCpanel) UAPI(ctx context.Context, module, function string, args Args, out interface{}) error {
	return c.api(ctx, "uapi", module, function, args, out)
}

func (c *remoteCpanel) API2(ctx context.Context, module, function string, args Args, out interface{}) error {
	// API2 responses need to be unwrapped from the outer cpanelresult hash
	// https://documentation.cpanel.net/display/DD/cPanel+API+2+-+Return+Data
	var resp struct {
		Result json.RawMessage `json:"cpanelresult"`
	}
	if err := c.api(ctx, "api2", module, function, args, &resp); err != nil {
		return err
	}
	return prettyJSONError(resp.Result, json.Unmarshal(resp.Result, out))
}

func (c *remoteCpanel) api(ctx context.Context, apiVersion, module, function string, args Args, out interface{}) error {
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

	req = req.WithContext(ctx)

	if c.cl == nil {
		c.cl = http.DefaultClient
	}

	if ctx.Value(LogRequestsAndResponses) != nil {
		log.Printf("%s:%s:%s request: %s", apiVersion, module, function, reqURL.String())
	}

	resp, err := c.cl.Do(req)
	if err != nil {
		return fmt.Errorf("API request %s:%s failed: %w", module, function, err)
	}
	defer resp.Body.Close()

	// Buffer the full response. This costs more memory but we want to to report the contents if
	// it's not valid JSON.
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read full API response: %w", err)
	}

	if ctx.Value(LogRequestsAndResponses) != nil {
		log.Printf("%s:%s:%s response: %q", apiVersion, module, function, buf)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request %s:%s failed: HTTP %s",
			module, function, resp.Status)
	}

	return prettyJSONError(buf, json.Unmarshal(buf, out))
}
