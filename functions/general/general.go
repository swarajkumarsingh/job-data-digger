package general

import "github.com/swarajkumarsingh/job-data-digger/model"

func IsModelEmpty(jobs []model.Job) bool {
	return len(jobs) == 0
}