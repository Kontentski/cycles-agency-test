package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/Kontentski/burgersDb/models"
	"github.com/Kontentski/burgersDb/storage"
	"github.com/gin-gonic/gin"
)

func GetBurgers(c *gin.Context) {
	rows, err := storage.DB.Query(context.Background(), "SELECT id, name, description, image_url FROM burgers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var burgers []models.Burgers
	for rows.Next() {
		var burger models.Burgers
		err := rows.Scan(&burger.ID, &burger.Name, &burger.Description, &burger.ImageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		burgers = append(burgers, burger)
	}

	c.JSON(http.StatusOK, burgers)
}

func GetBurgerById(c *gin.Context) {
	id := c.Param("id")
	var burger models.Burgers
	err := storage.DB.QueryRow(context.Background(), "SELECT id, name, description, image_url FROM burgers WHERE id=$1", id).Scan(&burger.ID, &burger.Name, &burger.Description, &burger.ImageURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Burger not found"})
		return
	}

	c.JSON(http.StatusOK, burger)
}

func GetBurgerByName(c *gin.Context) {
	name := c.Param("name")
	var burgers []models.Burgers

	rows, err := storage.DB.Query(context.Background(), "SELECT id, name, description, image_url FROM burgers WHERE LOWER(name) LIKE '%' || LOWER($1) || '%'", name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var burger models.Burgers
		err := rows.Scan(&burger.ID, &burger.Name, &burger.Description, &burger.ImageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		burgers = append(burgers, burger)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Row iteration failed"})
		return
	}
	if len(burgers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No burgers found"})
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
	var burgers []models.Burgers

	rows, err := storage.DB.Query(context.Background(), "SELECT id, name, description, image_url FROM burgers WHERE LOWER(name) LIKE $1", letter+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var burger models.Burgers
		err := rows.Scan(&burger.ID, &burger.Name, &burger.Description, &burger.ImageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		burgers = append(burgers, burger)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Row iteration failed"})
		return
	}
	if len(burgers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No burgers found"})
		return
	}

	c.JSON(http.StatusOK, burgers)
}

func GetBurgerByRandom(c *gin.Context) {
	var burger models.Burgers

	// Execute the query to get a random burger
	err := storage.DB.QueryRow(context.Background(), "SELECT id, name, description, image_url FROM burgers ORDER BY RANDOM() LIMIT 1").Scan(&burger.ID, &burger.Name, &burger.Description, &burger.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch random burger"})
		return
	}

	// Return the random burger as a JSON response
	c.JSON(http.StatusOK, burger)
}

func GetBurgersByRandom(c *gin.Context) {
	var burgers []models.Burgers

	// Execute the SQL query to fetch 10 random burgers
	rows, err := storage.DB.Query(context.Background(), "SELECT id, name, description, image_url FROM burgers ORDER BY RANDOM() LIMIT 10;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var burger models.Burgers
		err := rows.Scan(&burger.ID, &burger.Name, &burger.Description, &burger.ImageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan burgers"})
			return
		}
		burgers = append(burgers, burger)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows"})
		return
	}

	// Return the list of random burgers as JSON response
	c.JSON(http.StatusOK, burgers)
}

func GetLatestBurgers(c *gin.Context) {
    var burgers []models.Burgers

	rows, err := storage.DB.Query(context.Background(), "SELECT id, name, description, image_url, updated_at FROM burgers ORDER BY updated_at DESC LIMIT 10;")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
        return
    }
    defer rows.Close()

    for rows.Next() {
        var burger models.Burgers
        err := rows.Scan(&burger.ID, &burger.Name, &burger.Description, &burger.ImageURL, &burger.UpdatedAt)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan burgers"})
            return
        }
        burgers = append(burgers, burger)
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows"})
        return
    }

    c.JSON(http.StatusOK, burgers)
}


func CreateBurger(c *gin.Context) {
	var burger models.Burgers
	if err := c.ShouldBindJSON(&burger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := storage.DB.QueryRow(context.Background(), "INSERT INTO burgers (name, description, image_url) VALUES ($1, $2, $3) RETURNING id", burger.Name, burger.Description, burger.ImageURL).Scan(&burger.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, burger)
}
