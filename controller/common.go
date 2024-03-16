package controller

import (
	"encoding/base64"
	"gqlserver/graph/model"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spurtcms/pkgcore/auth"
	"github.com/spurtcms/pkgcore/member"
	"gorm.io/gorm"

	spurtcore "github.com/spurtcms/pkgcore"
)

type key string

const ContextKey key = "ginContext"

var(
	Mem member.MemberAuth
	Auth *auth.Authorization
	IST, _ = time.LoadLocation("Asia/Kolkata")
	ProfileImagePath = "Uploads/ProfileImages/"
	SpecialToken = "%$HEID$#PDGH*&MGEAFCC"
	SectionTypeId = 12
    MemberFieldTypeId = 14
)

func GetAuthorization(token string,db *gorm.DB)(*auth.Authorization) {

	auth := spurtcore.NewInstance(&auth.Option{DB: db, Token: token, Secret: os.Getenv("JWT_SECRET")})

	return &auth

}

func GetAuthorizationWithoutToken(db *gorm.DB)(*auth.Authorization){

	auth := spurtcore.NewInstance(&auth.Option{DB: db, Token: "", Secret: ""})

	return &auth
}

func StoreImageBase64ToLocal(imageData,storagePath,storingName string) (string,string,error) {

	extEndIndex := strings.Index(imageData, ";base64,")

	base64data := imageData[strings.IndexByte(imageData, ',')+1:]

	var ext = imageData[11:extEndIndex]

	randomNum := strconv.Itoa(rand.Intn(900000) + 100000)

	imageName := storingName +"-"+ randomNum + "." + ext

	err := os.MkdirAll(storagePath, 0755)

	if err!=nil{

		log.Println(err)

		return "","",err
	}

	storageDestination := storagePath + imageName
	
	decode, err := base64.StdEncoding.DecodeString(base64data)

	if err != nil {

		log.Println(err)

		return "","",err
	}

	file, err := os.Create(storageDestination)

	if err != nil {

		log.Println(err)

		return "","",err

	}
	if _, err := file.Write(decode); err != nil {

		log.Println(err)

		return "","",err

	}

	return imageName,storageDestination,nil
}

func SendMail(member model.Member, otp int,channel chan bool) {

	// Sender data.
	from := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")

	// Receiver email address.
	to := []string{
		member.Email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	subject := "Subject:Hello " + member.FirstName +" "+ member.LastName + "\n"

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	conv_otp := strconv.Itoa(otp)

	body := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>OTP Email</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				margin: 0;
				padding: 0;
			}
			.container {
				max-width: 600px;
				margin: 20px auto;
				background-color: #fff;
				padding: 20px;
				border-radius: 5px;
				box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			}
			h2 {
				color: #333;
			}
			p {
				color: #666;
			}
			.otp {
				font-size: 24px;
				font-weight: bold;
				color: #007bff;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>OTP Email</h2>
			<p>Your One-Time Password (OTP) is:</p>
			<p class="otp">`+conv_otp+`</p>
			<p>Please use this OTP to proceed to Login into ownDesk</p>
		</div>
	</body>
	</html>`

	msg := []byte(subject + mime + body)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)

	if err != nil {

		log.Println(err)

		channel <- false

		return
	}

	channel <- true
}
