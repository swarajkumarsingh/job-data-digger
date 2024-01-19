package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swarajkumarsingh/job-data-digger/conf"
	"github.com/swarajkumarsingh/job-data-digger/controller"
)

var version string = "1.0"

func enableCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Api-Key, token, User-Agent, Referer")
		c.Writer.Header().Set("AllowCredentials", "true")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		if c.Request.Method == "OPTIONS" {
			return
		}

		c.Next()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}

	fmt.Println("ENV: ", conf.ENV)

	if conf.ENV == conf.ENV_PROD {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.Use(gin.Recovery())
	r.Use(enableCORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "health ok",
		})
	})

	r.POST("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "health ok",
		})
	})

	r.GET("/scrape", controller.Scrape)

	log.Printf("Server Started, version: %s", version)
	r.Run(":8080")
}
