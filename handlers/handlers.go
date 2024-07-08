package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Kontentski/burgersDb/models"
	"github.com/Kontentski/burgersDb/storage"
	"github.com/gin-gonic/gin"
)

func CreateBurger(c *gin.Context) {
	var burger models.Burgers
	if err := c.ShouldBindJSON(&burger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Begin a transaction
	tx, err := storage.DB.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(context.Background())

	// Insert the burger
	err = tx.QueryRow(context.Background(),
		"INSERT INTO burgers (name, description, is_vegan, image_url, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
		burger.Name, burger.Description, burger.IsVegan, burger.ImageURL, time.Now()).Scan(&burger.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert burger", "details": err.Error()})
		return
	}

	// Insert ingredients and relationships
	for _, ingredient := range burger.Ingredients {
		var ingredientID int
		err := tx.QueryRow(context.Background(),
			"INSERT INTO ingredients (name, description) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description RETURNING id",
			ingredient.Name, ingredient.Description).Scan(&ingredientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert ingredient", "details": err.Error()})
			return
		}

		// Insert the relationship with measure
		_, err = tx.Exec(context.Background(),
			"INSERT INTO burger_ingredients (burger_id, ingredient_id, measure) VALUES ($1, $2, $3)",
			burger.ID, ingredientID, ingredient.Measure)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert burger-ingredient relationship", "details": err.Error()})
			return
		}
	}

	// Commit the transaction
	err = tx.Commit(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, burger)
}

func GetBurgers(c *gin.Context) {
	query := `
        SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
            i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
        FROM burgers b
        LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
        LEFT JOIN ingredients i ON bi.ingredient_id = i.id
    `

	burgers, err := fetchBurgers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, burgers)
}

func GetBurgerById(c *gin.Context) {
	id := c.Param("id")
	query := `
        SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
            i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
        FROM burgers b
        LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
        LEFT JOIN ingredients i ON bi.ingredient_id = i.id
        WHERE b.id=$1
    `

	burgers, err := fetchBurgers(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(burgers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Burger not found"})
		return
	}

	c.JSON(http.StatusOK, burgers[0]) // Assuming only one burger with the given ID
}

func GetBurgerByName(c *gin.Context) {
	name := c.Param("name")
	query := `
        SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
    		i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
        FROM burgers b
        LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
        LEFT JOIN ingredients i ON bi.ingredient_id = i.id
        WHERE LOWER(b.name) LIKE '%' || LOWER($1) || '%'
    `

	burgers, err := fetchBurgers(query, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, burgers)
}

func GetBurgerByLetter(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is required"})
		return
	}
	letter := strings.ToLower(string(name[0]))

	query := `
		SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
			i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
		FROM burgers b
		LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
		LEFT JOIN ingredients i ON bi.ingredient_id = i.id
		WHERE LOWER(b.name) LIKE $1
	`

	burgers, err := fetchBurgers(query, letter+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, burgers)
}

func GetBurgerByRandom(c *gin.Context) {
	query := `
		WITH random_burger AS (
			SELECT id
			FROM burgers
			ORDER BY RANDOM()
			LIMIT 1
		)
		SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
			i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
		FROM burgers b
		LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
		LEFT JOIN ingredients i ON bi.ingredient_id = i.id
		WHERE b.id = (SELECT id FROM random_burger)
	`

	burgers, err := fetchBurgers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(burgers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No burgers found"})
		return
	}

	c.JSON(http.StatusOK, burgers[0])
}

func GetBurgersByRandom(c *gin.Context) {
	query := `
		WITH random_burgers AS (
			SELECT id
			FROM burgers
			ORDER BY RANDOM()
			LIMIT 10
		)
		SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
			i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
		FROM burgers b
		LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
		LEFT JOIN ingredients i ON bi.ingredient_id = i.id
		WHERE b.id IN (SELECT id FROM random_burgers)
	`

	burgers, err := fetchBurgers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(burgers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No burgers found"})
		return
	}

	c.JSON(http.StatusOK, burgers)
}

func GetLatestBurgers(c *gin.Context) {
	query := `
		WITH latest_burgers AS (
			SELECT id
			FROM burgers
			ORDER BY updated_at DESC
			LIMIT 10
		)
		SELECT b.id, b.name, b.description, b.is_vegan, b.image_url, b.updated_at,
			i.id AS ingredient_id, i.name AS ingredient_name, i.description AS ingredient_description, bi.measure AS measure
		FROM burgers b
		LEFT JOIN burger_ingredients bi ON b.id = bi.burger_id
		LEFT JOIN ingredients i ON bi.ingredient_id = i.id
		WHERE b.id IN (SELECT id FROM latest_burgers)
		ORDER BY b.updated_at DESC

	`

	burgers, err := fetchBurgers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(burgers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No burgers found"})
		return
	}

	c.JSON(http.StatusOK, burgers)
}

func GetIngredientByName(c *gin.Context) {
	ingredientName := "%" + c.Param("name") + "%"
	query := `
        SELECT id, name, description
        FROM ingredients
        WHERE LOWER(name) LIKE LOWER($1)
        LIMIT 1
    `

	var ingredient models.Ingredients

	err := storage.DB.QueryRow(context.Background(), query, ingredientName).Scan(
		&ingredient.IngredientID, &ingredient.Name, &ingredient.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "ingredient not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ingredient)
}

func GetIngredientByID(c *gin.Context) {
	ingredientID := c.Param("id")
	query := `
        SELECT id, name, description
        FROM ingredients
        WHERE id = $1
        LIMIT 1
    `

	var ingredient models.Ingredients

	err := storage.DB.QueryRow(context.Background(), query, ingredientID).Scan(
		&ingredient.IngredientID, &ingredient.Name, &ingredient.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "ingredient not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ingredient)
}

func GetBurgersByIngredientName(c *gin.Context) {
	ingredientName := "%" + c.Param("name") + "%"
	query := `
        SELECT b.name, b.image_url, b.id
        FROM burgers b
        JOIN burger_ingredients bi ON b.id = bi.burger_id
        JOIN ingredients i ON bi.ingredient_id = i.id
        WHERE LOWER(i.name) LIKE LOWER($1)
    `
	fmt.Println("SQL Query:", query)

	rows, err := storage.DB.Query(context.Background(), query, ingredientName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var burgers []gin.H
	for rows.Next() {
		var burgerName, imageURL string
		var burgerID int
		err := rows.Scan(&burgerName, &imageURL, &burgerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		burger := gin.H{
			"name":      burgerName,
			"image_url": imageURL,
			"id":        burgerID,
		}
		burgers = append(burgers, burger)
	}

	if len(burgers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "burgers not found for ingredient"})
		return
	}

	c.JSON(http.StatusOK, burgers)
}

func GetBurgersByIngredients(c *gin.Context) {
	ingredientNamesString := c.Query("i") // Get comma-separated list of ingredient names from query parameter "i"
	if ingredientNamesString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no ingredients provided"})
		return
	}

	ingredientNames := strings.Split(ingredientNamesString, ",")
	if len(ingredientNames) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no ingredients provided"})
		return
	}

	// Build the WHERE clause dynamically based on the number of ingredients provided
	var placeholders []string
	var args []interface{}
	for _, name := range ingredientNames {
		placeholders = append(placeholders, "LOWER(i.name) LIKE LOWER($"+strconv.Itoa(len(args)+1)+")")
		args = append(args, "%"+name+"%")
	}
	whereClause := strings.Join(placeholders, " OR ")

	query := `
	SELECT b.id, b.name, b.image_url
	FROM burgers b
	JOIN burger_ingredients bi ON b.id = bi.burger_id
	JOIN ingredients i ON bi.ingredient_id = i.id
	WHERE ` + whereClause + `
	GROUP BY b.id, b.name, b.image_url
	HAVING COUNT(DISTINCT i.id) = $` + strconv.Itoa(len(args)+1) + `
	`

	// Add the number of ingredients to the arguments
	args = append(args, len(ingredientNames))

	// Debugging output
	fmt.Println("SQL Query:", query)
	fmt.Println("Arguments:", args)

	rows, err := storage.DB.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var burgers []gin.H
	for rows.Next() {
		var burgerID int
		var burgerName, imageURL string
		err := rows.Scan(&burgerID, &burgerName, &imageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		burger := gin.H{
			"id":        burgerID,
			"name":      burgerName,
			"image_url": imageURL,
		}
		burgers = append(burgers, burger)
	}

	if len(burgers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "burgers not found for provided ingredients"})
		return
	}

	c.JSON(http.StatusOK, burgers)
}

func fetchBurgers(query string, args ...interface{}) ([]models.Burgers, error) {
	var burgers []models.Burgers

	rows, err := storage.DB.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var burgerID int
		var burgerName, burgerDesc, burgerImageURL string
		var isVegan bool
		var updatedAt time.Time
		var ingredientID sql.NullInt64
		var ingredientName, ingredientDesc, measure sql.NullString

		err := rows.Scan(&burgerID, &burgerName, &burgerDesc, &isVegan, &burgerImageURL, &updatedAt,
			&ingredientID, &ingredientName, &ingredientDesc, &measure)
		if err != nil {
			return nil, err
		}

		var found bool
		for i := range burgers {
			if burgers[i].ID == burgerID {
				if ingredientID.Valid {
					ingredient := models.Ingredients{
						IngredientID: int(ingredientID.Int64),
						Name:         ingredientName.String,
						Description:  ingredientDesc.String,
					}
					if measure.Valid {
						ingredient.Measure = measure.String
					}
					burgers[i].Ingredients = append(burgers[i].Ingredients, ingredient)
				}
				found = true
				break
			}
		}

		if !found {
			burger := models.Burgers{
				ID:          burgerID,
				Name:        burgerName,
				Description: burgerDesc,
				IsVegan:     isVegan,
				ImageURL:    burgerImageURL,
				UpdatedAt:   updatedAt,
				Ingredients: []models.Ingredients{},
			}
			if ingredientID.Valid {
				ingredient := models.Ingredients{
					IngredientID: int(ingredientID.Int64),
					Name:         ingredientName.String,
					Description:  ingredientDesc.String,
				}
				if measure.Valid {
					ingredient.Measure = measure.String
				}
				burger.Ingredients = append(burger.Ingredients, ingredient)
			}
			burgers = append(burgers, burger)
		}
	}

	return burgers, nil
}

func GetBurgersByVeganStatus(c *gin.Context, isVegan bool) {
	query := `
        SELECT b.id, b.name, b.image_url
        FROM burgers b
        WHERE b.is_vegan = $1
    `

	rows, err := storage.DB.Query(context.Background(), query, isVegan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var burgers []gin.H
	for rows.Next() {
		var burgerID int
		var burgerName, imageURL string
		err := rows.Scan(&burgerID, &burgerName, &imageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		burger := gin.H{
			"id":        burgerID,
			"name":      burgerName,
			"image_url": imageURL,
		}
		burgers = append(burgers, burger)
	}

	if len(burgers) == 0 {
		if isVegan {
			c.JSON(http.StatusNotFound, gin.H{"error": "no vegan burgers found"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "no non-vegan burgers found"})
		}
		return
	}

	c.JSON(http.StatusOK, burgers)
}

func GetVeganBurgers(c *gin.Context) {
	GetBurgersByVeganStatus(c, true)
}

func GetNonVeganBurgers(c *gin.Context) {
	GetBurgersByVeganStatus(c, false)
}
