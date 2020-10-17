package cpanel

import "context"

// MkdirResponse is the response data from Fileman:mkdir
type MkdirResponse struct {
	BaseAPI2Response
}

// Mkdir invokes Fileman:mkdir
func (c *Client) Mkdir(ctx context.Context, path, name, permissions string) (*MkdirResponse, error) {
	resp := &MkdirResponse{}
	if err := c.a.API2(ctx, "Fileman", "mkdir", Args{
		"path":        path,
		"name":        name,
		"permissions": permissions,
	}, &resp); err != nil {
		return resp, err
	}
	return resp, resp.Error()
}
