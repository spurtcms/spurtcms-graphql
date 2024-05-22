package controller

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"spurtcms-graphql/graph/model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spurtcms/pkgcore/member"
	"gorm.io/gorm"
)

func MemberLogin(db *gorm.DB, ctx context.Context, email string) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	Mem.Auth = GetAuthorizationWithoutToken(db)

	member_details, err := Mem.GraphqlMemberLogin(email)

	if member_details.IsActive==0 && member_details.Id != 0{

		return false,ErrMemberInactive
	}

	if gorm.ErrRecordNotFound == err {

		var loginEnquiryTemplate  model.EmailTemplate

		if err := db.Debug().Table("tbl_email_templates").Where("is_deleted = 0 and template_name = ?",OwndeskLoginEnquiryTemplate).First(&loginEnquiryTemplate).Error;err!=nil{

			c.AbortWithError(http.StatusInternalServerError, err)
			
			return false,err
		}

		adminDetails, _ := Mem.GetAdminDetails(OwndeskChannelId)

		var admin_mail_data = MailConfig{Email: adminDetails.Email, MailUsername: os.Getenv("MAIL_USERNAME"), MailPassword: os.Getenv("MAIL_PASSWORD"), Subject: loginEnquiryTemplate.TemplateSubject}

		dataReplacer := strings.NewReplacer(
			"{OwndeskLogo}", EmailImagePath.Owndesk,
			"{Username}", adminDetails.Username,
			"{UnauthorizedMail}", email,
			"{CurrentTime}",time.Now().In(TimeZone).Format("02 Jan 2006 03:04 PM"),
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
			"<figure","<div",
			"</figure","</div",
			"&nbsp;","",          
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

		if err := tmpl.Execute(&template_buffers,gin.H{"body":htmlBody}); err != nil {

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

	}else if err!=nil {

		return false, err
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

	if err := db.Debug().Table("tbl_email_templates").Where("is_deleted=0 and template_name = ?",OwndeskLoginTemplate).First(&loginTemplate).Error;err!=nil{

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	dataReplacer := strings.NewReplacer(
		"{OwndeskLogo}", EmailImagePath.Owndesk,
		"{Username}", member_details.Username,
		"{CompanyName}", memberProfileData.CompanyName,
		"{Otp}",strconv.Itoa(otp),
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
		"<figure","<div",
		"</figure","</div",
		"&nbsp;","",          
	)

	integratedBody := dataReplacer.Replace(loginTemplate.TemplateMessage)

	htmlBody := template.HTML(integratedBody)

	if err !=nil{

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

	mail_data := MailConfig{Email: member_details.Email, MailUsername: os.Getenv("MAIL_USERNAME"), MailPassword: os.Getenv("MAIL_PASSWORD"), Subject: loginTemplate.TemplateSubject}

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

	memberDetails, token, err := Mem.VerifyLoginOtp(email, otp, currentTime)

	if err != nil {

		return &model.LoginDetails{}, err
	}

	var memberProfileDetails model.MemberProfile

	if err := db.Debug().Table("tbl_member_profiles").Select("tbl_member_profiles.*").Joins("inner join tbl_members on tbl_members.id = tbl_member_profiles.member_id").Where("tbl_member_profiles.is_deleted = 0 and tbl_members.is_deleted = 0 and  tbl_members.is_active =1 and tbl_member_profiles.member_id = ?", memberDetails.Id).First(&memberProfileDetails).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.LoginDetails{}, err
	}

	return &model.LoginDetails{MemberProfileData: memberProfileDetails, Token: token}, nil

}

func MemberRegister(db *gorm.DB, ctx context.Context, input model.MemberDetails, ecomModule *int) (bool, error) {

	if MemberRegisterPermission == "false"{

		return false, ErrMemberRegisterPerm
	}

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	Mem.Auth = GetAuthorizationWithoutToken(db)

	var (

		imageName, imagePath string

		err error

		ecomMod int = *ecomModule
	)

	if input.ProfileImage.IsSet() {

		imageName, imagePath, err = StoreImageBase64ToLocal(*input.ProfileImage.Value(), ProfileImagePath, "PROFILE")

		if err != nil {

			c.AbortWithError(http.StatusInternalServerError, err)

			return false, err
		}

	}

	var memberDetails member.MemberCreation

	if input.Mobile.IsSet() {

		memberDetails.MobileNo = *input.Mobile.Value()
	}

	if imageName != "" && imagePath != "" {

		memberDetails.ProfileImage = imageName

		memberDetails.ProfileImagePath = imagePath
	}

	if input.LastName.IsSet() {

		memberDetails.LastName = *input.LastName.Value()

	}

	if input.Username.IsSet() {

		memberDetails.Username = *input.Username.Value()

		_, isMemberExists, err := Mem.CheckUsernameInMember(0, *input.Username.Value())

		if isMemberExists || err == nil {

			err = errors.New("member already exists")

			c.AbortWithError(422, err)

			return isMemberExists, err
		}

		if ecomMod ==1{

			var count int64

			if err := db.Table("tbl_ecom_customers").Where("is_deleted = 0 and username = ?",*input.Username.Value()).Count(&count).Error;err!=nil{

				return false,err
			}

			if count > 0{

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

		if ecomMod ==1{

			var count int64

			if err := db.Table("tbl_ecom_customers").Where("is_deleted = 0 and email = ?",input.Email).Count(&count).Error;err!=nil{

				return false,err
			}

			if count > 0{

				err = errors.New("customer already exists")

			    c.AbortWithError(422, err)

			    return false, err
			}
		}

	}

	memberDetails.FirstName = input.FirstName

	memberDetails.Password = input.Password

	memberData,isRegistered, err := Mem.MemberRegister(memberDetails)

	if !isRegistered || err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return isRegistered, err
	}

	if isRegistered && ecomMod == 1{

		createdOn, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		var ecomCustomer = model.CustomerDetails{
			FirstName:        memberDetails.FirstName,
			LastName:         &memberDetails.LastName,
			Username:         memberDetails.Email,
			MobileNo:         memberDetails.MobileNo,
			Email:            input.Email,
			IsActive:         memberDetails.IsActive,
			ProfileImage:     &memberDetails.ProfileImage,
			ProfileImagePath: &memberDetails.ProfileImagePath,
			CreatedOn:        createdOn,
			Password:         HashingPassword(memberDetails.Password),
			MemberID:         &memberData.Id,
		}

		if err := db.Table("tbl_ecom_customers").Create(&ecomCustomer).Error;err!=nil{

			c.AbortWithError(http.StatusInternalServerError, err)
	
			return isRegistered, err
		}
	}

	return true, nil

}

func UpdateMember(db *gorm.DB, ctx context.Context, memberdata model.MemberDetails) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	Mem.Auth = GetAuthorization(token, db)

	var imageName, imagePath string

	var err error

	if memberdata.ProfileImage.IsSet() {

		imageName, imagePath, err = StoreImageBase64ToLocal(*memberdata.ProfileImage.Value(), ProfileImagePath, "PROFILE")

		if err != nil {

			return false, err
		}

	}

	memberDetails := member.MemberCreation{
		FirstName:        memberdata.FirstName,
		LastName:         *memberdata.LastName.Value(),
		MobileNo:         *memberdata.Mobile.Value(),
		Email:            memberdata.Email,
		Password:         memberdata.Password,
		ProfileImage:     imageName,
		ProfileImagePath: imagePath,
		// IsActive: *memberdata.IsActive,
		// Username: *memberdata.Username.Value(),
		// GroupId: *memberdata.GroupID,
	}

	isUpdated, err := Mem.MemberUpdate(memberDetails)

	if err != nil || !isUpdated {

		return isUpdated, err
	}

	return isUpdated, nil

}

func TemplateMemberLogin(db *gorm.DB, ctx context.Context, username, email *string, password string) (string, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	Mem.Auth = GetAuthorizationWithoutToken(db)

	var memberLogin member.MemberLogin

	if username != nil {

		memberLogin.Username = *username

	} else if email != nil {

		memberLogin.Emailid = *email
	}

	memberLogin.Password = password

	token, err := Mem.CheckMemberLogin(memberLogin, db, os.Getenv("JWT_SECRET"))

	if err != nil {

		c.AbortWithError(http.StatusUnauthorized, err)

		log.Println(err)
	}

	return token, err
}

func MemberProfileDetails(db *gorm.DB, ctx context.Context) (*model.MemberProfile, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return &model.MemberProfile{}, err

	}

	var memberProfile model.MemberProfile

	if err := db.Table("tbl_member_profiles").Where("is_deleted = 0 and member_id = ?", memberid).First(&memberProfile).Error; err != nil {

		c.AbortWithError(http.StatusUnprocessableEntity, err)

		return &model.MemberProfile{}, err
	}

	return &memberProfile, nil
}

func GetMemberProfileDetails(db *gorm.DB, ctx context.Context, id *int, profileSlug *string) (*model.MemberProfile, error) {

	var memberProfile member.TblMemberProfile

	query := db.Table("tbl_member_profiles").Where("is_deleted = 0")

	if id != nil {

		query = query.Where("member_id = ?", *id)

	} else if profileSlug != nil {

		query = query.Where("profile_slug = ?", *profileSlug)
	}

	if err := query.First(&memberProfile).Error; err != nil {

		return &model.MemberProfile{}, err
	}

	var profileLogo string

	if memberProfile.CompanyLogo!=""{

		profileLogo = PathUrl + strings.TrimPrefix(memberProfile.CompanyLogo,"/")
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
