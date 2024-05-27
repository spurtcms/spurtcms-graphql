package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"spurtcms-graphql/graph/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JobsList(db *gorm.DB, ctx context.Context, limit int, offset int, filter *model.JobFilter) (*model.JobsList, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	var jobs []model.Job

	var count int64

	listQuery := db.Debug().Table("tbl_jobs").Select("tbl_jobs.*,tbl_categories.id as CatId,tbl_categories.category_name,tbl_categories.category_slug").Joins("inner join tbl_categories on tbl_jobs.categories_id = tbl_categories.id").Where("tbl_jobs.is_deleted = 0 AND tbl_jobs.status = 1").Preload("Category")

	var (
		jobTitle, jobLocation, skill, keyWord, categorySlug, datePosted string

		minimumYears, maximumYears, categoryId                          int
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

		if filter.CategoryID.IsSet() {

			categoryId = *filter.CategoryID.Value()
		}

		if filter.CategorySlug.IsSet() {

			categorySlug = *filter.CategorySlug.Value()
		}

		if filter.MaximumYears.IsSet() {

			maximumYears = *filter.MaximumYears.Value()
		}

		if filter.MinimumYears.IsSet() {

			minimumYears = *filter.MinimumYears.Value()
		}

		if filter.DatePosted.IsSet() {

			datePosted = *filter.DatePosted.Value()
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

	if categorySlug != "" {

		listQuery = listQuery.Where("tbl_categories.category_slug = ?", categorySlug)
	}

	if categoryId != 0 {

		listQuery = listQuery.Where("categories_id = ?", categoryId)
	}

	if skill != "" {
		listQuery = listQuery.Where("skill = ?", skill)
	}

	if minimumYears != 0 && maximumYears != 0 {

		listQuery = listQuery.Where("minimum_years >= ? and maximum_years <= ?", minimumYears, maximumYears)

	} else if minimumYears != 0 {

		listQuery = listQuery.Where("minimum_years >= ?", minimumYears)

	} else if maximumYears != 0 {

		listQuery = listQuery.Where("maximum_years <= ?", maximumYears)
	}

	if datePosted != "" {

		var startDate,endDate time.Time

		var currentDate = time.Now().Local()

		if datePosted == "This Week" {

			currentDay := time.Now().Local().Weekday().String()

			switch currentDay {

			case "Monday":
				startDate = currentDate
				endDate = currentDate.AddDate(0, 0, 6)

			case "Tuesday":
				startDate = currentDate.AddDate(0, 0, -1)
				endDate = currentDate.AddDate(0, 0, 5)

			case "Wednesday":
				startDate = currentDate.AddDate(0, 0, -2)
				endDate = currentDate.AddDate(0, 0, 4)

			case "Thursday":
				startDate = currentDate.AddDate(0, 0, -3)
				endDate = currentDate.AddDate(0, 0, 3)

			case "Friday":
				startDate = currentDate.AddDate(0, 0, -4)
				endDate = currentDate.AddDate(0, 0, 2)

			case "Saturday":
				startDate = currentDate.AddDate(0, 0, -5)
				endDate = currentDate.AddDate(0, 0, 1)

			case "Sunday":
				startDate = currentDate.AddDate(0, 0, -6)
				endDate = currentDate
			}

		}

		if datePosted == "This Month" {

			startDate = time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())
			firstDayOfNxtMnth := startDate.AddDate(0, 1, 0)
			endDate = firstDayOfNxtMnth.Add(-time.Second)
		}

		if datePosted == "This Year" {

			startDate = time.Date(currentDate.Year(), time.January, 1, 0, 0, 0, 0, currentDate.Location())
			startofNxtYear := startDate.AddDate(1, 0, 0)
			endDate = startofNxtYear.Add(-time.Second)
		}

		if datePosted == "Today" {

			startDate = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
			nxtDay := startDate.AddDate(0, 0, 1)
			endDate = nxtDay.Add(-time.Second)
		}

		listQuery = listQuery.Where("posted_date between (?) and (?)", startDate, endDate)

	}

	listQuery = listQuery.Limit(limit).Offset(offset).Order("tbl_jobs.id desc").Find(&jobs)

	if listQuery.Error != nil {

		c.AbortWithError(http.StatusInternalServerError, listQuery.Error)

		return &model.JobsList{}, listQuery.Error
	}

	if len(jobs) <= 0 {

		c.AbortWithError(500, ErrRecordNotFound)

		return nil, ErrRecordNotFound
	}

	countQuery := listQuery.Count(&count)

	if countQuery.Error != nil {

		c.AbortWithError(http.StatusInternalServerError, countQuery.Error)

		return &model.JobsList{}, countQuery.Error
	}

	return &model.JobsList{Jobs: jobs, Count: int(count)}, nil
}

func JobDetail(db *gorm.DB, ctx context.Context, id *int, jobSlug *string) (*model.Job, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	var jobDetail *model.Job

	query := db.Debug().Table("tbl_jobs").Select("tbl_jobs.*,tbl_categories.id as CatId,tbl_categories.category_name,tbl_categories.category_slug").Joins("inner join tbl_categories on tbl_jobs.categories_id = tbl_categories.id").Where("tbl_jobs.is_deleted = 0").Preload("Category")

	if id != nil {

		query = query.Where("tbl_jobs.id = ?", id)

	} else if jobSlug != nil {

		query = query.Where("tbl_jobs.job_slug = ? ", jobSlug)

	}

	query = query.Find(&jobDetail)

	if query.Error != nil {

		c.AbortWithError(http.StatusInternalServerError, query.Error)

		return &model.Job{}, query.Error
	}

	return jobDetail, nil
}

func JobApplication(db *gorm.DB, ctx context.Context, applicationDetails model.ApplicationInput) (bool, error) {

	applicationInfo := applicationDetails

	applicantImage := applicationInfo.ApplicantImage

	resume := applicationInfo.Resume

	ImgBase64Data, err := io.ReadAll(applicantImage.File)

	resumeBase64Data, err := io.ReadAll(resume.File)

	if err != nil {

		return false, err
	}

	ImgTargetPath := filepath.Join("uploads/images", applicantImage.Filename)

	// Create the target file
	out, err := os.Create(ImgTargetPath)
	if err != nil {
		return false, errors.New("failed to create file")
	}
	defer out.Close()

	// Copy the uploaded file to the target file
	err = os.WriteFile(fmt.Sprintf("uploads/images/%v", applicantImage.Filename), ImgBase64Data, os.ModePerm)

	if err != nil {

		return false, errors.New("failed to copy file")
	}

	resumeTargetPath := filepath.Join("uploads/resumes", resume.Filename)

	out, err = os.Create(resumeTargetPath)

	if err != nil {

		return false, errors.New("failed to create file")
	}

	defer out.Close()

	err = os.WriteFile(fmt.Sprintf("uploads/resumes/%v", resume.Filename), resumeBase64Data, os.ModePerm)

	if err != nil {

		return false, errors.New("failed to copy file")
	}

	var newMember model.Member

	newMember.FirstName = applicationInfo.Name

	newMember.Email = applicationInfo.EmailID

	newMember.MobileNo = strconv.Itoa(applicationInfo.MobileNo)

	newMember.IsActive = 1

	newMember.MemberGroupID = 1

	newMember.ProfileImagePath = fmt.Sprintf("uploads/images/%v", applicantImage.Filename)

	newMember.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	result := db.Table("tbl_members").Create(&newMember)

	if result.Error != nil {
		
		return false, result.Error
	}

	return true, nil

}
