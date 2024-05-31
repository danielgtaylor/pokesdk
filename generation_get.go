package pokesdk

import "context"

type Names struct {
	Name     string    `json:"name"`
	Language NamedLink `json:"language"`
}

// Generation is a single generation of Pokemon games and all the Pokemon
// and associated moves within that generation.
type Generation struct {
	ID             int         `json:"id"`
	Name           string      `json:"name"`
	Abilities      []NamedLink `json:"abilities"`
	MainRegion     NamedLink   `json:"main_region"`
	Moves          []NamedLink `json:"moves"`
	Names          []Names     `json:"names"`
	PokemonSpecies []NamedLink `json:"pokemon_species"`
	Types          []NamedLink `json:"types"`
	VersionGroups  []NamedLink `json:"version_groups"`
}

// GetGeneration returns a single Generation from the API.
//
//	gen1, err := sdk.GetGeneration(ctx, "generation-i")
func (s *SDK) GetGeneration(ctx context.Context, name string) (*Generation, error) {
	return Follow[Generation](ctx, s, s.baseURL+"/api/v2/pokemon/"+name)
}
