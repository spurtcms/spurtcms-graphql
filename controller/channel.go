package controller

import (
	"context"
	"gqlserver/graph/model"
	"os"

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

	if channelID == nil{

		channelEntries, count, err = channelAuth.GetGraphqlAllChannelEntriesList(categoryId, limit, offset, pathUrl)

    }else{

		channelEntries, count, err = channelAuth.GetGraphqlChannelEntriesByChannelId(channelID, categoryId, limit, offset, pathUrl)
	}

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

		var sections []model.Section

		db.Table("tbl_group_fields").Select("tbl_fields.*,tbl_field_types.type_name").Joins("inner join tbl_fields on tbl_fields.id = tbl_group_fields.field_id").Joins("inner join tbl_field_types on tbl_field_types.id = tbl_fields.field_type_id").
			Where("tbl_fields.is_deleted = 0 and tbl_field_types.is_deleted = 0 and tbl_fields.field_type_id = ? and tbl_group_fields.channel_id = ?",SectionTypeId, entry.ChannelId).Find(&sections)

		var fields []model.Field

		db.Table("tbl_group_fields").Select("tbl_fields.*,tbl_field_types.type_name").Joins("inner join tbl_fields on tbl_fields.id = tbl_group_fields.field_id").Joins("inner join tbl_field_types on tbl_field_types.id = tbl_fields.field_type_id").
			Where("tbl_fields.is_deleted = 0 and tbl_field_types.is_deleted = 0 and tbl_fields.field_type_id != ? and tbl_group_fields.channel_id = ?",SectionTypeId,entry.ChannelId).Find(&fields)

		var final_fieldsList []model.Field

		for _, field := range fields {

			var fieldValue model.FieldValue

			db.Table("tbl_channel_entry_fields").Where("tbl_channel_entry_fields.field_id = ? and tbl_channel_entry_fields.channel_entry_id = ?", field.FieldID, entry.Id).First(&fieldValue)

			if fieldValue.ID != 0 {

				field.FieldValue = &fieldValue
			}

			var fieldOptions []model.FieldOptions

			db.Table("tbl_field_options").Where("tbl_field_options.is_deleted = 0 and tbl_field_options.field_id = ?", field.FieldID).Find(&fieldOptions)

			if len(fieldOptions) > 0 {

				field.FieldOptions = fieldOptions

			}

			final_fieldsList = append(final_fieldsList, field)
		}

		additionalFields := &model.AdditionalFields{Sections: sections, Fields: final_fieldsList}

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

	channelEntry, err := channelAuth.GetGraphqlChannelEntriesDetails(channelEntryId, channelId, categoryId, pathUrl)

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

	var sections []model.Section

	db.Table("tbl_group_fields").Select("tbl_fields.*,tbl_field_types.type_name").Joins("inner join tbl_fields on tbl_fields.id = tbl_group_fields.field_id").Joins("inner join tbl_field_types on tbl_field_types.id = tbl_fields.field_type_id").
		Where("tbl_fields.is_deleted = 0 and tbl_field_types.is_deleted = 0 and tbl_fields.field_type_id = ? and tbl_group_fields.channel_id = ?",SectionTypeId, channelEntry.ChannelId).Find(&sections)

	var fields []model.Field

	db.Table("tbl_group_fields").Select("tbl_fields.*,tbl_field_types.type_name").Joins("inner join tbl_fields on tbl_fields.id = tbl_group_fields.field_id").Joins("inner join tbl_field_types on tbl_field_types.id = tbl_fields.field_type_id").
		Where("tbl_fields.is_deleted = 0 and tbl_field_types.is_deleted = 0 and tbl_fields.field_type_id != ? and tbl_group_fields.channel_id = ?",SectionTypeId, channelEntry.ChannelId).Find(&fields)

	var final_fieldsList []model.Field

	for _, field := range fields {

		var fieldValue model.FieldValue

		db.Table("tbl_channel_entry_fields").Where("tbl_channel_entry_fields.field_id = ? and tbl_channel_entry_fields.channel_entry_id = ?", field.FieldID, channelEntry.Id).First(&fieldValue)

		if fieldValue.ID != 0 {

			field.FieldValue = &fieldValue
		}

		var fieldOptions []model.FieldOptions

		db.Table("tbl_field_options").Where("tbl_field_options.is_deleted = 0 and tbl_field_options.field_id = ?", field.FieldID).Find(&fieldOptions)

		if len(fieldOptions) > 0 {

			field.FieldOptions = fieldOptions

		}

		final_fieldsList = append(final_fieldsList, field)
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
	}

	return conv_channelEntry,nil

}
