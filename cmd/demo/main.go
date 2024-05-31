package main

import (
	"context"
	"fmt"
	"log"

	"github.com/danielgtaylor/pokesdk"
)

func main() {
	ctx := context.Background()
	sdk := pokesdk.New(pokesdk.Config{})

	// Print up to 50 pokemon names.
	iter, cancel := sdk.ListPokemon().AllWithCancel(ctx)
	for result := range iter {
		if result.Index > 50 {
			// We are done, so let the paginator know to stop.
			cancel()
			break
		}
		if result.Error != nil {
			log.Fatalf("Failed to list Pokemon: %v", result.Error)
		}
		fmt.Printf("Pokemon: %s\n", result.Value.Name)
	}

	// Print Pikachu's stat details.
	pika, err := sdk.GetPokemon(ctx, "pikachu")
	if err != nil {
		panic(err)
	}
	fmt.Println("Pikachu stats:")
	fmt.Println(pika.Stats)

	// Print generation names.
	for result := range sdk.ListGenerations().All(ctx) {
		if result.Error != nil {
			log.Fatalf("Failed to list Generation: %v", result.Error)
		}
		fmt.Printf("Generation: %+v\n", result.Value.Name)
	}
}
