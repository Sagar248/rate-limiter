package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"rate-limiter/internal/limiter"
	"rate-limiter/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// Set Gin mode from environment
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}

	// Initialize Redis
	rdb := storage.NewRedisClient()

	// Validate Redis connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("failed to connect to Redis: " + err.Error())
	}

	// Initialize limiter
	l := limiter.NewLimiter(rdb)

	// Setup router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // allow all (fine for demo)
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.SetTrustedProxies(nil)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Rate limited API
	r.GET("/api", func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "default-user"
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		allowed, remaining := l.Allow(ctx, userID)

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "allowed",
			"remaining": remaining,
		})
	})

	// Port configuration (important for Railway/Render)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	r.Run(":" + port)
}