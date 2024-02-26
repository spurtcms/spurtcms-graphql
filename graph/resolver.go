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

func Channellist(db *gorm.DB,ctx context.Context,limit, offset int) (model.ChannelDetails, error) {
	return controller.Channellist(db,ctx,limit,offset)
}

func ChannelEntriesList(db *gorm.DB,ctx context.Context, channelID, channelEntryID, categoryId *int, limit, offset *int) (model.ChannelEntryDetails, error) {
	return controller.ChannelEntriesList(db,ctx,channelID,channelEntryID,categoryId,limit,offset)
}

func MemberUpdate(db *gorm.DB,ctx context.Context, memberdata model.MemberDetails) (bool, error) {
	return controller.UpdateMember(db,ctx,memberdata)
}

func ChannelDetail(db *gorm.DB,ctx context.Context, channelID int) (model.TblChannel, error) {
	return controller.ChannelDetail(db,ctx,channelID)
}


