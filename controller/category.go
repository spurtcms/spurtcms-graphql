package controller

import (
	"context"
	"gqlserver/graph/model"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func CategoriesList(db *gorm.DB, ctx context.Context, limit, offset, categoryGroupId, hierarchyLevel *int) (model.CategoriesList, error) {

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

		outerlevel = ` and level = `+strconv.Itoa(*hierarchyLevel)

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

	if err := db.Raw(` ` + res + ` SELECT count(*) FROM cat_tree where is_deleted = 0 ` + selectGroupRemove + outerlevel + ` and parent_id != 0 group by id order by id desc`).Count(&count).Error; err != nil {

		return model.CategoriesList{}, err
	}

	var final_categoriesList []model.Category

	seenCategory := make(map[int]bool)

	for _, category := range categories {

		if !seenCategory[category.ID] {

			var modified_path string

			if category.ImagePath != "" {

				modified_path = PathUrl + strings.TrimPrefix(category.ImagePath, "/")
			}

			category.ImagePath = modified_path

			final_categoriesList = append(final_categoriesList, category)

			seenCategory[category.ID] = true
		}
	}

	return model.CategoriesList{Categories: final_categoriesList, Count: int(count)}, nil
}
