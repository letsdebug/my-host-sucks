package cpanel

import (
	"context"
	"encoding/json"
	"strings"
)

// DomainsDataResponse is the response data from DomainInfo:domains_data
type DomainsDataResponse struct {
	BaseUAPIResponse
	Data struct {
		Main       DomainsDataDomain   `json:"main_domain"`
		Addons     []DomainsDataDomain `json:"addon_domain"`
		Subdomains []DomainsDataDomain `json:"sub_domains"`
		// Parked domains are intentionally omitted here because it is made redundant
		// by Main.ServerAlias.
	} `json:"data"`
}

// DomainsDataDomain is an individual virtualhost entry within DomainsDataResponse
type DomainsDataDomain struct {
	Domain       string        `json:"domain"`
	DocumentRoot string        `json:"documentroot"`
	ServerName   string        `json:"servername"`
	ServerAlias  ServerAliases `json:"serveralias"`
}

// ServerAliases represents a list of ServerAlias within a DomainsDataDomain
type ServerAliases []string

// UnmarshalJSON implements custom JSON unmarshalling for ServerAliases
func (sa *ServerAliases) UnmarshalJSON(buf []byte) error {
	var s string
	if err := json.Unmarshal(buf, &s); err != nil {
		return err
	}
	*sa = strings.Split(s, " ")
	return nil
}

// DomainsData invokes DomainInfo:domains_data
func (c *Client) DomainsData(ctx context.Context) (*DomainsDataResponse, error) {
	var resp DomainsDataResponse
	if err := c.a.UAPI(ctx, "DomainInfo", "domains_data", nil, &resp); err != nil {
		return &resp, err
	}
	return &resp, resp.Error()
}
