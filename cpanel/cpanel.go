// Package cpanel implements local and remote API clients for userland cPanel
package cpanel

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// ContextKey provides a type that may be passed into API invocations
type ContextKey int

const (
	// LogRequestsAndResponses logs all requests and responses from the
	// cPanel API to the terminal
	LogRequestsAndResponses ContextKey = iota
)

// Args is a map container for arguments to the cPanel APIs
type Args map[string]interface{}

// Client implements a cPanel API client
type Client struct {
	a api
}

// API provides the basic cPanel API primitive operations
type api interface {
	UAPI(ctx context.Context, module, function string, args Args, out interface{}) error
	API2(ctx context.Context, module, function string, args Args, out interface{}) error
}

// BaseUAPIResponse is the inner UAPI response type.
type BaseUAPIResponse struct {
	StatusCode int      `json:"status"`
	Errors     []string `json:"errors"`
	Messages   []string `json:"messages"`
}

func (r BaseUAPIResponse) Error() error {
	if r.StatusCode == 1 {
		return nil
	}
	if len(r.Errors) == 0 {
		return errors.New("unknown error")
	}
	return errors.New(strings.Join(r.Errors, "\n"))
}

// Message merges each of the messages in the response into a single string.
func (r BaseUAPIResponse) Message() string {
	if r.Messages == nil || len(r.Messages) == 0 {
		return ""
	}
	return strings.Join(r.Messages, "\n")
}

// BaseAPI2Response is the inner API2 response type.
type BaseAPI2Response struct {
	Event struct {
		Result int    `json:"result"`
		Reason string `json:"reason"`
	} `json:"event"`
}

func (r BaseAPI2Response) Error() error {
	if r.Event.Result == 1 {
		return nil
	}
	if len(r.Event.Reason) == 0 {
		return errors.New("Unknown")
	}
	return errors.New(r.Event.Reason)
}

func prettyJSONError(buf []byte, err error) error {
	if err == nil {
		return nil
	}
	if len(buf) > 64 {
		buf = buf[:64]
	}
	return fmt.Errorf("JSON decoding error: %w (contents: %q...)", err, buf)
}
