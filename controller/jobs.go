package controller

import (
	"context"
	"net/http"
	"spurtcms-graphql/graph/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JobsList(db *gorm.DB, ctx context.Context, limit int, offset int, filter *model.JobFilter, sort *model.JobSort) (*model.JobsList, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	var jobs []model.Job
	var count int64

	listQuery := db.Debug().Table("tbl_jobs").Select("*")

	var (
		jobTitle, jobLocation, jobType, skill, keyWord string
		minimumYears, maximumYears                     int
	)

	if filter != nil {

		if filter.JobTitle.IsSet() {  

			jobTitle = *filter.JobTitle.Value()
		}

		if filter.KeyWord.IsSet() {

			keyWord = *filter.KeyWord.Value()
		}

		if filter.JobLocation.IsSet() {

			jobLocation = *filter.JobLocation.Value()
		}

		if filter.JobType.IsSet() {

			jobType = *filter.JobType.Value()
		}

		if filter.Skill.IsSet() {

			skill = *filter.Skill.Value()
		}

		if filter.MaximumYears.IsSet() {

			maximumYears = *filter.MaximumYears.Value()
		}

		if filter.MinimumYears.IsSet() {

			minimumYears = *filter.MinimumYears.Value()
		}
	}

	if jobTitle != "" {

		listQuery = listQuery.Where("job_title = ?", jobTitle)
	}

	if keyWord != "" {

		listQuery = listQuery.Where("LOWER(TRIM(job_title)) like LOWER(TRIM(?))", "%"+keyWord+"%")
	}

	if jobLocation != "" {

		listQuery = listQuery.Where("job_location = ?", jobLocation)
	}

	if jobType != "" {

		listQuery = listQuery.Where("job_type = ?", jobType)
	}

	if skill != "" {

		listQuery = listQuery.Where("skill = ?", skill)
	}

	if minimumYears != 0 {

		listQuery = listQuery.Where("minimum_years <= ?", minimumYears)
	}

	if maximumYears != 0 {

		listQuery = listQuery.Where("maximum_years <= ?", maximumYears)
	}

	// if sort != nil {

	if sort.Salary.Value() != nil && *sort.Salary.Value() != -1 {

		if *sort.Salary.Value() == 1 {

			listQuery = listQuery.Order("tbl_jobs.salary desc")

		} else if *sort.Salary.Value() == 0 {

			listQuery = listQuery.Order("tbl_jobs.salary")
		}

	} else if sort.PostedDate.Value() != nil && *sort.PostedDate.Value() != -1 {

		if *sort.PostedDate.Value() == 1 {

			listQuery = listQuery.Order("tbl_jobs.posted_date desc")

		} else if *sort.PostedDate.Value() == 0 {

			listQuery = listQuery.Order("tbl_jobs.posted_date")
		}

	} else {

		listQuery = listQuery.Order("tbl_jobs.id desc")
	}

	listQuery = listQuery.Limit(limit).Offset(offset).Find(&jobs)

	if listQuery.Error != nil {

		c.AbortWithError(http.StatusInternalServerError, listQuery.Error)

		return &model.JobsList{}, listQuery.Error
	}

	countQuery := listQuery.Count(&count)
	
	if countQuery.Error != nil {

		c.AbortWithError(http.StatusInternalServerError, countQuery.Error)

		return &model.JobsList{}, countQuery.Error
	}

	return &model.JobsList{Jobs: jobs, Count: int(count)}, nil
}
