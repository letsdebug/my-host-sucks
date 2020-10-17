// Package cpanel implements local and remote API clients for userland cPanel
package cpanel

import (
	"encoding/json"
	"errors"
	"strings"
)

// Args is a map container for arguments to the cPanel APIs
type Args map[string]interface{}

// API provides the basic cPanel API primitive operations
type API interface {
	UAPI(module, function string, args Args, out interface{}) error
	API2(module, function string, args Args, out interface{}) error
}

// BaseUAPIResponse is the inner UAPI response type.
type BaseUAPIResponse struct {
	BaseResult
	StatusCode int      `json:"status"`
	Errors     []string `json:"errors"`
	Messages   []string `json:"messages"`
}

func (r BaseUAPIResponse) Error() error {
	if r.StatusCode == 1 {
		return nil
	}
	if err := r.BaseResult.Error(); err != nil {
		return err
	}
	if len(r.Errors) == 0 {
		return errors.New("unknown error")
	}
	return errors.New(strings.Join(r.Errors, "\n"))
}

// UAPIResult is the outer UAPI response.
type UAPIResult struct {
	BaseResult
	Result json.RawMessage `json:"result"`
}

// BaseAPI2Response is the inner API2 response type.
type BaseAPI2Response struct {
	BaseResult
	Event struct {
		Result int    `json:"result"`
		Reason string `json:"reason"`
	} `json:"event"`
}

func (r BaseAPI2Response) Error() error {
	if r.Event.Result == 1 {
		return nil
	}
	err := r.BaseResult.Error()
	if err != nil {
		return err
	}
	if len(r.Event.Reason) == 0 {
		return errors.New("Unknown")
	}
	return errors.New(r.Event.Reason)
}

// API2Result is the outer API2 response type.
type API2Result struct {
	BaseResult
	Result json.RawMessage `json:"cpanelresult"`
}

// BaseResult is the basic cPanel API result type, common to both
// UAPI and API2, where it is contained within the respective response type.
type BaseResult struct {
	ErrorString string `json:"error"`
}

func (r BaseResult) Error() error {
	if r.ErrorString == "" {
		return nil
	}
	return errors.New(r.ErrorString)
}

// Message merges each of the messages in the response into a single string.
func (r BaseUAPIResponse) Message() string {
	if r.Messages == nil || len(r.Messages) == 0 {
		return ""
	}
	return strings.Join(r.Messages, "\n")
}
