// controller package
package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/swarajkumarsingh/job-data-digger/errorHandler"
	"github.com/swarajkumarsingh/job-data-digger/functions/logger"
	"github.com/swarajkumarsingh/job-data-digger/model"
)

func Scrape(r *gin.Context) {
	defer errorHandler.Recovery(r, http.StatusConflict)

	c := colly.NewCollector()
	jobs := []model.Job{}

	c.OnHTML("li.lLd3Je", func(e *colly.HTMLElement) {
		url := GetFullUrl(e)
		title := e.ChildText("h3.QJPWVe")
		location := e.ChildText("span.r0wTof")
		arr := e.ChildText("div.Xsxa1e ul li")
		qualification := SplitStringByDotAndInsertIntoArray(arr)
		j := model.Job{
			Title:         title,
			Location:      location,
			Link:          url,
			Qualification: qualification,
		}
		jobs = append(jobs, j)
	})

	c.OnError(func(_ *colly.Response, err error) {
		logger.WithRequest(r).Errorln("Error while fetch results...")
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Scraping url: ", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {})
	c.OnResponse(func(r *colly.Response) {})

	c.Visit("https://www.google.com/about/careers/applications/jobs/results/?location=India")
	c.Wait()

	// TODO: Add data in redis(ttl:24hr)
	r.JSON(http.StatusOK, map[string]interface{}{
		"error": false,
		"data":  jobs,
	})
}
