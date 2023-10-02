// controller package
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swarajkumarsingh/job-data-digger/errorHandler"
	"github.com/swarajkumarsingh/job-data-digger/functions/logger"
	"github.com/swarajkumarsingh/job-data-digger/errorCodes"
)

func Scrape(r *gin.Context) {
	defer errorHandler.Recovery(r, http.StatusConflict)

	if IsEmptyCache() {
		jobs := GetScrapeDataFromCache(r)
		r.JSON(errorcodes.STATUS_OK, gin.H{
			"error": false,
			"data":  jobs,
		})
		return
	}

	jobs, err := GetAllJobs(r)
	if err != nil {
		logger.WithRequest(r).Errorln(err)
	}

	err = AddedScrapeDataToRedis(r, jobs)
	if err != nil {
		logger.WithRequest(r).Errorln(err)
	}

	r.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  jobs,
	})
}
