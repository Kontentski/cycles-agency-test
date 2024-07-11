package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Kontentski/burgersDb/models"
	"github.com/Kontentski/burgersDb/queries"
	"github.com/Kontentski/burgersDb/storage"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var Domain = os.Getenv("DOMAIN")

func CreateBurger(c *gin.Context) {
	var req struct {
		Burger            models.Burgers             `json:"burger"`
		Ingredients       []models.Ingredients       `json:"ingredients"`
		BurgerIngredients []models.BurgerIngredients `json:"burgerIngredients"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := storage.DB.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback(context.Background())

	var burgerID int
	err = tx.QueryRow(context.Background(),
		`INSERT INTO burgers (name, description, is_vegan, image_url, updated_at) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.Burger.Name, req.Burger.Description, req.Burger.IsVegan, req.Burger.ImageURL, time.Now()).Scan(&burgerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ingredientIDMap := make(map[string]int)

	for _, ingredient := range req.Ingredients {
		var ingredientID int

		// Check if ingredient already exists
		err = tx.QueryRow(context.Background(),
			`SELECT id FROM ingredients WHERE name = $1`,
			ingredient.Name).Scan(&ingredientID)
		if err != nil {
			if err == pgx.ErrNoRows {
				err = tx.QueryRow(context.Background(),
					`INSERT INTO ingredients (name, description) VALUES ($1, $2) RETURNING id`,
					ingredient.Name, ingredient.Description).Scan(&ingredientID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		ingredientIDMap[ingredient.Name] = ingredientID
	}

	for _, bi := range req.BurgerIngredients {
		ingredientID, ok := ingredientIDMap[bi.IngredientName]
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ingredient not found"})
			return
		}

		_, err := tx.Exec(context.Background(),
			`INSERT INTO burger_ingredients (burger_id, ingredient_id, measure) VALUES ($1, $2, $3)`,
			burgerID, ingredientID, bi.Measure)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			log.Printf("%d, %d, %s", burgerID, ingredientID, bi.Measure)
			return
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"burger_id": burgerID})
}

func GetBurgers(c *gin.Context) {
	burgers, err := fetchBurgers(queries.GetBurgers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, burgers)
}

func GetBurgerById(c *gin.Context) {
	id := c.Param("id")

	burgers, err := fetchBurgers(queries.GetBurgerById, id)
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

	burgers, err := fetchBurgers(queries.GetBurgerByName, name)
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

	burgers, err := fetchBurgers(queries.GetBurgerByLetter, letter+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, burgers)
}

func GetBurgerByRandom(c *gin.Context) {
	burgers, err := fetchBurgers(queries.GetBurgerByRandom)
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
	burgers, err := fetchBurgers(queries.GetBurgersByRandom)
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
	burgers, err := fetchBurgers(queries.GetLatestBurgers)
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

	var ingredient models.Ingredients

	err := storage.DB.QueryRow(context.Background(), queries.GetIngredientByName, ingredientName).Scan(
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

	var ingredient models.Ingredients

	err := storage.DB.QueryRow(context.Background(), queries.GetIngredientByID, ingredientID).Scan(
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
	fmt.Print("SQL Query:", queries.GetBurgersByIngredientName)

	rows, err := storage.DB.Query(context.Background(), queries.GetBurgersByIngredientName, ingredientName)
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

		fullImageURL := Domain + imageURL

		burger := gin.H{
			"name":      burgerName,
			"image_url": fullImageURL,
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
	ingredientNamesString := c.Query("i")
	if ingredientNamesString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no ingredients provided"})
		return
	}

	ingredientNames := strings.Split(ingredientNamesString, ",")
	if len(ingredientNames) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no ingredients provided"})
		return
	}

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

	args = append(args, len(ingredientNames))

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

		fullImageURL := Domain + imageURL

		burger := gin.H{
			"id":        burgerID,
			"name":      burgerName,
			"image_url": fullImageURL,
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

		fullImageURL := Domain + burgerImageURL

		var found bool
		for i := range burgers {
			if burgers[i].ID == burgerID {
				if ingredientID.Valid {
					ingredient := models.Ingredients{
						IngredientID: int(ingredientID.Int64),
						Name:         ingredientName.String,
						Description:  ingredientDesc.String,
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
				ImageURL:    fullImageURL,
				UpdatedAt:   updatedAt,
				Ingredients: []models.Ingredients{},
			}
			if ingredientID.Valid {
				ingredient := models.Ingredients{
					IngredientID: int(ingredientID.Int64),
					Name:         ingredientName.String,
					Description:  ingredientDesc.String,
				}
				burger.Ingredients = append(burger.Ingredients, ingredient)
			}
			burgers = append(burgers, burger)
		}
	}

	return burgers, nil
}

func GetBurgersByVeganStatus(c *gin.Context, isVegan bool) {
	rows, err := storage.DB.Query(context.Background(), queries.GetBurgersByVeganStatus, isVegan)
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

		fullImageURL := Domain + imageURL

		burger := gin.H{
			"id":        burgerID,
			"name":      burgerName,
			"image_url": fullImageURL,
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
