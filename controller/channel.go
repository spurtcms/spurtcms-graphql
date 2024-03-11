package controller

import (
	"context"
	"gqlserver/graph/model"
	// "log"
	"os"
	"strconv"
	"strings"

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
func ChannelEntriesList(db *gorm.DB, ctx context.Context, channelID, categoryId *int, limit, offset int) (model.ChannelEntriesDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	channelAuth := channel.Channel{Authority: GetAuthorization(token.(string), db)}

	var pathUrl string

	if os.Getenv("DOMAIN_URL") != "" {

		pathUrl = os.Getenv("DOMAIN_URL")

	} else {

		pathUrl = os.Getenv("LOCAL_URL")
	}

	var channelEntries []channel.TblChannelEntries

	var count int64

	var err error

	channelEntries, count, err = channelAuth.GetGraphqlAllChannelEntriesList(channelID,categoryId, limit, offset,SectionTypeId, pathUrl)

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
			AuthorID: entry.AuthorDetail.AuthorID,
			FirstName: entry.AuthorDetail.FirstName,
			LastName: entry.AuthorDetail.LastName,
			Email: entry.AuthorDetail.Email,
			MobileNo: entry.AuthorDetail.MobileNo,
			IsActive: entry.AuthorDetail.IsActive,
			ProfileImage: entry.AuthorDetail.ProfileImage,
			ProfileImagePath: entry.AuthorDetail.ProfileImagePath,
			CreatedOn: entry.AuthorDetail.CreatedOn,
			CreatedBy: entry.AuthorDetail.CreatedBy,
		}

		var conv_sections []model.Section

	    for _,section :=  range entry.Sections{
			
		  conv_section:= model.Section{
			SectionID: &section.Id,
			SectionName: section.FieldName,
			SectionTypeID: section.FieldTypeId,
			CreatedOn: section.CreatedOn,
			CreatedBy: section.CreatedBy,
			ModifiedOn: &section.ModifiedOn,
			ModifiedBy: &section.ModifiedBy,
			OrderIndex: section.OrderIndex,
		  }

		  conv_sections = append(conv_sections, conv_section)
		}

		var conv_fields []model.Field

		for _,field := range entry.Fields{

			conv_field_value := model.FieldValue{
				ID: field.FieldValue.FieldId,
				FieldValue: field.FieldValue.FieldValue,
				CreatedOn: field.FieldValue.CreatedOn,
				CreatedBy: field.FieldValue.CreatedBy,
				ModifiedOn: &field.FieldValue.ModifiedOn,
				ModifiedBy: &field.FieldValue.ModifiedBy,
			}

			conv_field := model.Field{
				FieldID: field.Id,
				FieldName: field.FieldName,
				FieldTypeID: field.FieldTypeId,
				MandatoryField: field.MandatoryField,
				OptionExist: field.OptionExist,
				CreatedOn: field.CreatedOn,
				CreatedBy: field.CreatedBy,
				ModifiedOn: &field.ModifiedOn,
				ModifiedBy: &field.ModifiedBy,
				FieldDesc: field.FieldDesc,
				OrderIndex: field.OrderIndex,
				ImagePath: field.ImagePath,
				DatetimeFormat: &field.DatetimeFormat,
				TimeFormat: &field.TimeFormat,
				SectionParentID: &field.SectionParentId,
				CharacterAllowed: &field.CharacterAllowed,
				FieldTypeName: field.FieldTypeName,
				FieldValue: &conv_field_value,
			}

			conv_fields = append(conv_fields, conv_field)
		}

		additionalFields := &model.AdditionalFields{Sections: conv_sections, Fields: conv_fields}

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
			// MemberProfile: memberProfiles,
			AuthorDetails: authorDetails,

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

func ChannelEntryDetail(db *gorm.DB, ctx context.Context, channelEntryId int, channelId, categoryId *int) (model.ChannelEntries, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	channelAuth := channel.Channel{Authority: GetAuthorization(token.(string), db)}

	var pathUrl string

	if os.Getenv("DOMAIN_URL") != "" {

		pathUrl = os.Getenv("DOMAIN_URL")

	} else {

		pathUrl = os.Getenv("LOCAL_URL")
	}

	channelEntry, err := channelAuth.GetGraphqlChannelEntriesDetails(channelEntryId, channelId, categoryId, pathUrl,SectionTypeId)

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

	var memberids string

	var sections []model.Section

	db.Table("tbl_group_fields").Select("tbl_fields.*,tbl_field_types.type_name").Joins("inner join tbl_fields on tbl_fields.id = tbl_group_fields.field_id").Joins("inner join tbl_field_types on tbl_field_types.id = tbl_fields.field_type_id").
		Where("tbl_fields.is_deleted = 0 and tbl_field_types.is_deleted = 0 and tbl_fields.field_type_id = ? and tbl_group_fields.channel_id = ?", SectionTypeId, channelEntry.ChannelId).Find(&sections)

	var fields []model.Field

	db.Table("tbl_group_fields").Select("tbl_fields.*,tbl_field_types.type_name").Joins("inner join tbl_fields on tbl_fields.id = tbl_group_fields.field_id").Joins("inner join tbl_field_types on tbl_field_types.id = tbl_fields.field_type_id").
		Where("tbl_fields.is_deleted = 0 and tbl_field_types.is_deleted = 0 and tbl_fields.field_type_id != ? and tbl_group_fields.channel_id = ?", SectionTypeId, channelEntry.ChannelId).Find(&fields)

	var final_fieldsList []model.Field

	for _, field := range fields {

		var fieldValue model.FieldValue

		db.Table("tbl_channel_entry_fields").Where("tbl_channel_entry_fields.field_id = ? and tbl_channel_entry_fields.channel_entry_id = ?", field.FieldID, channelEntry.Id).First(&fieldValue)

		if fieldValue.ID != 0 {

			field.FieldValue = &fieldValue

			if field.FieldTypeID == 14 {

				memberids = fieldValue.FieldValue
			}
		}

		var fieldOptions []model.FieldOptions

		db.Table("tbl_field_options").Where("tbl_field_options.is_deleted = 0 and tbl_field_options.field_id = ?", field.FieldID).Find(&fieldOptions)

		if len(fieldOptions) > 0 {

			field.FieldOptions = fieldOptions

		}

		final_fieldsList = append(final_fieldsList, field)
	}

	var memberProfiles []model.MemberProfile

	MemIds := strings.Split(memberids, ",")

	for _, memberid := range MemIds {

		var memberProfile model.MemberProfile

		conv_memid,_ := strconv.Atoi(memberid)

		db.Debug().Table("tbl_member_profiles").Select("tbl_member_profiles.*").Joins("inner join tbl_members on tbl_members.id = tbl_member_profiles.member_id").Where("tbl_members.is_deleted = 0 and tbl_members.id = ?", conv_memid).First(&memberProfile)

		memberProfiles = append(memberProfiles, memberProfile)
	}

	additionalFields := &model.AdditionalFields{Sections: sections, Fields: final_fieldsList}

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
		MemberProfile: memberProfiles,
	}

	return conv_channelEntry, nil

}
