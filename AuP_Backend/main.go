package main

import (
	"fmt"
	"strings"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // Wichtig: Falls noch nicht geladen, im Terminal: go get github.com/golang-jwt/jwt/v5
	_ "github.com/go-sql-driver/mysql"
	
	"aup-backend/database"
	"aup-backend/handlers"
)

func GetGerichteByRestaurant(c *gin.Context) {
    restaurantID := c.Param("id")

    // 1. Spalte BildURL zur Query hinzugefügt
    query := `SELECT PK_ID_Gericht, Name, Preis, Vegetarisch, Vegan, BildURL 
              FROM Gericht 
              WHERE FK_ID_Restaurant = ?`
    
    rows, err := database.DB.Query(query, restaurantID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gerichte konnten nicht geladen werden"})
        return
    }
    defer rows.Close()

    var gerichte []gin.H
    for rows.Next() {
        var id int
        var name string
        var preis float64
        var vegetarisch []byte
        var vegan []byte
        var bildURL string // 2. Variable für den Pfad hinzugefügt

        // 3. Scan um BildURL erweitert
        if err := rows.Scan(&id, &name, &preis, &vegetarisch, &vegan, &bildURL); err == nil {
            gerichte = append(gerichte, gin.H{
                "id":          id,
                "name":        name,
                "preis":       preis,
                "vegetarisch": len(vegetarisch) > 0 && vegetarisch[0] == 1,
                "vegan":       len(vegan) > 0 && vegan[0] == 1,
                "bild_url":    bildURL, // 4. BildURL in das JSON-Objekt eingefügt
            })
        }
    }
    c.JSON(http.StatusOK, gerichte)
}

func main() {
	// 1. Datenbank-Verbindung beim Starten initialisieren
	database.InitDB() 
	defer database.DB.Close()

	r := gin.Default()
	r.Use(corsMiddleware())
	r.Static("/uploads", "./uploads")

	// =================================================================
	// ECHTE FRONTEND-ROUTEN (ABGESTIMMT AUF DIE CHECKLISTE)
	// =================================================================
	api := r.Group("/api")
	{
		// 1. KONTEN & LOGIN (Gruppe: /api/auth)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.RegisterKunde) // Registrierung speichern
			auth.POST("/login", handlers.LoginKunde)       // Login prüfen + JWT ausgeben
			r.GET("/api/restaurants/:id/gerichte", GetGerichteByRestaurant)
		}

		// 2. WARENKORB & CHECKOUT (Gruppe: /api/cart) -> JETZT MIT AUTH-MIDDLEWARE!
		cart := api.Group("/cart")
		cart.Use(jwtAuthMiddleware()) 
		{
			cart.POST("", handlers.AddToCart)              // POST /api/cart -> Gericht eintüten
			cart.GET("", handlers.GetCart)                // GET /api/cart -> Warenkorb abfragen & berechnen
			cart.POST("/checkout", handlers.CheckoutCart) // POST /api/cart/checkout -> In DB buchen
		}

		// 3. RESTAURANTS & SPEISEKARTEN (Direkt unter /api)
		api.GET("/restaurants", handlers.GetRestaurants) // Daten aus DB ausliefern
		api.GET("/gerichte", handlers.GetGerichte)       // Speisekarten-Daten aus DB ausliefern

		// 4. GUTSCHEINE (Gruppe: /api/gutscheine)
		gutscheine := api.Group("/gutscheine")
		{
			gutscheine.POST("/validate", handlers.ValidateGutschein) // Codes prüfen & Rabatt berechnen
		}

		// 5. BESTELLUNGEN & ZAHLUNG (Gruppe: /api/bestellungen)
		// Simuliert die erfolgreiche Zahlung für den Stripe/PayPal-Flow des Frontends
		api.POST("/bestellungen", func(c *gin.Context) {
			c.JSON(201, gin.H{
				"status":         "success",
				"message":        "Zahlung erfolgreich über Drittanbieter abgewickelt!",
				"transaction_id": "ch_fake_3N92xLzPq91",
			})
		})
	}

	r.Run(":8080")
}

// CORS-Middleware, damit ein Frontend (z.B. auf localhost:3000) mit dem Backend reden darf
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// JWT-Middleware, die das Token validiert und die kunde_id für die Handlers bereitstellt
func jwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Header fehlt"})
			c.Abort()
			return
		}

		// "Bearer <token>" splitten
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Ungültiges Header-Format (Format: Bearer <token>)"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Token parsen (Hier musst du deinen echten geheimen JWT-Schlüssel eintragen!)
		// HINWEIS: Falls dein Key ein String ist, passe []byte("dein_secret") an.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unerwartete Signaturmethode: %v", token.Header["alg"])
			}
			return []byte("dein_super_geheimes_zufaelliges_geheimnis"), nil 
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token ist ungültig oder abgelaufen"})
			c.Abort()
			return
		}

		// Claims auslesen und kunde_id in den Context pushen
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if kundeID, ok := claims["kunde_id"]; ok {
				// Speichert den Key exakt als "kunde_id", so wie handlers/cart.go es erwartet
				c.Set("kunde_id", int(kundeID.(float64))) 
			}
		}

		c.Next()
	}
}