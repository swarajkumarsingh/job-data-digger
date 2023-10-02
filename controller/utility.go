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
	"github.com/swarajkumarsingh/job-data-digger/functions/logger"
	redisUtils "github.com/swarajkumarsingh/job-data-digger/infra/redis"
	"github.com/swarajkumarsingh/job-data-digger/model"
)

func IsEmptyCache() bool {
	cacheData := GetScrapeDataListFromCache()
	return len(cacheData) != 0
}

func GetAllJobs(r *gin.Context) ([]model.Job, error) {
	// TODO: Add more job providers(microsoft, meta, amazon)
	return GetGoogleJobs(r)
}

func GetGoogleJobs(r *gin.Context) ([]model.Job, error) {
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
		logger.WithRequest(r).Errorln("Error while fetch results google jobs..")
	})
	c.OnRequest(OnRequest)
	c.OnScraped(OnScraped)
	c.OnResponse(OnResponse)

	c.Visit("https://www.google.com/about/careers/applications/jobs/results/?location=India")
	c.Wait()

	return jobs, nil
}

func OnRequest(r *colly.Request)   {}
func OnScraped(r *colly.Response)  {}
func OnResponse(r *colly.Response) {}

func GetScrapeDataListFromCache() []string {
	val, err := redisUtils.Rdb.LRange(context.Background(), constants.REDIS_JOBS_LIST_KEY, 0, -1).Result()
	if err != nil {
		return []string{}
	}
	return val
}

func GetScrapeDataFromCache(r *gin.Context) []model.Job {
	var jobs []model.Job
	rawData := GetScrapeDataListFromCache()

	if err := json.Unmarshal([]byte(rawData[0]), &jobs); err != nil {
		logger.WithRequest(r).Panicln(err)
	}

	return jobs
}

func AddedScrapeDataToRedis(r *gin.Context, jobs []model.Job) error {
	jsonData, err := json.Marshal(jobs)
	if err != nil {
		err = errors.New("Failed to marshal JSON data, cannot add data in cache: " + err.Error())
		return err
	}

	err = redisUtils.Rdb.LPush(context.Background(), "myList", jsonData).Err()
	if err != nil {
		err = errors.New("Failed to add data to the list in Redis: " + err.Error())
		return err
	}

	err = redisUtils.Rdb.Expire(context.Background(), "myList", time.Hour*24).Err()
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
