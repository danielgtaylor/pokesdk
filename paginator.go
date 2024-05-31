package pokesdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// DefaultPageBufferSize is the size of the channel holding page items. When
// this reaches the last `n` items a new request is made in the background to
// fetch more items. Increasing this value will increase memory usage but ensure
// items are more likely to already be available when requested.
var DefaultPageBufferSize = 10

// Page is a single page of results from the API. It may contain next/prev
// links for pagination, and a total count of items.
type Page[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

// Paginator is a helper for paginating through API results. It can be used to
// manually iterate over pages or to get a channel of all results.
type Paginator[T any] struct {
	sdk *SDK
	url string
}

// Next fetches the next page of results from the API. If there are no more
// pages, the `Next` field of the returned page will be empty.
func (p *Paginator[T]) Next(ctx context.Context) (*Page[T], error) {
	resp, err := p.sdk.Request(ctx, http.MethodGet, p.url, nil)
	if err != nil {
		return nil, err
	}

	var page *Page[T]
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	p.url = page.Next
	return page, nil
}

// IteratorResult is a single result from the paginator. It contains the page
// the result was found on, the value itself, and any error that occurred. This
// is used to send results oßver a channel so that errors can still be detected.
type IteratorResult[T any] struct {
	Page  *Page[T]
	Index int
	Value T
	Error error
}

// All returns a channel of all results from the paginator. This will fetch
// pages in the background as needed and close the channel when there are no
// more results. If an error occurs, the error will be sent on the channel and
// the channel will be closed.
//
//	for result := range paginator.All(ctx) {
//		if result.Error != nil {
//			return fmt.Errorf("failed to list items: %w", result.Error)
//		}
//		fmt.Printf("Item: %+v\n", result.Value)
//	}
func (p *Paginator[T]) All(ctx context.Context) chan IteratorResult[T] {
	iter, _ := p.AllWithCancel(ctx)
	return iter
}

// AllWithCancel returns a channel of all results from the paginator. This will fetch
// pages in the background as needed and close the channel when there are no
// more results. If an error occurs, the error will be sent on the channel and
// the channel will be closed.
//
//	iter, cancel := paginator.All(ctx)
//	for result := range iter {
//		if result.Error != nil {
//			return fmt.Errorf("failed to list items: %w", result.Error)
//		}
//		fmt.Printf("Item: %+v\n", result.Value)
//	}
func (p *Paginator[T]) AllWithCancel(ctx context.Context) (chan IteratorResult[T], func()) {
	ch := make(chan IteratorResult[T], DefaultPageBufferSize)
	done := make(chan struct{}, 1)
	closed := false

	go func() {
		defer close(ch)
		defer func() {
			if !closed {
				close(done)
			}
		}()

		index := 0
		for p.url != "" {
			page, err := p.Next(ctx)
			if err != nil {
				select {
				case <-done:
					return
				case ch <- IteratorResult[T]{Page: page, Index: index, Error: err}:
				}
				return
			}

			for _, v := range page.Results {
				select {
				case <-done:
					return
				case ch <- IteratorResult[T]{Page: page, Index: index, Value: v}:
				}
				index++
			}
		}
	}()

	return ch, func() {
		closed = true
		close(done)
	}
}
