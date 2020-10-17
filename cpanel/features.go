package cpanel

import "context"

// ListFeaturesResponse is the response from Features:list_features
type ListFeaturesResponse struct {
	BaseUAPIResponse
	Data map[string]int `json:"data"`
}

// HasFeature returns whether the requested feature is present
func (r *ListFeaturesResponse) HasFeature(feature string) bool {
	if r.Data == nil {
		return false
	}
	if v, ok := r.Data[feature]; ok && v == 1 {
		return true
	}
	return false
}

// ListFeatures invokes Features:list_features
func (c *Client) ListFeatures(ctx context.Context) (*ListFeaturesResponse, error) {
	var resp ListFeaturesResponse
	if err := c.a.UAPI(ctx, "Features", "list_features", nil, &resp); err != nil {
		return &resp, err
	}
	return &resp, resp.Error()
}
