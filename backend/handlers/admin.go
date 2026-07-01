package handlers

import (
	"net/http"

	"lieferino-backend/database"
	"lieferino-backend/models"

	"github.com/gin-gonic/gin"
)

// 🛠️ ADMIN-ENDPUNKTE (alle nur mit Admin-Rechten erreichbar)

// 📊 AdminStats liefert Kennzahlen fürs Dashboard.
func AdminStats(c *gin.Context) {
	var umsatz float64
	database.DB.Model(&models.Order{}).Select("COALESCE(SUM(summe),0)").Scan(&umsatz)

	var artikel int64
	database.DB.Model(&models.OrderItem{}).Select("COALESCE(SUM(menge),0)").Scan(&artikel)

	var bestellungen, nutzer, aktiveRestaurants int64
	database.DB.Model(&models.Order{}).Count(&bestellungen)
	database.DB.Model(&models.User{}).Count(&nutzer)
	database.DB.Model(&models.Restaurant{}).Where("aktiv = ?", true).Count(&aktiveRestaurants)

	c.JSON(http.StatusOK, gin.H{
		"umsatz":            umsatz,
		"verkaufteArtikel":  artikel,
		"bestellungen":      bestellungen,
		"nutzer":            nutzer,
		"aktiveRestaurants": aktiveRestaurants,
	})
}

// 👥 AdminUsers listet alle Nutzer (ohne Passwort-Hash).
func AdminUsers(c *gin.Context) {
	var users []models.User
	database.DB.Order("id asc").Find(&users)
	c.JSON(http.StatusOK, users)
}

type adminSperrenInput struct {
	Gesperrt bool `json:"gesperrt"`
}

// 🚫 AdminUserSperren sperrt/entsperrt einen Nutzer.
func AdminUserSperren(c *gin.Context) {
	id := c.Param("id")
	var in adminSperrenInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": err.Error()})
		return
	}
	if err := database.DB.Model(&models.User{}).Where("id = ?", id).
		Update("gesperrt", in.Gesperrt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"fehler": "Konnte nicht gespeichert werden"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// 🔓 AdminMfaReset setzt die MFA eines Nutzers zurück (muss neu eingerichtet werden).
func AdminMfaReset(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Model(&models.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{"mfa_aktiv": false, "mfa_secret": ""}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"fehler": "Konnte nicht gespeichert werden"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type adminAktivInput struct {
	Aktiv bool `json:"aktiv"`
}

// 🍽️ AdminRestaurantAktiv aktiviert/deaktiviert ein Restaurant (per Slug).
func AdminRestaurantAktiv(c *gin.Context) {
	slug := c.Param("slug")
	var in adminAktivInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": err.Error()})
		return
	}
	if err := database.DB.Model(&models.Restaurant{}).Where("slug = ?", slug).
		Update("aktiv", in.Aktiv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"fehler": "Konnte nicht gespeichert werden"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ⭐ AdminReviews listet alle Bewertungen (zum Moderieren).
func AdminReviews(c *gin.Context) {
	var reviews []models.Review
	database.DB.Order("created_at desc").Find(&reviews)
	c.JSON(http.StatusOK, reviews)
}

// 🗑️ AdminReviewLoeschen löscht eine Bewertung (per ID).
func AdminReviewLoeschen(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Review{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"fehler": "Konnte nicht gelöscht werden"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
