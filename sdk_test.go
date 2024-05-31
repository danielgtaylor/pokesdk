package pokesdk_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// mockTransport is a mock HTTP transport for testing. It takes a map of URLs to
// expected responses via `Expect` calls.
type mockTransport struct {
	responses map[string][]*http.Response
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Note: each call pops the first response off the list for the given URL.
	// This allows for multiple responses to be expected for the same URL.
	if r, ok := t.responses[req.URL.String()]; ok {
		if len(r) > 0 {
			res := r[0]
			t.responses[req.URL.String()] = r[1:]
			return res, nil
		}
	}
	return nil, fmt.Errorf("unexpected request: %s", req.URL.String())
}

// Expect adds an expected response for a given URL (including query params).
func (t *mockTransport) Expect(url string, status int, body string) {
	if t.responses == nil {
		t.responses = make(map[string][]*http.Response)
	}
	t.responses[url] = append(t.responses[url], &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
	})
}
