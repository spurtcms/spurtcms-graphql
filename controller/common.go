package controller

import (
	"encoding/base64"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
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
	TimeZone *time.Location
	ProfileImagePath,SpecialToken string
	SectionTypeId = 12
    MemberFieldTypeId = 14
)

type MailConfig struct{
	Email            string
	MailUsername     string
	MailPassword     string
	Subject          string
	AdditionalData   map[string]interface{}
}

func init(){
	
	err := godotenv.Load()

	if err != nil {

		log.Fatalf("Error loading .env file")
	}

	SpecialToken = "%$HEID$#PDGH*&MGEAFCC"

	TimeZone, _ = time.LoadLocation(os.Getenv("TIME_ZONE"))

	ProfileImagePath = "Uploads/ProfileImages/"

}

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

func SendMail(config MailConfig,html_content string,channel chan bool) {

	// Sender data.
	from := config.MailUsername
	password := config.MailPassword

	// Receiver email address.
	to := []string{
		config.Email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	subject := "Subject:"+config.Subject+" \n"

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte(subject + mime + html_content)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)

	if err != nil {

		log.Println(err)

		channel <- false

		return
	}

	channel <- true
}
