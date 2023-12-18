package graph

import (
	"context"
	"gqlserver/controller"
	"gqlserver/dbconfig"
	"gqlserver/graph/model"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	DB  *gorm.DB
}

func NewResolver() *Resolver {
	return &Resolver{DB: dbconfig.SetupDB()}
}

func MemberLogin(db *gorm.DB,input model.LoginCredentials)(string,error){
	return controller.MemberLogin(db,input)
}

func MemberRegister(db *gorm.DB, input model.MemberDetails)(bool,error){
	return controller.MemberRegister(db,input)
}

func Pagelist(db *gorm.DB, spaceid int) (model.PageAndPagegroups, error) {
	return controller.Pagelist(db,spaceid)
}

func Spacelist(db *gorm.DB, filter model.Filter) (model.SpacesDetails, error) {
	return controller.Spacelist(db,filter) 
}

func UpdateMember(db *gorm.DB,ctx context.Context ,memberdata model.MemberDetails)(bool,error){
	return controller.UpdateMember(db,ctx,memberdata)
}

func PageContent(db *gorm.DB,ctx context.Context, pageid int) (string, error) {
	return controller.PageContent(db,ctx,pageid)
}

func UpdateHighlights(db *gorm.DB,ctx context.Context, highlights model.Highlights)(bool,error){
	return controller.UpdateHighlights(db,ctx,highlights)
}

func UpdateNotes(db *gorm.DB,ctx context.Context, pageid int, notes string) (bool, error) {
	return controller.UpdateNotes(db,ctx,pageid,notes)
}

func GetNotesOrHighlights(db *gorm.DB, ctx context.Context, pageid int,contentType string) ([]model.MemberNotesHighlight, error) {
	return controller.GetNotesOrHighlights(db,ctx,pageid,contentType)
}

func DeleteNotesOrHighlights(db *gorm.DB,ctx context.Context, contentID int) (bool, error) {
	return controller.DeleteNotesOrHighlights(db,ctx,contentID)
}

func SendOtpToMail(db *gorm.DB, ctx context.Context, email string) (bool, error) {
	return controller.SendOtpToMail(db,ctx,email)
}

func ResetPassword(db *gorm.DB, otp int, newPassword, email string) (bool, error) {
	return controller.ResetPassword(db,otp,newPassword,email)
}


