package controller

import (
	"bytes"
	"context"
	"encoding/json"
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
	authPkg "github.com/spurtcms/auth"
	"github.com/spurtcms/pkgcore/member"
	"gorm.io/gorm"

	memberPkg "github.com/spurtcms/member"
	// ecomPkg "github.com/spurtcms/ecommerce"
)

func MemberLogin(db *gorm.DB, ctx context.Context, email string) (bool, error) {

	memberInstance := GetMemberInstanceWithoutAuth()

	memberSettings, err := memberInstance.GetMemberSettings()

	if err != nil{

		return false, err
	}

	if memberSettings.MemberLogin == "password" {

		return false, ErrMemberLoginPerm
	}

	c, _ := ctx.Value(ContextKey).(*gin.Context)


	memberDetails, err :=  memberInstance.GetMemberAndProfileData(0,email,0,"")

	if memberDetails.IsActive != 1{

		return false, ErrMemberInactive
	}

	if gorm.ErrRecordNotFound == err {

		var convIds []int

		adminIds := strings.Split(memberSettings.NotificationUsers, ",")

		for _, adminId := range adminIds {

			convId, _ := strconv.Atoi(adminId)

			convIds = append(convIds, convId)
		}

		_, notifyEmails, _ := GetNotifyAdminEmails(db, convIds)

		var loginEnquiryTemplate model.EmailTemplate

		if err := db.Debug().Table("tbl_email_templates").Where("is_deleted = 0 and template_name = ?", OwndeskLoginEnquiryTemplate).First(&loginEnquiryTemplate).Error; err != nil {

			c.AbortWithError(http.StatusInternalServerError, err)

			return false, err
		}

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

	if memberDetails.IsActive != 1{

		return false, ErrMemberInactive
	}

	channel := make(chan error)

	// rand.Seed(time.Now().UnixNano())

	otp := rand.Intn(900000) + 100000

	expiryTime,err := memberInstance.Auth.UpdateMemberOTP(authPkg.OTP{Length: 6,Duration: 5 * time.Minute,MemberId: memberDetails.Id})

	mail_expiry_time := expiryTime.In(TimeZone).Format("02 Jan 2006 03:04 PM")

	var loginTemplate model.EmailTemplate

	if err := db.Debug().Table("tbl_email_templates").Where("is_deleted=0 and template_name = ?", OwndeskLoginTemplate).First(&loginTemplate).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	dataReplacer := strings.NewReplacer(
		"{OwndeskLogo}", EmailImagePath.Owndesk,
		"{Username}", memberDetails.Username,
		"{CompanyName}", memberDetails.TblMemberProfile.CompanyName,
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

	sendMails = append(sendMails, memberDetails.Email)

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

	memberInstance := GetMemberInstanceWithoutAuth()

	member,err := memberInstance.Auth.CheckMemberLogin(authPkg.MemberLoginCheck{Email: email,OTP: otp,EmailWithOTP: true})

	if err != nil{

		return &model.LoginDetails{}, err
	}

	token, err := authPkg.CreateMemberToken(member.Id,member.MemberGroupId,os.Getenv("JWT_SECRET"),LocalLoginType)

	if err != nil{

		return &model.LoginDetails{}, err
	}

	memberProfile,err := memberInstance.GetMemberProfileByMemberId(member.Id)

	if memberProfile.CompanyLogo != "" {

		memberProfile.CompanyLogo = GetFilePathsRelatedToStorageTypes(db, memberProfile.CompanyLogo)
	}

	conv_memProfile := model.MemberProfile{
		ID:              &memberProfile.Id,
		MemberID:        &memberProfile.MemberId,
		ProfileName:     &memberProfile.ProfileName,
		ProfileSlug:     &memberProfile.ProfileSlug,
		ProfilePage:     &memberProfile.ProfilePage,
		MemberDetails:   &memberProfile.MemberDetails,
		CompanyName:     &memberProfile.CompanyName,
		CompanyLocation: &memberProfile.CompanyLocation,
		CompanyLogo:     &memberProfile.CompanyLogo,
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

	return &model.LoginDetails{MemberProfileData: conv_memProfile, Token: token}, nil

}

func MemberRegister(db *gorm.DB, ctx context.Context, input model.MemberDetails, ecomModule *int) (bool, error) {

	var (
		fileName, filePath string

		profileName  = input.FirstName

		err error

		ecomMod int = *ecomModule

		memberDetails memberPkg.MemberCreationUpdation

		memberProfile memberPkg.MemberprofilecreationUpdation
	)

	memberInstance := GetMemberInstance()

	memberSettings, err := memberInstance.GetMemberSettings()

	if err != nil{

		return false, err
	}

	if memberSettings.AllowRegistration == 0 {

		return false, ErrMemberRegisterPerm
	}

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	if input.Username.IsSet() && input.Username.Value() != nil {

		memberDetails.Username = *input.Username.Value()

		_, err := memberPkg.Membermodel.CheckNameInMember(0, *input.Username.Value(), db)

		if err == nil {

			err = errors.New("member username already exists")

			c.AbortWithError(400, err)

			return false, err
		}

	}

	if input.Email != "" {

		memberDetails.Email = input.Email

		err := memberPkg.Membermodel.CheckEmailInMember(&memberPkg.TblMember{}, input.Email, 0, db)

		if err == nil {

			err = errors.New("member email already exists")

			c.AbortWithError(http.StatusBadRequest, err)

			return false, err
		}

	}

	if input.Mobile.IsSet() && input.Mobile.Value() != nil {

		memberDetails.MobileNo = *input.Mobile.Value()

		err := memberPkg.Membermodel.CheckNumberInMember(&memberPkg.TblMember{}, *input.Mobile.Value(), 0, db)

		if err == nil {

			err = errors.New("member mobile number  already exists")

			c.AbortWithError(http.StatusBadRequest, err)

			return false, err
		}
	}

	if input.ProfileImage.IsSet() && input.ProfileImage.Value() != nil {

		storageType, _ := GetStorageType(db)

		fileName = input.ProfileImage.Value().Filename

		file := input.ProfileImage.Value().File

		if storageType.SelectedType == "aws" {

			fmt.Printf("aws-S3 storage selected\n")

			filePath = "member/" + fileName

			err = storage.UploadFileS3(storageType.Aws, input.ProfileImage.Value(), filePath)

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

	}

	if fileName != "" && filePath != "" {

		memberDetails.ProfileImage = fileName

		memberDetails.ProfileImagePath = filePath
	}

	if input.LastName.IsSet() && input.LastName.Value() != nil {

		memberDetails.LastName = *input.LastName.Value()

		profileName = input.FirstName + " " + *input.LastName.Value()

	}

	if input.Password.IsSet() && input.Password.Value() != nil {

		memberDetails.Password = *input.Password.Value()
	}

	if input.IsActive.IsSet() && input.IsActive.Value() != nil {

		memberDetails.IsActive = *input.IsActive.Value()
	}

	memberDetails.FirstName = input.FirstName

	memberDetails.GroupId = 1

	registeredMember, err := memberInstance.CreateMember(memberDetails)

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	memberProfile.ProfileName = profileName

	memberProfile.ProfileSlug = strings.ReplaceAll(profileName," ","-")

	memberProfile.MemberId = registeredMember.Id

	memberProfile.CreatedBy = registeredMember.Id

	err = memberInstance.CreateMemberProfile(memberProfile)

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false,err
	}

	if ecomMod == 1 {

		// ecomConfig := GetEcomInstanceWithoutAuth()

		// ecomConfig.CreateCustomer(ecomPkg.CreateCustomerReq{
		// 	MemberId: ,
		// })

		is_deleted := 0

		createdOn, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

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
			Password:         registeredMember.Password,
			MemberID:         &registeredMember.Id,
			IsDeleted:        &is_deleted,
		}

		if err := db.Table("tbl_ecom_customers").Create(&ecomCustomer).Error; err != nil {

			c.AbortWithError(http.StatusInternalServerError, err)

			return false,err
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

	memberInstance := GetMemberInstance()

	if err = memberInstance.MemberFlexibleUpdate(memberData, memberid, memberid); err != nil {

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

		ErrorLog.Printf("unauthorized error: %s", err)

		c.AbortWithError(http.StatusUnauthorized, err)

		return &model.MemberProfile{}, err

	}

	memberInstance := GetMemberInstance()

	memberProfile, err := memberInstance.GetMemberProfileByMemberId(memberid)

	if err != nil {

		ErrorLog.Printf("memberProfileDetails context error: %s", err)

		return &model.MemberProfile{}, err
	}

	if memberProfile.CompanyLogo != "" {

		memberProfile.CompanyLogo = GetFilePathsRelatedToStorageTypes(db, memberProfile.CompanyLogo)

	}

	conv_memProfile := &model.MemberProfile{
		ID:              &memberProfile.Id,
		MemberID:        &memberProfile.MemberId,
		ProfileName:     &memberProfile.ProfileName,
		ProfileSlug:     &memberProfile.ProfileSlug,
		ProfilePage:     &memberProfile.ProfilePage,
		MemberDetails:   &memberProfile.MemberDetails,
		CompanyName:     &memberProfile.CompanyName,
		CompanyLocation: &memberProfile.CompanyLocation,
		CompanyLogo:     &memberProfile.CompanyLogo,
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

	return conv_memProfile, nil
}

func GetMemberProfileDetails(db *gorm.DB, ctx context.Context, id *int, profileSlug *string) (*model.MemberProfile, error) {

	c, ok := ctx.Value(ContextKey).(*gin.Context)

	if !ok {

		ErrorLog.Printf("gin instance retrieval context error: %v", ok)

	}

	tokenType := c.GetString("tokenType")

	memberid := c.GetInt("memberid")

	memberInstance := GetMemberInstanceWithoutAuth()

	var memberDetailedProfile memberPkg.Tblmember

	var err error

	if id != nil && *id != 0 {

		memberDetailedProfile, err = memberInstance.GetMemberAndProfileData(0, "", *id, "")

	} else if profileSlug != nil && *profileSlug != "" {

		memberDetailedProfile, err = memberInstance.GetMemberAndProfileData(0, "", 0, *profileSlug)
	}

	if err != nil {

		return &model.MemberProfile{}, err
	}

	if memberDetailedProfile.IsActive == 0 && memberDetailedProfile.Id != 0 {

		if memberid != 0 && tokenType == LocalLoginType {

			return &model.MemberProfile{}, ErrMemberInactive

		} else if memberid == 0 {

			return &model.MemberProfile{}, ErrMemberInactive
		}
	}

	if memberDetailedProfile.TblMemberProfile.CompanyLogo != "" {

		memberDetailedProfile.TblMemberProfile.CompanyLogo = GetFilePathsRelatedToStorageTypes(db, memberDetailedProfile.TblMemberProfile.CompanyLogo)
	}

	MemberProfile := model.MemberProfile{
		ID:              &memberDetailedProfile.TblMemberProfile.Id,
		MemberID:        &memberDetailedProfile.TblMemberProfile.MemberId,
		ProfileName:     &memberDetailedProfile.TblMemberProfile.ProfileName,
		ProfileSlug:     &memberDetailedProfile.TblMemberProfile.ProfileSlug,
		ProfilePage:     &memberDetailedProfile.TblMemberProfile.ProfilePage,
		MemberDetails:   &memberDetailedProfile.TblMemberProfile.MemberDetails,
		CompanyName:     &memberDetailedProfile.TblMemberProfile.CompanyName,
		CompanyLocation: &memberDetailedProfile.TblMemberProfile.CompanyLocation,
		CompanyLogo:     &memberDetailedProfile.TblMemberProfile.CompanyLogo,
		About:           &memberDetailedProfile.TblMemberProfile.About,
		SeoTitle:        &memberDetailedProfile.TblMemberProfile.SeoTitle,
		SeoDescription:  &memberDetailedProfile.TblMemberProfile.SeoDescription,
		SeoKeyword:      &memberDetailedProfile.TblMemberProfile.SeoKeyword,
		CreatedBy:       &memberDetailedProfile.TblMemberProfile.CreatedBy,
		CreatedOn:       &memberDetailedProfile.TblMemberProfile.CreatedOn,
		ModifiedOn:      &memberDetailedProfile.TblMemberProfile.ModifiedOn,
		ModifiedBy:      &memberDetailedProfile.TblMemberProfile.ModifiedBy,
		Linkedin:        &memberDetailedProfile.TblMemberProfile.Linkedin,
		Twitter:         &memberDetailedProfile.TblMemberProfile.Twitter,
		Website:         &memberDetailedProfile.TblMemberProfile.Website,
		ClaimStatus:     &memberDetailedProfile.TblMemberProfile.ClaimStatus,
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

	memberInstance := GetMemberInstance()

	if err := memberInstance.MemberPasswordUpdate(newPassword, confirmPassword, oldPassword, memberId, memberId); err != nil {

		return false, err
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

	memberInstance := GetMemberInstance()

	memberDetails, err := memberInstance.GetMemberDetails(memberId)

	if err != nil {

		return &model.Member{}, err
	}

	if memberDetails.ProfileImagePath != "" {

		memberDetails.ProfileImagePath = GetFilePathsRelatedToStorageTypes(db, memberDetails.ProfileImagePath)

	}

	conv_Member := model.Member{
		ID:               memberDetails.Id,
		FirstName:        memberDetails.FirstName,
		LastName:         memberDetails.LastName,
		Email:            memberDetails.Email,
		MobileNo:         memberDetails.MobileNo,
		IsActive:         memberDetails.IsActive,
		ProfileImage:     memberDetails.ProfileImage,
		ProfileImagePath: memberDetails.ProfileImagePath,
		CreatedOn:        memberDetails.CreatedOn,
		CreatedBy:        memberDetails.CreatedBy,
		ModifiedOn:       &memberDetails.ModifiedOn,
		ModifiedBy:       &memberDetails.ModifiedBy,
		MemberGroupID:    memberDetails.MemberGroupId,
		Password:         &memberDetails.Password,
		Username:         &memberDetails.Username,
	}

	return &conv_Member, nil

}

func MemberProfileUpdate(db *gorm.DB, ctx context.Context, profiledata model.ProfileData) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		c.AbortWithStatus(http.StatusUnauthorized)

		return false, ErrLoginReq
	}

	companyData := make(map[string]interface{})

	var err error

	if profiledata.CompanyLogo.IsSet() && profiledata.CompanyLogo.Value() != nil {

		var fileName, filePath string

		storageType, _ := GetStorageType(db)

		fileName = profiledata.CompanyLogo.Value().Filename

		file := profiledata.CompanyLogo.Value().File

		if storageType.SelectedType == "aws" {

			fmt.Printf("aws-S3 storage selected\n")

			filePath = "member/" + fileName

			err = storage.UploadFileS3(storageType.Aws, profiledata.CompanyLogo.Value(), filePath)

			if err != nil {

				ErrorLog.Printf("company profile logo update failed in s3 error: %s", err)

				fmt.Printf("image upload failed %v\n", err)

				return false, ErrUpload

			}

		} else if storageType.SelectedType == "local" {

			fmt.Printf("local storage selected\n")

			b64Data, err := IoReadSeekerToBase64(file)

			if err != nil {

				ErrorLog.Printf("base64 conversion error: %s", err)

				return false, err
			}

			endpoint := "gqlSaveLocal"

			url := PathUrl + endpoint

			filePath, err = storage.UploadImageToAdminLocal(b64Data, fileName, url)

			if err != nil {

				ErrorLog.Printf("company profile logo upload failed in admin panel local error: %s", err)

				return false, ErrUpload
			}

		} else if storageType.SelectedType == "azure" {

			fmt.Printf("azure storage selected")

		} else if storageType.SelectedType == "drive" {

			fmt.Println("drive storage selected")
		}

		companyData["company_logo"] = filePath
	}

	companyData["company_name"] = profiledata.CompanyName

	companyData["profile_name"] = profiledata.ProfileName

	companyData["profile_slug"] = profiledata.ProfileSlug

	if profiledata.CompanyLocation.IsSet() && profiledata.CompanyLocation.Value() != nil {

		companyData["company_location"] = *profiledata.CompanyLocation.Value()
	}

	if profiledata.Website.IsSet() && profiledata.Website.Value() != nil {

		companyData["website"] = *profiledata.Website.Value()
	}

	if profiledata.Linkedin.IsSet() && profiledata.Linkedin.Value() != nil {

		companyData["linkedin"] = *profiledata.Linkedin.Value()
	}

	if profiledata.Twitter.IsSet() && profiledata.Twitter.Value() != nil {

		companyData["twitter"] = *profiledata.Twitter.Value()
	}

	if profiledata.SeoTitle.IsSet() && profiledata.SeoTitle.Value() != nil {

		companyData["seo_title"] = *profiledata.SeoTitle.Value()
	}

	if profiledata.SeoDescription.IsSet() && profiledata.SeoDescription.Value() != nil {

		companyData["seo_description"] = *profiledata.SeoDescription.Value()
	}

	if profiledata.SeoKeyword.IsSet() && profiledata.SeoKeyword.Value() != nil {

		companyData["seo_keyword"] = *profiledata.SeoKeyword.Value()
	}

	if profiledata.About.IsSet() && profiledata.About.Value() != nil {

		companyData["about"] = *profiledata.About.Value()
	}

	if profiledata.CompanyProfile.IsSet() && profiledata.CompanyProfile.Value() != nil {

		var jsonData map[string]interface{}

		err := json.Unmarshal([]byte(*profiledata.CompanyProfile.Value()), &jsonData)

		if err != nil {

			ErrorLog.Printf("company profile update error: %s", err)

			return false, err
		}

		companyData["member_details"] = jsonData

	}

	memberInstance := GetMemberInstance()

	if err = memberInstance.MemberProfileFlexibleUpdate(companyData, memberid, memberid); err != nil {

		return false, err
	}

	return true, nil
}

func VerifyProfileName(db *gorm.DB, ctx context.Context, profileSlug string, profileID int) (bool, error) {

	if profileSlug == "" || profileID < 0 {

		return false, nil
	}

	memberInstance := GetMemberInstanceWithoutAuth()

	slugPresence := memberInstance.CheckProfileSlug(profileSlug, profileID)

	return slugPresence, nil
}

func Memberclaimnow(db *gorm.DB, ctx context.Context, profileData model.ClaimData, profileId *int, profileSlug *string) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberInstance := GetMemberInstanceWithoutAuth()

	verify_chan := make(chan error)

	var(

	 MemberDetails memberPkg.Tblmember

	 err error

	)

	if *profileId != 0{

		MemberDetails, err = memberInstance.GetMemberAndProfileData(0,"",*profileId,"")

	}else if *profileSlug != ""{

		MemberDetails, err = memberInstance.GetMemberAndProfileData(0,"",0,*profileSlug)
	}

	if MemberDetails.TblMemberProfile.ClaimStatus == 1 {

		return false, ErrclaimAlready
	}

	if MemberDetails.IsActive != 1 {

		return false, ErrMemberInactive
	}

	var memberSettings model.MemberSettings

	if err := db.Debug().Table("tbl_member_settings").First(&memberSettings).Error; err != nil {

		return false, err
	}

	var convIds []int

	adminIds := strings.Split(memberSettings.NotificationUsers, ",")

	for _, adminId := range adminIds {

		convId, _ := strconv.Atoi(adminId)

		convIds = append(convIds, convId)
	}

	_, notifyEmails, _ := GetNotifyAdminEmails(db, convIds)

	var claimTemplate model.EmailTemplate

	if err := db.Debug().Table("tbl_email_templates").Where("is_deleted=0 and template_name = ?", OwndeskClaimnowTemplate).First(&claimTemplate).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	dataReplacer := strings.NewReplacer(
		"{OwndeskLogo}", EmailImagePath.Owndesk,
		"{Username}", "Admin",
		"{CompanyName}", MemberDetails.TblMemberProfile.CompanyName,
		"{ProfileName}", profileData.ProfileName,
		"{ProfileSlug}", profileData.ProfileSlug,
		"{WorkMail}", profileData.WorkMail,
		"{CompanyNumber}", profileData.CompanyNumber,
		"{PersonName}", profileData.PersonName,
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

	integratedBody := dataReplacer.Replace(claimTemplate.TemplateMessage)

	htmlBody := template.HTML(integratedBody)

	tmpl, err := template.ParseFiles("./view/email/email-template.html")

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	var template_buffers bytes.Buffer

	if err := tmpl.Execute(&template_buffers, gin.H{"body": htmlBody}); err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	modifiedSubject := strings.TrimSuffix(claimTemplate.TemplateSubject, "{CompanyName}") + MemberDetails.TblMemberProfile.CompanyName

	mail_data := MailConfig{Emails: notifyEmails, MailUsername: os.Getenv("MAIL_USERNAME"), MailPassword: os.Getenv("MAIL_PASSWORD"), Subject: modifiedSubject}

	html_content := template_buffers.String()

	go SendMail(mail_data, html_content, verify_chan)

	if <-verify_chan == nil {

		return true, nil

	} else {

		c.AbortWithError(500, <-verify_chan)

		return false, <-verify_chan
	}
}
