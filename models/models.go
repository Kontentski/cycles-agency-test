package models

import "time"

type Ingredients struct {
	IngredientID int    `json:"ingredient_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Measure      string `json:"measure"`
}

type Burgers struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	IsVegan     bool          `json:"is_vegan"`
	ImageURL    string        `json:"image_url"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Ingredients []Ingredients `json:"Ingredients"`
}

type BurgerIngredients struct {
	BurgerID     int `json:"burger_id"`
	IngredientID int `json:"ingredient_id"`
	Measure      int `json:"measure"`
}
