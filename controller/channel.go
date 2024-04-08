package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"os"
	"spurtcms-graphql/graph/model"
	"strconv"
	"time"

	"log"
	// "time"

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

	memberid := c.GetInt("memberid")

	channelAuth := channel.Channel{Authority: GetAuthorization(token.(string), db)}

	var channelEntries []channel.TblChannelEntries

	var count int64

	var err error

	channelEntries, count, err = channelAuth.GetGraphqlAllChannelEntriesList(channelID, categoryId, limit, offset, SectionTypeId, MemberFieldTypeId, PathUrl, title, categoryChildId, categorySlug, categoryChildSlug)

	if err != nil {

		return &model.ChannelEntriesDetails{}, err
	}

	conv_channelEntries := make([]model.ChannelEntries, len(channelEntries))

	type data interface{} 

	for index, entry := range channelEntries {

		conv_categories := make([][]model.Category, len(entry.Categories))

		for cat_index, categories := range entry.Categories {

			conv_categoryz := make([]model.Category, len(categories))

			for i, category := range categories {

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

				conv_categoryz[i] = conv_category

			}

			conv_categories[cat_index] = conv_categoryz
		}

		conv_channelEntries[index].Categories = conv_categories

		authorDetails := model.Author{
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

		conv_channelEntries[index].AuthorDetails = authorDetails

		conv_sections := make([]model.Section, len(entry.Sections))

		for section_index, section := range entry.Sections {

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

			conv_sections[section_index] = conv_section

		}

		conv_fields := make([]model.Field, len(entry.Fields))

		for field_index, field := range entry.Fields {

			conv_field_value := model.FieldValue{
				ID:         field.FieldValue.FieldId,
				FieldValue: field.FieldValue.FieldValue,
				CreatedOn:  field.FieldValue.CreatedOn,
				CreatedBy:  field.FieldValue.CreatedBy,
				ModifiedOn: &field.FieldValue.ModifiedOn,
				ModifiedBy: &field.FieldValue.ModifiedBy,
			}

			conv_fieldOptions := make([]model.FieldOptions, len(field.FieldOptions))

			for option_index, field_option := range field.FieldOptions {

				conv_fieldOption := model.FieldOptions{
					ID:          field_option.Id,
					OptionName:  field_option.OptionName,
					OptionValue: field_option.OptionValue,
					CreatedOn:   field_option.CreatedOn,
					CreatedBy:   field_option.CreatedBy,
					ModifiedOn:  &field_option.ModifiedOn,
					ModifiedBy:  &field_option.ModifiedBy,
				}

				conv_fieldOptions[option_index] = conv_fieldOption
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

			conv_fields[field_index] = conv_field

		}

		conv_channelEntries[index].Fields = conv_fields

		additionalFields := model.AdditionalFields{Sections: conv_sections, Fields: conv_fields}

		conv_channelEntries[index].AdditionalFields = &additionalFields

		MemberProfile := model.MemberProfile{
			ID:              entry.MemberProfile.Id,
			MemberID:        entry.MemberProfile.MemberId,
			ProfileName:     entry.MemberProfile.ProfileName,
			ProfileSlug:     entry.MemberProfile.ProfileSlug,
			ProfilePage:     entry.MemberProfile.ProfilePage,
			MemberDetails:   entry.MemberProfile.MemberDetails,
			CompanyName:     &entry.MemberProfile.CompanyName,
			CompanyLocation: &entry.MemberProfile.CompanyLocation,
			CompanyLogo:     &entry.MemberProfile.CompanyLogo,
			About:           &entry.MemberProfile.About,
			SeoTitle:        &entry.MemberProfile.SeoTitle,
			SeoDescription:  &entry.MemberProfile.SeoDescription,
			SeoKeyword:      &entry.MemberProfile.SeoKeyword,
			CreatedBy:       &entry.MemberProfile.CreatedBy,
			CreatedOn:       &entry.MemberProfile.CreatedOn,
			ModifiedOn:      &entry.MemberProfile.ModifiedOn,
			ModifiedBy:      &entry.MemberProfile.ModifiedBy,
			Linkedin:        &entry.MemberProfile.Linkedin,
			Twitter:         &entry.MemberProfile.Twitter,
			Website:         &entry.MemberProfile.Website,
			ClaimStatus:     &entry.MemberProfile.ClaimStatus,
		}

		conv_channelEntries[index].MemberProfile = MemberProfile

		var claimStatus bool

		if entry.MemberProfile.ClaimStatus == 1 && memberid == entry.MemberProfile.MemberId {

			claimStatus = true

		} else {

			claimStatus = false
		}

		conv_channelEntries[index].ClaimStatus = claimStatus

		conv_channelEntries[index].Author = &entry.Author

		conv_channelEntries[index].CategoriesID = entry.CategoriesId

		conv_channelEntries[index].ChannelID = entry.ChannelId

		conv_channelEntries[index].CoverImage = entry.CoverImage

		conv_channelEntries[index].CreateDate = &entry.CreateDate

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

		conv_channelEntries[index].ModifiedBy = &entry.ModifiedBy

		conv_channelEntries[index].ModifiedOn = &entry.ModifiedOn

		conv_channelEntries[index].PublishedTime = &entry.PublishedTime

		conv_channelEntries[index].ReadingTime = &entry.ReadingTime

		conv_channelEntries[index].RelatedArticles = entry.RelatedArticles

		conv_channelEntries[index].Slug = entry.Slug

		conv_channelEntries[index].SortOrder = &entry.SortOrder

		conv_channelEntries[index].Status = entry.Status

		conv_channelEntries[index].Tags = &entry.Tags

		conv_channelEntries[index].ThumbnailImage = entry.ThumbnailImage

		conv_channelEntries[index].Title = entry.Title

		conv_channelEntries[index].UserID = entry.UserId

		conv_channelEntries[index].ViewCount = entry.ViewCount

	}

	for _,val := range conv_channelEntries{

		log.Println("final",val.ModifiedOn)
	}

	channelEntryDetails := model.ChannelEntriesDetails{ChannelEntriesList: conv_channelEntries, Count: int(count)}

	return &channelEntryDetails, nil

}

func ChannelDetail(db *gorm.DB, ctx context.Context, channelID int) (*model.Channel, error) {

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

		return &model.Channel{}, err
	}

	return &conv_channel, nil
}

func ChannelEntryDetail(db *gorm.DB, ctx context.Context, channelEntryId, channelId, categoryId *int, slug *string) (*model.ChannelEntries, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	memberid := c.GetInt("memberid")

	channelAuth := channel.Channel{Authority: GetAuthorization(token.(string), db)}

	channelEntry, err := channelAuth.GetGraphqlChannelEntriesDetails(channelEntryId, channelId, categoryId, PathUrl, SectionTypeId, MemberFieldTypeId, slug)

	if err != nil {

		return &model.ChannelEntries{}, err
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

	additionalFields := model.AdditionalFields{Sections: conv_sections, Fields: conv_fields}

	MemberProfile := model.MemberProfile{
		ID:              channelEntry.MemberProfile.Id,
		MemberID:        channelEntry.MemberProfile.MemberId,
		ProfileName:     channelEntry.MemberProfile.ProfileName,
		ProfileSlug:     channelEntry.MemberProfile.ProfileSlug,
		ProfilePage:     channelEntry.MemberProfile.ProfilePage,
		MemberDetails:   channelEntry.MemberProfile.MemberDetails,
		CompanyName:     &channelEntry.MemberProfile.CompanyName,
		CompanyLocation: &channelEntry.MemberProfile.CompanyLocation,
		CompanyLogo:     &channelEntry.MemberProfile.CompanyLogo,
		About:           &channelEntry.MemberProfile.About,
		SeoTitle:        &channelEntry.MemberProfile.SeoTitle,
		SeoDescription:  &channelEntry.MemberProfile.SeoDescription,
		SeoKeyword:      &channelEntry.MemberProfile.SeoKeyword,
		CreatedBy:       &channelEntry.MemberProfile.CreatedBy,
		CreatedOn:       &channelEntry.MemberProfile.CreatedOn,
		ModifiedOn:      &channelEntry.MemberProfile.ModifiedOn,
		ModifiedBy:      &channelEntry.MemberProfile.ModifiedBy,
		Linkedin:        &channelEntry.MemberProfile.Linkedin,
		Twitter:         &channelEntry.MemberProfile.Twitter,
		Website:         &channelEntry.MemberProfile.Website,
		ClaimStatus:     &channelEntry.MemberProfile.ClaimStatus,
	}

	var claimStatus bool

	if *MemberProfile.ClaimStatus==1 && memberid == MemberProfile.MemberID {

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
		AdditionalFields: &additionalFields,
		MemberProfile:    MemberProfile,
		AuthorDetails:    authorDetails,
		FeaturedEntry:    channelEntry.Feature,
		ViewCount:        channelEntry.ViewCount,
		ClaimStatus:      claimStatus,
		Author:           &channelEntry.Author,
		SortOrder:        &channelEntry.SortOrder,
		CreateDate:       &channelEntry.CreateDate,
		PublishedTime:    &channelEntry.PublishedTime,
		ReadingTime:      &channelEntry.ReadingTime,
		Tags:             &channelEntry.Tags,
		Excerpt:          &channelEntry.Excerpt,
	}

	return &conv_channelEntry, nil

}

func Memberclaimnow(db *gorm.DB, ctx context.Context, profileData model.ClaimData, entryId int) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	memberid := c.GetInt("memberid")

	verify_chan := make(chan bool)

	channelAuth := channel.Channel{Authority: GetAuthorization(token, db)}

	channelEntry, err := channelAuth.GetGraphqlChannelEntriesDetails(&entryId, nil, nil, PathUrl, SectionTypeId, MemberFieldTypeId, nil)

	if err != nil {

		return false, err
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

	additionalFields := model.AdditionalFields{Sections: conv_sections, Fields: conv_fields}

	MemberProfile := model.MemberProfile{
		ID:              channelEntry.MemberProfile.Id,
		MemberID:        channelEntry.MemberProfile.MemberId,
		ProfileName:     channelEntry.MemberProfile.ProfileName,
		ProfileSlug:     channelEntry.MemberProfile.ProfileSlug,
		ProfilePage:     channelEntry.MemberProfile.ProfilePage,
		MemberDetails:   channelEntry.MemberProfile.MemberDetails,
		CompanyName:     &channelEntry.MemberProfile.CompanyName,
		CompanyLocation: &channelEntry.MemberProfile.CompanyLocation,
		CompanyLogo:     &channelEntry.MemberProfile.CompanyLogo,
		About:           &channelEntry.MemberProfile.About,
		SeoTitle:        &channelEntry.MemberProfile.SeoTitle,
		SeoDescription:  &channelEntry.MemberProfile.SeoDescription,
		SeoKeyword:      &channelEntry.MemberProfile.SeoKeyword,
		CreatedBy:       &channelEntry.MemberProfile.CreatedBy,
		CreatedOn:       &channelEntry.MemberProfile.CreatedOn,
		ModifiedOn:      &channelEntry.MemberProfile.ModifiedOn,
		ModifiedBy:      &channelEntry.MemberProfile.ModifiedBy,
		Linkedin:        &channelEntry.MemberProfile.Linkedin,
		Twitter:         &channelEntry.MemberProfile.Twitter,
		Website:         &channelEntry.MemberProfile.Website,
		ClaimStatus:     &channelEntry.MemberProfile.ClaimStatus,
	}

	var claimStatus bool

	if *MemberProfile.ClaimStatus == 1 && memberid == MemberProfile.MemberID {

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
		AdditionalFields: &additionalFields,
		MemberProfile:    MemberProfile,
		AuthorDetails:    authorDetails,
		FeaturedEntry:    channelEntry.Feature,
		ViewCount:        channelEntry.ViewCount,
		ClaimStatus:      claimStatus,
		Author:           &channelEntry.Author,
		SortOrder:        &channelEntry.SortOrder,
		CreateDate:       &channelEntry.CreateDate,
		PublishedTime:    &channelEntry.PublishedTime,
		ReadingTime:      &channelEntry.ReadingTime,
		Tags:             &channelEntry.Tags,
		Excerpt:          &channelEntry.Excerpt,
	}

	data := map[string]interface{}{"claimData": profileData, "authorDetails": conv_channelEntry.AuthorDetails, "entry": conv_channelEntry, "additionalData": AdditionalData, "link": PathUrl + "member/updatemember?id=" + strconv.Itoa(conv_channelEntry.MemberProfile.MemberID)}

	log.Println("maildata", data["link"], conv_channelEntry.ClaimStatus, data["additionalData"])

	tmpl, _ := template.ParseFiles("view/email/claim-template.html")

	var template_buff bytes.Buffer

	err = tmpl.Execute(&template_buff, data)

	if err != nil {

		return false, err
	}

	mail_data := MailConfig{Email: conv_channelEntry.AuthorDetails.Email, MailUsername: os.Getenv("MAIL_USERNAME"), MailPassword: os.Getenv("MAIL_PASSWORD"), Subject: "My Claim Request for " + conv_channelEntry.Title}

	html_content := template_buff.String()

	go SendMail(mail_data, html_content, verify_chan)

	if <-verify_chan {

		return true, nil

	} else {

		return false, nil
	}
}

func MemberProfileUpdate(db *gorm.DB, ctx context.Context, profiledata model.ProfileData, entryId int) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	memberid := c.GetInt("memberid")

	if token == SpecialToken || token == "" || memberid == 0 {

		return false, errors.New("login required")
	}

	var memberProfile model.MemberProfile

	if err := db.Debug().Table("tbl_channel_entry_fields").Select("tbl_member_profiles.*").Joins("inner join tbl_fields on tbl_fields.id = tbl_channel_entry_fields.field_id").Joins("inner join tbl_members on tbl_members.id = tbl_channel_entry_fields.field_value::integer").Joins("inner join tbl_member_profiles on tbl_member_profiles.member_id = tbl_members.id").Where("tbl_fields.is_deleted = 0 and tbl_members.is_deleted = 0 and tbl_member_profiles.is_deleted = 0 and tbl_member_profiles.claim_status = 1 and tbl_fields.field_type_id = ? and tbl_channel_entry_fields.channel_entry_id = ?", MemberFieldTypeId, entryId).Find(&memberProfile).Error; err != nil {

		return false, err
	}

	var jsonData map[string]interface{}

	err := json.Unmarshal([]byte(profiledata.MemberProfile), &jsonData)

	if err != nil {

		return false, err
	}

	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().In(TimeZone).Format("2006-01-02 15:04:05"))

	memberProfileDetails := model.MemberProfile{
		MemberDetails: profiledata.MemberProfile,
		Linkedin:      profiledata.Linkedin.Value(),
		Twitter:       profiledata.Twitter.Value(),
		Website:       profiledata.Website.Value(),
		ModifiedOn:    &currentTime,
	}

	if memberid != memberProfile.MemberID {

		return false, errors.New("authorized member id mismatched in member profile")
	}

	if err := db.Debug().Table("tbl_member_profiles").Where("is_deleted = 0 and claim_status = 1 and member_id = ?", memberProfile.MemberID).UpdateColumns(map[string]interface{}{"member_details": memberProfileDetails.MemberDetails, "linkedin": memberProfileDetails.Linkedin, "twitter": memberProfileDetails.Twitter, "website": memberProfileDetails.Website, "modified_on": memberProfileDetails.ModifiedOn}).Error; err != nil {

		return false, err
	}

	return true, nil
}

func VerifyProfileName(db *gorm.DB, ctx context.Context, profileName string) (bool, error) {

	if profileName == "" {

		return false, errors.New("empty values not allowed")
	}

	var count int64

	if err := db.Debug().Table("tbl_member_profiles").Where("tbl_member_profiles.is_deleted = 0 and tbl_member_profiles.claim_status = 1 and tbl_member_profiles.profile_name = ?", profileName).Count(&count).Error; err != nil {

		return false, err
	}

	if count > 0 {

		return false, errors.New("profile name already exists")
	}

	return true, nil
}
