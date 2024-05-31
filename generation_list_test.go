package pokesdk_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/danielgtaylor/pokesdk"
)

const listGenPage1 = `{
	"count": 123,
	"next": "https://pokeapi.co/api/v2/generation?offset=20&limit=20",
	"previous": null,
	"results": [
		{"name": "generation-i", "url": "https://pokeapi.co/api/v2/generation/1/"},
		{"name": "generation-ii", "url": "https://pokeapi.co/api/v2/generation/2/"},
		{"name": "generation-iii", "url": "https://pokeapi.co/api/v2/generation/3/"}
	]
}`

const listGenPage2 = `{
	"count": 123,
	"next": null,
	"previous": null,
	"results": [
		{"name": "generation-iv", "url": "https://pokeapi.co/api/v2/generation/4/"},
		{"name": "generation-v", "url": "https://pokeapi.co/api/v2/generation/5/"},
		{"name": "generation-vi", "url": "https://pokeapi.co/api/v2/generation/6/"}
	]
}`

func TestListGenerations(t *testing.T) {
	ctx := context.Background()

	transport := &mockTransport{}
	transport.Expect("https://pokeapi.co/api/v2/generation", http.StatusOK, listGenPage1)
	transport.Expect("https://pokeapi.co/api/v2/generation?offset=20&limit=20", http.StatusOK, listGenPage2)

	sdk := pokesdk.New(pokesdk.Config{
		Client: &http.Client{Transport: transport},
	})

	listed := []pokesdk.NamedLink{}
	names := make([]string, 0, 6)
	for result := range sdk.ListGenerations().All(ctx) {
		if result.Error != nil {
			t.Fatalf("failed to list pokemon: %v", result.Error)
		}
		listed = append(listed, result.Value)
		names = append(names, result.Value.Name)
	}

	if len(listed) != 6 {
		t.Errorf("expected 6 results, got %d", len(listed))
	}

	if !reflect.DeepEqual(names, []string{"generation-i", "generation-ii", "generation-iii", "generation-iv", "generation-v", "generation-vi"}) {
		t.Errorf("expected names to match")
	}
}
