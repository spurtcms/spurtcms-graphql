package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"spurtcms-graphql/graph/model"
	"spurtcms-graphql/storage"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spurtcms/pkgcore/member"
	"gorm.io/gorm"
)

func MemberLogin(db *gorm.DB, ctx context.Context, email string) (bool, error) {

	var memberSettings model.MemberSettings

	if err := db.Debug().Table("tbl_member_settings").First(&memberSettings).Error; err != nil {

		return false, err
	}

	if memberSettings.MemberLogin == "password" {

		return false, ErrMemberLoginPerm
	}

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	Mem.Auth = GetAuthorizationWithoutToken(db)

	member_details, err := Mem.GraphqlMemberLogin(email)

	if gorm.ErrRecordNotFound == err {

		var loginEnquiryTemplate model.EmailTemplate

		if err := db.Debug().Table("tbl_email_templates").Where("is_deleted = 0 and template_name = ?", OwndeskLoginEnquiryTemplate).First(&loginEnquiryTemplate).Error; err != nil {

			c.AbortWithError(http.StatusInternalServerError, err)

			return false, err
		}

		var convIds []int

		adminIds := strings.Split(memberSettings.NotificationUsers, ",")

		for _, adminId := range adminIds {

			convId, _ := strconv.Atoi(adminId)

			convIds = append(convIds, convId)
		}

		_, notifyEmails, _ := GetNotifyAdminEmails(db, convIds)

		var admin_mail_data = MailConfig{Emails: notifyEmails, MailUsername: os.Getenv("MAIL_USERNAME"), MailPassword: os.Getenv("MAIL_PASSWORD"), Subject: loginEnquiryTemplate.TemplateSubject}

		dataReplacer := strings.NewReplacer(
			"{OwndeskLogo}", EmailImagePath.Owndesk,
			"{Username}", "Admin",
			"{UnauthorizedMail}", email,
			"{CurrentTime}", time.Now().In(TimeZone).Format("02 Jan 2006 03:04 PM"),
			"{OwndeskFacebookLink}", SocialMediaLinks.Facebook,
			"{OwndeskLinkedinLink}", SocialMediaLinks.Linkedin,
			"{OwndeskTwitterLink}", SocialMediaLinks.Twitter,
			"{OwndeskYoutubeLink}", SocialMediaLinks.Youtube,
			"{OwndeskInstagramLink}", SocialMediaLinks.Instagram,
			"{FacebookLogo}", EmailImagePath.Facebook,
			"{LinkedinLogo}", EmailImagePath.LinkedIn,
			"{TwitterLogo}", EmailImagePath.Twitter,
			"{YoutubeLogo}", EmailImagePath.Youtube,
			"{InstagramLogo}", EmailImagePath.Instagram,
			"<figure", "<div",
			"</figure", "</div",
			"&nbsp;", "",
		)

		integratedBody := dataReplacer.Replace(loginEnquiryTemplate.TemplateMessage)

		htmlBody := template.HTML(integratedBody)

		tmpl, err := template.ParseFiles("./view/email/email-template.html")

		if err != nil {

			c.AbortWithError(http.StatusInternalServerError, err)

			return false, err
		}

		channel := make(chan error)

		var template_buffers bytes.Buffer

		if err := tmpl.Execute(&template_buffers, gin.H{"body": htmlBody}); err != nil {

			c.AbortWithError(http.StatusInternalServerError, err)

			return false, err
		}

		admin_content := template_buffers.String()

		go SendMail(admin_mail_data, admin_content, channel)

		if <-channel != nil {

			c.AbortWithError(http.StatusInternalServerError, <-channel)

			return false, <-channel

		}

		return false, ErrInvalidMail

	} else if err != nil {

		return false, ErrInvalidMail
	}

	if member_details.IsActive == 0 && member_details.Id != 0 {

		return false, ErrMemberInactive
	}

	var memberProfileData member.TblMemberProfile

	if err := db.Debug().Table("tbl_member_profiles").Where("is_deleted = 0 and member_id = ?", member_details.Id).First(&memberProfileData).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	channel := make(chan error)

	// rand.Seed(time.Now().UnixNano())

	otp := rand.Intn(900000) + 100000

	current_time := time.Now()

	otp_expiry_time := current_time.UTC().Add(5 * time.Minute).Format("2006-01-02 15:04:05")

	mail_expiry_time := current_time.In(TimeZone).Add(5 * time.Minute).Format("02 Jan 2006 03:04 PM")

	err = Mem.StoreGraphqlMemberOtp(otp, member_details.Id, otp_expiry_time)

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	var loginTemplate model.EmailTemplate

	if err := db.Debug().Table("tbl_email_templates").Where("is_deleted=0 and template_name = ?", OwndeskLoginTemplate).First(&loginTemplate).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	dataReplacer := strings.NewReplacer(
		"{OwndeskLogo}", EmailImagePath.Owndesk,
		"{Username}", member_details.Username,
		"{CompanyName}", memberProfileData.CompanyName,
		"{Otp}", strconv.Itoa(otp),
		"{OtpExpiryTime}", mail_expiry_time,
		"{OwndeskFacebookLink}", SocialMediaLinks.Facebook,
		"{OwndeskLinkedinLink}", SocialMediaLinks.Linkedin,
		"{OwndeskTwitterLink}", SocialMediaLinks.Twitter,
		"{OwndeskYoutubeLink}", SocialMediaLinks.Youtube,
		"{OwndeskInstagramLink}", SocialMediaLinks.Instagram,
		"{FacebookLogo}", EmailImagePath.Facebook,
		"{LinkedinLogo}", EmailImagePath.LinkedIn,
		"{TwitterLogo}", EmailImagePath.Twitter,
		"{YoutubeLogo}", EmailImagePath.Youtube,
		"{InstagramLogo}", EmailImagePath.Instagram,
		"<figure", "<div",
		"</figure", "</div",
		"&nbsp;", "",
	)

	integratedBody := dataReplacer.Replace(loginTemplate.TemplateMessage)

	htmlBody := template.HTML(integratedBody)

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	tmpl, err := template.ParseFiles("./view/email/login-template.html")

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	var template_buffer bytes.Buffer

	if err := tmpl.Execute(&template_buffer, gin.H{"body": htmlBody}); err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	var sendMails []string

	sendMails = append(sendMails, member_details.Email)

	mail_data := MailConfig{Emails: sendMails, MailUsername: os.Getenv("MAIL_USERNAME"), MailPassword: os.Getenv("MAIL_PASSWORD"), Subject: loginTemplate.TemplateSubject}

	html_content := template_buffer.String()

	go SendMail(mail_data, html_content, channel)

	if <-channel != nil {

		c.AbortWithError(http.StatusServiceUnavailable, <-channel)

		return false, <-channel

	}

	return true, nil
}

func VerifyMemberOtp(db *gorm.DB, ctx context.Context, email string, otp int) (*model.LoginDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	Mem.Auth = GetAuthorizationWithoutToken(db)

	currentTime := time.Now().UTC()

	memberDetails, token, err := Mem.VerifyLoginOtp(email, otp, currentTime, LocalLoginType)

	if err != nil {

		return &model.LoginDetails{}, err
	}

	var memberProfileDetails model.MemberProfile

	if err := db.Debug().Table("tbl_member_profiles").Select("tbl_member_profiles.*").Joins("inner join tbl_members on tbl_members.id = tbl_member_profiles.member_id").Where("tbl_member_profiles.is_deleted = 0 and tbl_members.is_deleted = 0 and  tbl_members.is_active =1 and tbl_member_profiles.member_id = ?", memberDetails.Id).First(&memberProfileDetails).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.LoginDetails{}, err
	}

	if memberProfileDetails.CompanyLogo != nil && *memberProfileDetails.CompanyLogo != "" {

		logoPath := PathUrl + strings.TrimPrefix(*memberProfileDetails.CompanyLogo, "/")

		memberProfileDetails.CompanyLogo = &logoPath
	}

	return &model.LoginDetails{MemberProfileData: memberProfileDetails, Token: token}, nil

}

func MemberRegister(db *gorm.DB, ctx context.Context, input model.MemberDetails, ecomModule *int) (bool, error) {

	var memberSettings model.MemberSettings

	if err := db.Debug().Table("tbl_member_settings").First(&memberSettings).Error; err != nil {

		return false, err
	}

	if memberSettings.AllowRegistration == 0 {

		return false, ErrMemberRegisterPerm
	}

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	Mem.Auth = GetAuthorizationWithoutToken(db)

	var (
		// imageName, imagePath string

		err error

		ecomMod int = *ecomModule
	)

	// if input.ProfileImage.IsSet() {

	// 	imageName, imagePath, err = StoreImageBase64ToLocal(*input.ProfileImage.Value(), ProfileImagePath, "PROFILE")

	// 	if err != nil {

	// 		c.AbortWithError(http.StatusInternalServerError, err)

	// 		return false, err
	// 	}

	// }

	var memberDetails member.MemberCreation

	if input.Mobile.IsSet() {

		memberDetails.MobileNo = *input.Mobile.Value()
	}

	// if imageName != "" && imagePath != "" {

	// 	memberDetails.ProfileImage = imageName

	// 	memberDetails.ProfileImagePath = imagePath
	// }

	if input.LastName.IsSet() {

		memberDetails.LastName = *input.LastName.Value()

	}

	var hashpass string

	if input.Password.IsSet() {

		hashpass, err = HashingPassword(*input.Password.Value())

		if err != nil {

			return false, ErrPassHash
		}

		memberDetails.Password = hashpass
	}

	if input.Username.IsSet() {

		memberDetails.Username = *input.Username.Value()

		_, isMemberExists, err := Mem.CheckUsernameInMember(0, *input.Username.Value())

		if isMemberExists || err == nil {

			err = errors.New("member already exists")

			c.AbortWithError(422, err)

			return isMemberExists, err
		}

		if ecomMod == 1 {

			var count int64

			if err := db.Table("tbl_ecom_customers").Where("is_deleted = 0 and username = ?", *input.Username.Value()).Count(&count).Error; err != nil {

				return false, err
			}

			if count > 0 {

				err = errors.New("customer already exists")

				c.AbortWithError(422, err)

				return false, err
			}
		}

	}

	if input.Email != "" {

		memberDetails.Email = input.Email

		_, isMemberExists, err := Mem.CheckEmailInMember(0, input.Email)

		if isMemberExists || err == nil {

			err = errors.New("member already exists")

			c.AbortWithError(http.StatusBadRequest, err)

			return isMemberExists, err
		}

		if ecomMod == 1 {

			var count int64

			if err := db.Table("tbl_ecom_customers").Where("is_deleted = 0 and email = ?", input.Email).Count(&count).Error; err != nil {

				return false, err
			}

			if count > 0 {

				err = errors.New("customer already exists")

				c.AbortWithError(422, err)

				return false, err
			}
		}

	}

	memberDetails.FirstName = input.FirstName

	memberData, isRegistered, err := Mem.MemberRegister(memberDetails)

	if !isRegistered || err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return isRegistered, err
	}

	if isRegistered && ecomMod == 1 {

		createdOn, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		is_deleted := 0

		var ecomCustomer = model.CustomerDetails{
			FirstName:        memberDetails.FirstName,
			LastName:         &memberDetails.LastName,
			Username:         memberDetails.Username,
			MobileNo:         memberDetails.MobileNo,
			Email:            input.Email,
			IsActive:         memberDetails.IsActive,
			ProfileImage:     &memberDetails.ProfileImage,
			ProfileImagePath: &memberDetails.ProfileImagePath,
			CreatedOn:        createdOn,
			Password:         hashpass,
			MemberID:         &memberData.Id,
			IsDeleted:        &is_deleted,
		}

		if err := db.Table("tbl_ecom_customers").Create(&ecomCustomer).Error; err != nil {

			c.AbortWithError(http.StatusInternalServerError, err)

			return isRegistered, err
		}
	}

	return true, nil

}

func UpdateMember(db *gorm.DB, ctx context.Context, memberdata model.MemberDetails) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return false, err

	}

	memberData := make(map[string]interface{})

	var err error

	if memberdata.ProfileImage.IsSet() && memberdata.ProfileImage.Value() != nil {

		var fileName, filePath string

		storageType, _ := GetStorageType(db)

		fileName = memberdata.ProfileImage.Value().Filename

		file := memberdata.ProfileImage.Value().File

		if storageType.SelectedType == "aws" {

			fmt.Printf("aws-S3 storage selected\n")

			filePath = "member/" + fileName

			err = storage.UploadFileS3(storageType.Aws, memberdata.ProfileImage.Value(), filePath)

			if err != nil {

				fmt.Printf("image upload failed %v\n", err)

				return false, ErrUpload

			}

		} else if storageType.SelectedType == "local" {

			fmt.Printf("local storage selected\n")

			b64Data, err := IoReadSeekerToBase64(file)

			if err != nil {

				return false, err
			}

			endpoint := "gqlSaveLocal"

			url := PathUrl + endpoint

			filePath, err = storage.UploadImageToAdminLocal(b64Data, fileName, url)

			if err != nil {

				return false, ErrUpload
			}

			log.Printf("local stored path: %v\n", filePath)

		} else if storageType.SelectedType == "azure" {

			fmt.Printf("azure storage selected")

		} else if storageType.SelectedType == "drive" {

			fmt.Println("drive storage selected")
		}

		memberData["profile_image"] = fileName

		memberData["profile_image_path"] = filePath

	}

	memberData["first_name"] = memberdata.FirstName

	memberData["email"] = memberdata.Email

	if memberdata.Mobile.IsSet() && memberdata.Mobile.Value() != nil {

		memberData["mobile_no"] = *memberdata.Mobile.Value()

	}

	if memberdata.GroupID.IsSet() && memberdata.GroupID.Value() != nil && *memberdata.GroupID.Value() != 0 {

		memberData["member_group_id"] = *memberdata.GroupID.Value()

	}

	if memberdata.Password.IsSet() && memberdata.Password.Value() != nil && *memberdata.Password.Value() != "" {

		hashpass, err := HashingPassword(*memberdata.Password.Value())

		if err != nil {

			return false, ErrPassHash
		}

		memberData["password"] = &hashpass
	}

	if memberdata.LastName.IsSet() && memberdata.LastName.Value() != nil {

		memberData["last_name"] = *memberdata.LastName.Value()

	}

	if memberdata.Username.IsSet() && memberdata.Username.Value() != nil {

		memberData["username"] = memberdata.Username.Value()

	}

	if memberdata.IsActive.IsSet() && memberdata.IsActive.Value() != nil {

		memberData["is_active"] = *memberdata.IsActive.Value()

	}

	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	memberData["modified_on"] = &currentTime

	memberData["modified_by"] = &memberid

	query := db.Debug().Table("tbl_members").Where("is_deleted = 0 and id = ?", memberid)

	if err = query.UpdateColumns(memberData).Error; err != nil {

		return false, err
	}

	return true, nil

}

func TemplateMemberLogin(db *gorm.DB, ctx context.Context, username, email *string, password string) (string, error) {

	var memberSettings model.MemberSettings

	if err := db.Debug().Table("tbl_member_settings").First(&memberSettings).Error; err != nil {

		return "", err
	}

	if memberSettings.MemberLogin == "otp" {

		return "", ErrMemberLoginPerm
	}

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	Mem.Auth = GetAuthorizationWithoutToken(db)

	var memberLogin member.MemberLogin

	if username != nil {

		memberLogin.Username = *username

	} else if email != nil {

		memberLogin.Emailid = *email
	}

	memberLogin.Password = password

	token, err := Mem.CheckMemberLogin(memberLogin, db, os.Getenv("JWT_SECRET"), LocalLoginType)

	if err != nil {

		c.AbortWithError(http.StatusUnauthorized, err)

		log.Println(err)
	}

	return token, err
}

func MemberProfileDetails(db *gorm.DB, ctx context.Context) (*model.MemberProfile, error) {

	c, ok := ctx.Value(ContextKey).(*gin.Context)

	if !ok {

		ErrorLog.Printf("memberProfileDetails context error: %v", ok)
	}

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		ErrorLog.Printf("memberProfileDetails context error: %s", err)

		c.AbortWithError(http.StatusUnauthorized, err)

		return &model.MemberProfile{}, err

	}

	var memberProfile model.MemberProfile

	if err := db.Debug().Table("tbl_member_profiles").Select("tbl_member_profiles.*,tbl_members.is_active").Joins("inner join tbl_members on tbl_members.id = tbl_member_profiles.member_id").Where("tbl_members.is_deleted = 0 and tbl_member_profiles.is_deleted = 0 and tbl_member_profiles.member_id = ?", memberid).First(&memberProfile).Error; err != nil {

		ErrorLog.Printf("get memberProfileDetails data error: %s", err)

		c.AbortWithError(http.StatusUnprocessableEntity, err)

		return &model.MemberProfile{}, err
	}

	if memberProfile.CompanyLogo != nil && *memberProfile.CompanyLogo != "" {

		var logoPath string

		storageType, _ := GetStorageType(db)

		awsCreds := storageType.Aws

		isExist, _ := storage.CheckS3FileExistence(awsCreds, *memberProfile.CompanyLogo)

		if isExist {

			s3FileServeEndpoint := "image-resize"

			logoPath = PathUrl + s3FileServeEndpoint + "?name=" + strings.TrimPrefix(*memberProfile.CompanyLogo, "/")

		} else {

			logoPath = PathUrl + strings.TrimPrefix(*memberProfile.CompanyLogo, "/")
		}

		memberProfile.CompanyLogo = &logoPath
	}

	return &memberProfile, nil
}

func GetMemberProfileDetails(db *gorm.DB, ctx context.Context, id *int, profileSlug *string) (*model.MemberProfile, error) {

	c, ok := ctx.Value(ContextKey).(*gin.Context)

	if !ok {

		ErrorLog.Printf("getmemberProfileDetails context error: %v", ok)

	}

	tokenType := c.GetString("tokenType")

	memberid := c.GetInt("memberid")

	var memberProfile member.TblMemberProfile

	query := db.Select("tbl_member_profiles.*").Table("tbl_member_profiles").Joins("inner join tbl_members on tbl_members.id = tbl_member_profiles.member_id").Where("tbl_members.is_deleted = 0 and tbl_member_profiles.is_deleted = 0")

	if id != nil {

		query = query.Where("tbl_member_profiles.member_id = ?", *id)

	} else if profileSlug != nil {

		query = query.Where("tbl_member_profiles.profile_slug = ?", *profileSlug)
	}

	if err := query.First(&memberProfile).Error; err != nil {

		ErrorLog.Printf("getmemberProfileDetails data error: %s", err)

		return &model.MemberProfile{}, err
	}

	var memberDetails model.Member

	if err := db.Debug().Table("tbl_members").Where("is_deleted = 0 and id = ?", memberProfile.MemberId).First(&memberDetails).Error; err != nil {

		return &model.MemberProfile{}, err
	}

	if memberDetails.IsActive == 0 && memberDetails.ID != 0 {

		if memberid != 0 && tokenType == LocalLoginType {

			return &model.MemberProfile{}, ErrMemberInactive

		} else if memberid == 0 {

			return &model.MemberProfile{}, ErrMemberInactive
		}
	}

	var profileLogo string

	if memberProfile.CompanyLogo != "" {

		storageType, _ := GetStorageType(db)

		awsCreds := storageType.Aws

		isExist, _ := storage.CheckS3FileExistence(awsCreds, memberProfile.CompanyLogo)

		if isExist {

			s3FileServeEndpoint := "image-resize"

			profileLogo = PathUrl + s3FileServeEndpoint + "?name=" + strings.TrimPrefix(memberProfile.CompanyLogo, "/")

		} else {

			profileLogo = PathUrl + strings.TrimPrefix(memberProfile.CompanyLogo, "/")
		}

	}

	MemberProfile := model.MemberProfile{
		ID:              &memberProfile.Id,
		MemberID:        &memberProfile.MemberId,
		ProfileName:     &memberProfile.ProfileName,
		ProfileSlug:     &memberProfile.ProfileSlug,
		ProfilePage:     &memberProfile.ProfilePage,
		MemberDetails:   &memberProfile.MemberDetails,
		CompanyName:     &memberProfile.CompanyName,
		CompanyLocation: &memberProfile.CompanyLocation,
		CompanyLogo:     &profileLogo,
		About:           &memberProfile.About,
		SeoTitle:        &memberProfile.SeoTitle,
		SeoDescription:  &memberProfile.SeoDescription,
		SeoKeyword:      &memberProfile.SeoKeyword,
		CreatedBy:       &memberProfile.CreatedBy,
		CreatedOn:       &memberProfile.CreatedOn,
		ModifiedOn:      &memberProfile.ModifiedOn,
		ModifiedBy:      &memberProfile.ModifiedBy,
		Linkedin:        &memberProfile.Linkedin,
		Twitter:         &memberProfile.Twitter,
		Website:         &memberProfile.Website,
		ClaimStatus:     &memberProfile.ClaimStatus,
	}

	return &MemberProfile, nil
}

func MemberPasswordUpdate(db *gorm.DB, ctx context.Context, oldPassword string, newPassword string, confirmPassword string) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberId := c.GetInt("memberid")

	if memberId == 0 {

		err := errors.New("unauthorized access")

		ErrorLog.Printf("memberProfileDetails context error: %s", err)

		c.AbortWithError(http.StatusUnauthorized, err)

		return false, err

	}

	var loggedInMember member.TblMember

	result := db.Where("is_deleted = 0 and id = ?", memberId).First(&loggedInMember)

	if result.Error != nil {

		c.AbortWithError(404, result.Error)

		return false, result.Error
	}

	err := CompareBcryptPassword(loggedInMember.Password, oldPassword)

	if err != nil {

		return false, ErrOldPass

	}

	if newPassword != confirmPassword {

		return false, ErrConfirmPass
	}

	hasspass, err := HashingPassword(confirmPassword)

	if err != nil {

		return false, ErrPassHash
	}

	err = CompareBcryptPassword(hasspass, oldPassword)

	if err == nil {

		return false, ErrSamePass
	}

	updateQuery := db.Model(&member.TblMember{}).Where("id = ?", memberId).Update("password", hasspass)

	if updateQuery.Error != nil {

		c.AbortWithError(404, updateQuery.Error)

		return false, updateQuery.Error

	}

	return true, nil
}

func GetMemberDetails(db *gorm.DB, ctx context.Context) (*model.Member, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberId := c.GetInt("memberid")

	if memberId == 0 {

		err := errors.New("unauthorized access")

		ErrorLog.Printf("memberProfileDetails context error: %s", err)

		c.AbortWithError(http.StatusUnauthorized, err)

		return &model.Member{}, err

	}

	var memberDetails model.Member

	result := db.Table("tbl_members").Where("is_deleted = 0 and id = ?", memberId).First(&memberDetails)

	if result.Error != nil {

		c.AbortWithError(404, result.Error)

		return &model.Member{}, result.Error
	}

	if memberDetails.ProfileImagePath != "" {

		storageType, _ := GetStorageType(db)

		awsCreds := storageType.Aws

		isExist, _ := storage.CheckS3FileExistence(awsCreds, memberDetails.ProfileImagePath)

		var modifiedPath string

		if isExist {

			s3FileServeEndpoint := "image-resize"

			modifiedPath = PathUrl + s3FileServeEndpoint + "?name=" + strings.TrimPrefix(memberDetails.ProfileImagePath, "/")

		} else {

			modifiedPath = PathUrl + strings.TrimPrefix(memberDetails.ProfileImagePath, "/")
		}

		memberDetails.ProfileImagePath = modifiedPath

	}

	return &memberDetails, nil

}
