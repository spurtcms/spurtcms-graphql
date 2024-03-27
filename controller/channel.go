package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gqlserver/graph/model"
	// "log"
	"os"
	"time"

	"html/template"

	"github.com/gin-gonic/gin"
	channel "github.com/spurtcms/pkgcontent/channels"
	"gorm.io/gorm"
)

func Channellist(db *gorm.DB, ctx context.Context, limit, offset int) (model.ChannelDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	channelAuth := channel.Channel{Authority: GetAuthorization(token.(string), db)}

	channelList, count, err := channelAuth.GetGraphqlChannelList(limit, offset)

	if err != nil {

		return model.ChannelDetails{}, err
	}

	var conv_channelList []model.Channel

	for _, channel := range channelList {

		conv_channel := model.Channel{
			ID:                 channel.Id,
			ChannelName:        channel.ChannelName,
			ChannelDescription: channel.ChannelDescription,
			SlugName:           channel.SlugName,
			FieldGroupID:       channel.FieldGroupId,
			IsActive:           channel.IsActive,
			CreatedOn:          channel.CreatedOn,
			ModifiedOn:         &channel.ModifiedOn,
			CreatedBy:          channel.CreatedBy,
			ModifiedBy:         &channel.ModifiedBy,
		}

		conv_channelList = append(conv_channelList, conv_channel)
	}

	return model.ChannelDetails{Channellist: conv_channelList, Count: int(count)}, nil

}

// this function provides the published channel entries list under a channel and channel entry details for a particular channeel entry by using its id
func ChannelEntriesList(db *gorm.DB, ctx context.Context, channelID, categoryId *int, limit, offset int, title *string) (model.ChannelEntriesDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	memberid := c.GetInt("memberid")

	channelAuth := channel.Channel{Authority: GetAuthorization(token.(string), db)}

	var channelEntries []channel.TblChannelEntries

	var count int64

	var err error

	channelEntries, count, err = channelAuth.GetGraphqlAllChannelEntriesList(channelID, categoryId, limit, offset, SectionTypeId, MemberFieldTypeId, PathUrl, title)

	if err != nil {

		return model.ChannelEntriesDetails{}, err
	}

	var conv_channelEntries []model.ChannelEntries

	for _, entry := range channelEntries {

		var conv_categories [][]model.Category

		for _, categories := range entry.Categories {

			var conv_categoryz []model.Category

			for _, category := range categories {

				conv_category := model.Category{
					ID:           category.Id,
					CategoryName: category.CategoryName,
					CategorySlug: category.CategorySlug,
					Description:  category.Description,
					ImagePath:    category.ImagePath,
					CreatedOn:    category.CreatedOn,
					CreatedBy:    category.CreatedBy,
					ModifiedOn:   &category.ModifiedOn,
					ModifiedBy:   &category.ModifiedBy,
					ParentID:     category.ParentId,
				}

				conv_categoryz = append(conv_categoryz, conv_category)

			}

			conv_categories = append(conv_categories, conv_categoryz)

		}

		authorDetails := &model.Author{
			AuthorID:         entry.AuthorDetail.AuthorID,
			FirstName:        entry.AuthorDetail.FirstName,
			LastName:         entry.AuthorDetail.LastName,
			Email:            entry.AuthorDetail.Email,
			MobileNo:         entry.AuthorDetail.MobileNo,
			IsActive:         entry.AuthorDetail.IsActive,
			ProfileImagePath: entry.AuthorDetail.ProfileImagePath,
			CreatedOn:        entry.AuthorDetail.CreatedOn,
			CreatedBy:        entry.AuthorDetail.CreatedBy,
		}

		var conv_sections []model.Section

		for _, section := range entry.Sections {

			conv_section := model.Section{
				SectionID:     &section.Id,
				SectionName:   section.FieldName,
				SectionTypeID: section.FieldTypeId,
				CreatedOn:     section.CreatedOn,
				CreatedBy:     section.CreatedBy,
				ModifiedOn:    &section.ModifiedOn,
				ModifiedBy:    &section.ModifiedBy,
				OrderIndex:    section.OrderIndex,
			}

			conv_sections = append(conv_sections, conv_section)

		}

		var conv_fields []model.Field

		for _, field := range entry.Fields {

			conv_field_value := model.FieldValue{
				ID:         field.FieldValue.FieldId,
				FieldValue: field.FieldValue.FieldValue,
				CreatedOn:  field.FieldValue.CreatedOn,
				CreatedBy:  field.FieldValue.CreatedBy,
				ModifiedOn: &field.FieldValue.ModifiedOn,
				ModifiedBy: &field.FieldValue.ModifiedBy,
			}

			var conv_fieldOptions []model.FieldOptions

			for _, field_option := range field.FieldOptions {

				conv_fieldOption := model.FieldOptions{
					ID:          field_option.Id,
					OptionName:  field_option.OptionName,
					OptionValue: field_option.OptionValue,
					CreatedOn:   field_option.CreatedOn,
					CreatedBy:   field_option.CreatedBy,
					ModifiedOn:  &field_option.ModifiedOn,
					ModifiedBy:  &field_option.ModifiedBy,
				}

				conv_fieldOptions = append(conv_fieldOptions, conv_fieldOption)
			}

			conv_field := model.Field{
				FieldID:          field.Id,
				FieldName:        field.FieldName,
				FieldTypeID:      field.FieldTypeId,
				MandatoryField:   field.MandatoryField,
				OptionExist:      field.OptionExist,
				CreatedOn:        field.CreatedOn,
				CreatedBy:        field.CreatedBy,
				ModifiedOn:       &field.ModifiedOn,
				ModifiedBy:       &field.ModifiedBy,
				FieldDesc:        field.FieldDesc,
				OrderIndex:       field.OrderIndex,
				ImagePath:        field.ImagePath,
				DatetimeFormat:   &field.DatetimeFormat,
				TimeFormat:       &field.TimeFormat,
				SectionParentID:  &field.SectionParentId,
				CharacterAllowed: &field.CharacterAllowed,
				FieldTypeName:    field.FieldTypeName,
				FieldValue:       &conv_field_value,
				FieldOptions:     conv_fieldOptions,
			}

			conv_fields = append(conv_fields, conv_field)

		}

		additionalFields := &model.AdditionalFields{Sections: conv_sections, Fields: conv_fields}

		var conv_memberProfiles []model.MemberProfile

		claimStatus := false

		for _, memberProfile := range entry.MemberProfiles {

			if memberid == memberProfile.MemberId && memberProfile.ClaimStatus == 1 {

				claimStatus = true
	
			}

			conv_MemberProfile := model.MemberProfile{
				ID:              &memberProfile.Id,
				MemberID:        &memberProfile.MemberId,
				ProfileName:     &memberProfile.ProfileName,
				ProfileSlug:     &memberProfile.ProfileSlug,
				ProfilePage:     &memberProfile.ProfilePage,
				MemberDetails:   memberProfile.MemberDetails,
				CompanyName:     &memberProfile.CompanyName,
				CompanyLocation: &memberProfile.CompanyLocation,
				CompanyLogo:     &memberProfile.CompanyLogo,
				About:           &memberProfile.About,
				SeoTitle:        &memberProfile.SeoTitle,
				SeoDescription:  &memberProfile.SeoDescription,
				SeoKeyword:      &memberProfile.SeoKeyword,
				CreatedBy:       &memberProfile.CreatedBy,
				CreatedOn:       &memberProfile.CreatedOn,
				ModifiedOn:      &memberProfile.ModifiedOn,
				ModifiedBy:      &memberProfile.ModifiedBy,
				Linkedin:        &memberProfile.Linkedin,
				Twitter:         &memberProfile.Twitter,
				Website:         &memberProfile.Website,
				ClaimStatus:     &memberProfile.ClaimStatus,
			}

			conv_memberProfiles = append(conv_memberProfiles, conv_MemberProfile)

		}

		conv_channelEntry := model.ChannelEntries{
			ID:               entry.Id,
			Title:            entry.Title,
			Slug:             entry.Slug,
			Description:      entry.Description,
			UserID:           entry.UserId,
			ChannelID:        entry.ChannelId,
			Status:           entry.Status,
			IsActive:         entry.IsActive,
			CreatedOn:        entry.CreatedOn,
			CreatedBy:        entry.CreatedBy,
			ModifiedBy:       &entry.ModifiedBy,
			ModifiedOn:       &entry.ModifiedOn,
			CoverImage:       entry.CoverImage,
			ThumbnailImage:   entry.ThumbnailImage,
			MetaTitle:        entry.MetaTitle,
			MetaDescription:  entry.MetaDescription,
			Keyword:          entry.Keyword,
			CategoriesID:     entry.CategoriesId,
			RelatedArticles:  entry.RelatedArticles,
			Categories:       conv_categories,
			AdditionalFields: additionalFields,
			MemberProfile:    conv_memberProfiles,
			AuthorDetails:    authorDetails,
			FeaturedEntry:    entry.Feature,
			ViewCount:        entry.ViewCount,
			ClaimStatus:      claimStatus,
		}

		conv_channelEntries = append(conv_channelEntries, conv_channelEntry)

	}

	channelEntryDetails := model.ChannelEntriesDetails{ChannelEntriesList: conv_channelEntries, Count: int(count)}

	return channelEntryDetails, nil

}

func ChannelDetail(db *gorm.DB, ctx context.Context, channelID int) (model.Channel, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	channelAuth := channel.Channel{Authority: GetAuthorization(token.(string), db)}

	channel, err := channelAuth.GetGraphqlChannelDetails(channelID)

	conv_channel := model.Channel{
		ID:                 channel.Id,
		ChannelName:        channel.ChannelName,
		ChannelDescription: channel.ChannelDescription,
		SlugName:           channel.SlugName,
		FieldGroupID:       channel.FieldGroupId,
		IsActive:           channel.IsActive,
		CreatedOn:          channel.CreatedOn,
		ModifiedOn:         &channel.ModifiedOn,
		CreatedBy:          channel.CreatedBy,
		ModifiedBy:         &channel.ModifiedBy,
	}

	if err != nil {

		return model.Channel{}, err
	}

	return conv_channel, nil
}

func ChannelEntryDetail(db *gorm.DB, ctx context.Context, channelEntryId, channelId, categoryId *int, slug *string) (model.ChannelEntries, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	memberid := c.GetInt("memberid")

	channelAuth := channel.Channel{Authority: GetAuthorization(token.(string), db)}

	channelEntry, err := channelAuth.GetGraphqlChannelEntriesDetails(channelEntryId, channelId, categoryId, PathUrl, SectionTypeId, MemberFieldTypeId, slug)

	if err != nil {

		return model.ChannelEntries{}, err
	}

	var conv_categories [][]model.Category

	for _, categories := range channelEntry.Categories {

		var conv_categoryz []model.Category

		for _, category := range categories {

			conv_category := model.Category{
				ID:           category.Id,
				CategoryName: category.CategoryName,
				CategorySlug: category.CategorySlug,
				Description:  category.Description,
				ImagePath:    category.ImagePath,
				CreatedOn:    category.CreatedOn,
				CreatedBy:    category.CreatedBy,
				ModifiedOn:   &category.ModifiedOn,
				ModifiedBy:   &category.ModifiedBy,
				ParentID:     category.ParentId,
			}

			conv_categoryz = append(conv_categoryz, conv_category)

		}

		conv_categories = append(conv_categories, conv_categoryz)

	}
	authorDetails := &model.Author{
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

		conv_section := model.Section{
			SectionID:     &section.Id,
			SectionName:   section.FieldName,
			SectionTypeID: section.FieldTypeId,
			CreatedOn:     section.CreatedOn,
			CreatedBy:     section.CreatedBy,
			ModifiedOn:    &section.ModifiedOn,
			ModifiedBy:    &section.ModifiedBy,
			OrderIndex:    section.OrderIndex,
		}

		conv_sections = append(conv_sections, conv_section)
	}

	var conv_fields []model.Field

	for _, field := range channelEntry.Fields {

		conv_field_value := model.FieldValue{
			ID:         field.FieldValue.FieldId,
			FieldValue: field.FieldValue.FieldValue,
			CreatedOn:  field.FieldValue.CreatedOn,
			CreatedBy:  field.FieldValue.CreatedBy,
			ModifiedOn: &field.FieldValue.ModifiedOn,
			ModifiedBy: &field.FieldValue.ModifiedBy,
		}

		var conv_fieldOptions []model.FieldOptions

		for _, field_option := range field.FieldOptions {

			conv_fieldOption := model.FieldOptions{
				ID:          field_option.Id,
				OptionName:  field_option.OptionName,
				OptionValue: field_option.OptionValue,
				CreatedOn:   field_option.CreatedOn,
				CreatedBy:   field_option.CreatedBy,
				ModifiedOn:  &field_option.ModifiedOn,
				ModifiedBy:  &field_option.ModifiedBy,
			}

			conv_fieldOptions = append(conv_fieldOptions, conv_fieldOption)
		}

		conv_field := model.Field{
			FieldID:          field.Id,
			FieldName:        field.FieldName,
			FieldTypeID:      field.FieldTypeId,
			MandatoryField:   field.MandatoryField,
			OptionExist:      field.OptionExist,
			CreatedOn:        field.CreatedOn,
			CreatedBy:        field.CreatedBy,
			ModifiedOn:       &field.ModifiedOn,
			ModifiedBy:       &field.ModifiedBy,
			FieldDesc:        field.FieldDesc,
			OrderIndex:       field.OrderIndex,
			ImagePath:        field.ImagePath,
			DatetimeFormat:   &field.DatetimeFormat,
			TimeFormat:       &field.TimeFormat,
			SectionParentID:  &field.SectionParentId,
			CharacterAllowed: &field.CharacterAllowed,
			FieldTypeName:    field.FieldTypeName,
			FieldValue:       &conv_field_value,
			FieldOptions:     conv_fieldOptions,
		}

		conv_fields = append(conv_fields, conv_field)
	}

	additionalFields := &model.AdditionalFields{Sections: conv_sections, Fields: conv_fields}

	var conv_memberProfiles []model.MemberProfile

	claimStatus := false

	for _, memberProfile := range channelEntry.MemberProfiles {

		if memberid == memberProfile.MemberId && memberProfile.ClaimStatus == 1 {

			claimStatus = true

		}	

		conv_MemberProfile := model.MemberProfile{
			ID:              &memberProfile.Id,
			MemberID:        &memberProfile.MemberId,
			ProfileName:     &memberProfile.ProfileName,
			ProfileSlug:     &memberProfile.ProfileSlug,
			ProfilePage:     &memberProfile.ProfilePage,
			MemberDetails:   memberProfile.MemberDetails,
			CompanyName:     &memberProfile.CompanyName,
			CompanyLocation: &memberProfile.CompanyLocation,
			CompanyLogo:     &memberProfile.CompanyLogo,
			About:           &memberProfile.About,
			SeoTitle:        &memberProfile.SeoTitle,
			SeoDescription:  &memberProfile.SeoDescription,
			SeoKeyword:      &memberProfile.SeoKeyword,
			CreatedBy:       &memberProfile.CreatedBy,
			CreatedOn:       &memberProfile.CreatedOn,
			ModifiedOn:      &memberProfile.ModifiedOn,
			ModifiedBy:      &memberProfile.ModifiedBy,
			Linkedin:        &memberProfile.Linkedin,
			Twitter:         &memberProfile.Twitter,
			Website:         &memberProfile.Website,
			ClaimStatus:     &memberProfile.ClaimStatus,
		}

		conv_memberProfiles = append(conv_memberProfiles, conv_MemberProfile)
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
		MemberProfile:    conv_memberProfiles,
		AuthorDetails:    authorDetails,
		FeaturedEntry:    channelEntry.Feature,
		ViewCount:        channelEntry.ViewCount,
		ClaimStatus:      claimStatus,
	}

	return conv_channelEntry, nil

}

func Memberclaimnow(db *gorm.DB, ctx context.Context, profileData model.ClaimData, entryId int) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	if token == SpecialToken {

		return false, errors.New("login required")
	}

	verify_chan := make(chan bool)

	var channelEntry model.ChannelEntries

	if err := db.Table("tbl_channel_entries").Where("is_deleted = 0 and id = ?", entryId).First(&channelEntry).Error; err != nil {

		return false, err
	}

	var AuthorDetails model.Author

	if err := db.Table("tbl_users").Where("is_deleted = 0 and id = ?", channelEntry.CreatedBy).First(&AuthorDetails).Error; err != nil {

		return false, err
	}

	data := map[string]interface{}{"claimData": profileData, "authorDetails": AuthorDetails, "entry": channelEntry, "additionalData": AdditionalData}

	tmpl, _ := template.ParseFiles("view/email/claim-template.html")

	var template_buff bytes.Buffer

	err := tmpl.Execute(&template_buff, data)

	if err != nil {

		return false, err
	}

	mail_data := MailConfig{Email: AuthorDetails.Email, MailUsername: os.Getenv("MAIL_USERNAME"), MailPassword: os.Getenv("MAIL_PASSWORD"), Subject: "OwnDesk - Member Profile Claim Request"}

	html_content := template_buff.String()

	go SendMail(mail_data, html_content, verify_chan)

	if <-verify_chan {

		return true, nil

	} else {

		return false, nil
	}
}

func MemberProfileUpdate(db *gorm.DB, ctx context.Context, profiledata model.ProfileData, entryId int, updateExactMemberProfileOnly bool) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	if token == SpecialToken {

		return false, errors.New("login required")
	}

	memberid := c.GetInt("memberid")

	channelAuth := channel.Channel{Authority: GetAuthorization(token, db)}

	entryDetails, err := channelAuth.GetGraphqlChannelEntriesDetails(&entryId, nil, nil, PathUrl, SectionTypeId, MemberFieldTypeId, nil)

	if err != nil {

		return false, err
	}

	var claimedMembers []int

	for _, memberProfile := range entryDetails.MemberProfiles {

		if updateExactMemberProfileOnly {

			if memberProfile.MemberId == memberid {

				claimedMembers = append(claimedMembers, memberProfile.MemberId)

			}

		} else {

			claimedMembers = append(claimedMembers, memberProfile.MemberId)

		}

	}

	var jsonData map[string]interface{}

	err = json.Unmarshal([]byte(profiledata.MemberProfile), &jsonData)

	if err != nil {

		return false, err
	}

	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().In(TimeZone).Format("2006-01-02 15:04:05"))

	memberProfileDetails := model.MemberProfile{
		MemberDetails: profiledata.MemberProfile,
		Linkedin:      profiledata.Linkedin,
		Twitter:       profiledata.Twitter,
		Website:       profiledata.Website,
		ModifiedOn:    &currentTime,
	}

	if err := db.Table("tbl_member_profiles").Where("is_deleted = 0 and claim_status = 1 and member_id in (?)", claimedMembers).UpdateColumns(map[string]interface{}{"member_details": memberProfileDetails.MemberDetails, "linkedin": memberProfileDetails.Linkedin, "twitter": memberProfileDetails.Twitter, "website": memberProfileDetails.Website, "modified_on": memberProfileDetails.ModifiedOn}).Error; err != nil {

		return false, err
	}

	return true, nil
}
