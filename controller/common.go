package controller

import (
	"encoding/base64"
	"log"
	"math/rand"
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



