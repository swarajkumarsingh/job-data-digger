package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

type Job struct {
	Title         string   `json:"title"`
	Location      string   `json:"location"`
	Link          string   `json:"link"`
	Qualification []string `json:"qualification"`
}

func splitStringByDotAndInsertIntoArray(arr string) []string {
	var words []string
	for _, word := range strings.Split(arr, ".") {
		if word == "" {continue}
		word = strings.TrimSpace(word)
		words = append(words, word)
	}
	return words
}

func main() {
	c := colly.NewCollector()

	jobs := []Job{}

	c.OnHTML("li.lLd3Je", func(e *colly.HTMLElement) {

		url := fmt.Sprintf("https://www.google.com/about/careers/applications/%s", e.ChildAttr("a.WpHeLc", "href"))

		a := e.ChildText("div.Xsxa1e ul li")

		qualification := splitStringByDotAndInsertIntoArray(a)

		j := Job{
			Title:         e.ChildText("h3.QJPWVe"),
			Location:      e.ChildText("span.r0wTof"),
			Link:          url,
			Qualification: qualification,
		}
		jobs = append(jobs, j)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL)
	})

	c.Visit("https://www.google.com/about/careers/applications/jobs/results/?location=India")

	c.Wait()

	fmt.Println(jobs[0].Title)
	fmt.Println(jobs[0].Location)
	fmt.Println(jobs[0].Link)

	for _, v := range jobs[0].Qualification {
		fmt.Println(v)
	}
}

func main1() {
	//   gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(enableCORS())

	// Load environment variables from a .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//   TODO: ADD DB
	//   TODO: ADD Redis

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "health ok",
		})
	})

	r.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "health ok",
		})
	})

	r.GET("/scrape", controller.Scrape)

	log.Printf("Server Started, version: %s", version)
	r.Run("localhost:8080")
}

// c := colly.NewCollector()

// jobs := []Job{}

// c.OnHTML("h3.QJPWVe", func(e *colly.HTMLElement) {
// 	fmt.Println("Job Title", e.Text)
// })

// c.OnHTML("span.r0wTof", func(e *colly.HTMLElement) {
// 	fmt.Println("Job Location", e.Text)
// })

// c.OnHTML("a.WpHeLc", func(e *colly.HTMLElement) {
// 	url := fmt.Sprintf("https://www.google.com/about/careers/applications/%s", e.Attr("href"))
// 	fmt.Println("Job Share Link", url)
// })

// c.OnRequest(func(r *colly.Request) {
// 	fmt.Println(r.URL)
// })

// c.Visit("https://www.google.com/about/careers/applications/jobs/results/?location=India")
