package controller

import (
	"context"
	"gqlserver/graph/model"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CategoriesList(db *gorm.DB, ctx context.Context, limit, offset, categoryGroupId, hierarchyLevel, checkEntriesPresence *int) (model.CategoriesList, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	var categories []model.Category

	var count int64

	category_string := ""

	selectGroupRemove := ""

	if categoryGroupId != nil {

		category_string = `WHERE id = ` + strconv.Itoa(*categoryGroupId)

		selectGroupRemove = `AND id != ` + strconv.Itoa(*categoryGroupId)
	}

	hierarchy_string := ""

	fromhierarchy_string := ""

	selecthierarchy_string := ""

	outerlevel := ""

	if hierarchyLevel != nil {

		hierarchy_string = ` WHERE CAT_TREE.LEVEL < ` + strconv.Itoa(*hierarchyLevel)

		fromhierarchy_string = `,CAT_TREE.LEVEL + 1`

		selecthierarchy_string = `,0 AS LEVEL`

		outerlevel = ` and level = ` + strconv.Itoa(*hierarchyLevel)

	}

	limit_offString := ""

	if limit != nil && offset != nil {

		limit_offString = `limit ` + strconv.Itoa(*limit) + ` offset ` + strconv.Itoa(*offset)
	}

	res := `WITH RECURSIVE cat_tree AS (
		SELECT id, category_name, category_slug,image_path, parent_id,created_on,modified_on,is_deleted` + selecthierarchy_string + `
		FROM tbl_categories ` + category_string + `
		UNION ALL
		SELECT cat.id, cat.category_name, cat.category_slug, cat.image_path ,cat.parent_id,cat.created_on,cat.modified_on,
		cat.is_deleted` + fromhierarchy_string + `
		FROM tbl_categories AS cat
		JOIN cat_tree ON cat.parent_id = cat_tree.id ` + hierarchy_string + ` )`

	if err := db.Debug().Raw(` ` + res + `SELECT cat_tree.* FROM cat_tree where is_deleted = 0 ` + selectGroupRemove + outerlevel + ` and parent_id != 0 order by id desc ` + limit_offString).Find(&categories).Error; err != nil {

		return model.CategoriesList{}, err
	}

	if err := db.Raw(` ` + res + ` SELECT count(*) FROM cat_tree where is_deleted = 0 ` + selectGroupRemove + outerlevel + ` and parent_id != 0  group by id order by id desc`).Count(&count).Error; err != nil {

		return model.CategoriesList{}, err
	}

	var final_categoriesList []model.Category

	seenCategory := make(map[int]bool)

	for _, category := range categories {

		if !seenCategory[category.ID] {

			var categoryIds string

			if checkEntriesPresence != nil && *checkEntriesPresence > 0 && *hierarchyLevel > 0 {

				Query := db.Table("tbl_channel_entries").Select("tbl_channel_entries.categories_id").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1").Where(`` + strconv.Itoa(category.ID) + `= any(string_to_array(tbl_channel_entries.categories_id,',')::integer[])`)

				if memberid > 0 {

					innerSubQuery := db.Table("tbl_channel_entries").Select("tbl_channel_entries.id").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1").Where(`` + strconv.Itoa(category.ID) + `= any(string_to_array(tbl_channel_entries.categories_id,',')::integer[])`)

					subquery := db.Table("tbl_access_control_pages").Select("tbl_access_control_pages.entry_id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
						Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id")

					subquery = subquery.Where("tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?", memberid).Where("tbl_access_control_pages.entry_id in (?)", innerSubQuery)

					Query = Query.Where("tbl_channel_entries.id not in (?)", subquery)
				}

				err := Query.Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1").Where(`` + strconv.Itoa(category.ID) + `= any(string_to_array(tbl_channel_entries.categories_id,',')::integer[])`).Find(&categoryIds).Error

				if err != nil {

					return model.CategoriesList{}, err
				}

				if categoryIds != "" {

					log.Println("categoryIds", categoryIds)

					var modified_path string

					if category.ImagePath != "" {

						modified_path = PathUrl + strings.TrimPrefix(category.ImagePath, "/")
					}

					category.ImagePath = modified_path

					final_categoriesList = append(final_categoriesList, category)
				}

			}else {

				var modified_path string
	
				if category.ImagePath != "" {
	
					modified_path = PathUrl + strings.TrimPrefix(category.ImagePath, "/")
				}
	
				category.ImagePath = modified_path
	
				final_categoriesList = append(final_categoriesList, category)
	
			}
	
			seenCategory[category.ID] = true

		} 

	}

	if checkEntriesPresence != nil && *checkEntriesPresence > 0 && *hierarchyLevel > 0 {

		count = int64(len(final_categoriesList))
	}

	return model.CategoriesList{Categories: final_categoriesList, Count: int(count)}, nil
}
