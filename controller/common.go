package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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
	"spurtcms-graphql/models"
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
	"gorm.io/gorm"

	spurtcore "github.com/spurtcms/pkgcore"

	memberpkg "github.com/spurtcms/member"

	teampkg "github.com/spurtcms/team"
)

type key string

const ContextKey key = "ginContext"

var (
	DB                             *gorm.DB
	Model                          models.ModelConfig
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
)

var (
	ErrInvalidMail                = errors.New("your email is not yet registered in our owndesk platform")
	ErrSendMail                   = errors.New("failed to send unauthorized login attempt mail to admin")
	ErrclaimAlready               = errors.New("member profile is already claimed")
	ErrEmptyProfileSlug           = errors.New("profile slug should not be empty")
	ErrProfileSlugExist           = errors.New("profile slug already exists")
	ErrMandatory                  = errors.New("missing mandatory fields")
	ErrMemberRegisterPerm         = errors.New("member register permission denied")
	ErrMemberInactive             = errors.New("inactive member")
	ErrMemberLoginPerm            = errors.New("member login permission denied")
	ErrRecordNotFound             = errors.New("record not found")
	ErrPassHash                   = errors.New("hasing password failed")
	ErrUpload                     = errors.New("file upload failed")
	ErrOldPass                    = errors.New("old password mismatched")
	ErrConfirmPass                = errors.New("new passowrd and confirmation password mismatched")
	ErrSamePass                   = errors.New("old password and new password should not be same")
	ErrLoginReq                   = errors.New("login required")
	ErrUnauthorize                = errors.New("unauthorized access")
	ErrGinInstance                = errors.New("gin instance retrieval context error")
	ErrMemberSettings             = errors.New("failed to fetch member settings")
	ErrFetchMailConfig            = errors.New("failed to fetch email configurations")
	ErrInactiveTemplate           = errors.New("mail template is inactive")
	ErrParsingHtmlTemplate        = errors.New("failed to parse the html template")
	ErrExecutingHtmlTemplate      = errors.New("failed to execute html template")
	ErrNoMemberDetails            = errors.New("failed to get the member details")
	ErrNoOtpUpdate                = errors.New("unable to update otp")
	ErrFetchAdmin                 = errors.New("failed to fetch the admin details")
	ErrCreatingToken              = errors.New("unable to create token")
	ErrLoginDataCheck             = errors.New("failed to get member login data")
	ErrMailExist                  = errors.New("email already exists")
	ErrMobileExist                = errors.New("mobile number already exists")
	ErrFetchStorageType           = errors.New("failed to fetch the storage type details")
	ErrIllegalB64                 = errors.New("illegal base64 data")
	ErrJsonUnMarshal              = errors.New("failed to unmarshall the json")
	ErrClaimMail                  = errors.New("failed to send claim request mail to the admin")
	ErrClaimSubmitMail            = errors.New("failed to send claim request submission status mail to the user")
	ErrLoginClaimMail             = errors.New("current login email sholuld not be used in another claim")
	ErrLoginClaimMob              = errors.New("current login mobile number sholuld not be used in another claim")
	ErrFetchEntries               = errors.New("failed to fetch the channel Entries data")
	ErrFetchJobsList              = errors.New("failed to fetch jobs list")
	ErrFetchJobDetails            = errors.New("failed to fetch details of the particular job")
	ErrUnauthorizedAccess         = errors.New("unauthorized access")
	ErrGettingMemberId            = errors.New("couldn't able to get member id from request")
	ErrApplicantNotFound          = errors.New("applicant not found")
	ErrGettingContext             = errors.New("error getting context")
	ErrGettingApplicantDetail     = errors.New("unable to fetch the applicant Details")
	ErrCheckingAlreadyRegistered  = errors.New("unable to check whether the applicant is already applied to this job or not")
	ErrApplicantAlreadyRegistered = errors.New("this applicant has already applied for this job using this emailid ")
	EmailImagePath                 models.MailImages
	SocialMediaLinks               models.SocialMedias
	OwndeskLoginEnquiryTemplate    = "owndeskloginenquiry"
	OwndeskLoginTemplate           = "owndesklogin"
	OwndeskClaimnowTemplate        = "owndeskclaimrequest"
	OwndeskClaimSubmitTemplate     = "owndeskclaimsubmit"
	LocalLoginType                 = "member"
	TokenExpiryTime                = 1
	ErrorLog                       *log.Logger
	WarnLog                        *log.Logger
)

func init() {

	err := godotenv.Load()

	if err != nil {

		log.Fatalf("Error loading .env file")
	}

	DB = dbconfig.SetupDB()

	Model = models.ModelConfig{DB: DB}

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

	SmtpHost = os.Getenv("SMTP_HOST")

	SmtpPort = os.Getenv("SMTP_PORT")

	EmailImageUrlPrefix = os.Getenv("EMAIL_IMAGE_PREFIX_URL")

	EmailImagePath = models.MailImages{
		Owndesk:   EmailImageUrlPrefix + strings.TrimPrefix("/view/img/own-desk-logo.png", "/"),
		Twitter:   EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons3.png", "/"),
		Facebook:  EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons1.png", "/"),
		LinkedIn:  EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons2.png", "/"),
		Youtube:   EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons4.png", "/"),
		Instagram: EmailImageUrlPrefix + strings.TrimPrefix("/view/img/social-media-icons5.png", "/"),
	}

	SocialMediaLinks = models.SocialMedias{
		Linkedin:  os.Getenv("LINKEDIN"),
		Twitter:   os.Getenv("TWITTER"),
		Facebook:  os.Getenv("FACEBOOK"),
		Instagram: os.Getenv("INSTAGRAM"),
		Youtube:   os.Getenv("YOUTUBE"),
	}

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

func SendMail(config models.MailConfig, html_content string, channel chan error) {

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

func HashingPassword(pass string) (string, error) {

	passbyte, err := bcrypt.GenerateFromPassword([]byte(pass), 14)

	if err != nil {

		return "", err
	}

	return string(passbyte), nil
}

func GetNotifyAdminEmails(db *gorm.DB, adminIds []int) ([]teampkg.TblUser, []string, error) {

	_, adminDetails, err := TeamInstance.GetUserById(0, adminIds)

	if err != nil {

		return []teampkg.TblUser{}, []string{}, err
	}

	var adminEmails []string

	for _, admin := range adminDetails {

		adminEmails = append(adminEmails, admin.Email)
	}

	return adminDetails, adminEmails, nil
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

func ConvertByteToJson(byteData []byte) (map[string]interface{}, error) {

	var jsonMap map[string]interface{}

	err := json.Unmarshal(byteData, &jsonMap)

	if err != nil {

		return map[string]interface{}{}, err
	}

	return jsonMap, nil

}

func GetEmailConfigurations() (models.MailConfig, error) {

	var email_configs models.EmailConfiguration

	err := Model.GetEmailConfig(&email_configs)

	if err != nil{

		return models.MailConfig{},err
	}

	var sendMailData models.MailConfig

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

func ImageResize(c *gin.Context) {

	fileName := c.Query("name")

	filePath := c.Query("path")

	extension := path.Ext(fileName)

	var storageType models.StorageType

	err := Model.GetStorageType(&storageType)

	if err != nil {

		fmt.Println(err)

		c.AbortWithError(500, fmt.Errorf("%v-%v", ErrGetAwsCreds, err))

		return
	}

	var byteData []byte

	rawObject, err := storage.GetObjectFromS3(storageType.Aws, filePath+fileName)

	if err != nil {

		fmt.Println(err)

		c.AbortWithError(500, fmt.Errorf("%v-%v", ErrGetImage, err))

		return
	}

	buf := new(bytes.Buffer)

	buf.ReadFrom(rawObject.Body)

	byteData = buf.Bytes()

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
