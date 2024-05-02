package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"os"
	"spurtcms-graphql/graph/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	channel "github.com/spurtcms/pkgcontent/channels"
	"gorm.io/gorm"
)

func Channellist(db *gorm.DB, ctx context.Context, limit, offset int) (*model.ChannelDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	channelAuth := channel.Channel{Authority: GetAuthorization(token.(string), db)}

	channelList, count, err := channelAuth.GetGraphqlChannelList(limit, offset)

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.ChannelDetails{}, err
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

	return &model.ChannelDetails{Channellist: conv_channelList, Count: int(count)}, nil

}

// this function provides the published channel entries list under a channel and channel entry details for a particular channeel entry by using its id
func ChannelEntriesList(db *gorm.DB, ctx context.Context, channelID, categoryId *int, limit, offset int, title *string, categoryChildId *int, categorySlug, categoryChildSlug *string) (*model.ChannelEntriesDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	// memberid := c.GetInt("memberid")

	channelAuth := channel.Channel{Authority: GetAuthorization(token.(string), db)}

	var channelEntries []channel.TblChannelEntries

	var count int64

	var err error

	channelEntries, count, err = channelAuth.GetGraphqlAllChannelEntriesList(channelID, categoryId, limit, offset, SectionTypeId, MemberFieldTypeId, PathUrl, title, categoryChildId, categorySlug, categoryChildSlug)

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.ChannelEntriesDetails{}, err
	}

	conv_channelEntries := make([]model.ChannelEntries, len(channelEntries))

	for index, entry := range channelEntries {

		conv_categories := make([][]model.Category, len(entry.Categories))

		for cat_index, categories := range entry.Categories {

			conv_categoryz := make([]model.Category, len(categories))

			for i, category := range categories {

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

				conv_categoryz[i] = conv_category

			}

			conv_categories[cat_index] = conv_categoryz
		}

		conv_channelEntries[index].Categories = conv_categories

		authorMobnumber := entry.AuthorDetail.MobileNo

		authorIsActive := entry.AuthorDetail.IsActive

		authorProfileImage := entry.AuthorDetail.ProfileImagePath

		authorDetails := model.Author{
			AuthorID:         entry.AuthorDetail.AuthorID,
			FirstName:        entry.AuthorDetail.FirstName,
			LastName:         entry.AuthorDetail.LastName,
			Email:            entry.AuthorDetail.Email,
			MobileNo:         &authorMobnumber,
			IsActive:         &authorIsActive,
			ProfileImagePath: &authorProfileImage,
			CreatedOn:        entry.AuthorDetail.CreatedOn,
			CreatedBy:        entry.AuthorDetail.CreatedBy,
		}

		conv_channelEntries[index].AuthorDetails = authorDetails

		conv_sections := make([]model.Section, len(entry.Sections))

		for section_index, section := range entry.Sections {

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

			conv_sections[section_index] = conv_section

		}

		conv_fields := make([]model.Field, len(entry.Fields))

		for field_index, field := range entry.Fields {

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

			conv_fieldOptions := make([]model.FieldOptions, len(field.FieldOptions))

			for option_index, field_option := range field.FieldOptions {

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

				conv_fieldOptions[option_index] = conv_fieldOption
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

			conv_fields[field_index] = conv_field

		}

		additionalFields := model.AdditionalFields{Sections: conv_sections, Fields: conv_fields}

		conv_channelEntries[index].AdditionalFields = &additionalFields

		memberProfileId := entry.MemberProfile.Id
		memberProfileMemId := entry.MemberProfile.MemberId
		memberProfileName := entry.MemberProfile.ProfileName
		memberProfileSlug := entry.MemberProfile.ProfileSlug
		memberProfilePage := entry.MemberProfile.ProfilePage
		memberProfileMemDetails := entry.MemberProfile.MemberDetails
		memberProfileComName := entry.MemberProfile.CompanyName
		memberProfileComLocation := entry.MemberProfile.CompanyLocation
		memberProfileComLogo := entry.MemberProfile.CompanyLogo
		memberProfileAbout := entry.MemberProfile.About
		memberProfileSeoTitle := entry.MemberProfile.SeoTitle
		memberProfileSeoDesc := entry.MemberProfile.SeoDescription
		memberProfileSeoKey := entry.MemberProfile.SeoKeyword
		memberProfileCreateBy := entry.MemberProfile.CreatedBy
		memberProfileCreateOn := entry.MemberProfile.CreatedOn
		memberProfileModon := entry.MemberProfile.ModifiedOn
		memberProfileModBy := entry.MemberProfile.ModifiedBy
		memberProfileLinkedin := entry.MemberProfile.Linkedin
		memberProfileTwitter := entry.MemberProfile.Twitter
		memberProfileWeb := entry.MemberProfile.Website
		memberProfileClaim := entry.MemberProfile.ClaimStatus

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

		conv_channelEntries[index].MemberProfile = MemberProfile

		conv_channelEntries[index].Author = &entry.Author

		conv_channelEntries[index].CategoriesID = entry.CategoriesId

		conv_channelEntries[index].ChannelID = entry.ChannelId

		conv_channelEntries[index].CoverImage = entry.CoverImage

		conv_channelEntries[index].CreateTime = &entry.CreateTime

		conv_channelEntries[index].CreatedBy = entry.CreatedBy

		conv_channelEntries[index].CreatedOn = entry.CreatedOn

		conv_channelEntries[index].Description = entry.Description

		conv_channelEntries[index].Excerpt = &entry.Excerpt

		conv_channelEntries[index].FeaturedEntry = entry.Feature

		conv_channelEntries[index].ID = entry.Id

		conv_channelEntries[index].IsActive = entry.IsActive

		conv_channelEntries[index].Keyword = entry.Keyword

		conv_channelEntries[index].MetaDescription = entry.MetaDescription

		conv_channelEntries[index].MetaTitle = entry.MetaTitle

		modifiedBy := entry.ModifiedBy

		conv_channelEntries[index].ModifiedBy = &modifiedBy

		modifiedOn := entry.ModifiedOn

		conv_channelEntries[index].ModifiedOn = &modifiedOn

		publishedOn := entry.PublishedTime

		conv_channelEntries[index].PublishedTime = &publishedOn

		readingTime := entry.ReadingTime

		conv_channelEntries[index].ReadingTime = &readingTime

		conv_channelEntries[index].RelatedArticles = entry.RelatedArticles

		conv_channelEntries[index].Slug = entry.Slug

		sortOrder := entry.SortOrder

		conv_channelEntries[index].SortOrder = &sortOrder

		conv_channelEntries[index].Status = entry.Status

		tags := entry.Tags

		conv_channelEntries[index].Tags = &tags

		conv_channelEntries[index].ThumbnailImage = entry.ThumbnailImage

		conv_channelEntries[index].Title = entry.Title

		conv_channelEntries[index].UserID = entry.UserId

		conv_channelEntries[index].ViewCount = entry.ViewCount

		imageAltTag := entry.ImageAltTag

		conv_channelEntries[index].ImageAltTag = &imageAltTag

	}

	channelEntryDetails := model.ChannelEntriesDetails{ChannelEntriesList: conv_channelEntries, Count: int(count)}

	return &channelEntryDetails, nil

}

func ChannelDetail(db *gorm.DB, ctx context.Context, channelID *int, channelSlug *string) (*model.Channel, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	channelAuth := channel.Channel{Authority: GetAuthorization(token.(string), db)}

	channel, err := channelAuth.GetGraphqlChannelDetails(channelID,channelSlug)

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.Channel{}, err
	}

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

	return &conv_channel, nil
}

func ChannelEntryDetail(db *gorm.DB, ctx context.Context, channelEntryId, channelId, categoryId *int, slug, profileSlug *string) (*model.ChannelEntries, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	// memberid := c.GetInt("memberid")

	channelAuth := channel.Channel{Authority: GetAuthorization(token.(string), db)}

	channelEntry, err := channelAuth.GetGraphqlChannelEntriesDetails(channelEntryId, channelId, categoryId, PathUrl, SectionTypeId, MemberFieldTypeId, slug, profileSlug)

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.ChannelEntries{}, err
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

	authorMobnumber := channelEntry.AuthorDetail.MobileNo

	authorIsActive := channelEntry.AuthorDetail.IsActive

	authorProfileImage := channelEntry.AuthorDetail.ProfileImagePath

	authorDetails := model.Author{
		AuthorID:         channelEntry.AuthorDetail.AuthorID,
		FirstName:        channelEntry.AuthorDetail.FirstName,
		LastName:         channelEntry.AuthorDetail.LastName,
		Email:            channelEntry.AuthorDetail.Email,
		MobileNo:         &authorMobnumber,
		IsActive:         &authorIsActive,
		ProfileImagePath: &authorProfileImage,
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
		Author:           &channelEntry.Author,
		SortOrder:        &channelEntry.SortOrder,
		CreateTime:       &channelEntry.CreateTime,
		PublishedTime:    &channelEntry.PublishedTime,
		ReadingTime:      &channelEntry.ReadingTime,
		Tags:             &channelEntry.Tags,
		Excerpt:          &channelEntry.Excerpt,
		ImageAltTag:      &channelEntry.ImageAltTag,
	}

	return &conv_channelEntry, nil

}

func Memberclaimnow(db *gorm.DB, ctx context.Context, profileData model.ClaimData, entryId *int, profileSlug *string) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	// memberid := c.GetInt("memberid")

	verify_chan := make(chan error)

	channelAuth := channel.Channel{Authority: GetAuthorization(token, db)}

	var channelEntry channel.TblChannelEntries

	var err error

	if entryId != nil {

		channelEntry, err = channelAuth.GetGraphqlChannelEntriesDetails(entryId, nil, nil, PathUrl, SectionTypeId, MemberFieldTypeId, nil, nil)

	} else {

		channelEntry, err = channelAuth.GetGraphqlChannelEntriesDetails(nil, nil, nil, PathUrl, SectionTypeId, MemberFieldTypeId, nil, profileSlug)

	}

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
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

	authorMobnumber := channelEntry.AuthorDetail.MobileNo

	authorIsActive := channelEntry.AuthorDetail.IsActive

	authorProfileImage := channelEntry.AuthorDetail.ProfileImagePath

	authorDetails := model.Author{
		AuthorID:         channelEntry.AuthorDetail.AuthorID,
		FirstName:        channelEntry.AuthorDetail.FirstName,
		LastName:         channelEntry.AuthorDetail.LastName,
		Email:            channelEntry.AuthorDetail.Email,
		MobileNo:         &authorMobnumber,
		IsActive:         &authorIsActive,
		ProfileImagePath: &authorProfileImage,
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
		Author:           &channelEntry.Author,
		SortOrder:        &channelEntry.SortOrder,
		CreateTime:       &channelEntry.CreateTime,
		PublishedTime:    &channelEntry.PublishedTime,
		ReadingTime:      &channelEntry.ReadingTime,
		Tags:             &channelEntry.Tags,
		Excerpt:          &channelEntry.Excerpt,
		ImageAltTag:      &channelEntry.ImageAltTag,
	}

	data := map[string]interface{}{"claimData": profileData, "authorDetails": conv_channelEntry.AuthorDetails, "entry": conv_channelEntry, "additionalData": AdditionalData, "link": PathUrl + "member/updatemember?id=" + strconv.Itoa(*conv_channelEntry.MemberProfile.MemberID)}

	tmpl, err := template.ParseFiles("view/email/claim-template.html")

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	var template_buff bytes.Buffer

	err = tmpl.Execute(&template_buff, data)

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	mail_data := MailConfig{Email: conv_channelEntry.AuthorDetails.Email, MailUsername: os.Getenv("MAIL_USERNAME"), MailPassword: os.Getenv("MAIL_PASSWORD"), Subject: "My Claim Request for " + conv_channelEntry.Title}

	html_content := template_buff.String()

	go SendMail(mail_data, html_content, verify_chan)

	if <-verify_chan == nil {

		return true, nil

	} else {

		return false, nil
	}
}

func MemberProfileUpdate(db *gorm.DB, ctx context.Context, profiledata model.ProfileData, entryId *int, profileSlug *string) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	memberid := c.GetInt("memberid")

	if token == SpecialToken || token == "" || memberid == 0 {

		c.AbortWithStatus(http.StatusUnauthorized)

		return false, errors.New("login required")
	}

	var memberProfile model.MemberProfile

	query := db.Debug().Table("tbl_channel_entry_fields").Select("tbl_member_profiles.*").Joins("inner join tbl_fields on tbl_fields.id = tbl_channel_entry_fields.field_id").Joins("inner join tbl_members on tbl_members.id = tbl_channel_entry_fields.field_value::integer").Joins("inner join tbl_member_profiles on tbl_member_profiles.member_id = tbl_members.id").Where("tbl_fields.is_deleted = 0 and tbl_members.is_deleted = 0 and tbl_member_profiles.is_deleted = 0 and tbl_fields.field_type_id = ?", MemberFieldTypeId)

	if memberid > 0 {

		query = query.Where("tbl_member_profiles.member_id = ?", memberid)
	}

	if entryId != nil {

		query = query.Where("tbl_channel_entry_fields.channel_entry_id = ?", entryId)

	} else if profileSlug != nil {

		profileSubQuery := db.Table("tbl_channel_entries").Select("tbl_channel_entries.id").Joins("inner join tbl_channel_entry_fields on tbl_channel_entry_fields.channel_entry_id = tbl_channel_entries.id").Joins("inner join tbl_fields on tbl_fields.id = tbl_channel_entry_fields.field_id").Joins("inner join tbl_member_profiles on tbl_member_profiles.member_id = any(string_to_array(tbl_channel_entry_fields.field_value,',')::integer[])").Where("tbl_channel_entries.is_deleted = 0  and tbl_member_profiles.claim_status = 1 and tbl_channel_entries.status = 1 and tbl_fields.is_deleted = 0 and tbl_member_profiles.is_deleted = 0 and tbl_fields.field_type_id = ? and tbl_member_profiles.profile_slug = ?", MemberFieldTypeId, profileSlug).Limit(1)

		if memberid > 0 {

			profileSubQuery = profileSubQuery.Where("tbl_member_profiles.member_id = ?", memberid)
		}

		query = query.Where("tbl_channel_entry_fields.channel_entry_id = (?)", profileSubQuery)
	}

	if err := query.Find(&memberProfile).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	var jsonData map[string]interface{}

	err := json.Unmarshal([]byte(profiledata.MemberProfile), &jsonData)

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	currentTime, err := time.Parse("2006-01-02 15:04:05", time.Now().In(TimeZone).Format("2006-01-02 15:04:05"))

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	memberProfileDetails := model.MemberProfile{
		MemberDetails: profiledata.MemberProfile,
		Linkedin:      profiledata.Linkedin.Value(),
		Twitter:       profiledata.Twitter.Value(),
		Website:       profiledata.Website.Value(),
		ModifiedOn:    &currentTime,
	}

	if err := db.Debug().Table("tbl_member_profiles").Where("is_deleted = 0 and claim_status = 1 and member_id = ?", memberProfile.MemberID).UpdateColumns(map[string]interface{}{"member_details": memberProfileDetails.MemberDetails, "linkedin": memberProfileDetails.Linkedin, "twitter": memberProfileDetails.Twitter, "website": memberProfileDetails.Website, "modified_on": memberProfileDetails.ModifiedOn}).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	return true, nil
}

func VerifyProfileName(db *gorm.DB, ctx context.Context, profileName string) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	if profileName == "" {

		return false, errors.New("empty values not allowed")
	}

	var count int64

	if err := db.Debug().Table("tbl_member_profiles").Where("tbl_member_profiles.is_deleted = 0 and tbl_member_profiles.claim_status = 1 and tbl_member_profiles.profile_name = ?", profileName).Count(&count).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	if count > 0 {

		return false, errors.New("profile name already exists")
	}

	return true, nil
}

func UpdateChannelEntryViewCount(db *gorm.DB, ctx context.Context, entryId *int, slug *string) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	chanAuth := channel.Channel{Authority: GetAuthorizationWithoutToken(db)}

	isUpdated, err := chanAuth.UpdateChannelEntryViewCount(entryId, slug)

	if !isUpdated || err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	return true, nil
}
