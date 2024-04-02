package controller

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"log"
	"math/rand"
	"os"
	"spurtcms-graphql/graph/model"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/spurtcms/pkgcontent/channels"
	"github.com/spurtcms/pkgcore/member"
	"gorm.io/gorm"
)

func MemberLogin(db *gorm.DB, ctx context.Context, email string) (bool, error) {

	Mem.Auth = GetAuthorizationWithoutToken(db)

	member_details, err := Mem.GraphqlMemberLogin(email)

	if err != nil {

		return false, err
	}

	conv_member := model.Member{
		ID: member_details.Id,
		FirstName: member_details.FirstName,
		LastName: member_details.LastName,
		Email: member_details.Email,
		MobileNo: member_details.MobileNo,
		IsActive: member_details.IsActive,
		ProfileImagePath: member_details.ProfileImagePath,
		CreatedOn: member_details.CreatedOn,
		CreatedBy: member_details.CreatedBy,
		ModifiedOn: &member_details.ModifiedOn,
		ModifiedBy: &member_details.ModifiedBy,
	}

    channel := make(chan bool)

	// rand.Seed(time.Now().UnixNano())

    otp := rand.Intn(900000) + 100000 

	current_time := time.Now().In(TimeZone)

	otp_expiry_time := current_time.Add(5*time.Minute).Format("2006-01-02 15:04:05")

	mail_expiry_time := current_time.Add(5*time.Minute).Format("02 Jan 2006 03:04 PM")

	err = Mem.StoreGraphqlMemberOtp(otp,conv_member.ID,otp_expiry_time)

	if err!=nil{

		return false, err
	}

	data := map[string]interface{}{"otp": otp,"expiryTime": mail_expiry_time,"member": conv_member,"additionalData": AdditionalData}

	tmpl, err := template.ParseFiles("view/email/login-template.html")

	if err != nil {
		
		return false, err
	}

	var template_buffer bytes.Buffer
	
	if err := tmpl.Execute(&template_buffer,data); err != nil {

		return false,err
	}

	mail_data := MailConfig{Email: conv_member.Email,MailUsername: os.Getenv("MAIL_USERNAME"),MailPassword: os.Getenv("MAIL_PASSWORD"),Subject: "OwnDesk - Login Otp Confirmation"}

	html_content := template_buffer.String()

	go SendMail(mail_data,html_content,channel)

	if <-channel{

		return true,nil

	}else{

		return false,nil
	}
}

func VerifyMemberOtp(db *gorm.DB,ctx context.Context,email string,otp int)(string,error){

	Mem.Auth = GetAuthorizationWithoutToken(db)

	currentTime := time.Now().In(TimeZone).Unix()

	memberDetails,token,err := Mem.VerifyLoginOtp(email,otp,currentTime)

	if err!=nil{

		return "",err
	}

	log.Println("memberdetails",memberDetails)

	var channelEntryDetails model.ChannelEntries

	if err := db.Debug().Table("tbl_channel_entries").Select("tbl_channel_entries.*").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id ").Joins("inner join tbl_channel_entry_fields on tbl_channel_entry_fields.channel_entry_id = tbl_channel_entries.id").Joins("inner join tbl_fields on tbl_fields.id = tbl_channel_entry_fields.field_id").Joins("inner join tbl_field_types on tbl_field_types.id = tbl_fields.field_type_id").Joins("inner join tbl_members on tbl_members.id = any(string_to_array(tbl_channel_entry_fields.field_value,',')::integer[])").
	Joins("inner join tbl_member_profiles on tbl_members.id = tbl_member_profiles.member_id").Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_field_types.is_deleted = 0 and tbl_fields.is_deleted = 0 and tbl_members.is_deleted = 0 and tbl_member_profiles.is_deleted = 0 and tbl_member_profiles.claim_status = 1 and tbl_field_types.id = ? and tbl_members.id = ?",MemberFieldTypeId,memberDetails.Id).First(&channelEntryDetails).Error;err!=nil{

		return "",err
	}

	log.Println("chkking",channelEntryDetails)

	// channelAuth := channels.Channel{Authority: GetAuthorization(token,db)}

	// channelEntry,err :=  channelAuth.GetGraphqlChannelEntriesDetails(&channelEntryDetails.ID,&channelEntryDetails.ChannelID,nil,PathUrl,SectionTypeId,MemberFieldTypeId,nil)

	return  token,nil
	
}

func MemberRegister(db *gorm.DB, input model.MemberDetails) (bool, error) {

	Mem.Auth = GetAuthorizationWithoutToken(db)

	var imageName, imagePath string

	var err error

	if input.ProfileImage != nil {

		imageName, imagePath, err = StoreImageBase64ToLocal(*input.ProfileImage, ProfileImagePath, "PROFILE")

		if err != nil {

			return false, err
		}

	}

	memberDetails := member.MemberCreation{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		MobileNo:  input.Mobile,
		Email:     input.Email,
		Password:  input.Password,
		//  Username: *input.Username,
		ProfileImage:     imageName,
		ProfileImagePath: imagePath,
	}

	_, isMemberExists, err := Mem.CheckEmailInMember(0, input.Email)

	if isMemberExists {

		return isMemberExists, errors.New("Member already exists!")
	}

	isRegistered, err := Mem.MemberRegister(memberDetails)

	if !isRegistered || err != nil {

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

	if memberdata.ProfileImage != nil {

		imageName, imagePath, err = StoreImageBase64ToLocal(*memberdata.ProfileImage, ProfileImagePath, "PROFILE")

		if err != nil {

			return false, err
		}

	}

	memberDetails := member.MemberCreation{
		FirstName:        memberdata.FirstName,
		LastName:         memberdata.LastName,
		MobileNo:         memberdata.Mobile,
		Email:            memberdata.Email,
		Password:         memberdata.Password,
		ProfileImage:     imageName,
		ProfileImagePath: imagePath,
		// IsActive: *memberdata.IsActive,
		// Username: *memberdata.Username,
		// GroupId: *memberdata.GroupID,
	}

	isUpdated, err := Mem.MemberUpdate(memberDetails)

	if err != nil || !isUpdated {

		return isUpdated, err
	}

	return isUpdated, nil

}


