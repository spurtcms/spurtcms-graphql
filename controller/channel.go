package controller

import (
	"context"
	"errors"
	"gqlserver/graph/model"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Channellist(db *gorm.DB,ctx context.Context,limit,offset int)(model.ChannelDetails,error){

	c,_ := ctx.Value(ContextKey).(*gin.Context)

	token ,_ := c.Get("token")

	var channellist []model.TblChannel

	var count int64

	if token == SpecialToken{

	   listerr := db.Table("tbl_channels").Select("distinct on (tbl_channels.id) tbl_channels.*").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
	   Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
	   Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
	   Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0").
	   Order("tbl_channels.id desc").Limit(limit).Offset(offset).Find(&channellist).Error

	   counterr := db.Table("tbl_channels").Distinct("tbl_channels.id").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
	   Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
	   Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
	   Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0").
	   Count(&count).Error

	   if listerr == nil && counterr == nil{

		  channelDetails := model.ChannelDetails{Channellist: channellist,Count: int(count)}

		  return channelDetails,nil

	    }

	   var final_error error

	   if listerr!=nil{

	     final_error = listerr

	   }else{

		 final_error = counterr

	   }

	   return model.ChannelDetails{},final_error

	}else{

		memberid := c.GetInt("memberid")

	    listerr := db.Table("tbl_channels").Select("distinct on (tbl_channels.id) tbl_channels.*").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
	    Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
	    Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
	    Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?",memberid).
	    Order("tbl_channels.id desc").Limit(limit).Offset(offset).Find(&channellist).Error

	    counterr := db.Table("tbl_channels").Distinct("tbl_channels.id").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
	    Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
	    Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
	    Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?",memberid).
	    Count(&count).Error

	    if listerr == nil && counterr == nil{

		  channelDetails := model.ChannelDetails{Channellist: channellist,Count: int(count)}

		  return channelDetails,nil

	    }

	    var final_error error

	    if listerr!=nil{

		  final_error = listerr

	    }else{

		  final_error = counterr

	    }

	   return model.ChannelDetails{},final_error

	}

}

func ChannelEntriesList(db *gorm.DB,ctx context.Context, channelID *int, channelEntryID *int, limit ,offset *int) (model.ChannelEntryDetails, error) {

	c,_ := ctx.Value(ContextKey).(*gin.Context)

	token ,_ := c.Get("token")

	if token == SpecialToken{

		if channelEntryID!=nil || (channelID!=nil && channelEntryID!=nil) && limit==nil && offset==nil {
	
			var channelEntry model.TblChannelEntries
	
			var entryerr error
	
			if channelID!=nil{
	
			   entryerr = db.Table("tbl_channel_entries").Where("tbl_channel_entries.status = 1 and tbl_channel_entries.channel_id = ? and tbl_channel_entries.id = ?",channelID,channelEntryID).
			   First(&channelEntry).Error
	
			}else{
	
			  entryerr = db.Table("tbl_channel_entries").Where("tbl_channel_entries.status = 1 and tbl_channel_entries.id = ?",channelEntryID).
			  First(&channelEntry).Error

			}
	
			if entryerr==nil{
	
				splittedArr := strings.Split(channelEntry.CategoriesID, ",")
	
				var parentCatId int
	
				var indivCategories [][]model.TblCategory
	
				for _, catId := range splittedArr{
	
					var indivCategory []model.TblCategory
	
					conv_id,_ := strconv.Atoi(catId)
	
					var category model.TblCategory
	
					caterr := db.Table("tbl_categories").Where("is_deleted = 0 and id = ?",conv_id).First(&category).Error
	
					if caterr==nil{
	
						indivCategory = append(indivCategory, category)
	
						parentCatId = category.ParentID
	
						if parentCatId!=0{
	
							var count int 
	
							LOOP:
							  
							   for{
	
								  count = count + 1 //count increment used to check how many times the loop gets executed
	
								  var parentCategory model.TblCategory
	
								   caterr = db.Table("tbl_categories").Where("is_deleted = 0 and id = ?",parentCatId).First(&parentCategory).Error
	
								   if caterr==nil{
									  
									   indivCategory = append(indivCategory, parentCategory)
	
									   parentCatId = parentCategory.ParentID
	
									   if parentCatId!=0{ //mannuall condition to break the loop in overlooping situations
										  
										   goto LOOP
	
									   }else if count>49{
	
										break  //use to break the loop if infinite loop doesn't break ,So forcing the loop to break at overlooping conditions
	
									   }else{
	
										  break
									   }
									   
								   }else{
	
									   indivCategory = append(indivCategory, model.TblCategory{})
	
									   break
	
								   }
							   }
						}
	
					}else{
	
						indivCategory = append(indivCategory, model.TblCategory{})
					}
	
					indivCategories = append(indivCategories, indivCategory)
	
				}
	
				channelEntry.Categories = indivCategories
	
				channelEntryDetails := model.ChannelEntryDetails{ChannelEntryList: &model.ChannelEntries{},ChannelEntry: &channelEntry}
	
				return channelEntryDetails,nil
	
			}
	
			return model.ChannelEntryDetails{},entryerr	
	
		} else if channelEntryID==nil && channelID!=nil && limit!=nil && offset!=nil{
	
			var channelEntries []model.TblChannelEntries
	
			var count int64
	
			entrieserr := db.Table("tbl_channel_entries").Where("tbl_channel_entries.status = 1 and tbl_channel_entries.channel_id = ?",channelID).Limit(*limit).Offset(*offset).Order("tbl_channel_entries.id desc").
			Find(&channelEntries).Error
	
			counterr := db.Table("tbl_channel_entries").Where("tbl_channel_entries.status = 1 and tbl_channel_entries.channel_id = ?",channelID).
			Count(&count).Error
	
			if entrieserr==nil && counterr==nil{
	
				var final_entries_list []model.TblChannelEntries
	
				for _,entry := range channelEntries{
	
					var indivCategories [][]model.TblCategory
	
					splittedArr := strings.Split(entry.CategoriesID, ",")
		
					var parentCatId int
		
					for _, catId := range splittedArr{
		
						var indivCategory []model.TblCategory
		
						conv_id,_ := strconv.Atoi(catId)
		
						var category model.TblCategory
		
						caterr := db.Table("tbl_categories").Where("is_deleted = 0 and id = ?",conv_id).First(&category).Error
		
						if caterr==nil{
		
							indivCategory = append(indivCategory, category)
		
							parentCatId = category.ParentID
		
							if parentCatId!=0{
	
								var count int
		
								LOOP1:
								  
								   for{
	 
									  count = count + 1 //count increment used to check how many times the loop gets executed
		
									  var parentCategory model.TblCategory
		
									   caterr = db.Table("tbl_categories").Where("is_deleted = 0 and id = ?",parentCatId).First(&parentCategory).Error
		
									   if caterr==nil{
										  
										   indivCategory = append(indivCategory, parentCategory)
		
										   parentCatId = parentCategory.ParentID
		
										   if parentCatId!=0{
											  
											   goto LOOP1
		
										   }else if count>49{ //mannuall condition to break the loop in overlooping situations
	
											  break //use to break the loop if infinite loop doesn't break , So forcing the loop to break at overlooping conditions 
	
										   }else{
		
											  break
		
										   }
										   
									   }else{
		
										   indivCategory = append(indivCategory, model.TblCategory{})
		
										   break
		
									   }
								   }
							}
		
						}else{
		
							indivCategory = append(indivCategory, model.TblCategory{})
						}
		
						indivCategories = append(indivCategories, indivCategory)
		
					}
	
					entry.Categories = indivCategories
	
					final_entries_list = append(final_entries_list, entry)
	
				}
	
				channelEntrieslist := &model.ChannelEntries{ChannelEntryList: final_entries_list,Count: int(count)}
	
				channelEntriesDetails := model.ChannelEntryDetails{ChannelEntryList: channelEntrieslist,ChannelEntry: nil}
	
				return  channelEntriesDetails,nil
	
			}
	
			var final_error error
	
			if entrieserr!=nil{
	
				final_error = entrieserr
	
			}else{
	
				final_error = counterr
				
			}
	
			return model.ChannelEntryDetails{},final_error
	
		}else if channelID==nil && channelEntryID==nil && limit!=nil && offset!=nil{
	
			var channelEntries []model.TblChannelEntries
	
			var count int64
	
			entrieserr := db.Table("tbl_channel_entries").Where("tbl_channel_entries.status = 1").Limit(*limit).Offset(*offset).Order("tbl_channel_entries.id desc").
			Find(&channelEntries).Error
	
			counterr := db.Table("tbl_channel_entries").Where("tbl_channel_entries.status = 1").Count(&count).Error
	
			if entrieserr==nil && counterr==nil{
	
				var final_entries_list []model.TblChannelEntries
	
				for _,entry := range channelEntries{
	
					var indivCategories [][]model.TblCategory
	
					splittedArr := strings.Split(entry.CategoriesID, ",")
		
					var parentCatId int
		
					for _, catId := range splittedArr{
		
						var indivCategory []model.TblCategory
		
						conv_id,_ := strconv.Atoi(catId)
		
						var category model.TblCategory
		
						caterr := db.Table("tbl_categories").Where("is_deleted = 0 and id = ?",conv_id).First(&category).Error
		
						if caterr==nil{
		
							indivCategory = append(indivCategory, category)
		
							parentCatId = category.ParentID
		
							if parentCatId!=0{
	
								var count int
		
								LOOP2:
								  
								   for{
	 
									  count = count + 1 //count increment used to check how many times the loop gets executed
		
									  var parentCategory model.TblCategory
		
									   caterr = db.Table("tbl_categories").Where("is_deleted = 0 and id = ?",parentCatId).First(&parentCategory).Error
		
									   if caterr==nil{
										  
										   indivCategory = append(indivCategory, parentCategory)
		
										   parentCatId = parentCategory.ParentID
		
										   if parentCatId!=0{
											  
											   goto LOOP2
		
										   }else if count>49{ //mannuall condition to break the loop in overlooping situations
	
											  break //use to break the loop if infinite loop doesn't break , So forcing the loop to break at overlooping conditions 
	
										   }else{
		
											  break
		
										   }
										   
									   }else{
		
										   indivCategory = append(indivCategory, model.TblCategory{})
		
										   break
		
									   }
								   }
							}
		
						}else{
		
							indivCategory = append(indivCategory, model.TblCategory{})
						}
		
						indivCategories = append(indivCategories, indivCategory)
		
					}
	
					entry.Categories = indivCategories
	
					final_entries_list = append(final_entries_list, entry)
	
				}
	
				channelEntrieslist := &model.ChannelEntries{ChannelEntryList: final_entries_list,Count: int(count)}
	
				channelEntriesDetails := model.ChannelEntryDetails{ChannelEntryList: channelEntrieslist,ChannelEntry: nil}
	
				return  channelEntriesDetails,nil
			}
	
			var final_error error
	
			if entrieserr!=nil{
	
				final_error = entrieserr
	
			}else{
	
				final_error = counterr
	
			}
	
			return model.ChannelEntryDetails{},final_error
			
		}
	
		return model.ChannelEntryDetails{},errors.New("unable to fetch data")

	}else{

		memberid := c.GetInt("memberid")

		if channelEntryID!=nil || (channelID!=nil && channelEntryID!=nil) && limit==nil && offset==nil {
	
			var channelEntry model.TblChannelEntries
	
			var entryerr error
	
			if channelID!=nil{
	
			   entryerr = db.Table("tbl_channel_entries").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			   Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			   Where("tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and tbl_channel_entries.channel_id = ? and tbl_channel_entries.id = ?",memberid,channelID,channelEntryID).
			   First(&channelEntry).Error
	
			}else{
	
			  entryerr = db.Table("tbl_channel_entries").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			  Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			  Where("tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and tbl_channel_entries.id = ?",memberid,channelEntryID).
			  First(&channelEntry).Error
			}
	
			if entryerr==nil{
	
				splittedArr := strings.Split(channelEntry.CategoriesID, ",")
	
				var parentCatId int
	
				var indivCategories [][]model.TblCategory
	
				for _, catId := range splittedArr{
	
					var indivCategory []model.TblCategory
	
					conv_id,_ := strconv.Atoi(catId)
	
					var category model.TblCategory
	
					caterr := db.Table("tbl_categories").Where("is_deleted = 0 and id = ?",conv_id).First(&category).Error
	
					if caterr==nil{
	
						indivCategory = append(indivCategory, category)
	
						parentCatId = category.ParentID
	
						if parentCatId!=0{
	
							var count int 
	
							LOOP3:
							  
							   for{
	
								  count = count + 1 //count increment used to check how many times the loop gets executed
	
								  var parentCategory model.TblCategory
	
								   caterr = db.Table("tbl_categories").Where("is_deleted = 0 and id = ?",parentCatId).First(&parentCategory).Error
	
								   if caterr==nil{
									  
									   indivCategory = append(indivCategory, parentCategory)
	
									   parentCatId = parentCategory.ParentID
	
									   if parentCatId!=0{ //mannuall condition to break the loop in overlooping situations
										  
										   goto LOOP3
	
									   }else if count>49{
	
										break  //use to break the loop if infinite loop doesn't break ,So forcing the loop to break at overlooping conditions
	
									   }else{
	
										  break
									   }
									   
								   }else{
	
									   indivCategory = append(indivCategory, model.TblCategory{})
	
									   break
	
								   }
							   }
						}
	
					}else{
	
						indivCategory = append(indivCategory, model.TblCategory{})
					}
	
					indivCategories = append(indivCategories, indivCategory)
	
				}
	
				channelEntry.Categories = indivCategories
	
				channelEntryDetails := model.ChannelEntryDetails{ChannelEntryList: &model.ChannelEntries{},ChannelEntry: &channelEntry}
	
				return channelEntryDetails,nil
	
			}
	
			return model.ChannelEntryDetails{},entryerr	
	
		} else if channelEntryID==nil && channelID!=nil && limit!=nil && offset!=nil{
	
			var channelEntries []model.TblChannelEntries
	
			var count int64
	
			entrieserr := db.Table("tbl_channel_entries").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and tbl_channel_entries.channel_id = ?",memberid,channelID).Limit(*limit).Offset(*offset).Order("tbl_channel_entries.id desc").
			Find(&channelEntries).Error
	
			counterr := db.Table("tbl_channel_entries").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and tbl_channel_entries.channel_id = ?",memberid,channelID).
			Count(&count).Error
	
			if entrieserr==nil && counterr==nil{
	
				var final_entries_list []model.TblChannelEntries
	
				for _,entry := range channelEntries{
	
					var indivCategories [][]model.TblCategory
	
					splittedArr := strings.Split(entry.CategoriesID, ",")
		
					var parentCatId int
		
					for _, catId := range splittedArr{
		
						var indivCategory []model.TblCategory
		
						conv_id,_ := strconv.Atoi(catId)
		
						var category model.TblCategory
		
						caterr := db.Table("tbl_categories").Where("is_deleted = 0 and id = ?",conv_id).First(&category).Error
		
						if caterr==nil{
		
							indivCategory = append(indivCategory, category)
		
							parentCatId = category.ParentID
		
							if parentCatId!=0{
	
								var count int
		
								LOOP4:
								  
								   for{
	 
									  count = count + 1 //count increment used to check how many times the loop gets executed
		
									  var parentCategory model.TblCategory
		
									   caterr = db.Table("tbl_categories").Where("is_deleted = 0 and id = ?",parentCatId).First(&parentCategory).Error
		
									   if caterr==nil{
										  
										   indivCategory = append(indivCategory, parentCategory)
		
										   parentCatId = parentCategory.ParentID
		
										   if parentCatId!=0{
											  
											   goto LOOP4
		
										   }else if count>49{ //mannuall condition to break the loop in overlooping situations
	
											  break //use to break the loop if infinite loop doesn't break , So forcing the loop to break at overlooping conditions 
	
										   }else{
		
											  break
		
										   }
										   
									   }else{
		
										   indivCategory = append(indivCategory, model.TblCategory{})
		
										   break
		
									   }
								   }
							}
		
						}else{
		
							indivCategory = append(indivCategory, model.TblCategory{})
						}
		
						indivCategories = append(indivCategories, indivCategory)
		
					}
	
					entry.Categories = indivCategories
	
					final_entries_list = append(final_entries_list, entry)
	
				}
	
				channelEntrieslist := &model.ChannelEntries{ChannelEntryList: final_entries_list,Count: int(count)}
	
				channelEntriesDetails := model.ChannelEntryDetails{ChannelEntryList: channelEntrieslist,ChannelEntry: nil}
	
				return  channelEntriesDetails,nil
	
			}
	
			var final_error error
	
			if entrieserr!=nil{
	
				final_error = entrieserr
	
			}else{
	
				final_error = counterr
				
			}
	
			return model.ChannelEntryDetails{},final_error
	
		}else if channelID==nil && channelEntryID==nil && limit!=nil && offset!=nil{
	
			var channelEntries []model.TblChannelEntries
	
			var count int64
	
			entrieserr := db.Debug().Table("tbl_channel_entries").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?",memberid).Limit(*limit).Offset(*offset).Order("tbl_channel_entries.id desc").
			Find(&channelEntries).Error
	
			counterr := db.Table("tbl_channel_entries").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?",memberid).
			Count(&count).Error
	
			if entrieserr==nil && counterr==nil{
	
				var final_entries_list []model.TblChannelEntries
	
				for _,entry := range channelEntries{
	
					var indivCategories [][]model.TblCategory
	
					splittedArr := strings.Split(entry.CategoriesID, ",")
		
					var parentCatId int
		
					for _, catId := range splittedArr{
		
						var indivCategory []model.TblCategory
		
						conv_id,_ := strconv.Atoi(catId)
		
						var category model.TblCategory
		
						caterr := db.Table("tbl_categories").Where("is_deleted = 0 and id = ?",conv_id).First(&category).Error
		
						if caterr==nil{
		
							indivCategory = append(indivCategory, category)
		
							parentCatId = category.ParentID
		
							if parentCatId!=0{
	
								var count int
		
								LOOP5:
								  
								   for{
	 
									  count = count + 1 //count increment used to check how many times the loop gets executed
		
									  var parentCategory model.TblCategory
		
									   caterr = db.Table("tbl_categories").Where("is_deleted = 0 and id = ?",parentCatId).First(&parentCategory).Error
		
									   if caterr==nil{
										  
										   indivCategory = append(indivCategory, parentCategory)
		
										   parentCatId = parentCategory.ParentID
		
										   if parentCatId!=0{
											  
											   goto LOOP5
		
										   }else if count>49{ //mannuall condition to break the loop in overlooping situations
	
											  break //use to break the loop if infinite loop doesn't break , So forcing the loop to break at overlooping conditions 
	
										   }else{
		
											  break
		
										   }
										   
									   }else{
		
										   indivCategory = append(indivCategory, model.TblCategory{})
		
										   break
		
									   }
								   }
							}
		
						}else{
		
							indivCategory = append(indivCategory, model.TblCategory{})
						}
		
						indivCategories = append(indivCategories, indivCategory)
		
					}
	
					entry.Categories = indivCategories
	
					final_entries_list = append(final_entries_list, entry)
	
				}
	
				channelEntrieslist := &model.ChannelEntries{ChannelEntryList: final_entries_list,Count: int(count)}
	
				channelEntriesDetails := model.ChannelEntryDetails{ChannelEntryList: channelEntrieslist,ChannelEntry: nil}
	
				return  channelEntriesDetails,nil
			}
	
			var final_error error
	
			if entrieserr!=nil{
	
				final_error = entrieserr
	
			}else{
	
				final_error = counterr
	
			}
	
			return model.ChannelEntryDetails{},final_error
			
		}
	
		return model.ChannelEntryDetails{},errors.New("unable to fetch data")

	}

}

func ChannelDetail(db *gorm.DB,ctx context.Context, channelID int) (model.TblChannel, error) {

	c,_ := ctx.Value(ContextKey).(*gin.Context)

	token ,_ := c.Get("token")

	var channel model.TblChannel

	if token == SpecialToken{

	    if err := db.Table("tbl_channels").Select("distinct on (tbl_channels.id) tbl_channels.*").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
	    Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and id = ?",channelID).First(&channel).Error;err!=nil{

		  return model.TblChannel{},err
	    }

	    return channel,nil

	}else{

		memberid := c.GetInt("memberid")

	    var channel model.TblChannel

	    if err := db.Table("tbl_channels").Select("distinct on (tbl_channels.id) tbl_channels.*").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
	    Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
	    Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
	    Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and  id = ?",memberid,channelID).First(&channel).Error;err!=nil{

		  return model.TblChannel{},err
	    }

	    return channel,nil
	}

}


