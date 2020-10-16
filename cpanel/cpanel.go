// Package cpanel implements local and remote API clients for userland cPanel
package cpanel

// API provides the basic cPanel API primitive operations
type API interface {
	UAPI(module, function string, args map[string]string, out interface{}) error
	API2(module, function string, args map[string]string, out interface{}) error
}

// DomainsDataResponse is the response data from DomainInfo:domains_data
type DomainsDataResponse struct {
}

// DomainsData invokes DomainInfo:domains_data
func DomainsData(api API) (*DomainsDataResponse, error) {
	resp := &DomainsDataResponse{}
	return resp, api.UAPI("DomainInfo", "domains_data", nil, resp)
}
