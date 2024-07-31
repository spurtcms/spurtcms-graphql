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
	"github.com/spurtcms/jobs"
	"gorm.io/gorm"
)

func JobsList(db *gorm.DB, ctx context.Context, limit int, offset int, filter *model.JobFilter) (*model.JobsList, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	var (
		jobsLocal  []model.Job
		job        model.Job
		count      int64
		jobsFilter jobs.Filter
		err        error
	)

	if filter != nil {

		if filter.JobTitle.IsSet() {

			jobsFilter.JobTitle = *filter.JobTitle.Value()
		}

		if filter.KeyWord.IsSet() {

			jobsFilter.KeyWord = *filter.KeyWord.Value()
		}

		if filter.JobLocation.IsSet() {

			jobsFilter.JobLocation = *filter.JobLocation.Value()
		}

		if filter.CategoryID.IsSet() {

			jobsFilter.CategoryId = *filter.CategoryID.Value()
		}

		if filter.CategorySlug.IsSet() {

			jobsFilter.CategorySlug = *filter.CategorySlug.Value()
		}

		if filter.MaximumYears.IsSet() {

			jobsFilter.MaximumYears = *filter.MaximumYears.Value()
		}

		if filter.MinimumYears.IsSet() {

			jobsFilter.MinimumYears = *filter.MinimumYears.Value()
		}

		if filter.DatePosted.IsSet() {

			jobsFilter.DatePosted = *filter.DatePosted.Value()
		}
	}

	jobsList, count, err := JobsInstance.GetJobsList(limit, offset, jobsFilter)
	if err != nil {

		ErrorLog.Printf("%v: %v", ErrFetchJobsList, err)

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.JobsList{}, err
	}

	for _, jobList := range jobsList {
		job.CategoriesID = jobList.CategoriesId
		job.Category.CategoryName = jobList.Category.CategoryName
		job.Category.CategorySlug = jobList.Category.CategorySlug
		job.Category.CreatedBy = jobList.Category.CreatedBy
		job.Category.CreatedOn = jobList.Category.CreatedOn
		job.Category.Description = jobList.Category.Description
		job.Category.ID = jobList.Category.Id
		job.Category.ImagePath = jobList.Category.ImagePath
		job.Category.ModifiedBy = &jobList.Category.ModifiedBy
		job.Category.ModifiedOn = &jobList.Category.ModifiedOn
		job.Category.ParentID = jobList.Category.ParentId
		job.CreatedBy = jobList.CreatedBy
		job.CreatedOn = jobList.CreatedOn
		job.DeletedBy = &jobList.DeletedBy
		job.DeletedOn = &jobList.DeletedOn
		job.Department = &jobList.Department
		job.Education = jobList.Education
		job.Experience = &jobList.Experience
		job.ID = jobList.Id
		job.IsDeleted = &jobList.IsDeleted
		job.JobDescription = jobList.JobDescription
		job.JobLocation = jobList.JobLocation
		job.JobSlug = jobList.JobSlug
		job.JobTitle = jobList.JobTitle
		job.JobType = jobList.JobType
		job.Keyword = &jobList.Keywords
		job.MaximumYears = jobList.MaximumYears
		job.MinimumYears = jobList.MinimumYears
		job.ModifiedBy = &jobList.ModifiedBy
		job.ModifiedOn = &jobList.ModifiedOn
		job.PostedDate = jobList.PostedDate
		job.Salary = jobList.Salary
		job.Skill = jobList.Skill
		job.Status = jobList.Status
		job.ValidThrough = jobList.ValidThrough

		jobsLocal = append(jobsLocal, job)

	}

	return &model.JobsList{Jobs: jobsLocal, Count: int(count)}, nil
}

func JobDetail(db *gorm.DB, ctx context.Context, id *int, jobSlug *string) (*model.Job, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	var (
		jobDetailLocal model.Job
		jobId          int
		slug           string
	)

	if id != nil {

		jobId = *id
	} else if jobSlug != nil {

		slug = *jobSlug
	}

	jobDetails, err := JobsInstance.GetJobDetails(jobId, slug)
	if err != nil {

		ErrorLog.Printf("%v: %v", ErrFetchJobDetails, err)

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.Job{}, err
	}

	jobDetailLocal.CategoriesID = jobDetails.CategoriesId
	jobDetailLocal.Category.CategoryName = jobDetails.Category.CategoryName
	jobDetailLocal.Category.CategorySlug = jobDetails.Category.CategorySlug
	jobDetailLocal.Category.CreatedBy = jobDetails.Category.CreatedBy
	jobDetailLocal.Category.CreatedOn = jobDetails.Category.CreatedOn
	jobDetailLocal.Category.Description = jobDetails.Category.Description
	jobDetailLocal.Category.ID = jobDetails.Category.Id
	jobDetailLocal.Category.ImagePath = jobDetails.Category.ImagePath
	jobDetailLocal.Category.ModifiedBy = &jobDetails.Category.ModifiedBy
	jobDetailLocal.Category.ModifiedOn = &jobDetails.Category.ModifiedOn
	jobDetailLocal.Category.ParentID = jobDetails.Category.ParentId
	jobDetailLocal.CreatedBy = jobDetails.CreatedBy
	jobDetailLocal.CreatedOn = jobDetails.CreatedOn
	jobDetailLocal.DeletedBy = &jobDetails.DeletedBy
	jobDetailLocal.DeletedOn = &jobDetails.DeletedOn
	jobDetailLocal.Department = &jobDetails.Department
	jobDetailLocal.Education = jobDetails.Education
	jobDetailLocal.Experience = &jobDetails.Experience
	jobDetailLocal.ID = jobDetails.Id
	jobDetailLocal.IsDeleted = &jobDetails.IsDeleted
	jobDetailLocal.JobDescription = jobDetails.JobDescription
	jobDetailLocal.JobLocation = jobDetails.JobLocation
	jobDetailLocal.JobSlug = jobDetails.JobSlug
	jobDetailLocal.JobTitle = jobDetails.JobTitle
	jobDetailLocal.JobType = jobDetails.JobType
	jobDetailLocal.Keyword = &jobDetails.Keywords
	jobDetailLocal.MaximumYears = jobDetails.MaximumYears
	jobDetailLocal.MinimumYears = jobDetails.MinimumYears
	jobDetailLocal.ModifiedBy = &jobDetails.ModifiedBy
	jobDetailLocal.ModifiedOn = &jobDetails.ModifiedOn
	jobDetailLocal.PostedDate = jobDetails.PostedDate
	jobDetailLocal.Salary = jobDetails.Salary
	jobDetailLocal.Skill = jobDetails.Skill
	jobDetailLocal.Status = jobDetails.Status
	jobDetailLocal.ValidThrough = jobDetails.ValidThrough

	return &jobDetailLocal, nil
}

func ApplicantDetails(db *gorm.DB, ctx context.Context, jobId int, emailId string) (*model.ApplicantDetails, error) {

	c, ok := ctx.Value(ContextKey).(*gin.Context)

	if !ok {
		ErrorLog.Printf("%v", ErrGettingContext)

		return &model.ApplicantDetails{}, ErrGettingContext
	}

	memberId := c.GetInt("memberid")

	if memberId == 0 {

		ErrorLog.Printf("%v", ErrGettingMemberId)

		c.AbortWithError(http.StatusUnauthorized, ErrUnauthorizedAccess)

		return &model.ApplicantDetails{}, ErrUnauthorizedAccess
	}

	var (
		finalApplicantDetails model.ApplicantDetails
		err                   error
		applicantDetails      jobs.ApplicantDetails
		imagePath, resumePath string
	)

	applicantDetails, err = JobsAuthInstance.GetApplicantDetails(jobId, memberId, emailId)
	if err != nil {

		ErrorLog.Printf("%v: %v", ErrGettingApplicantDetail, err)

		c.AbortWithError(http.StatusUnauthorized, ErrGettingApplicantDetail)

		return &model.ApplicantDetails{}, err
	}

	finalApplicantDetails.ApplicantID = &applicantDetails.ApplicantID
	finalApplicantDetails.CompanyName = &applicantDetails.CompanyName
	finalApplicantDetails.CreatedBy = &applicantDetails.CreatedBy
	finalApplicantDetails.CreatedOn = &applicantDetails.CreatedOn
	finalApplicantDetails.CurrentSalary = &applicantDetails.CurrentSalary
	finalApplicantDetails.DeletedBy = &applicantDetails.DeletedBy
	finalApplicantDetails.DeletedOn = &applicantDetails.DeletedOn
	finalApplicantDetails.Education = &applicantDetails.Education
	finalApplicantDetails.EmailID = &applicantDetails.EmailID
	finalApplicantDetails.ExpectedSalary = &applicantDetails.ExpectedSalary
	finalApplicantDetails.Experience = &applicantDetails.Experience
	finalApplicantDetails.Gender = &applicantDetails.Gender
	finalApplicantDetails.Graduation = &applicantDetails.Graduation
	finalApplicantDetails.ID = &applicantDetails.ID
	finalApplicantDetails.Image = &applicantDetails.Image
	finalApplicantDetails.ImagePath = &applicantDetails.ImagePath
	finalApplicantDetails.IsDeleted = &applicantDetails.IsDeleted
	finalApplicantDetails.JobID = &applicantDetails.JobID
	finalApplicantDetails.JobType = &applicantDetails.JobType
	finalApplicantDetails.Location = &applicantDetails.Location
	finalApplicantDetails.MobileNo = &applicantDetails.MobileNo
	finalApplicantDetails.ModifiedBy = &applicantDetails.ModifiedBy
	finalApplicantDetails.ModifiedOn = &applicantDetails.ModifiedOn
	finalApplicantDetails.Name = &applicantDetails.Name
	finalApplicantDetails.ResumeName = &applicantDetails.ResumeName
	finalApplicantDetails.ResumePath = &applicantDetails.ResumePath
	finalApplicantDetails.Skills = &applicantDetails.Skills
	finalApplicantDetails.Status = &applicantDetails.Status
	finalApplicantDetails.StorageType = &applicantDetails.StorageType

	if finalApplicantDetails.StorageType != nil && *finalApplicantDetails.StorageType == "aws" && finalApplicantDetails.ImagePath != nil && *finalApplicantDetails.ImagePath != "" {

		imagePath = "image-resize?name=" + *finalApplicantDetails.ImagePath

	} else if finalApplicantDetails.StorageType != nil && *finalApplicantDetails.StorageType == "local" && finalApplicantDetails.ImagePath != nil && *finalApplicantDetails.ImagePath != "" {

		imagePath = *finalApplicantDetails.ImagePath

	} else {

		imagePath = ""
	}

	finalApplicantDetails.ImagePath = &imagePath

	if finalApplicantDetails.StorageType != nil && *finalApplicantDetails.StorageType == "aws" && finalApplicantDetails.ResumePath != nil && *finalApplicantDetails.ResumePath != "" {

		resumePath = "image-resize?name=" + *finalApplicantDetails.ResumePath

	} else if finalApplicantDetails.StorageType != nil && *finalApplicantDetails.StorageType == "local" && finalApplicantDetails.ResumePath != nil && *finalApplicantDetails.ResumePath != "" {

		resumePath = *finalApplicantDetails.ResumePath
	} else {

		resumePath = ""
	}

	finalApplicantDetails.ResumePath = &resumePath

	return &finalApplicantDetails, nil
}

func JobApplication(db *gorm.DB, ctx context.Context, applicationDetails model.ApplicationInput) (bool, error) {

	c, ok := ctx.Value(ContextKey).(*gin.Context)

	if !ok {

		ErrorLog.Printf("%v", ErrGettingContext)

		return false, ErrGettingContext
	}

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		ErrorLog.Printf("%v: %s", ErrGettingMemberId, ErrUnauthorizedAccess)

		c.AbortWithError(http.StatusUnauthorized, ErrUnauthorizedAccess)

		return false, ErrUnauthorizedAccess

	}

	var (
		applicantDetails                                                    jobs.ApplicantDetails
		applicationData                                                     jobs.ApplicantDetails
		imageName, imagePath, resumeName, resumePath, base64Data, extension string
		storageType                                                         StorageType
		err                                                                 error
		isValidBase64                                                       bool
		registeredApplicant                                                 int64
	)

	storageType, err = GetStorageType(db)
	if err != nil {
		ErrorLog.Printf("%v: %v", ErrFetchStorageType, err)

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	applicantDetails, err = JobsInstance.GetApplicantDetails(0, memberid, "")
	if err != nil {

		ErrorLog.Printf("%v: %v", ErrGettingApplicantDetail, err)

		c.AbortWithError(http.StatusInternalServerError, ErrGettingApplicantDetail)

		return false, err
	}

	if applicationDetails.EmailID == applicantDetails.EmailID {

		applicationData.JobID = applicationDetails.JobID

		applicationData.EmailID = applicationDetails.EmailID

		registeredApplicant, err = JobsInstance.CheckAlreadyRegistered(applicationData.JobID, applicationData.EmailID)
		if err != nil {

			ErrorLog.Printf("%v: %v", ErrCheckingAlreadyRegistered, err)

			c.AbortWithError(http.StatusInternalServerError, ErrCheckingAlreadyRegistered)

			return false, err
		}

		if registeredApplicant > 0 {

			ErrorLog.Printf("%v", ErrApplicantAlreadyRegistered)

			c.AbortWithError(http.StatusInternalServerError, ErrApplicantAlreadyRegistered)

			return false, ErrApplicantAlreadyRegistered
		}

		applicationData.ApplicantID = applicantDetails.ID

		applicationData.CreatedBy = applicantDetails.ID

		currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		applicationData.CreatedOn = currentTime

		applicationData.Name = applicationDetails.Name

		applicationData.MobileNo = applicationDetails.MobileNo

		if applicationDetails.JobType.IsSet() && applicationDetails.JobType.Value() != nil {

			applicationData.JobType = *applicationDetails.JobType.Value()
		}

		applicationData.Location = applicationDetails.Location

		applicationData.Education = applicationDetails.Education

		applicationData.Graduation = applicationDetails.Graduation

		applicationData.Gender = applicationDetails.Gender

		if applicationDetails.CompanyName.IsSet() && applicationDetails.CompanyName.Value() != nil {

			applicationData.CompanyName = *applicationDetails.CompanyName.Value()
		}

		applicationData.Experience = applicationDetails.Experience

		applicationData.Skills = applicationDetails.Skills

		if applicationDetails.Image != "" {

			isValidBase64, base64Data, extension = IsValidBase64(applicationDetails.Image)

			if isValidBase64 && base64Data != "" {

				rand_num := strconv.Itoa(int(time.Now().Unix()))

				if extension == "msword" {

					imageName = "IMG-" + rand_num + "." + "doc"

				} else if extension == "vnd.openxmlformats-officedocument.wordprocessingml.document" {

					imageName = "IMG-" + rand_num + "." + "docx"

				} else {

					imageName = "IMG-" + rand_num + "." + extension

				}

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

				c.AbortWithError(http.StatusNotAcceptable, errors.New("illegal base64 data "))

				return false, errors.New("illegal base64 data ")

			}

			applicationData.ImagePath = imagePath

			applicationData.Image = imageName

		}

		isDeleted := 0

		applicationData.IsDeleted = isDeleted

		if applicationDetails.CurrentSalary.IsSet() && applicationDetails.CurrentSalary.Value() != nil {

			applicationData.CurrentSalary = *applicationDetails.CurrentSalary.Value()
		}

		if applicationDetails.ExpectedSalary.IsSet() && applicationDetails.ExpectedSalary.Value() != nil {

			applicationData.ExpectedSalary = *applicationDetails.ExpectedSalary.Value()
		}

		if applicationDetails.Resume != "" {

			isValidBase64, base64Data, extension = IsValidBase64(applicationDetails.Resume)

			if isValidBase64 && base64Data != "" {

				rand_num := strconv.Itoa(int(time.Now().Unix()))

				if extension == "msword" {

					resumeName = "RES-" + rand_num + "." + "doc"

				} else if extension == "vnd.openxmlformats-officedocument.wordprocessingml.document" {

					resumeName = "RES-" + rand_num + "." + "docx"

				} else {

					resumeName = "RES-" + rand_num + "." + extension

				}

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

				c.AbortWithError(http.StatusNotAcceptable, errors.New("illegal base64 data "))

				return false, errors.New("illegal base64 data ")

			}

			applicationData.ResumePath = resumePath

			applicationData.ResumeName = resumeName

		}

		applicationData.StorageType = storageType.SelectedType

		err = JobsAuthInstance.CreateJobApplication(applicationData)
		if err != nil {

			return false, err
		}

	} else {

		ErrorLog.Printf("%v", "Please use the registered email")

		return false, errors.New("please use the registered email")

	}

	return true, nil

}
