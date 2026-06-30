package handlers

import (
	"net/http"

	"lieferino-backend/database"
	"lieferino-backend/models"

	"github.com/gin-gonic/gin"
)

// ❤️ ListFavorites liefert die Slugs der Lieblings-Restaurants des Nutzers.
func ListFavorites(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var favs []models.Favorite
	database.DB.Where("user_id = ?", userID).Find(&favs)

	slugs := []string{}
	for _, f := range favs {
		slugs = append(slugs, f.RestaurantSlug)
	}
	c.JSON(http.StatusOK, slugs)
}

// ➕ AddFavorite fügt ein Restaurant zu den Favoriten hinzu (ohne Dopplung).
func AddFavorite(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	slug := c.Param("slug")

	var vorhanden models.Favorite
	err := database.DB.Where("user_id = ? AND restaurant_slug = ?", userID, slug).First(&vorhanden).Error
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"ok": true}) // war schon Favorit
		return
	}

	database.DB.Create(&models.Favorite{UserID: userID, RestaurantSlug: slug})
	c.JSON(http.StatusCreated, gin.H{"ok": true})
}

// 🗑️ RemoveFavorite entfernt ein Restaurant aus den Favoriten.
func RemoveFavorite(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	slug := c.Param("slug")

	database.DB.Where("user_id = ? AND restaurant_slug = ?", userID, slug).Delete(&models.Favorite{})
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
