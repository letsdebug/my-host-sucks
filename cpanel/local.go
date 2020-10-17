package cpanel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
)

const (
	uapiPath = "/bin/uapi"
	api2Path = "/bin/cpapi2"
)

// IsLocal determines whether the local environment is a non-root cPanel one
func IsLocal() bool {
	s, _ := os.Stat(uapiPath)
	return s != nil && os.Getuid() != 0
}

// NewLocalAPI provides an API client for the local environment
func NewLocalAPI() API {
	return &localCpanel{}
}

type localCpanel struct {
}

func (c *localCpanel) execAndUnmarshal(ctx context.Context, binary, module, function string, args Args, out interface{}) error {
	encodedArgs := []string{"--output=json", module, function}
	for k, v := range args {
		encodedArgs = append(encodedArgs, k+"="+url.QueryEscape(fmt.Sprintf("%v", v)))
	}

	if ctx.Value(LogRequestsAndResponses) != nil {
		log.Printf("%s:%s:%s request: %v", binary, module, function, encodedArgs)
	}

	buf, err := exec.CommandContext(ctx, binary, encodedArgs...).Output()
	if err != nil {
		return fmt.Errorf("%s:%s failed: %q", module, function, buf)
	}

	if ctx.Value(LogRequestsAndResponses) != nil {
		log.Printf("%s:%s:%s response: %q", binary, module, function, buf)
	}

	return prettyJSONError(buf, json.Unmarshal(buf, out))
}

func (c *localCpanel) UAPI(ctx context.Context, module, function string, args Args, out interface{}) error {
	// The local UAPI responses have an extra `result` wrapper compared to the remote responses
	var resp struct {
		Result json.RawMessage `json:"result"`
	}
	if err := c.execAndUnmarshal(ctx, uapiPath, module, function, args, &resp); err != nil {
		return err
	}
	return prettyJSONError(resp.Result, json.Unmarshal(resp.Result, out))
}

func (c *localCpanel) API2(ctx context.Context, module, function string, args Args, out interface{}) error {
	// API2 responses need to be unwrapped from the outer cpanelresult hash
	// https://documentation.cpanel.net/display/DD/cPanel+API+2+-+Return+Data
	var resp struct {
		Result json.RawMessage `json:"cpanelresult"`
	}
	if err := c.execAndUnmarshal(ctx, api2Path, module, function, args, &resp); err != nil {
		return err
	}
	return prettyJSONError(resp.Result, json.Unmarshal(resp.Result, out))
}
