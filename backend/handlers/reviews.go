package handlers

import (
	"net/http"

	"lieferino-backend/database"
	"lieferino-backend/models"

	"github.com/gin-gonic/gin"
)

// ⭐ ListReviews liefert alle Bewertungen eines Restaurants (öffentlich).
func ListReviews(c *gin.Context) {
	slug := c.Param("slug")

	var reviews []models.Review
	database.DB.Where("restaurant_slug = ?", slug).
		Order("created_at desc").
		Find(&reviews)

	c.JSON(http.StatusOK, reviews)
}

type reviewInput struct {
	Name           string `json:"name" binding:"required"`
	Sterne         int    `json:"sterne" binding:"required,min=1,max=5"`
	Text           string `json:"text" binding:"required"`
	RestaurantName string `json:"restaurantName" binding:"required"`
}

// ✍️ AddReview legt eine Bewertung an – aber NUR, wenn der Nutzer dort schon
// etwas bestellt hat (serverseitig erzwungen).
func AddReview(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	slug := c.Param("slug")

	var in reviewInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": err.Error()})
		return
	}

	// Hat der Nutzer bei diesem Restaurant (per Name) schon bestellt?
	var anzahl int64
	database.DB.Model(&models.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.user_id = ? AND order_items.restaurant = ?", userID, in.RestaurantName).
		Count(&anzahl)

	if anzahl == 0 {
		c.JSON(http.StatusForbidden, gin.H{"fehler": "Du kannst erst bewerten, wenn du hier bestellt hast"})
		return
	}

	review := models.Review{
		UserID:         userID,
		RestaurantSlug: slug,
		RestaurantName: in.RestaurantName,
		Name:           in.Name,
		Sterne:         in.Sterne,
		Text:           in.Text,
	}
	if err := database.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"fehler": "Bewertung konnte nicht gespeichert werden"})
		return
	}
	c.JSON(http.StatusCreated, review)
}
