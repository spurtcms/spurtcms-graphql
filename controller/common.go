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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	spurtcore "github.com/spurtcms/pkgcore"
)

type key string

const ContextKey key = "ginContext"

type MailConfig struct {
	Emails         []string
	MailUsername   string
	MailPassword   string
	Subject        string
	AdditionalData map[string]interface{}
}

type MailImages struct {
	Owndesk   string
	Twitter   string
	Facebook  string
	LinkedIn  string
	Youtube   string
	Instagram string
}

type SocialMedias struct {
	Linkedin  string
	Twitter   string
	Facebook  string
	Instagram string
	Youtube   string
}

var (
	Mem                            member.MemberAuth
	Auth                           *auth.Authorization
	TimeZone                       *time.Location
	ProfileImagePath, SpecialToken string
	MemberRegisterPermission       string
	SectionTypeId                  = 12
	MemberFieldTypeId              = 14
	PathUrl                        string
	EmailImageUrlPrefix            string
	SmtpPort, SmtpHost             string
	// OwndeskChannelId               = 108
	EmailImagePath                 MailImages
	SocialMediaLinks               SocialMedias
	OwndeskLoginEnquiryTemplate    = "OwndeskLoginEnquiry"
	OwndeskLoginTemplate           = "OwndeskLogin"
	OwndeskClaimnowTemplate        = "OwndeskClaimRequest"
	LocalLoginType                 = "member"
)

var (
	ErrInvalidMail         = errors.New("your email is not yet registered in our owndesk platform")
	ErrSendMail            = errors.New("failed to send unauthorized login attempt mail to admin")
	ErrclaimAlready        = errors.New("member profile is already claimed")
	ErrEmptyProfileSlug    = errors.New("profile slug should not be empty")
	ErrProfileSlugExist    = errors.New("profile slug already exists")
	ErrMandatory           = errors.New("missing mandatory fields")
	ErrMemberRegisterPerm  = errors.New("member register permission denied")
	ErrMemberInactive      = errors.New("inactive member")
	ErrMemberLoginPerm     = errors.New("member login permission denied")
)

func init() {

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

	EmailImagePath = MailImages{
		Owndesk:   EmailImageUrlPrefix + strings.TrimPrefix("/view/img/own-desk-logo.png", "/"),
		Twitter:   EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons3.png", "/"),
		Facebook:  EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons1.png", "/"),
		LinkedIn:  EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons2.png", "/"),
		Youtube:   EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons4.png", "/"),
		Instagram: EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons5.png", "/"),
	}

	SocialMediaLinks = SocialMedias{
		Linkedin:  os.Getenv("LINKEDIN"),
		Twitter:   os.Getenv("TWITTER"),
		Facebook:  os.Getenv("FACEBOOK"),
		Instagram: os.Getenv("INSTAGRAM"),
		Youtube:   os.Getenv("YOUTUBE"),
	}

	MemberRegisterPermission = os.Getenv("ALLOW_MEMBER_REGISTER")

}

func GetAuthorization(token string, db *gorm.DB) *auth.Authorization {

	auth := spurtcore.NewInstance(&auth.Option{DB: db, Token: token, Secret: os.Getenv("JWT_SECRET")})

	return &auth

}

func GetAuthorizationWithoutToken(db *gorm.DB) *auth.Authorization {

	auth := spurtcore.NewInstance(&auth.Option{DB: db, Token: "", Secret: os.Getenv("JWT_SECRET")})

	return &auth
}

func StoreImageBase64ToLocal(imageData, storagePath, storingName string) (string, string, error) {

	extEndIndex := strings.Index(imageData, ";base64,")

	base64data := imageData[strings.IndexByte(imageData, ',')+1:]

	var ext = imageData[11:extEndIndex]

	randomNum := strconv.Itoa(rand.Intn(900000) + 100000)

	imageName := storingName + "-" + randomNum + "." + ext

	err := os.MkdirAll(storagePath, 0755)

	if err != nil {

		log.Println(err)

		return "", "", err
	}

	storageDestination := storagePath + imageName

	decode, err := base64.StdEncoding.DecodeString(base64data)

	if err != nil {

		log.Println(err)

		return "", "", err
	}

	file, err := os.Create(storageDestination)

	if err != nil {

		log.Println(err)

		return "", "", err

	}
	if _, err := file.Write(decode); err != nil {

		log.Println(err)

		return "", "", err

	}

	return imageName, storageDestination, nil
}

func SendMail(config MailConfig, html_content string, channel chan error) {

	// Sender data
	from := config.MailUsername
	password := config.MailPassword

	// Receiver email address
	to := config.Emails

	// Authentication
	auth := smtp.PlainAuth("", from, password, SmtpHost)

	subject := "Subject:" + config.Subject + " \n"

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

func HashingPassword(pass string) string {

	passbyte, err := bcrypt.GenerateFromPassword([]byte(pass), 14)

	if err != nil {

		panic(err)

	}

	return string(passbyte)
}

func GetNotifyAdminEmails(db *gorm.DB, adminIds []int) ([]auth.TblUser,[]string,error){

	Mem.Auth = GetAuthorizationWithoutToken(db)

	adminDetails,err := Mem.GetAdminDetails(adminIds)

	if err != nil{

		return []auth.TblUser{},[]string{},err
	}

	var adminEmails []string

	for _,admin := range adminDetails{

		adminEmails = append(adminEmails, admin.Email)
	}

	return adminDetails,adminEmails,nil
}

