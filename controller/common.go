package controller

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"sync"
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

func GenerateEmail(email, subject, message string, wg *sync.WaitGroup) error {

	data := map[string]interface{}{

		"Body": template.HTML(message),
	}

	t, err2 := template.ParseFiles("view/email/email-template.html")

	if err2 != nil {

		fmt.Println(err2)
	}

	var htmlBuffer bytes.Buffer

	if err1 := t.Execute(&htmlBuffer, data); err1 != nil {

		log.Println(err1)
	}

	result := htmlBuffer.String()

	defer wg.Done()

	from := os.Getenv("MAIL_USERNAME")

	smtpHost := "smtp.gmail.com"

	smtpPort := "587"

	contentType := "text/html"

	// Set up the SMTP server configuration.
	auth := smtp.PlainAuth("", os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"), smtpHost)

	// Compose the email.
	emailBody := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: %s; charset=UTF-8\r\n\r\n%s", from, email, subject, contentType, result)

	// Connect to the SMTP server and send the email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(emailBody))

	if err != nil {

		fmt.Println("Failed to send email:", err)

		return err

	} else {

		fmt.Println("Email sent successfully to:", email)
		
		return nil
	}

}