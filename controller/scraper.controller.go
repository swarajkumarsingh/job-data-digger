// controller package
package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/swarajkumarsingh/job-data-digger/errorHandler"
	"github.com/swarajkumarsingh/job-data-digger/functions/logger"
	redisUtils "github.com/swarajkumarsingh/job-data-digger/infra/redis"
	"github.com/swarajkumarsingh/job-data-digger/model"
)

func Scrape(r *gin.Context) {
	defer errorHandler.Recovery(r, http.StatusConflict)

	listData, _ := redisUtils.Rdb.LRange(context.Background(), "myList", 0, -1).Result()
	if len(listData) != 0 {
		var jobs []model.Job

		if err := json.Unmarshal([]byte(listData[0]), &jobs); err != nil {
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

	// err := redisUtils.SetStructWithLongTTL("google-data", jobs[1])
	// err := redisUtils.Set("key", "simu dude", time.Hour*24)

	jsonData, err := json.Marshal(jobs)
	if err != nil {
		r.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON data"})
		return
	}

	// Use the Redis LPUSH command to add data to the list
	err = redisUtils.Rdb.LPush(context.Background(), "myList", jsonData).Err()
	if err != nil {
		r.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add data to the list in Redis"})
		return
	}

	r.JSON(http.StatusOK, map[string]interface{}{
		"error": false,
		"data":  jobs,
	})
}
