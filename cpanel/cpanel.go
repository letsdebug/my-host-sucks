// Package cpanel implements local and remote API clients for userland cPanel
package cpanel

// Args is a map container for arguments to the cPanel APIs
type Args map[string]interface{}

// API provides the basic cPanel API primitive operations
type API interface {
	UAPI(module, function string, args Args, out interface{}) error
	API2(module, function string, args Args, out interface{}) error
}

// DomainsDataResponse is the response data from DomainInfo:domains_data
type DomainsDataResponse struct {
}

// DomainsData invokes DomainInfo:domains_data
func DomainsData(api API) (*DomainsDataResponse, error) {
	resp := &DomainsDataResponse{}
	return resp, api.UAPI("DomainInfo", "domains_data", nil, resp)
}

// MkdirResponse is the response data from Fileman:mkdir
type MkdirResponse struct {
}

// Mkdir invokes Fileman:mkdir
func Mkdir(api API, path, name, permissions string) (*MkdirResponse, error) {
	resp := &MkdirResponse{}
	return resp, api.API2("Fileman", "mkdir", Args{
		"path":        path,
		"name":        name,
		"permissions": permissions,
	}, &resp)
}
