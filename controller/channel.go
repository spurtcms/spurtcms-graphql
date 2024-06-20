package controller

import (
	"context"
	"net/http"
	"spurtcms-graphql/graph/model"

	"github.com/gin-gonic/gin"
	channel "github.com/spurtcms/pkgcontent/channels"

	// "github.com/spurtcms/pkgcore/member"
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
func ChannelEntriesList(db *gorm.DB, ctx context.Context, channelID, categoryId *int, limit, offset int, title *string, categoryChildId *int, categorySlug, categoryChildSlug *string, requireData *model.RequireData) (*model.ChannelEntriesDetails, error) {

	channelAuth := channel.Channel{Authority: GetAuthorizationWithoutToken(db)}

	var channelEntries []channel.TblChannelEntries

	var count int64

	var err error

	var memberprofileflg, authorflg, categoriesflg, fieldsflg bool

	if requireData != nil {

		if requireData.MemberProfile.IsSet() {

			memberprofileflg = *requireData.MemberProfile.Value()
		}

		if requireData.AuthorDetails.IsSet() {

			authorflg = *requireData.AuthorDetails.Value()
		}

		if requireData.Categories.IsSet() {

			categoriesflg = *requireData.Categories.Value()
		}

		if requireData.AdditionalFields.IsSet() {

			fieldsflg = *requireData.AdditionalFields.Value()
		}
	}

	channelEntries, count, err = channelAuth.GetGraphqlAllChannelEntriesList(channelID, categoryId, limit, offset, SectionTypeId, MemberFieldTypeId, title, categoryChildId, categorySlug, categoryChildSlug, authorflg, memberprofileflg, categoriesflg, fieldsflg)

	if err != nil {

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

		memberProfileComLogo := entry.MemberProfile.CompanyLogo
		memberProfileId := entry.MemberProfile.Id
		memberProfileMemId := entry.MemberProfile.MemberId
		memberProfileName := entry.MemberProfile.ProfileName
		memberProfileSlug := entry.MemberProfile.ProfileSlug
		memberProfilePage := entry.MemberProfile.ProfilePage
		memberProfileMemDetails := entry.MemberProfile.MemberDetails
		memberProfileComName := entry.MemberProfile.CompanyName
		memberProfileComLocation := entry.MemberProfile.CompanyLocation
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

	channel, err := channelAuth.GetGraphqlChannelDetails(channelID, channelSlug)

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

func UpdateChannelEntryViewCount(db *gorm.DB, ctx context.Context, entryId *int, slug *string) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	chanAuth := channel.Channel{Authority: GetAuthorizationWithoutToken(db)}

	isUpdated, err := chanAuth.UpdateChannelEntryViewCount(entryId, slug)

	if !isUpdated || err != nil {

		ErrorLog.Printf("channel entry update view count error: %s", err)

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	return true, nil
}
