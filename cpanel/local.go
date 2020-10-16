package cpanel

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
)

const (
	uapiPath = "/bin/uapi"
	api2Path = "/bin/cpapi2"
)

var (
	_ API = &localCpanel{}
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

func (c *localCpanel) execAndUnmarshal(binary, module, function string, args map[string]string, out interface{}) error {
	encodedArgs := []string{"--output=json", module, function}
	for k, v := range args {
		encodedArgs = append(encodedArgs, k+"="+url.QueryEscape(v))
	}

	buf, err := exec.Command(binary, encodedArgs...).Output()
	if err != nil {
		return fmt.Errorf("%s:%s failed: %q", module, function, buf)
	}

	return json.Unmarshal(buf, out)
}

func (c *localCpanel) UAPI(module, function string, args map[string]string, out interface{}) error {
	return c.execAndUnmarshal(uapiPath, module, function, args, out)
}

func (c *localCpanel) API2(module, function string, args map[string]string, out interface{}) error {
	return c.execAndUnmarshal(api2Path, module, function, args, out)
}
