package pokesdk_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielgtaylor/pokesdk"
)

func TestGetGeneration(t *testing.T) {
	ctx := context.Background()

	transport := &mockTransport{}
	transport.Expect(
		"https://pokeapi.co/api/v2/generation/generation-i",
		http.StatusOK,
		`{"name":"generation-i"}`, // TODO: add more fields...
	)

	sdk := pokesdk.New(pokesdk.Config{
		Client: &http.Client{Transport: transport},
	})

	gen, err := sdk.GetGeneration(ctx, "generation-i")
	if err != nil {
		t.Fatalf("failed to get gen1: %v", err)
	}

	if gen.Name != "generation-i" {
		t.Errorf("expected generation-i, got %s", gen.Name)
	}
}
