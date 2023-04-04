package api

import (
	"ip_manager/config"
	"ip_manager/db"
	"ip_manager/types"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func LoadGin(config *config.Config, db *db.MySQLDatabase) {
	// Create a channel to signal when the database has finished loading
	done := make(chan struct{})

	svc := types.NewService()

	go loadService(svc, db, config, done)

	r := gin.Default()
	r.Use(authMiddleware(config))

	<-done

	loadReadDef(r, svc, db)
	loadWriteDef(r, svc, db)
	r.Run(config.Listen)
}

func loadService(svc *types.Service, db *db.MySQLDatabase, cfg *config.Config, done chan struct{}) {
	// Set up a ticker to load the tables on a regular interval
	ticker := time.NewTicker(time.Duration(cfg.DbSaveInterval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			svc.Delete()
			err := db.LoadService(svc)
			if err != nil {
				log.Printf("Error loading database: %v", err)
			}
		case <-done:
			return
		}
	}
}

func authMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header from the request
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			// If the header is missing, return a 401 Unauthorized error
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Extract the token from the header string
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			// If the header format is invalid, return a 401 Unauthorized error
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := authHeaderParts[1]

		// Check if the token is present in the authorized tokens list
		for _, authorizedToken := range cfg.AuthorizedTokens {
			if tokenString == authorizedToken {
				// Call the next middleware function in the chain
				c.Next()
				return
			}
		}

		// If the token is not authorized, return a 401 Unauthorized error
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
