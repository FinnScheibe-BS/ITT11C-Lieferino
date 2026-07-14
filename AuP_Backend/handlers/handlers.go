package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetProdukte(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hier kommen bald die Produkte"})
}

func RegisterUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Registrierung folgt"})
}

func LoginUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login folgt"})
}
