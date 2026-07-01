package handlers

import (
	"net/http"

	"lieferino-backend/database"
	"lieferino-backend/models"

	"github.com/gin-gonic/gin"
)

// 🏪 ListRestaurants liefert alle Verkäufer inkl. ihrer Produkte (öffentlich).
func ListRestaurants(c *gin.Context) {
	var restaurants []models.Restaurant
	database.DB.Preload("Produkte").Find(&restaurants)
	c.JSON(http.StatusOK, restaurants)
}

// 🍕 ListProducts liefert die Produkte eines Restaurants (per Slug, öffentlich).
func ListProducts(c *gin.Context) {
	slug := c.Param("slug")
	var produkte []models.Product
	database.DB.Where("restaurant_slug = ?", slug).Find(&produkte)
	c.JSON(http.StatusOK, produkte)
}
