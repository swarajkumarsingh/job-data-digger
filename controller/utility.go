package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/swarajkumarsingh/job-data-digger/constants"
	"github.com/swarajkumarsingh/job-data-digger/functions/general"
	"github.com/swarajkumarsingh/job-data-digger/functions/logger"
	redisUtils "github.com/swarajkumarsingh/job-data-digger/infra/redis"
	"github.com/swarajkumarsingh/job-data-digger/model"
)

func IsCacheDataPresent() bool {
	cacheData := getScrapeDataListFromCache()
	return len(cacheData) != 0
}

func GetAllJobs(r *gin.Context) ([]model.Job, error) {

	jobs := []model.Job{}

	data, err := googleJobsProvider(r)
	if err != nil {
		logger.WithRequest(r).Errorln("Error while fetching data")
	}

	jobs = append(jobs, data...)
	if general.IsModelEmpty(jobs) {
		return jobs, errors.New("data not found")
	}

	return jobs, nil
}

func googleJobsProvider(r *gin.Context) ([]model.Job, error) {
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

	c.OnError(func(c *colly.Response, err error) {
		onError(r, c)
	})
	c.OnRequest(onRequest)
	c.OnScraped(onScraped)
	c.OnResponse(onResponse)

	c.Visit(constants.GOOGLE_CAREER_PAGE_URL)
	c.Wait()

	return jobs, nil
}

func onRequest(r *colly.Request) {
	logger.Log.Debug()
}
func onScraped(r *colly.Response) {}
func onError(r *gin.Context, _ *colly.Response) {
	logger.WithRequest(r).Panicln("Error while fetch results google jobs..")

}
func onResponse(r *colly.Response) {
	logger.Log.Errorln("Received response from ", r.Request.URL)
}

func getScrapeDataListFromCache() []string {
	val, err := redisUtils.Rdb.LRange(context.Background(), constants.REDIS_SCRAPE_DATA_KEY, 0, -1).Result()
	if err != nil {
		return []string{}
	}
	return val
}

func GetScrapeDataFromCache(r *gin.Context) []model.Job {
	var jobs []model.Job
	rawData := getScrapeDataListFromCache()

	if err := json.Unmarshal([]byte(rawData[0]), &jobs); err != nil {
		logger.WithRequest(r).Panicln(err)
	}

	return jobs
}

func AddScrapeDataToRedis(r *gin.Context, jobs []model.Job) error {
	jsonData, err := json.Marshal(jobs)
	if err != nil {
		err = errors.New("Failed to marshal JSON data, cannot add data in cache: " + err.Error())
		return err
	}

	err = redisUtils.Rdb.LPush(context.Background(), constants.REDIS_SCRAPE_DATA_KEY, jsonData).Err()
	if err != nil {
		err = errors.New("Failed to add data to the list in Redis: " + err.Error())
		return err
	}

	err = redisUtils.Rdb.Expire(context.Background(), constants.REDIS_SCRAPE_DATA_KEY, time.Hour*24).Err()
	if err != nil {
		err = errors.New("Error while setting ttl: " + err.Error())
		return err
	}

	return nil
}

func SplitStringByDotAndInsertIntoArray(arr string) []string {
	var words []string
	for _, word := range strings.SplitN(arr, ".", 1) {
		if word == "" {
			continue
		}
		word = strings.TrimSpace(word)
		words = append(words, word)
	}
	return words
}

func GetFullUrl(e *colly.HTMLElement) string {
	return fmt.Sprintf("https://www.google.com/about/careers/applications/%s", e.ChildAttr("a.WpHeLc", "href"))
}
