package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"spurtcms-graphql/graph/model"
	"spurtcms-graphql/storage"
	"strconv"
	"strings"
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

		minimumYears, maximumYears, categoryId int
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

		listQuery = listQuery.Where("LOWER(TRIM(job_location)) = LOWER(TRIM(?))", jobLocation)
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

		var startDate, endDate time.Time

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

		c.AbortWithError(http.StatusInternalServerError, ErrRecordNotFound)

		return &model.JobsList{}, ErrRecordNotFound
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

func ApplicantDetails(db *gorm.DB, ctx context.Context) (*model.ApplicantDetails, error) {

	c, ok := ctx.Value(ContextKey).(*gin.Context)

	if !ok {
		ErrorLog.Printf("Applicant Details Context error: %v", ok)
	}

	memberId := c.GetInt("memberid")

	if memberId == 0 {

		ErrorLog.Printf("Applicant Details context error: %s", ErrUnauthorizedAccess)

		c.AbortWithError(http.StatusUnauthorized, ErrUnauthorizedAccess)
	}

	var (
		applicantDetails      model.ApplicantDetails
		imagePath, resumePath string
	)

	if err := db.Debug().Table("tbl_jobs_applicants").Where("is_deleted = 0 and status = 1 and member_id = ?", memberId).First(&applicantDetails).Error; err != nil {

		ErrorLog.Printf("%s: %s", ErrApplicantNotFound, err)

		c.AbortWithError(http.StatusUnprocessableEntity, err)

		return &model.ApplicantDetails{}, err
	}

	fmt.Println("applicantImage", *applicantDetails.ImagePath)

	if *applicantDetails.StorageType == "aws" && applicantDetails.ImagePath != nil && applicantDetails.StorageType != nil && *applicantDetails.ImagePath != "" {

		imagePath = "image-resize?name=" + *applicantDetails.ImagePath

	} else if *applicantDetails.StorageType == "local" && applicantDetails.ImagePath != nil && applicantDetails.StorageType != nil && *applicantDetails.ImagePath != "" {

		imagePath = *applicantDetails.ImagePath

	} else {

		imagePath = ""
	}

	applicantDetails.ImagePath = &imagePath

	if *applicantDetails.StorageType == "aws" && applicantDetails.ResumePath != nil && applicantDetails.StorageType != nil && *applicantDetails.ResumePath != "" {

		resumePath = "image-resize?name=" + *applicantDetails.ResumePath

	} else if *applicantDetails.StorageType == "local" && applicantDetails.ResumePath != nil && applicantDetails.StorageType != nil && *applicantDetails.ResumePath != "" {

		resumePath = *applicantDetails.ResumePath
	} else {

		resumePath = ""
	}

	applicantDetails.ResumePath = &resumePath

	return &applicantDetails, nil
}

func JobApplication(db *gorm.DB, ctx context.Context, applicationDetails model.ApplicationInput) (bool, error) {

	c, ok := ctx.Value(ContextKey).(*gin.Context)

	if !ok {

		ErrorLog.Printf("job Application context error: %v", ok)
	}

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		ErrorLog.Printf("job Application context error: %s", ErrUnauthorizedAccess)

		c.AbortWithError(http.StatusUnauthorized, ErrUnauthorizedAccess)

		return false, ErrUnauthorizedAccess

	}

	var (
		applicantDetails                                                    model.ApplicantDetails
		result                                                              *gorm.DB
		applicationData                                                     model.ApplicantDetails
		imageName, imagePath, resumeName, resumePath, base64Data, extension string
		storageType                                                         StorageType
		err                                                                 error
		isValidBase64                                                       bool
	)

	storageType, err = GetStorageType(db)

	if err != nil {

		return false, err
	}

	result = db.Debug().Table("tbl_jobs_applicants").Where("is_deleted = 0 and member_id = ? and status = 1", memberid).First(&applicantDetails)
	if result.Error != nil {

		return false, result.Error
	}

	applicationData.JobID = &applicationDetails.JobID

	applicationData.ApplicantID = applicantDetails.ID

	applicationData.CreatedBy = applicantDetails.ID

	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	applicationData.CreatedOn = &currentTime

	applicationData.Name = &applicationDetails.Name

	applicationData.EmailID = &applicationDetails.EmailID

	applicationData.MobileNo = &applicationDetails.MobileNo

	if applicationDetails.JobType.IsSet() && applicationDetails.JobType.Value() != nil {

		applicationData.JobType = applicationDetails.JobType.Value()
	}

	applicationData.Location = &applicationDetails.Location

	applicationData.Education = &applicationDetails.Education

	applicationData.Graduation = &applicationDetails.Graduation

	if applicationDetails.CompanyName.IsSet() && applicationDetails.CompanyName.Value() != nil {

		applicationData.CompanyName = applicationDetails.CompanyName.Value()
	}

	applicationData.Experience = &applicationDetails.Experience

	applicationData.Skills = &applicationDetails.Skills

	if applicationDetails.Image != "" {

		isValidBase64, base64Data, extension = IsValidBase64(applicationDetails.Image)

		if isValidBase64 && base64Data != "" {

			rand_num := strconv.Itoa(int(time.Now().Unix()))

			imageName = "IMG-" + rand_num + "." + extension

			if storageType.SelectedType == "aws" {

				fmt.Printf("aws-S3 storage selected\n")

				imagePath = "member/" + imageName

				err = storage.UploadFileS3(storageType.Aws, nil, base64Data, imagePath)
				if err != nil {

					fmt.Printf("image upload failed %v\n", err)

					return false, ErrUpload

				}

			} else if storageType.SelectedType == "azure" {

				fmt.Printf("azure storage selected")

			} else if storageType.SelectedType == "drive" {

				fmt.Println("drive storage selected")
			}
		} else if strings.Contains(applicationDetails.Image, "image-resize?name") {

			imagePath = strings.ReplaceAll(applicationDetails.Image, "image-resize?name=", "")

		} else {

			ErrorLog.Printf("%v", "illegal base64 data")

			return false, errors.New("illegal base64 data ")

		}

		applicationData.ImagePath = &imagePath

		applicationData.Image = &imageName

	}

	isDeleted := 0

	applicationData.IsDeleted = &isDeleted

	if applicationDetails.CurrentSalary.IsSet() && applicationDetails.CurrentSalary.Value() != nil {

		applicationData.CurrentSalary = applicationDetails.CurrentSalary.Value()
	}

	if applicationDetails.ExpectedSalary.IsSet() && applicationDetails.ExpectedSalary.Value() != nil {

		applicationData.ExpectedSalary = applicationDetails.ExpectedSalary.Value()
	}

	if applicationDetails.Resume != "" {

		isValidBase64, base64Data, extension = IsValidBase64(applicationDetails.Resume)

		if isValidBase64 && base64Data != "" {

			rand_num := strconv.Itoa(int(time.Now().Unix()))

			resumeName = "RES-" + rand_num + "." + extension

			if storageType.SelectedType == "aws" {

				fmt.Printf("aws-S3 storage selected\n")

				resumePath = "member/" + resumeName

				err = storage.UploadFileS3(storageType.Aws, nil, base64Data, resumePath)
				if err != nil {

					fmt.Printf("image upload failed %v\n", err)

					return false, ErrUpload

				}

			} else if storageType.SelectedType == "azure" {

				fmt.Printf("azure storage selected")

			} else if storageType.SelectedType == "drive" {

				fmt.Println("drive storage selected")
			}

		} else if strings.Contains(applicationDetails.Resume, "image-resize?name") {

			resumePath = strings.ReplaceAll(applicationDetails.Resume, "image-resize?name=", "")

		} else {

			ErrorLog.Printf("%v", "illegal base64 data")

			return false, errors.New("illegal base64 data ")

		}

		applicationData.ResumePath = &resumePath

		applicationData.ResumeName = &resumeName

	}

	applicationData.StorageType = &storageType.SelectedType

	result = db.Debug().Table("tbl_jobs_registers").Create(&applicationData).Omit("member_id")
	if result.Error != nil {

		return false, result.Error
	}

	return true, nil

}
