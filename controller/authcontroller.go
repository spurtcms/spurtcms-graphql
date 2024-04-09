package controller

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"log"
	"math/rand"
	"os"
	"spurtcms-graphql/graph/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spurtcms/pkgcontent/channels"
	"github.com/spurtcms/pkgcore/member"
	"gorm.io/gorm"
)

func MemberLogin(db *gorm.DB, ctx context.Context, email string) (bool, error) {

	Mem.Auth = GetAuthorizationWithoutToken(db)

	member_details, err := Mem.GraphqlMemberLogin(email)

	if err != nil {

		return false, err
	}

	conv_member := model.Member{
		ID:               member_details.Id,
		FirstName:        member_details.FirstName,
		LastName:         member_details.LastName,
		Email:            member_details.Email,
		MobileNo:         member_details.MobileNo,
		IsActive:         member_details.IsActive,
		ProfileImagePath: member_details.ProfileImagePath,
		CreatedOn:        member_details.CreatedOn,
		CreatedBy:        member_details.CreatedBy,
		ModifiedOn:       &member_details.ModifiedOn,
		ModifiedBy:       &member_details.ModifiedBy,
	}

	channel := make(chan bool)

	// rand.Seed(time.Now().UnixNano())

	otp := rand.Intn(900000) + 100000

	current_time := time.Now().In(TimeZone)

	otp_expiry_time := current_time.Add(5 * time.Minute).Format("2006-01-02 15:04:05")

	mail_expiry_time := current_time.Add(5 * time.Minute).Format("02 Jan 2006 03:04 PM")

	err = Mem.StoreGraphqlMemberOtp(otp, conv_member.ID, otp_expiry_time)

	if err != nil {

		return false, err
	}

	data := map[string]interface{}{"otp": otp, "expiryTime": mail_expiry_time, "member": conv_member, "additionalData": AdditionalData}

	tmpl, err := template.ParseFiles("view/email/login-template.html")

	if err != nil {

		return false, err
	}

	var template_buffer bytes.Buffer

	if err := tmpl.Execute(&template_buffer, data); err != nil {

		return false, err
	}

	mail_data := MailConfig{Email: conv_member.Email, MailUsername: os.Getenv("MAIL_USERNAME"), MailPassword: os.Getenv("MAIL_PASSWORD"), Subject: "OwnDesk - Login Otp Confirmation"}

	html_content := template_buffer.String()

	go SendMail(mail_data, html_content, channel)

	if <-channel {

		return true, nil

	} else {

		return false, errors.New("mail sent unsuccess")
	}
}

func VerifyMemberOtp(db *gorm.DB, ctx context.Context, email string, otp int) (*model.LoginDetails, error) {

	Mem.Auth = GetAuthorizationWithoutToken(db)

	currentTime := time.Now().In(TimeZone).Unix()

	memberDetails, token, err := Mem.VerifyLoginOtp(email, otp, currentTime)

	if err != nil {

		return &model.LoginDetails{}, err
	}

	log.Println("memberdetails", memberDetails)

	var channelEntryDetails model.ChannelEntries

	if err := db.Debug().Table("tbl_channel_entries").Select("tbl_channel_entries.*").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id ").Joins("inner join tbl_channel_entry_fields on tbl_channel_entry_fields.channel_entry_id = tbl_channel_entries.id").Joins("inner join tbl_fields on tbl_fields.id = tbl_channel_entry_fields.field_id").Joins("inner join tbl_field_types on tbl_field_types.id = tbl_fields.field_type_id").Joins("inner join tbl_members on tbl_members.id = any(string_to_array(tbl_channel_entry_fields.field_value,',')::integer[])").
		Joins("inner join tbl_member_profiles on tbl_members.id = tbl_member_profiles.member_id").Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_field_types.is_deleted = 0 and tbl_fields.is_deleted = 0 and tbl_members.is_deleted = 0 and tbl_member_profiles.is_deleted = 0 and tbl_member_profiles.claim_status = 1 and tbl_field_types.id = ? and tbl_members.id = ?", MemberFieldTypeId, memberDetails.Id).First(&channelEntryDetails).Error; err != nil {

		return &model.LoginDetails{}, err
	}

	channelAuth := channels.Channel{Authority: GetAuthorization(token, db)}

	channelEntry, err := channelAuth.GetGraphqlChannelEntriesDetails(&channelEntryDetails.ID, &channelEntryDetails.ChannelID, nil, PathUrl, SectionTypeId, MemberFieldTypeId, nil)

	if err != nil {

		return &model.LoginDetails{}, err
	}

	var conv_categories [][]model.Category

	for _, categories := range channelEntry.Categories {

		var conv_categoryz []model.Category

		for _, category := range categories {

			categoryModon := category.ModifiedOn

			categoryModBy := category.ModifiedBy

			conv_category := model.Category{
				ID:           category.Id,
				CategoryName: category.CategoryName,
				CategorySlug: category.CategorySlug,
				Description:  category.Description,
				ImagePath:    category.ImagePath,
				CreatedOn:    category.CreatedOn,
				CreatedBy:    category.CreatedBy,
				ModifiedOn:   &categoryModon,
				ModifiedBy:   &categoryModBy,
				ParentID:     category.ParentId,
			}

			conv_categoryz = append(conv_categoryz, conv_category)

		}

		conv_categories = append(conv_categories, conv_categoryz)

	}
	authorDetails := model.Author{
		AuthorID:         channelEntry.AuthorDetail.AuthorID,
		FirstName:        channelEntry.AuthorDetail.FirstName,
		LastName:         channelEntry.AuthorDetail.LastName,
		Email:            channelEntry.AuthorDetail.Email,
		MobileNo:         channelEntry.AuthorDetail.MobileNo,
		IsActive:         channelEntry.AuthorDetail.IsActive,
		ProfileImagePath: channelEntry.AuthorDetail.ProfileImagePath,
		CreatedOn:        channelEntry.AuthorDetail.CreatedOn,
		CreatedBy:        channelEntry.AuthorDetail.CreatedBy,
	}

	var conv_sections []model.Section

	for _, section := range channelEntry.Sections {

		sectionId := section.Id

		sectionModon := section.ModifiedOn

		sectionModBy := section.ModifiedBy

		conv_section := model.Section{
			SectionID:     &sectionId,
			SectionName:   section.FieldName,
			SectionTypeID: section.FieldTypeId,
			CreatedOn:     section.CreatedOn,
			CreatedBy:     section.CreatedBy,
			ModifiedOn:    &sectionModon,
			ModifiedBy:    &sectionModBy,
			OrderIndex:    section.OrderIndex,
		}

		conv_sections = append(conv_sections, conv_section)
	}

	var conv_fields []model.Field

	for _, field := range channelEntry.Fields {

		fieldValueModon := field.FieldValue.ModifiedOn

		fieldValueModBy := field.FieldValue.ModifiedBy

		conv_field_value := model.FieldValue{
			ID:         field.FieldValue.FieldId,
			FieldValue: field.FieldValue.FieldValue,
			CreatedOn:  field.FieldValue.CreatedOn,
			CreatedBy:  field.FieldValue.CreatedBy,
			ModifiedOn: &fieldValueModon,
			ModifiedBy: &fieldValueModBy,
		}

		var conv_fieldOptions []model.FieldOptions

		for _, field_option := range field.FieldOptions {

			optionModOn := field_option.ModifiedOn

			optionModBy := field_option.ModifiedBy

			conv_fieldOption := model.FieldOptions{
				ID:          field_option.Id,
				OptionName:  field_option.OptionName,
				OptionValue: field_option.OptionValue,
				CreatedOn:   field_option.CreatedOn,
				CreatedBy:   field_option.CreatedBy,
				ModifiedOn:  &optionModOn,
				ModifiedBy:  &optionModBy,
			}

			conv_fieldOptions = append(conv_fieldOptions, conv_fieldOption)
		}

		fieldModon := field.ModifiedOn

		fieldModBy := field.ModifiedBy

		fieldDateTime := field.DatetimeFormat

		fieldTime := field.TimeFormat

		fieldSectionParentId := field.SectionParentId

		fieldCharAllowed := field.CharacterAllowed

		conv_field := model.Field{
			FieldID:          field.Id,
			FieldName:        field.FieldName,
			FieldTypeID:      field.FieldTypeId,
			MandatoryField:   field.MandatoryField,
			OptionExist:      field.OptionExist,
			CreatedOn:        field.CreatedOn,
			CreatedBy:        field.CreatedBy,
			ModifiedOn:       &fieldModon,
			ModifiedBy:       &fieldModBy,
			FieldDesc:        field.FieldDesc,
			OrderIndex:       field.OrderIndex,
			ImagePath:        field.ImagePath,
			DatetimeFormat:   &fieldDateTime,
			TimeFormat:       &fieldTime,
			SectionParentID:  &fieldSectionParentId,
			CharacterAllowed: &fieldCharAllowed,
			FieldTypeName:    field.FieldTypeName,
			FieldValue:       &conv_field_value,
			FieldOptions:     conv_fieldOptions,
		}

		conv_fields = append(conv_fields, conv_field)
	}

	additionalFields := &model.AdditionalFields{Sections: conv_sections, Fields: conv_fields}

	    memberProfileId := channelEntry.MemberProfile.Id
		memberProfileMemId := channelEntry.MemberProfile.MemberId
        memberProfileName := channelEntry.MemberProfile.ProfileName
		memberProfileSlug := channelEntry.MemberProfile.ProfileSlug
		memberProfilePage := channelEntry.MemberProfile.ProfilePage
		memberProfileMemDetails := channelEntry.MemberProfile.MemberDetails
		memberProfileComName := channelEntry.MemberProfile.CompanyName
		memberProfileComLocation := channelEntry.MemberProfile.CompanyLocation
		memberProfileComLogo := channelEntry.MemberProfile.CompanyLogo
		memberProfileAbout := channelEntry.MemberProfile.About
		memberProfileSeoTitle := channelEntry.MemberProfile.SeoTitle
		memberProfileSeoDesc := channelEntry.MemberProfile.SeoDescription
		memberProfileSeoKey := channelEntry.MemberProfile.SeoKeyword
		memberProfileCreateBy := channelEntry.MemberProfile.CreatedBy
		memberProfileCreateOn := channelEntry.MemberProfile.CreatedOn
		memberProfileModon := channelEntry.MemberProfile.ModifiedOn
		memberProfileModBy := channelEntry.MemberProfile.ModifiedBy
		memberProfileLinkedin := channelEntry.MemberProfile.Linkedin
		memberProfileTwitter := channelEntry.MemberProfile.Twitter
		memberProfileWeb := channelEntry.MemberProfile.Website
		memberProfileClaim := channelEntry.MemberProfile.ClaimStatus


		MemberProfile := model.MemberProfile{
			ID:              &memberProfileId,
			MemberID:        &memberProfileMemId,
			ProfileName:     &memberProfileName,
			ProfileSlug:     &memberProfileSlug,
			ProfilePage:     &memberProfilePage,
			MemberDetails:   &memberProfileMemDetails,
			CompanyName:     &memberProfileComName,
			CompanyLocation: &memberProfileComLocation,
			CompanyLogo:     &memberProfileComLogo,
			About:           &memberProfileAbout,
			SeoTitle:        &memberProfileSeoTitle,
			SeoDescription:  &memberProfileSeoDesc,
			SeoKeyword:      &memberProfileSeoKey,
			CreatedBy:       &memberProfileCreateBy,
			CreatedOn:       &memberProfileCreateOn,
			ModifiedOn:      &memberProfileModon,
			ModifiedBy:      &memberProfileModBy,
			Linkedin:        &memberProfileLinkedin,
			Twitter:         &memberProfileTwitter,
			Website:         &memberProfileWeb,
			ClaimStatus:     &memberProfileClaim,
		}

	var claimStatus bool

	if channelEntry.MemberProfile.ClaimStatus == 1 {

		claimStatus = true

	} else {

		claimStatus = false
	}

	conv_channelEntry := model.ChannelEntries{
		ID:               channelEntry.Id,
		Title:            channelEntry.Title,
		Slug:             channelEntry.Slug,
		Description:      channelEntry.Description,
		UserID:           channelEntry.UserId,
		ChannelID:        channelEntry.ChannelId,
		Status:           channelEntry.Status,
		IsActive:         channelEntry.IsActive,
		CreatedOn:        channelEntry.CreatedOn,
		CreatedBy:        channelEntry.CreatedBy,
		ModifiedBy:       &channelEntry.ModifiedBy,
		ModifiedOn:       &channelEntry.ModifiedOn,
		CoverImage:       channelEntry.CoverImage,
		ThumbnailImage:   channelEntry.ThumbnailImage,
		MetaTitle:        channelEntry.MetaTitle,
		MetaDescription:  channelEntry.MetaDescription,
		Keyword:          channelEntry.Keyword,
		CategoriesID:     channelEntry.CategoriesId,
		RelatedArticles:  channelEntry.RelatedArticles,
		Categories:       conv_categories,
		AdditionalFields: additionalFields,
		MemberProfile:    MemberProfile,
		AuthorDetails:    authorDetails,
		FeaturedEntry:    channelEntry.Feature,
		ViewCount:        channelEntry.ViewCount,
		ClaimStatus:      claimStatus,
		Author:           &channelEntry.Author,
		SortOrder:        &channelEntry.SortOrder,
		CreateTime:       &channelEntry.CreateTime,
		PublishedTime:    &channelEntry.PublishedTime,
		ReadingTime:      &channelEntry.ReadingTime,
		Tags:             &channelEntry.Tags,
		Excerpt:          &channelEntry.Excerpt,
	}

	return &model.LoginDetails{ClaimEntryDetails: conv_channelEntry, Token: token}, nil

}

func MemberRegister(db *gorm.DB, input model.MemberDetails) (bool, error) {

	Mem.Auth = GetAuthorizationWithoutToken(db)

	var imageName, imagePath string

	var err error

	if input.ProfileImage.IsSet() {

		imageName, imagePath, err = StoreImageBase64ToLocal(*input.ProfileImage.Value(), ProfileImagePath, "PROFILE")

		if err != nil {

			return false, err
		}

	}

	memberDetails := member.MemberCreation{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		MobileNo:  input.Mobile,
		Email:     input.Email,
		Password:  input.Password,
		//  Username: *input.Username,
		ProfileImage:     imageName,
		ProfileImagePath: imagePath,
	}

	_, isMemberExists, err := Mem.CheckEmailInMember(0, input.Email)

	if isMemberExists || err != nil {

		return isMemberExists, errors.New("Member already exists!")
	}

	isRegistered, err := Mem.MemberRegister(memberDetails)

	if !isRegistered || err != nil {

		return isRegistered, err
	}

	return isRegistered, nil

}

func UpdateMember(db *gorm.DB, ctx context.Context, memberdata model.MemberDetails) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	Mem.Auth = GetAuthorization(token, db)

	var imageName, imagePath string

	var err error

	if memberdata.ProfileImage.IsSet() {

		imageName, imagePath, err = StoreImageBase64ToLocal(*memberdata.ProfileImage.Value(), ProfileImagePath, "PROFILE")

		if err != nil {

			return false, err
		}

	}

	memberDetails := member.MemberCreation{
		FirstName:        memberdata.FirstName,
		LastName:         memberdata.LastName,
		MobileNo:         memberdata.Mobile,
		Email:            memberdata.Email,
		Password:         memberdata.Password,
		ProfileImage:     imageName,
		ProfileImagePath: imagePath,
		// IsActive: *memberdata.IsActive,
		// Username: *memberdata.Username,
		// GroupId: *memberdata.GroupID,
	}

	isUpdated, err := Mem.MemberUpdate(memberDetails)

	if err != nil || !isUpdated {

		return isUpdated, err
	}

	return isUpdated, nil

}

func TemplateMemberLogin(db *gorm.DB, ctx context.Context, username string, password string) (string, error) {

	Mem.Auth = GetAuthorizationWithoutToken(db)

	member_details, err := Mem.CheckMemberLogin(member.MemberLogin{Username: username, Password: password}, db, os.Getenv("JWT_SECRET"))

	if err != nil {

		log.Println(err)
	}

	return member_details, err
}
