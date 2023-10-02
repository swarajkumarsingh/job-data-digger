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

func GetGoogleJobs() {
	
}

func OnRequest(r *colly.Request)   {}
func OnScraped(r *colly.Response)  {}
func OnResponse(r *colly.Response) {}
func OnError(_ *colly.Response, err error) {
	logger.Log.Errorln("Error while fetch results google jobs...")
}

func GetScrapeDataList() ([]string, error) {
	return redisUtils.Rdb.LRange(context.Background(), constants.REDIS_JOBS_LIST_KEY, 0, -1).Result()
}

func GetScrapeData(r *gin.Context, listData string) ([]model.Job, error) {
	var jobs []model.Job

	if err := json.Unmarshal([]byte(listData), &jobs); err != nil {
		logger.WithRequest(r).Panicln(err)
	}

	return jobs, nil
}

func AddedScrapeDataToRedis(r *gin.Context, jobs []model.Job) error {
	jsonData, err := json.Marshal(jobs)
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to marshal JSON data, cannot add data in cache: %s", err))
		return err
	}

	err = redisUtils.Rdb.LPush(context.Background(), "myList", jsonData).Err()
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to add data to the list in Redis: %s", err))
		return err
	}

	err = redisUtils.Rdb.Expire(context.Background(), "myList", time.Hour*24).Err()
	if err != nil {
		err = errors.New(fmt.Sprintf("Error while setting ttl: %s", err))
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
