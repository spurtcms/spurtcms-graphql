package graph

import (
	"context"
	"spurtcms-graphql/controller"
	"spurtcms-graphql/graph/model"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *gorm.DB
}

func NewResolver() *Resolver {
	return &Resolver{DB: controller.DB}
}

func MemberLogin(db *gorm.DB, ctx context.Context, email string) (bool, error) {
	return controller.MemberLogin(db, ctx, email)
}

func MemberRegister(db *gorm.DB, ctx context.Context, input model.MemberDetails, ecomModule *int) (bool, error) {
	return controller.MemberRegister(db, ctx, input, ecomModule)
}

func Channellist(db *gorm.DB, ctx context.Context, limit, offset int) (*model.ChannelDetails, error) {
	return controller.Channellist(db, ctx, limit, offset)
}

func ChannelEntriesList(db *gorm.DB, ctx context.Context, channelID, categoryId *int, limit, offset int, title *string, categoryChildId *int, categorySlug, categoryChildSlug *string, reduireData *model.RequireData) (*model.ChannelEntriesDetails, error) {
	return controller.ChannelEntriesList(db, ctx, channelID, categoryId, limit, offset, title, categoryChildId, categorySlug, categoryChildSlug, reduireData)
}

func MemberUpdate(db *gorm.DB, ctx context.Context, memberdata model.MemberDetails) (bool, error) {
	return controller.UpdateMember(db, ctx, memberdata)
}

func ChannelDetail(db *gorm.DB, ctx context.Context, channelID *int, channelSlug *string) (*model.Channel, error) {
	return controller.ChannelDetail(db, ctx, channelID, channelSlug)
}

func SpaceList(db *gorm.DB, ctx context.Context, limit, offset int, categoryId *int) (*model.SpaceDetails, error) {
	return controller.SpaceList(db, ctx, limit, offset, categoryId)
}

func SpaceDetails(db *gorm.DB, ctx context.Context, spaceId *int, spaceSlug *string) (*model.Space, error) {
	return controller.SpaceDetails(db, ctx, spaceId, spaceSlug)
}

func PagesAndPageGroupsUnderSpace(db *gorm.DB, ctx context.Context, spaceID int) (*model.PageAndPageGroups, error) {
	return controller.PagesAndPageGroupsBySpaceId(db, ctx, spaceID)
}

func CategoriesList(db *gorm.DB, ctx context.Context, limit, offset, categoryGroupId *int, categoryGroupSlug *string, hierarchyLevel, checkEntriesPresence *int) (*model.CategoriesList, error) {
	return controller.CategoriesList(db, ctx, limit, offset, categoryGroupId, categoryGroupSlug, hierarchyLevel, checkEntriesPresence)
}

func ChannelEntryDetail(db *gorm.DB, ctx context.Context, channelID *int, categoryID *int, channelEntryID *int, slug, profileSlug *string) (*model.ChannelEntries, error) {
	return controller.ChannelEntryDetail(db, ctx, channelEntryID, channelID, categoryID, slug, profileSlug)
}

func MemberProfileUpdate(db *gorm.DB, ctx context.Context, profiledata model.ProfileData) (bool, error) {
	return controller.MemberProfileUpdate(db, ctx, profiledata)
}

func VerifyMemberOtp(db *gorm.DB, ctx context.Context, email string, otp int) (*model.LoginDetails, error) {
	return controller.VerifyMemberOtp(db, ctx, email, otp)
}

func Memberclaimnow(db *gorm.DB, ctx context.Context, input model.ClaimData, profileId *int, profileSlug *string) (bool, error) {
	return controller.Memberclaimnow(db, ctx, input, profileId, profileSlug)
}

func EcommerceProductList(db *gorm.DB, ctx context.Context, limit int, offset int, filter *model.ProductFilter, sort *model.ProductSort) (*model.EcommerceProducts, error) {
	return controller.EcommerceProductList(db, ctx, limit, offset, filter, sort)
}

func VerifyProfileName(db *gorm.DB, ctx context.Context, profileSlug string, profileID int) (bool, error) {
	return controller.VerifyProfileName(db, ctx, profileSlug, profileID)
}

func TemplateMemberLogin(db *gorm.DB, ctx context.Context, username, email *string, password string, ecomModule *int) (string, error) {
	return controller.TemplateMemberLogin(db, ctx, username, email, password, ecomModule)
}

func EcommerceProductDetails(db *gorm.DB, ctx context.Context, productID *int, productSlug *string) (*model.EcommerceProduct, error) {
	return controller.EcommerceProductDetails(db, ctx, productID, productSlug)
}

func EcommerceCartList(db *gorm.DB, ctx context.Context, limit, offset int) (*model.EcommerceCartDetails, error) {
	return controller.EcommerceCartList(db, ctx, limit, offset)
}

func EcommerceAddToCart(db *gorm.DB, ctx context.Context, productID *int, productSlug *string, quantity int) (bool, error) {
	return controller.EcommerceAddToCart(db, ctx, productID, productSlug, quantity)
}

func UpdateChannelEntryViewCount(db *gorm.DB, ctx context.Context, entryId *int, slug *string) (bool, error) {
	return controller.UpdateChannelEntryViewCount(db, ctx, entryId, slug)
}

func RemoveProductFromCartlist(db *gorm.DB, ctx context.Context, productID int) (bool, error) {
	return controller.RemoveProductFromCartlist(db, ctx, productID)
}

func MemberProfileDetails(db *gorm.DB, ctx context.Context) (*model.MemberProfile, error) {
	return controller.MemberProfileDetails(db, ctx)
}

func EcommerceProductOrdersList(db *gorm.DB, ctx context.Context, limit int, offset int, filter *model.OrderFilter, sort *model.OrderSort) (*model.EcommerceProducts, error) {
	return controller.EcommerceProductOrdersList(db, ctx, limit, offset, filter, sort)
}

func EcommerceProductOrderDetails(db *gorm.DB, ctx context.Context, productID *int, productSlug *string, orderId int) (*model.EcomOrderedProductDetails, error) {
	return controller.EcommerceProductOrderDetails(db, ctx, productID, productSlug, orderId)
}

func EcommerceOrderPlacement(db *gorm.DB, ctx context.Context, paymentMode string, shippingAddress string, orderProducts []model.OrderProduct, orderSummary *model.OrderSummary) (bool, error) {
	return controller.EcommerceOrderPlacement(db, ctx, paymentMode, shippingAddress, orderProducts, orderSummary)
}

func GetMemberProfileDetails(db *gorm.DB, ctx context.Context, id *int, profileSlug *string) (*model.MemberProfile, error) {
	return controller.GetMemberProfileDetails(db, ctx, id, profileSlug)
}

func EcommerceCustomerDetails(db *gorm.DB, ctx context.Context) (*model.CustomerDetails, error) {
	return controller.EcommerceCustomerDetails(db, ctx)
}

func CustomerProfileUpdate(db *gorm.DB, ctx context.Context, customerDetails model.CustomerInput) (bool, error) {
	return controller.CustomerProfileUpdate(db, ctx, customerDetails)
}

func UpdateProductViewCount(db *gorm.DB, ctx context.Context, productID *int, productSlug *string) (bool, error) {
	return controller.UpdateProductViewCount(db, ctx, productID, productSlug)
}

func JobsList(db *gorm.DB, ctx context.Context, limit int, offset int, filter *model.JobFilter) (*model.JobsList, error) {
	return controller.JobsList(db, ctx, limit, offset, filter)
}

func JobDetail(db *gorm.DB, ctx context.Context, id *int, jobSlug *string) (*model.Job, error) {
	return controller.JobDetail(db, ctx, id, jobSlug)
}

func JobApplication(db *gorm.DB, ctx context.Context, applicationDetails model.ApplicationInput) (bool, error) {
	return controller.JobApplication(db, ctx, applicationDetails)
}

func MemberPasswordUpdate(db *gorm.DB, ctx context.Context, oldPassword string, newPassword string, confirmPassword string) (bool, error) {
	return controller.MemberPasswordUpdate(db, ctx, oldPassword, newPassword, confirmPassword)
}

func GetMemberDetails(db *gorm.DB, ctx context.Context) (*model.Member, error) {
	return controller.GetMemberDetails(db, ctx)
}

func EcommerceOrderStatusNames(db *gorm.DB, ctx context.Context) ([]model.OrderStatusNames, error) {
	return controller.EcommerceOrderStatusNames(db, ctx)
}
