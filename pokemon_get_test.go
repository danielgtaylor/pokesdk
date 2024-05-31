package pokesdk_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/danielgtaylor/pokesdk"
)

func TestGetPokemon(t *testing.T) {
	ctx := context.Background()

	transport := &mockTransport{}
	transport.Expect(
		"https://pokeapi.co/api/v2/pokemon/pikachu",
		http.StatusOK,
		`{"name":"pikachu"}`, // TODO: add more fields...
	)

	sdk := pokesdk.New(pokesdk.Config{
		Client: &http.Client{Transport: transport},
	})

	pika, err := sdk.GetPokemon(ctx, "pikachu")
	if err != nil {
		t.Fatalf("failed to get pikachu: %v", err)
	}

	if pika.Name != "pikachu" {
		t.Errorf("expected pikachu, got %s", pika.Name)
	}
}

func TestGetPokemon404(t *testing.T) {
	ctx := context.Background()

	transport := &mockTransport{}
	transport.Expect(
		"https://pokeapi.co/api/v2/pokemon/pikachu",
		http.StatusNotFound,
		"",
	)

	sdk := pokesdk.New(pokesdk.Config{
		Client: &http.Client{Transport: transport},
	})

	_, err := sdk.GetPokemon(ctx, "pikachu")
	if !errors.Is(err, pokesdk.APIError) {
		t.Fatalf("expected error from 404 response")
	}
}
