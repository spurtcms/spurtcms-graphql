package controller

import (
	"bytes"
	"context"

	// "encoding/base64"
	"errors"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"spurtcms-graphql/graph/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spurtcms/pkgcore/member"
	"gorm.io/gorm"
)

func MemberLogin(db *gorm.DB, ctx context.Context, email string) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	Mem.Auth = GetAuthorizationWithoutToken(db)

	member_details, err := Mem.GraphqlMemberLogin(email)

	if gorm.ErrRecordNotFound == err {

		adminDetails, _ := Mem.GetAdminDetails(OwndeskChannelId)

		var admin_mail_data = MailConfig{Email: adminDetails.Email, MailUsername: os.Getenv("MAIL_USERNAME"), MailPassword: os.Getenv("MAIL_PASSWORD"), Subject: "Notification: New User Attempted Login to OwnDesk Platform"}

		channel := make(chan error)

		tmpls, err := template.ParseFiles("view/email/admin-loginenquiry.html")

		if err != nil {

			c.AbortWithError(http.StatusInternalServerError, err)

			return false, err
		}

		var template_buffers bytes.Buffer

		if err := tmpls.Execute(&template_buffers, map[string]interface{}{"adminDetails": adminDetails, "unauthorizedMail": email, "currentTime": time.Now().In(TimeZone).Format("02 Jan 2006 03:04 PM")}); err != nil {

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

	conv_member := model.Member{
		ID:               member_details.Id,
		FirstName:        member_details.FirstName,
		LastName:         member_details.LastName,
		Email:            member_details.Email,
		MobileNo:         member_details.MobileNo,
		IsActive:         member_details.IsActive,
		ProfileImagePath: member_details.ProfileImagePath,
		CreatedOn:        member_details.CreatedOn,
		CreatedBy:        member_details.CreatedBy,
		ModifiedOn:       &member_details.ModifiedOn,
		ModifiedBy:       &member_details.ModifiedBy,
	}

	var memberProfileData model.MemberProfile

	if err := db.Debug().Table("tbl_member_profiles").Where("is_deleted = 0 and member_id = ?", conv_member.ID).First(&memberProfileData).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	channel := make(chan error)

	// rand.Seed(time.Now().UnixNano())

	otp := rand.Intn(900000) + 100000

	current_time := time.Now().In(TimeZone)

	otp_expiry_time := current_time.Add(5 * time.Minute).Format("2006-01-02 15:04:05")

	mail_expiry_time := current_time.Add(5 * time.Minute).Format("02 Jan 2006 03:04 PM")

	err = Mem.StoreGraphqlMemberOtp(otp, conv_member.ID, otp_expiry_time)

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	data := map[string]interface{}{"otp": otp, "expiryTime": mail_expiry_time, "member": conv_member, "additionalData": AdditionalData, "memberProfile": memberProfileData}

	tmpl, err := template.ParseFiles("view/email/login-template.html")

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	var template_buffer bytes.Buffer

	if err := tmpl.Execute(&template_buffer, data); err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	mail_data := MailConfig{Email: conv_member.Email, MailUsername: os.Getenv("MAIL_USERNAME"), MailPassword: os.Getenv("MAIL_PASSWORD"), Subject: "OwnDesk - Login Otp Confirmation"}

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

	currentTime := time.Now().In(TimeZone).Unix()

	memberDetails, token, err := Mem.VerifyLoginOtp(email, otp, currentTime)

	if err != nil {

		return &model.LoginDetails{}, err
	}

	var memberProfileDetails model.MemberProfile

	if err := db.Debug().Table("tbl_member_profiles").Select("tbl_member_profiles.*").Where("tbl_member_profiles.is_deleted = 0 and tbl_member_profiles.member_id = ? and tbl_member_profiles.claim_status = 1", memberDetails.Id).First(&memberProfileDetails).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.LoginDetails{}, err
	}

	return &model.LoginDetails{MemberProfileData: memberProfileDetails, Token: token}, nil

}

func MemberRegister(db *gorm.DB, ctx context.Context, input model.MemberDetails) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	Mem.Auth = GetAuthorizationWithoutToken(db)

	var imageName, imagePath string

	var err error

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

			c.AbortWithError(http.StatusBadRequest, err)

			return isMemberExists, err
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
	}

	memberDetails.FirstName = input.FirstName

	memberDetails.Password = input.Password

	isRegistered, err := Mem.MemberRegister(memberDetails)

	if !isRegistered || err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return isRegistered, err
	}

	return isRegistered, nil

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

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	var memberProfile member.TblMemberProfile

	query := db.Table("tbl_member_profiles").Where("is_deleted = 0")

	if id != nil {

		query = query.Where("member_id = ?", *id)

	} else if profileSlug != nil {

		query = query.Where("profile_slug = ?", *profileSlug)
	}

	if err := query.First(&memberProfile).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.MemberProfile{}, err
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
		CompanyLogo:     &memberProfile.About,
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
