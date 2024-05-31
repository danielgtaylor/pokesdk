package pokesdk

// ListGenerations returns a paginator for listing Generations in the API. You
// can manually iterate over pages via `Next(ctx)` or use the `All(ctx)` method
// to get a channel of all results.
//
//	for result := range sdk.ListGenerations().All(ctx) {
//	  if result.Error != nil {
//	    return fmt.Errorf("failed to list generation: %w", result.Error)
//	  }
//	  fmt.Printf("Generation: %s\n", result.Value.Name)
//	}
func (s *SDK) ListGenerations() *Paginator[NamedLink] {
	return &Paginator[NamedLink]{
		sdk: s,
		url: s.baseURL + "/api/v2/generation",
	}
}
