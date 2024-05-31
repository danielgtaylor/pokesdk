package pokesdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// APIError is an error type for API errors, such as 404 not found responses.
var APIError = errors.New("API error")

// NamedLink is a common structure for named links in the API that contain a
// name and a URL.
type NamedLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Config provides optional configuration for the Pokemon API SDK.
type Config struct {
	BaseURL string
	Client  *http.Client
	// TODO: add auth if desired...
}

// SDK is the Pokemon API SDK.
type SDK struct {
	baseURL string
	client  *http.Client
}

// New returns a new instance of the Pokemon API SDK.
func New(config Config) *SDK {
	if config.BaseURL == "" {
		config.BaseURL = "https://pokeapi.co"
	}

	if config.Client == nil {
		config.Client = http.DefaultClient
	}

	return &SDK{
		baseURL: config.BaseURL,
		client:  config.Client,
	}
}

// Request makes an HTTP request to the given URL with the given method and
// body using the SDK's client. It returns the response or an error.
func (s *SDK) Request(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// TODO: auth headers could be inserted here.

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

// Follow is a helper function to follow a URL and decode the response into a
// pointer of the given type. This is useful for following links in API
// responses hypermedia-style.
//
//	thing, err := Follow[Thing](ctx, sdk, "https://example.com/things/123")
func Follow[T any](ctx context.Context, sdk *SDK, url string) (*T, error) {
	resp, err := sdk.Request(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status %d response: %w", resp.StatusCode, APIError)
	}

	// TODO: content negotiation could be added here to support more formats.
	var value *T
	if err := json.NewDecoder(resp.Body).Decode(&value); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return value, nil
}
