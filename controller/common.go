package controller

import (
	"encoding/base64"
	"errors"
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

type MailConfig struct{
	Email            string
	MailUsername     string
	MailPassword     string
	Subject          string
	AdditionalData   map[string]interface{}
}

var(
	Mem member.MemberAuth
	Auth *auth.Authorization
	TimeZone *time.Location
	ProfileImagePath,SpecialToken string
	SectionTypeId = 12
    MemberFieldTypeId = 14
	PathUrl string
	EmailImageUrlPrefix string
	SmtpPort,SmtpHost string
	OwndeskChannelId int = 108
	AdditionalData map[string]interface{}
)

var(
	ErrInvalidMail = errors.New("your email is not yet registered in our owndesk platform")
	ErrSendMail = errors.New("failed to send unauthorized login attempt mail to admin")
	ErrclaimAlready = errors.New("member profile is already claimed")
	ErrEmptyProfileName = errors.New("profile name should not be empty")
	ErrEmptyProfileSlug = errors.New("profile slug should not be empty")
)

func init(){
	
	err := godotenv.Load()

	if err != nil {

		log.Fatalf("Error loading .env file")
	}

	SpecialToken = "%$HEID$#PDGH*&MGEAFCC"

	TimeZone, _ = time.LoadLocation(os.Getenv("TIME_ZONE"))

	ProfileImagePath = "Uploads/ProfileImages/"

	if os.Getenv("DOMAIN_URL") != "" {

		PathUrl = os.Getenv("DOMAIN_URL")

	} else {

		PathUrl = os.Getenv("LOCAL_URL")
	}

	SmtpHost = os.Getenv("SMTP_HOST")

	SmtpPort = os.Getenv("SMTP_PORT")

	EmailImageUrlPrefix = os.Getenv("EMAIL_IMAGE_PREFIX_URL")

	EmailImagePath := struct{
		Owndesk    string
		Twitter    string
		Facebook   string
		LinkedIn   string
		Youtube    string
		Instagram  string
	}{
		Owndesk  :  EmailImageUrlPrefix + strings.TrimPrefix("/view/img/own-desk-logo.png","/"),
		Twitter  :  EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons3.png","/"),
		Facebook :  EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons1.png","/"),
		LinkedIn :  EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons2.png","/"),
		Youtube  :  EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons4.png","/"),
		Instagram:  EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons5.png","/"),
	}

	SocialMediaLinks := struct{
		Linkedin    string
		Twitter     string
		Facebook    string
		Instagram   string
		Youtube     string
	}{
		Linkedin: os.Getenv("LINKEDIN"),
		Twitter: os.Getenv("TWITTER"),
		Facebook: os.Getenv("FACEBOOK"),
		Instagram: os.Getenv("INSTAGRAM"),
		Youtube: os.Getenv("YOUTUBE"),
	}

	AdditionalData = map[string]interface{}{"emailImagePath": EmailImagePath,"socialMediaLinks": SocialMediaLinks}
}

func GetAuthorization(token string,db *gorm.DB)(*auth.Authorization) {

	auth := spurtcore.NewInstance(&auth.Option{DB: db, Token: token, Secret: os.Getenv("JWT_SECRET")})

	return &auth

}

func GetAuthorizationWithoutToken(db *gorm.DB)(*auth.Authorization){

	auth := spurtcore.NewInstance(&auth.Option{DB: db, Token: "", Secret: os.Getenv("JWT_SECRET")})

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

func SendMail(config MailConfig,html_content string,channel chan error) {

	// Sender data
	from := config.MailUsername
	password := config.MailPassword

	// Receiver email address
	to := []string{
		config.Email,
	}

	// Authentication
	auth := smtp.PlainAuth("", from, password, SmtpHost)

	subject := "Subject:"+config.Subject+" \n"

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte(subject + mime + html_content)

	// Sending email
	err := smtp.SendMail(SmtpHost+":"+SmtpPort, auth, from, to, msg)

	if err != nil {

		log.Println(err)

		channel <- err

		return
	}

	channel <- nil
}
