package pokesdk

// ListPokemon returns a paginator for listing Pokemon in the API. You can
// manually iterate over pages via `Next(ctx)` or use the `All(ctx)` method to
// get a channel of all results.
//
//	for result := range sdk.ListPokemon().All(ctx) {
//		if result.Error != nil {
//			return fmt.Errorf("failed to list pokemon: %w", result.Error)
//		}
//		fmt.Printf("Pokemon: %s\n", result.Value.Name)
//	}
func (s *SDK) ListPokemon() *Paginator[NamedLink] {
	return &Paginator[NamedLink]{
		sdk: s,
		url: s.baseURL + "/api/v2/pokemon",
	}
}
