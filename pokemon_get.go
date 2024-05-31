package pokesdk

import "context"

type Abilities struct {
	IsHidden bool      `json:"is_hidden"`
	Slot     int       `json:"slot"`
	Ability  NamedLink `json:"ability"`
}

type GameIndices struct {
	GameIndex int       `json:"game_index"`
	Version   NamedLink `json:"version"`
}

type VersionDetails struct {
	Rarity  int       `json:"rarity"`
	Version NamedLink `json:"version"`
}

type HeldItems struct {
	Item           NamedLink        `json:"item"`
	VersionDetails []VersionDetails `json:"version_details"`
}

type VersionGroupDetails struct {
	LevelLearnedAt  int       `json:"level_learned_at"`
	VersionGroup    NamedLink `json:"version_group"`
	MoveLearnMethod NamedLink `json:"move_learn_method"`
}

type Moves struct {
	Move                NamedLink             `json:"move"`
	VersionGroupDetails []VersionGroupDetails `json:"version_group_details"`
}

type Sprites struct {
	BackDefault      string                               `json:"back_default"`
	BackFemale       string                               `json:"back"`
	BackShiny        string                               `json:"back_shiny"`
	BackShinyFemale  string                               `json:"back_shiny_female"`
	FrontDefault     string                               `json:"front_default"`
	FrontFemale      string                               `json:"front"`
	FrontShiny       string                               `json:"front_shiny"`
	FrontShinyFemale string                               `json:"front_shiny_female"`
	Other            map[string]map[string]string         `json:"other"`
	Versions         map[string]map[string]map[string]any `json:"versions"`
}

type Stats struct {
	BaseStat int       `json:"base_stat"`
	Effort   int       `json:"effort"`
	Stat     NamedLink `json:"stat"`
}

type Types struct {
	Slot int       `json:"slot"`
	Type NamedLink `json:"type"`
}

type PastTypes struct {
	Generation NamedLink `json:"generation"`
	Types      []Types   `json:"types"`
}

// Pokemon is a single Pokemon and all its associated data.
type Pokemon struct {
	ID                     int               `json:"id"`
	Name                   string            `json:"name"`
	BaseExperience         int               `json:"base_experience"`
	Height                 int               `json:"height"`
	IsDefault              bool              `json:"is_default"`
	Order                  int               `json:"order"`
	Weight                 int               `json:"weight"`
	Abilities              []Abilities       `json:"abilities"`
	Forms                  []NamedLink       `json:"forms"`
	GameIndices            []GameIndices     `json:"game_indices"`
	HeldItems              []HeldItems       `json:"held_items"`
	LocationAreaEncounters string            `json:"location_area_encounters"`
	Moves                  []Moves           `json:"moves"`
	Species                NamedLink         `json:"species"`
	Sprites                Sprites           `json:"sprites"`
	Cries                  map[string]string `json:"cries"`
	Stats                  []Stats           `json:"stats"`
	Types                  []Types           `json:"types"`
	PastTypes              []PastTypes       `json:"past_types"`
}

// GetPokemon returns a single Pokemon from the API.
//
//	pikachu, err := sdk.GetPokemon(ctx, "pikachu")
func (s *SDK) GetPokemon(ctx context.Context, name string) (*Pokemon, error) {
	return Follow[Pokemon](ctx, s, s.baseURL+"/api/v2/pokemon/"+name)
}
