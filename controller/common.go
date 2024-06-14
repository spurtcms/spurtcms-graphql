package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/smtp"
	"os"
	"spurtcms-graphql/dbconfig"
	"spurtcms-graphql/logger"
	"spurtcms-graphql/storage"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spurtcms/pkgcore/auth"
	"github.com/spurtcms/pkgcore/member"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	spurtcore "github.com/spurtcms/pkgcore"

	memberpkg "github.com/spurtcms/member"
)

type key string

const ContextKey key = "ginContext"

type MailConfig struct {
	Emails         []string
	MailUsername   string
	MailPassword   string
	SmtpPort       string
	SmtpHost       string
	Subject        string
	TimeOut        time.Duration
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

type StorageType struct {
	Id           int
	Local        string
	Aws          datatypes.JSONMap `gorm:"type:jsonb"`
	Azure        datatypes.JSONMap `gorm:"type:jsonb"`
	Drive        datatypes.JSONMap `gorm:"type:jsonb"`
	SelectedType string
}

type EmailConfiguration struct {
	Id           int
	StmpConfig   datatypes.JSONMap `gorm:"type:jsonb"`
	SelectedType string
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
	EmailImagePath                 MailImages
	SocialMediaLinks               SocialMedias
	OwndeskLoginEnquiryTemplate    = "owndeskloginenquiry"
	OwndeskLoginTemplate           = "owndesklogin"
	OwndeskClaimnowTemplate        = "owndeskclaimrequest"
	OwndeskClaimSubmitTemplate     = "owndeskclaimsubmit"
	LocalLoginType                 = "member"
	ErrorLog                       *log.Logger
	WarnLog                        *log.Logger
	DB                             *gorm.DB
)

var (
	ErrInvalidMail        = errors.New("your email is not yet registered in our owndesk platform")
	ErrclaimAlready       = errors.New("member profile is already claimed")
	ErrEmptyProfileSlug   = errors.New("profile slug should not be empty")
	ErrProfileSlugExist   = errors.New("profile slug already exists")
	ErrMandatory          = errors.New("missing mandatory fields")
	ErrMemberRegisterPerm = errors.New("member register permission denied")
	ErrMemberInactive     = errors.New("inactive member")
	ErrMemberLoginPerm    = errors.New("member login permission denied")
	ErrRecordNotFound     = errors.New("record not found")
	ErrPassHash           = errors.New("hasing password failed")
	ErrUpload             = errors.New("file upload failed")
	ErrOldPass            = errors.New("old password mismatched")
	ErrConfirmPass        = errors.New("new passowrd and confirmation password mismatched")
	ErrSamePass           = errors.New("old password and new password should not be same")
	ErrLoginReq           = errors.New("login required")
	ErrSendMailFail       = errors.New("failed to send mail with smtp configurations")
	ErrClaimMail          = errors.New("failed to send claim request mail to the admin")
	ErrClaimSubmitMail    = errors.New("failed to send claim request submission status mail to the user")
	ErrInactiveTemplate   = errors.New("mail template is inactive")

)

func init() {

	err := godotenv.Load()

	if err != nil {

		log.Fatalf("Error loading .env file")
	}

	DB = dbconfig.SetupDB()

	SpecialToken = "%$HEID$#PDGH*&MGEAFCC"

	TimeZone, _ = time.LoadLocation(os.Getenv("TIME_ZONE"))

	ErrorLog = logger.ErrorLOG()

	WarnLog = logger.WarnLOG()

	ProfileImagePath = "Uploads/ProfileImages/"

	if os.Getenv("DOMAIN_URL") != "" {

		PathUrl = os.Getenv("DOMAIN_URL")

	} else {

		PathUrl = os.Getenv("LOCAL_URL")
	}

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

func GetMemberPackageSetup(db *gorm.DB) *memberpkg.Member {

	memberConfig := memberpkg.Config{DB: db}

	memberSetup := memberpkg.MemberSetup(memberConfig)

	return memberSetup

}

func GetAuthorization(token string, db *gorm.DB) *auth.Authorization {

	auth := spurtcore.NewInstance(&auth.Option{DB: db, Token: token, Secret: os.Getenv("JWT_SECRET")})

	return &auth

}

func GetAuthorizationWithoutToken(db *gorm.DB) *auth.Authorization {

	auth := spurtcore.NewInstance(&auth.Option{DB: db, Token: "", Secret: os.Getenv("JWT_SECRET")})

	return &auth
}

func SendMail(config MailConfig, html_content string, channel chan error) {

	// Sender data
	from := config.MailUsername
	password := config.MailPassword

	// Receiver email address
	to := config.Emails

	// Authentication
	auth := smtp.PlainAuth("", from, password, config.SmtpHost)

	subject := "Subject:" + config.Subject + " \n"

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte(subject + mime + html_content)

	// Sending email
	err := smtp.SendMail(config.SmtpHost+":"+config.SmtpPort, auth, from, to, msg)

	if err != nil {

		channel <- err

		return
	}

	channel <- nil
}

func HashingPassword(pass string) (string, error) {

	passbyte, err := bcrypt.GenerateFromPassword([]byte(pass), 14)

	if err != nil {

		return "", err
	}

	return string(passbyte), nil
}

func GetNotifyAdminEmails(db *gorm.DB, adminIds []int) ([]auth.TblUser, []string, error) {

	Mem.Auth = GetAuthorizationWithoutToken(db)

	adminDetails, err := Mem.GetAdminDetails(adminIds)

	if err != nil {

		return []auth.TblUser{}, []string{}, err
	}

	var adminEmails []string

	for _, admin := range adminDetails {

		adminEmails = append(adminEmails, admin.Email)
	}

	return adminDetails, adminEmails, nil
}

func GetStorageType(db *gorm.DB) (StorageType, error) {

	var storageType StorageType

	if err := db.Debug().Table("tbl_storage_types").First(&storageType).Error; err != nil {

		return StorageType{}, err
	}

	return storageType, nil
}

func IoReadSeekerToBase64(file io.ReadSeeker) (string, error) {

	_, err := file.Seek(0, io.SeekStart)

	if err != nil {

		return "", err
	}

	// Read the data into a buffer
	var buf bytes.Buffer

	_, err = io.Copy(&buf, file)

	if err != nil {

		return "", err
	}

	// Encode the buffer to a base64 string
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64Str, nil
}

func CompareBcryptPassword(hashpass, oldpass string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(oldpass))

	if err != nil {

		return err
	}

	return nil
}

func GetFilePathsRelatedToStorageTypes(db *gorm.DB, path string) string {

	storageType, _ := GetStorageType(db)

	awsCreds := storageType.Aws

	isExist, _ := storage.CheckS3FileExistence(awsCreds, path)

	if isExist {

		s3FileServeEndpoint := "image-resize"

		s3Path := PathUrl + s3FileServeEndpoint + "?name=" + strings.TrimPrefix(path, "/")

		return s3Path

	}

	localPath := PathUrl + strings.TrimPrefix(path, "/")

	return localPath
}

func ConvertByteToJson(byteData []byte) (map[string]interface{}, error) {

	var jsonMap map[string]interface{}

	err := json.Unmarshal(byteData, &jsonMap)

	if err != nil {

		return map[string]interface{}{}, err
	}

	return jsonMap, nil

}

func GetEmailConfigurations(db *gorm.DB) (MailConfig, error) {

	var email_configs EmailConfiguration

	if err := db.Debug().Table("tbl_email_configurations").First(&email_configs).Error; err != nil {

		return MailConfig{}, err
	}

	var sendMailData MailConfig

	if email_configs.SelectedType == "environment" {

		sendMailData.MailUsername = os.Getenv("MAIL_USERNAME")

		sendMailData.MailPassword = os.Getenv("MAIL_PASSWORD")

		sendMailData.SmtpHost = os.Getenv("SMTP_HOST")

		sendMailData.SmtpPort = os.Getenv("SMTP_PORT")

	} else if email_configs.SelectedType == "smtp" {

		sendMailData.MailUsername = email_configs.StmpConfig["Mail"].(string)

		sendMailData.MailPassword = email_configs.StmpConfig["Password"].(string)

		sendMailData.SmtpHost = email_configs.StmpConfig["Host"].(string)

		sendMailData.SmtpPort = email_configs.StmpConfig["Port"].(string)

	}

	sendMailData.TimeOut = 5 * time.Second

	return sendMailData, nil

}
