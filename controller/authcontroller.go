package controller

import (
	"context"
	"encoding/json"
	"errors"
	"gqlserver/graph/model"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spurtcms/pkgcore/member"

	// "gorm.io/datatypes"
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

	otp_expiry_time := time.Now().Add(5*time.Minute).In(TimeZone).Format("2006-01-02 15:04:05")

	err = Mem.StoreGraphqlMemberOtp(otp,conv_member.ID,otp_expiry_time)

	if err!=nil{

		return false, err
	}

	go SendMail(conv_member,otp,channel)

	if <-channel{

		return true,nil

	}else{

		return false,nil
	}
}

func VerifyMemberOtp(db *gorm.DB,ctx context.Context,otp int)(string,error){

	Mem.Auth = GetAuthorizationWithoutToken(db)

	currentTime := time.Now().In(TimeZone).Unix()

	token,err := Mem.VerifyLoginOtp(otp,currentTime,os.Getenv("JWT_SECRET"))

	if err!=nil{

		return "",err
	}

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

func MemberProfileUpdate(db *gorm.DB, ctx context.Context, profiledata model.ProfileData) (bool, error) {

	var jsonData map[string]interface{}

	err := json.Unmarshal([]byte(profiledata.MemberProfile), &jsonData)

	if err != nil {

		log.Println(err)
	}

	if err := db.Debug().Table("tbl_member_profiles").Where("id = ?", profiledata.MemberProfileID).UpdateColumns(map[string]interface{}{"member_details": jsonData}).Error;err!=nil{

		return false,err
	}

	return true, nil
}

