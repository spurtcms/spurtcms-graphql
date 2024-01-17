package controller

import (
	"context"
	"gqlserver/graph/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Channellist(db *gorm.DB,ctx context.Context,limit,offset int)(model.ChannelDetails,error){

	c,_ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	var channellist []model.TblChannel

	var count int64

	listerr := db.Debug().Table("tbl_channels").Select("distinct on (tbl_channels.id) tbl_channels.*").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
	Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
	Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
	Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?",memberid).
	Order("tbl_channels.id desc").Limit(limit).Offset(offset).Find(&channellist).Error

	counterr := db.Debug().Table("tbl_channels").Distinct("tbl_channels.id").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
	Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
	Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
	Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?",memberid).
	Count(&count).Error

	if listerr == nil && counterr == nil{

		channelDetails := model.ChannelDetails{Channellist: channellist,Count: int(count)}

		return channelDetails,nil

	}else{

		if listerr!=nil{

			return model.ChannelDetails{},listerr

		}else{

			return model.ChannelDetails{},counterr

		}
	}
}

func ChannelEntriesList(db *gorm.DB,ctx context.Context, channelID *int, channelEntryID *int, limit ,offset int) (model.ChannelEntryDetails, error) {

	c,_ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if channelEntryID!=nil{

		var channelEntry model.TblChannelEntries

		entryerr := db.Debug().Table("tbl_channel_Entries").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
		Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
		Where("tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and tbl_channel_entries.channel_id = ? and tbl_channel_entries.id = ?",memberid,channelID,channelEntryID).
		First(&channelEntry).Error

		var categories []model.TblCategory

		raw := `WITH RECURSIVE cat_tree AS (
			SELECT * FROM tbl_categories WHERE id = ?
			UNION ALL
			SELECT cat.*
			FROM tbl_categories AS cat
			JOIN cat_tree ON cat.parent_id = cat_tree.id )`

		categoryerr := db.Debug().Raw(``+raw+`SELECT * FROM cat_tree where is_deleted = 0 order by id desc `, channelEntry.CategoriesID).Find(&categories).Error

		if entryerr==nil && categoryerr==nil{

			channelEntryDetails := model.ChannelEntryDetails{ChannelEntryList: &model.ChannelEntries{},ChannelEntry: &channelEntry}

			return channelEntryDetails,nil

		}else{

			if entryerr!=nil{

				return model.ChannelEntryDetails{},entryerr

			}else{

				return model.ChannelEntryDetails{},categoryerr

			}
		}
		

	}else{

		var channelEntries []model.TblChannelEntries

	    var count int64

		entrieserr := db.Debug().Table("tbl_channel_entries").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
		Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
		Where("tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and tbl_channel_entries.channel_id = ?",memberid,channelID).
		Find(&channelEntries).Error

		counterr := db.Debug().Table("tbl_channel_entries").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
		Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
		Where("tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and tbl_channel_entries.channel_id = ?",memberid,channelID).
		Count(&count).Error

		var final_entries_list []model.TblChannelEntries

		for _,entry := range channelEntries{

			var categories []model.TblCategory

			raw := `WITH RECURSIVE cat_tree AS (
				SELECT * FROM tbl_categories WHERE id = ?
				UNION ALL
				SELECT cat.*
				FROM tbl_categories AS cat
				JOIN cat_tree ON cat.parent_id = cat_tree.id )`
	
			categoryerr := db.Debug().Raw(``+raw+`SELECT * FROM cat_tree where is_deleted = 0 order by id desc `, entry.CategoriesID).Scan(&categories).Error

			if categoryerr!=nil{

				entry.Categories = []model.TblCategory{}
			}
			
			final_entries_list = append(final_entries_list, entry)

		}

		if entrieserr==nil && counterr==nil{

			channelEntrieslist := &model.ChannelEntries{ChannelEntryList: final_entries_list,Count: int(count)}

			channelEntriesDetails := model.ChannelEntryDetails{ChannelEntryList: channelEntrieslist,ChannelEntry: &model.TblChannelEntries{}}

			return  channelEntriesDetails,nil

		}else{

			if entrieserr!=nil{

				return model.ChannelEntryDetails{},entrieserr

			}else{

				return model.ChannelEntryDetails{},counterr
			}

		}

	}

}


