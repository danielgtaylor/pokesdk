package pokesdk_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/danielgtaylor/pokesdk"
)

const listResultPage1 = `{
	"count": 123,
	"next": "https://pokeapi.co/api/v2/pokemon?offset=20&limit=20",
	"previous": null,
	"results": [
		{"name": "bulbasaur", "url": "https://pokeapi.co/api/v2/pokemon/1/"},
		{"name": "ivysaur", "url": "https://pokeapi.co/api/v2/pokemon/2/"},
		{"name": "venusaur", "url": "https://pokeapi.co/api/v2/pokemon/3/"}
	]
}`

const listResultPage2 = `{
	"count": 123,
	"next": null,
	"previous": null,
	"results": [
		{"name": "charmander", "url": "https://pokeapi.co/api/v2/pokemon/4/"},
		{"name": "charmeleon", "url": "https://pokeapi.co/api/v2/pokemon/5/"},
		{"name": "charizard", "url": "https://pokeapi.co/api/v2/pokemon/6/"}
	]
}`

func TestListPokemon(t *testing.T) {
	ctx := context.Background()

	transport := &mockTransport{}
	transport.Expect("https://pokeapi.co/api/v2/pokemon", http.StatusOK, listResultPage1)
	transport.Expect("https://pokeapi.co/api/v2/pokemon?offset=20&limit=20", http.StatusOK, listResultPage2)

	sdk := pokesdk.New(pokesdk.Config{
		Client: &http.Client{Transport: transport},
	})

	listed := []pokesdk.NamedLink{}
	names := make([]string, 0, 6)
	for result := range sdk.ListPokemon().All(ctx) {
		if result.Error != nil {
			t.Fatalf("failed to list pokemon: %v", result.Error)
		}
		listed = append(listed, result.Value)
		names = append(names, result.Value.Name)
	}

	if len(listed) != 6 {
		t.Errorf("expected 6 results, got %d", len(listed))
	}

	if !reflect.DeepEqual(names, []string{"bulbasaur", "ivysaur", "venusaur", "charmander", "charmeleon", "charizard"}) {
		t.Errorf("expected names to match")
	}
}

func TestListPokemonFailure(t *testing.T) {
	ctx := context.Background()

	transport := &mockTransport{}
	transport.Expect("https://pokeapi.co/api/v2/pokemon", http.StatusOK, listResultPage1)
	transport.Expect("https://pokeapi.co/api/v2/pokemon?offset=20&limit=20", http.StatusInternalServerError, "")

	sdk := pokesdk.New(pokesdk.Config{
		Client: &http.Client{Transport: transport},
	})

	var err error
	for result := range sdk.ListPokemon().All(ctx) {
		if result.Error != nil {
			err = result.Error
			break
		}
	}

	if err == nil {
		t.Errorf("expected error")
	}
}

func TestListPokemonCancel(t *testing.T) {
	ctx := context.Background()

	// Since we are going to cancel, we only need to set up the first page. Any
	// additional request that goes out if the code doesn't clean up properly
	// will result in an error which will cause the test to fail.
	transport := &mockTransport{}
	transport.Expect("https://pokeapi.co/api/v2/pokemon", http.StatusOK, listResultPage1)

	sdk := pokesdk.New(pokesdk.Config{
		Client: &http.Client{Transport: transport},
	})

	var err error
	iter, cancel := sdk.ListPokemon().AllWithCancel(ctx)
	for result := range iter {
		if result.Error != nil {
			err = result.Error
			break
		}

		// We are done, so let the paginator know to stop.
		cancel()
		break
	}

	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}
}
