package controller

import (
	"context"
	"gqlserver/graph/model"
	"os"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SpaceList(db *gorm.DB, ctx context.Context, limit, offset int) (model.SpaceDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	memberid := c.GetInt("memberid")

	var spacelist []model.Space

	var count int64

	if token == SpecialToken {

		db.Debug().Table("tbl_spaces_aliases").Select("tbl_spaces_aliases.id,tbl_spaces_aliases.spaces_name,tbl_spaces_aliases.spaces_slug,tbl_spaces_aliases.spaces_description,tbl_spaces_aliases.image_path,tbl_spaces_aliases.language_id,tbl_spaces_aliases.created_on,tbl_spaces_aliases.created_by,tbl_spaces_aliases.modified_on,tbl_spaces_aliases.modified_by,tbl_spaces_aliases.is_deleted,tbl_spaces_aliases.deleted_on,tbl_spaces_aliases.deleted_by,tbl_spaces.page_category_id").
			Joins("inner join tbl_spaces on tbl_spaces.id = tbl_spaces_aliases.spaces_id").Where("tbl_spaces.is_deleted = 0 and tbl_spaces_aliases.is_deleted = 0").
			Order("tbl_spaces.id desc").Limit(limit).Offset(offset).Find(&spacelist)

		db.Debug().Table("tbl_spaces_aliases").Joins("inner join tbl_spaces on tbl_spaces.id = tbl_spaces_aliases.spaces_id").Where("tbl_spaces.is_deleted = 0 and tbl_spaces_aliases.is_deleted = 0").Count(&count)

	} else {

		db.Debug().Table("tbl_spaces_aliases").Select("distinct on (tbl_spaces.id) tbl_spaces_aliases.id,tbl_spaces_aliases.spaces_name,tbl_spaces_aliases.spaces_slug,tbl_spaces_aliases.spaces_description,tbl_spaces_aliases.image_path,tbl_spaces_aliases.language_id,tbl_spaces_aliases.created_on,tbl_spaces_aliases.created_by,tbl_spaces_aliases.modified_on,tbl_spaces_aliases.modified_by,tbl_spaces_aliases.is_deleted,tbl_spaces_aliases.deleted_on,tbl_spaces_aliases.deleted_by,tbl_spaces.page_category_id").
			Joins("inner join tbl_spaces on tbl_spaces.id = tbl_spaces_aliases.spaces_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.spaces_id = tbl_spaces.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_spaces.is_deleted = 0 and tbl_spaces_aliases.is_deleted = 0 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?", memberid).Order("tbl_spaces.id desc").Limit(limit).Offset(offset).Find(&spacelist)

		db.Debug().Table("tbl_spaces_aliases").Distinct("tbl_spaces.id").Joins("inner join tbl_spaces on tbl_spaces.id = tbl_spaces_aliases.spaces_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.spaces_id = tbl_spaces.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_spaces.is_deleted = 0 and tbl_spaces_aliases.is_deleted = 0 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?", memberid).Count(&count)
	}

	var final_spacelist []model.Space

	for _, space := range spacelist {

		modified_path := os.Getenv("DOMAIN_URL") + strings.TrimPrefix(space.ImagePath, "/")

		space.ImagePath = modified_path

		var categories []model.TblCategory

		var parent_category model.TblCategory

		db.Table("tbl_categories").Where("tbl_categories.is_deleted = 0 and tbl_categories.id = ?", space.CategoryID).First(&parent_category)

		if parent_category.ID != 0 {

			categories = append(categories, parent_category)
		}

		parentCatId := parent_category.ParentID

		if parentCatId != 0 {

		LOOP:

			count := 0

			for {

				count++ //count increment used to check how many times the loop gets executed

				var loopParentCategory model.TblCategory

				db.Table("tbl_categories").Where("tbl_categories.is_deleted = 0 and tbl_categories.id = ?", parentCatId).First(&loopParentCategory)

				if loopParentCategory.ID != 0 {

					categories = append(categories, loopParentCategory)
				}

				parentCatId = loopParentCategory.ParentID

				if parentCatId != 0 {

					goto LOOP

				} else if count > 49 { //mannuall condition to break the loop in overlooping situations

					break //use to break the loop if infinite loop doesn't break ,So forcing the loop to break at overlooping conditions

				} else {

					break

				}

			}

		}

		if len(categories) > 0 {

			sort.SliceStable(categories, func(i, j int) bool {

				return categories[i].ID < categories[j].ID

			})

			space.Categories = categories

		}

		final_spacelist = append(final_spacelist, space)

	}

	return model.SpaceDetails{Spacelist: final_spacelist, Count: int(count)}, nil
}

func SpaceDetails(db *gorm.DB, ctx context.Context, spaceId int) (model.Space, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	memberid := c.GetInt("memberid")

	var space model.Space

	if token == SpecialToken {

		db.Debug().Table("tbl_spaces_aliases").Select("tbl_spaces_aliases.id,tbl_spaces_aliases.spaces_name,tbl_spaces_aliases.spaces_slug,tbl_spaces_aliases.spaces_description,tbl_spaces_aliases.image_path,tbl_spaces_aliases.language_id,tbl_spaces_aliases.created_on,tbl_spaces_aliases.created_by,tbl_spaces_aliases.modified_on,tbl_spaces_aliases.modified_by,tbl_spaces_aliases.is_deleted,tbl_spaces_aliases.deleted_on,tbl_spaces_aliases.deleted_by,tbl_spaces.page_category_id").
			Joins("inner join tbl_spaces on tbl_spaces.id = tbl_spaces_aliases.spaces_id").Where("tbl_spaces.is_deleted = 0 and tbl_spaces_aliases.is_deleted = 0 and tbl_spaces.id = ?", spaceId).First(&space)

	} else {

		db.Debug().Table("tbl_spaces_aliases").Select("distinct on (tbl_spaces_aliases.id) tbl_spaces_aliases.id,tbl_spaces_aliases.spaces_name,tbl_spaces_aliases.spaces_slug,tbl_spaces_aliases.spaces_description,tbl_spaces_aliases.image_path,tbl_spaces_aliases.language_id,tbl_spaces_aliases.created_on,tbl_spaces_aliases.created_by,tbl_spaces_aliases.modified_on,tbl_spaces_aliases.modified_by,tbl_spaces_aliases.is_deleted,tbl_spaces_aliases.deleted_on,tbl_spaces_aliases.deleted_by,tbl_spaces.page_category_id").
			Joins("inner join tbl_spaces on tbl_spaces.id = tbl_spaces_aliases.spaces_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.spaces_id = tbl_spaces.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_spaces.is_deleted = 0 and tbl_spaces_aliases.is_deleted = 0 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_spaces.id = ? and tbl_members.id = ?", spaceId, memberid).First(&space)
	}

	modified_path := os.Getenv("DOMAIN_URL") + strings.TrimPrefix(space.ImagePath, "/")

	space.ImagePath = modified_path

	var categories []model.TblCategory

	var parent_category model.TblCategory

	db.Table("tbl_categories").Where("tbl_categories.is_deleted = 0 and tbl_categories.id = ?", space.CategoryID).First(&parent_category)

	if parent_category.ID != 0 {

		categories = append(categories, parent_category)
	}

	parentCatId := parent_category.ParentID

	if parentCatId != 0 {

	LOOP:

		count := 0

		for {

			count++ //count increment used to check how many times the loop gets executed

			var loopParentCategory model.TblCategory

			db.Table("tbl_categories").Where("tbl_categories.is_deleted = 0 and tbl_categories.id = ?", parentCatId).First(&loopParentCategory)

			if loopParentCategory.ID != 0 {

				categories = append(categories, loopParentCategory)
			}

			parentCatId = loopParentCategory.ParentID

			if parentCatId != 0 {

				goto LOOP

			} else if count > 49 { //mannuall condition to break the loop in overlooping situations

				break //use to break the loop if infinite loop doesn't break ,So forcing the loop to break at overlooping conditions

			} else {

				break

			}

		}

	}

	if len(categories) > 0 {

		sort.SliceStable(categories, func(i, j int) bool {

			return categories[i].ID < categories[j].ID

		})

		space.Categories = categories

	}

	return space,nil

}

// func PageAndPageGrouplistBySpaceId(db *gorm.DB,ctx context.Context,limit, offset int)(model,error){

// 	c, _ := ctx.Value(ContextKey).(*gin.Context)

// 	token, _ := c.Get("token")

// 	memberid := c.GetInt("memberid")

// }
