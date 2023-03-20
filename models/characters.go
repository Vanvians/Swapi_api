package models

type Character struct {
	Name string `json:"name"`
	Height float64 `json:"height"`
	Gender string `json:"gender"`
}

type CharacterMetadata struct {
	TotalCharacters int `json:"total_characters"`
	TotalHeightCm float64 `json:"total_height_cm"`
	TotalHeightFt string `json:"total_height_ft"`
}