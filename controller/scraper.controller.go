// controller package
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/swarajkumarsingh/job-data-digger/errorHandler"
	"github.com/swarajkumarsingh/job-data-digger/functions/logger"
	"github.com/swarajkumarsingh/job-data-digger/model"
)

func Scrape(r *gin.Context) {
	defer errorHandler.Recovery(r, http.StatusConflict)

	listData, _ := GetScrapeDataList()
	if len(listData) != 0 {
		jobs, err := GetScrapeData(r, listData[0])
		if err != nil {
			logger.WithRequest(r).Panicln(err)
		}

		r.JSON(http.StatusOK, map[string]interface{}{
			"error": false,
			"data":  jobs,
		})
		return
	}

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

	c.OnError(OnError)
	c.OnRequest(OnRequest)
	c.OnScraped(OnScraped)
	c.OnResponse(OnResponse)

	c.Visit("https://www.google.com/about/careers/applications/jobs/results/?location=India")
	c.Wait()

	err := AddedScrapeDataToRedis(r, jobs)
	if err != nil {
		logger.WithRequest(r).Errorln(err)
	}

	r.JSON(http.StatusOK, map[string]interface{}{
		"error": false,
		"data":  jobs,
	})
}
