package models

import "time"

type Ingredients struct {
	IngredientID int    `json:"ingredient_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type Burgers struct {
	ID                int                 `json:"id"`
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	IsVegan           bool                `json:"is_vegan"`
	ImageURL          string              `json:"image_url"`
	UpdatedAt         time.Time           `json:"updated_at"`
	Ingredients       []Ingredients       `json:"ingredients,omitempty"`
	BurgerIngredients []BurgerIngredients `json:"burger_ingredients,omitempty"`
}

type BurgerIngredients struct {
	BurgerID       int    `json:"burger_id"`
	IngredientID   int    `json:"ingredient_id"`
	IngredientName string `json:"ingredient_name"`
	Measure        string `json:"measure"`
}
