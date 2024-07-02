package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/smtp"
	"os"
	"path"
	"spurtcms-graphql/dbconfig"
	"spurtcms-graphql/logger"
	"spurtcms-graphql/storage"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
	"github.com/spurtcms/pkgcore/auth"
	"github.com/spurtcms/pkgcore/member"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	spurtcore "github.com/spurtcms/pkgcore"
)

type key string

const ContextKey key = "ginContext"

type MailConfig struct {
	Emails       []string
	MailUsername string
	MailPassword string
	SmtpPort     string
	SmtpHost     string
	Subject      string
	TimeOut      time.Duration
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
	SmtpConfig   datatypes.JSONMap `gorm:"type:jsonb"`
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
	ErrGetImage           = errors.New("failed to retrieve the file")
	ErrGetAwsCreds        = errors.New("failed to retrieve the aws credentials")
	ErrDecodeImg          = errors.New("failed to parse the image data")
	ErrImageResize        = errors.New("failed to resize the image")
	ErrMailExist          = errors.New("email already exists")
	ErrLoginClaimMail     = errors.New("current login email sholuld not be used in another claim")
	ErrLoginClaimMob      = errors.New("current login mobile number sholuld not be used in another claim")
	ErrMobileExist        = errors.New("mobile number already exists")
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

		sendMailData.MailUsername = email_configs.SmtpConfig["Mail"].(string)

		sendMailData.MailPassword = email_configs.SmtpConfig["Password"].(string)

		sendMailData.SmtpHost = email_configs.SmtpConfig["Host"].(string)

		sendMailData.SmtpPort = email_configs.SmtpConfig["Port"].(string)

	}

	sendMailData.TimeOut = 5 * time.Second

	return sendMailData, nil

}

func GenerateFileName(imageData string) (fileName string) {

	extEndIndex := strings.Index(imageData, ";base64,")

	var ext = imageData[11:extEndIndex]

	rand_num := strconv.Itoa(int(time.Now().Unix()))

	imageName := "IMG-" + rand_num + "." + ext

	return imageName
}

func ImageResize(c *gin.Context) {

	fileName := c.Query("filename")

	filePath := c.Query("path")

	extension := path.Ext(fileName)

	storageType, err := GetStorageType(DB)

	if err != nil {

		fmt.Println(err)

		c.AbortWithError(500, fmt.Errorf("%v-%v", ErrGetAwsCreds, err))

		return
	}

	var byteData []byte

	if storageType.SelectedType == "aws" {

		rawObject, err := storage.GetObjectFromS3(storageType.Aws, filePath+fileName)

		if err != nil {

			fmt.Println(err)

			c.AbortWithError(500, fmt.Errorf("%v-%v", ErrGetImage, err))

			return
		}

		buf := new(bytes.Buffer)

		buf.ReadFrom(rawObject.Body)

		byteData = buf.Bytes()

	} else if storageType.SelectedType == "azure" {

		fmt.Printf("azure storage selected")

		return

	} else if storageType.SelectedType == "drive" {

		fmt.Println("drive storage selected")

		return
	}

	extType := strings.Trim(extension, ".")

	if c.Query("width") == "" || c.Query("height") == "" {

		if extType == "svg" {

			extType = "svg+xml"
		}

		c.Data(200, "image/"+extType, byteData)

		return
	}

	width, _ := strconv.ParseUint(c.Query("width"), 10, 64)

	height, _ := strconv.ParseUint(c.Query("height"), 10, 64)

	Image, _, err := image.Decode(bytes.NewReader(byteData))

	if err != nil {

		fmt.Println(err)

		c.AbortWithError(500, fmt.Errorf("%v-%v", ErrDecodeImg, err))

		return
	}

	newImage := resize.Resize(uint(width), uint(height), Image, resize.Lanczos3)

	if extension == ".png" {

		err = png.Encode(c.Writer, newImage)

		if err != nil {

			fmt.Println(err)

			c.AbortWithError(500, fmt.Errorf("%v-%v", ErrImageResize, err))

			return
		}
	}

	if extension == ".jpeg" || extension == ".jpg" {

		err = jpeg.Encode(c.Writer, newImage, nil)

		if err != nil {

			fmt.Println(err)

			c.AbortWithError(500, fmt.Errorf("%v-%v", ErrImageResize, err))

			return
		}

	}

}

func IsValidBase64(input string) (isvalid bool, base64Data string, extension string) {

	if !strings.Contains(input, "data:image/png;base64") && !strings.Contains(input, "data:image/jpeg;base64") && !strings.Contains(input, "data:image/jpg;base64") && !strings.Contains(input, "data:image/svg;base64") {

		return false, "", ""
	}

	base64Data = input[strings.IndexByte(input, ',')+1:]

	_, err := base64.StdEncoding.DecodeString(base64Data)

	if err != nil {
		return false, "", ""
	}

	extEndIndex := strings.Index(input, ";base64,")

	var ext = input[11:extEndIndex]

	return true, base64Data, ext
}
