package cpanel

import "context"

// MkdirResponse is the response data from Fileman:mkdir
type MkdirResponse struct {
	BaseAPI2Response
}

// Mkdir invokes Fileman:mkdir
func Mkdir(ctx context.Context, api API, path, name, permissions string) (*MkdirResponse, error) {
	resp := &MkdirResponse{}
	if err := api.API2(ctx, "Fileman", "mkdir", Args{
		"path":        path,
		"name":        name,
		"permissions": permissions,
	}, &resp); err != nil {
		return resp, err
	}
	return resp, resp.Error()
}
