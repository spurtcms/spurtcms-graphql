package graph

import (
	"context"
	"gqlserver/controller"
	"gqlserver/dbconfig"
	"gqlserver/graph/model"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	DB  *gorm.DB
}

func NewResolver() *Resolver {
	return &Resolver{DB: dbconfig.SetupDB()}
}

func MemberLogin(db *gorm.DB,ctx context.Context, email string)(bool,error){
	return controller.MemberLogin(db,ctx,email)
}

func MemberRegister(db *gorm.DB, input model.MemberDetails)(bool,error){
	return controller.MemberRegister(db,input)
}

func Channellist(db *gorm.DB,ctx context.Context,limit, offset int) (model.ChannelDetails, error) {
	return controller.Channellist(db,ctx,limit,offset)
}

func ChannelEntriesList(db *gorm.DB,ctx context.Context, channelID, categoryId *int, limit, offset int, title *string,categoryChildId *int) (model.ChannelEntriesDetails, error) {
	return controller.ChannelEntriesList(db,ctx,channelID,categoryId,limit,offset,title,categoryChildId)
}

func MemberUpdate(db *gorm.DB,ctx context.Context, memberdata model.MemberDetails) (bool, error) {
	return controller.UpdateMember(db,ctx,memberdata)
}

func ChannelDetail(db *gorm.DB,ctx context.Context, channelID int) (model.Channel, error) {
	return controller.ChannelDetail(db,ctx,channelID)
}

func SpaceList(db *gorm.DB,ctx context.Context,limit, offset int) (model.SpaceDetails, error){
	return controller.SpaceList(db,ctx,limit,offset)
}

func SpaceDetails(db *gorm.DB,ctx context.Context,spaceId int)(model.Space,error){
	return controller.SpaceDetails(db,ctx,spaceId)
}

func PagesAndPageGroupsUnderSpace(db *gorm.DB,ctx context.Context, spaceID int) (model.PageAndPageGroups, error) {
	return controller.PagesAndPageGroupsBySpaceId(db,ctx,spaceID)
}

func CategoriesList(db *gorm.DB,ctx context.Context, limit, offset, categoryGroupId, hierarchyLevel,checkEntriesPresence *int) (model.CategoriesList, error) {
	return controller.CategoriesList(db,ctx,limit,offset,categoryGroupId,hierarchyLevel,checkEntriesPresence)
}

func ChannelEntryDetail(db *gorm.DB,ctx context.Context, channelID *int,categoryID *int, channelEntryID *int,slug *string) (model.ChannelEntries, error) {
	return controller.ChannelEntryDetail(db,ctx,channelEntryID,channelID,categoryID,slug)
}

func MemberProfileUpdate(db *gorm.DB,ctx context.Context, profiledata model.ProfileData,entryId int,updateExactMemberProfileOnly bool) (bool, error) {
	return controller.MemberProfileUpdate(db,ctx,profiledata,entryId,updateExactMemberProfileOnly)
}

func VerifyMemberOtp(db *gorm.DB,ctx context.Context,email string ,otp int) (string, error) {
	return controller.VerifyMemberOtp(db,ctx, email,otp)
}

func Memberclaimnow(db *gorm.DB,ctx context.Context, input model.ClaimData, entryId int) (bool, error) {
	return controller.Memberclaimnow(db,ctx,input,entryId)
}

func EcommerceProductList(db *gorm.DB,ctx context.Context, limit int, offset int, filter *model.ProductFilter, sort *model.ProductSort) (model.EcommerceProducts, error) {
	return controller.EcommerceProductList(db,ctx,limit,offset,filter,sort)
}

