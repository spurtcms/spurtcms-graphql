package controller

import (
	"context"
	"errors"
	"gqlserver/graph/model"
	"math/rand"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spurtcms/pkgcore/member"
	"gorm.io/gorm"
)

func MemberLogin(db *gorm.DB ,input model.LoginCredentials) (string, error) {

	token, err := Mem.CheckMemberLogin(member.MemberLogin{Emailid: input.Email, Password: input.Password}, db, os.Getenv("JWT_SECRET"))

	if err!=nil{

		return "",err
	}
	
	Auth = GetAuthorization(token,db)

	return token,nil
}

func MemberRegister(db *gorm.DB,input model.MemberDetails)(bool,error){

	Mem.Auth = GetAuthorizationWithoutToken(db)

	var imageName,imagePath string 

	var err error

	if input.ProfileImage != nil {

		imageName,imagePath,err = StoreImageBase64ToLocal(*input.ProfileImage,ProfileImagePath,"PROFILE")

		if err!=nil{

			return false,err
		}
		
	}

	memberDetails := member.MemberCreation{
		             FirstName: input.FirstName,
					 LastName: input.LastName,
					 MobileNo: input.Mobile,
					 Email: input.Email,
					 Password: input.Password,
					//  Username: *input.Username,
					 ProfileImage: imageName,
					 ProfileImagePath: imagePath,

	}

	_,isMemberExists,err := Mem.CheckEmailInMember(0,input.Email)

	if isMemberExists{

		return isMemberExists,errors.New("Member already exists!")
	}

	isRegistered, err := Mem.MemberRegister(memberDetails)

	if !isRegistered || err!= nil{

		return isRegistered,err
	}

	return isRegistered,nil

}

func UpdateMember(db *gorm.DB ,ctx context.Context,memberdata model.MemberDetails)(bool,error){

	c,_ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	Mem.Auth = GetAuthorization(token,db)

	var imageName,imagePath string 

	var err error

	if memberdata.ProfileImage!=nil{

		imageName,imagePath,err = StoreImageBase64ToLocal(*memberdata.ProfileImage,ProfileImagePath,"PROFILE")

		if err!=nil{
	
			return false,err
		}

	}

	memberDetails := member.MemberCreation{
		FirstName: memberdata.FirstName,
		LastName: memberdata.LastName,
		MobileNo: memberdata.Mobile,
		Email: memberdata.Email,
		Password: memberdata.Password,
		ProfileImage: imageName,
		ProfileImagePath: imagePath,
		// IsActive: *memberdata.IsActive,
		// Username: *memberdata.Username,
		// GroupId: *memberdata.GroupID,
    }

	isUpdated, err := Mem.MemberUpdate(memberDetails)

	if err!=nil || !isUpdated{

		return isUpdated,err
	}

	return isUpdated,nil

}

func SendOtpToMail(db *gorm.DB, ctx context.Context, email string) (bool, error) {

	Mem.Auth = GetAuthorizationWithoutToken(db)

	member,isValidMember,err := Mem.CheckEmailInMember(0,email)

	if err!=nil || !isValidMember{

		return isValidMember,err
	}

	randNum := rand.Intn(900000) + 100000

	otp := strconv.Itoa(randNum)

	isOtpStored,err := Mem.UpdateOtp(randNum,member.Id)

	if !isOtpStored || err!=nil{

		return isOtpStored,errors.New("Failed to generate otp!")
	}

	ch := make(chan error)

	subject := "Forgot Password OTP Verification"

	body := `<html>
	        <head>
		         <meta charset="UTF-8">
		         <meta name="viewport" content="width=device-width, initial-scale=1.0">
		    <style>
			     body {
				       font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
				       background-color: #f4f4f4;
				       margin: 0;
				       padding: 0;
				       display: flex;
				       align-items: center;
				       justify-content: center;
				       height: 100vh;
			         }
	
			     .container {
				            max-width: 400px;
				            background-color: #fff;
				            padding: 30px;
				            border-radius: 10px;
				            box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
			               }
	
				   h2{
						text-align: center;
						color: #3498db;
					}
	
			       p {
				      color: #555;
				      line-height: 1.6;
				      margin-bottom: 20px;
			         }
	
			      .otp {
				        text-align: center;
				        font-size: 36px;
				        font-weight: bold;
				        margin: 20px 0;
				        color: #3498db;
			          }
	
			      .note {
				        text-align: center;
				        color: #888;
			            }
	
			      .button {
				            display: inline-block;
				            padding: 10px 20px;
				            background-color: #3498db;
				            color: #fff;
				            text-decoration: none;
			             	border-radius: 5px;
			            }
	
			           .button:hover {
				                background-color: #2980b9;
			                         }
		       </style
	           </head>
	        <body>
		          <div class="container">
			        <p>Hello `+member.FirstName+` `+member.LastName+`,</p>
                    <p>Your OTP for password reset is:</p>
                    <div class="otp">[`+otp+`]</div>
                    <p class="note">Note: This OTP is valid for only 5 minutes. Please do not share it with anyone.</p>
                    <p class="note">If you didn't request a password reset, please ignore this email.</p>
                 </div>
            </body>
            </html>`

	go SendEmail(member.Email,subject,body,ch)

	send_err := <-ch

	if send_err!=nil {

		return false,send_err
	}

	return true,nil
}

func ResetPassword(db *gorm.DB, otp int, newPassword, email string) (bool, error) {

	Mem.Auth = GetAuthorizationWithoutToken(db)

	member,isvalidMail,err := Mem.CheckEmailInMember(0,email)

	if !isvalidMail || err!=nil{

		return isvalidMail,err
	}

	isPswdChanged,err := Mem.ChangePassword(otp,member.Id,newPassword)

	if !isPswdChanged || err!=nil{

		return isPswdChanged,err
	}

	return isPswdChanged,nil
}

