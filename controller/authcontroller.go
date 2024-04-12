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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spurtcms/pkgcore/member"
	"gorm.io/gorm"
)

func MemberLogin(db *gorm.DB, ctx context.Context, email string) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	Mem.Auth = GetAuthorizationWithoutToken(db)

	member_details, err := Mem.GraphqlMemberLogin(email)

	if err != nil {

		adminDetails, _ := Mem.GetAdminDetails(OwndeskChannelId)

		var admin_mail_data = MailConfig{Email: adminDetails.Email, MailUsername: os.Getenv("MAIL_USERNAME"), MailPassword: os.Getenv("MAIL_PASSWORD"), Subject: "Notification: New User Attempted Login to OwnDesk Platform"}

		channel := make(chan error)

		tmpls, err := template.ParseFiles("view/email/admin-loginenquiry.html")

		if err != nil {

			err = errors.New("failed to send unauthorized login attempt mail to admin")

			c.AbortWithError(http.StatusUnauthorized, err)

			return false, err
		}

		var template_buffers bytes.Buffer

		if err := tmpls.Execute(&template_buffers, map[string]interface{}{"adminDetails": adminDetails, "unauthorizedMail": email,"currentTime": time.Now().In(TimeZone).Format("02 Jan 2006 03:04 PM")}); err != nil {

			err = errors.New("failed to send unauthorized login attempt mail to admin")

			c.AbortWithError(http.StatusUnauthorized, err)

			return false, err
		}

		admin_content := template_buffers.String()

		go SendMail(admin_mail_data, admin_content, channel)

		if <-channel != nil {

			err = errors.New("failed to send unauthorized login attempt mail to admin")

			c.AbortWithError(http.StatusUnauthorized, err)

			return false, err

		}

		c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized access"))

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

	if <-channel == nil {

		return true, nil

	} else {

		c.AbortWithError(http.StatusServiceUnavailable, <-channel)

		return false, <-channel
	}
}

func VerifyMemberOtp(db *gorm.DB, ctx context.Context, email string, otp int) (*model.LoginDetails, error) {

	Mem.Auth = GetAuthorizationWithoutToken(db)

	currentTime := time.Now().In(TimeZone).Unix()

	memberDetails, token, err := Mem.VerifyLoginOtp(email, otp, currentTime)

	if err != nil {

		return &model.LoginDetails{}, err
	}

	var channelEntryDetails model.ChannelEntries

	if err := db.Debug().Table("tbl_channel_entries").Select("tbl_channel_entries.*").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id ").Joins("inner join tbl_channel_entry_fields on tbl_channel_entry_fields.channel_entry_id = tbl_channel_entries.id").Joins("inner join tbl_fields on tbl_fields.id = tbl_channel_entry_fields.field_id").Joins("inner join tbl_field_types on tbl_field_types.id = tbl_fields.field_type_id").Joins("inner join tbl_members on tbl_members.id = any(string_to_array(tbl_channel_entry_fields.field_value,',')::integer[])").
		Joins("inner join tbl_member_profiles on tbl_members.id = tbl_member_profiles.member_id").Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_field_types.is_deleted = 0 and tbl_fields.is_deleted = 0 and tbl_members.is_deleted = 0 and tbl_member_profiles.is_deleted = 0 and tbl_member_profiles.claim_status = 1 and tbl_field_types.id = ? and tbl_members.id = ? and tbl_member_profiles.claim_status = 1", MemberFieldTypeId, memberDetails.Id).First(&channelEntryDetails).Error; err != nil {

		return &model.LoginDetails{}, err
	}

	var memberProfileDetails model.MemberProfile

	if err := db.Debug().Table("tbl_member_profiles").Select("tbl_member_profiles.*").Joins("inner join tbl_members on tbl_members.id = tbl_member_profiles.member_id").Joins("INNER JOIN TBL_CHANNEL_ENTRY_FIELDS ON TBL_MEMBERS.ID::text = tbl_channel_entry_fields.field_value").Joins("inner join tbl_channel_entries on tbl_channel_entry_fields.channel_entry_id = tbl_channel_entries.id").Joins("inner join tbl_fields on tbl_fields.id = tbl_channel_entry_fields.field_id").Joins("inner join tbl_field_types on tbl_field_types.id = tbl_fields.field_type_id").
	          Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id ").Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_field_types.is_deleted = 0 and tbl_fields.is_deleted = 0 and tbl_members.is_deleted = 0 and tbl_member_profiles.is_deleted = 0 and tbl_member_profiles.claim_status = 1 and tbl_field_types.id = ? and tbl_members.id = ? and tbl_member_profiles.claim_status = 1", MemberFieldTypeId, memberDetails.Id).First(&memberProfileDetails).Error; err != nil {

		return &model.LoginDetails{}, err
	}

	return &model.LoginDetails{ClaimEntryDetails: channelEntryDetails,MemberProfileData: memberProfileDetails, Token: token}, nil

}

func MemberRegister(db *gorm.DB,ctx context.Context, input model.MemberDetails) (bool, error) {

	c,_:= ctx.Value(ContextKey).(*gin.Context)

	Mem.Auth = GetAuthorizationWithoutToken(db)

	var imageName, imagePath string

	var err error

	if input.ProfileImage.IsSet() {

		imageName, imagePath, err = StoreImageBase64ToLocal(*input.ProfileImage.Value(), ProfileImagePath, "PROFILE")

		if err != nil {

			c.AbortWithError(http.StatusInternalServerError,err)

			return false, err
		}

	}

	var memberDetails member.MemberCreation

	if input.Mobile.IsSet(){

		memberDetails.MobileNo = *input.Mobile.Value()
	}

	if imageName!="" && imagePath!=""{

		memberDetails.ProfileImage = imageName

		memberDetails.ProfileImagePath = imagePath
	}

	memberDetails.FirstName = input.FirstName

	memberDetails.LastName = *input.LastName.Value()

	memberDetails.Email = input.Email

	memberDetails.Password = input.Password

	memberDetails.Username = input.Username

	_, isMemberExists, err := Mem.CheckUsernameInMember(0, input.Username)

	if isMemberExists || err == nil {

		err = errors.New("member already exists") 

		c.AbortWithError(http.StatusBadRequest,err)

		return isMemberExists, err
	}

	isRegistered, err := Mem.MemberRegister(memberDetails)

	if !isRegistered || err != nil {

		c.AbortWithError(http.StatusInternalServerError,err)

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

func TemplateMemberLogin(db *gorm.DB, ctx context.Context, username string, password string) (string, error) {

	Mem.Auth = GetAuthorizationWithoutToken(db)

	member_details, err := Mem.CheckMemberLogin(member.MemberLogin{Username: username, Password: password}, db, os.Getenv("JWT_SECRET"))

	if err != nil {

		log.Println(err)
	}

	return member_details, err
}
