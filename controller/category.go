package controller

import (
	"context"
	"gqlserver/graph/model"
	"strconv"

	"os"
	"strings"
	// "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CategoriesList(db *gorm.DB, ctx context.Context, limit, offset int) (model.CategoriesList, error) {

	// c,_ := ctx.Value(ContextKey).(*gin.Context)

	// token := c.GetString("token")

	// memberid := c.GetInt("memberid")

	var pathUrl string

	if os.Getenv("DOMAIN_URL") != "" {

		pathUrl = os.Getenv("DOMAIN_URL")

	} else {

		pathUrl = os.Getenv("LOCAL_URL")
	}

	var categories []model.TblCategory

	var count int64

	res := `WITH RECURSIVE cat_tree AS (
		SELECT id, category_name, category_slug,image_path, parent_id,created_on,modified_on,is_deleted
		FROM tbl_categories
		WHERE id = 1
		UNION ALL
		SELECT cat.id, cat.category_name, cat.category_slug, cat.image_path ,cat.parent_id,cat.created_on,cat.modified_on,
		cat.is_deleted
		FROM tbl_categories AS cat
		JOIN cat_tree ON cat.parent_id = cat_tree.id )`

	if err := db.Debug().Raw(` ` + res + ` SELECT cat_tree.* FROM cat_tree where is_deleted = 0 and id not in (1) and parent_id =1  order by id desc limit ` + strconv.Itoa(limit) + ` offset ` + strconv.Itoa(offset)).Find(&categories).Error; err != nil {

		return model.CategoriesList{}, err
	}

	if err := db.Debug().Raw(` ` + res + ` SELECT count(*) FROM cat_tree where is_deleted = 0 and id not in (1) and parent_id =1 group by id order by id desc`).Count(&count).Error; err != nil {

		return model.CategoriesList{}, err
	}

	var final_categoriesList []model.TblCategory

	for _, parentCat := range categories {

		modified_path := pathUrl + strings.TrimPrefix(parentCat.ImagePath, "/")

		parentCat.ImagePath = modified_path

		var childCategories []model.TblCategory

		err := db.Debug().Raw(` ` + res + ` SELECT cat_tree.* FROM cat_tree where is_deleted = 0 and id not in (` + strconv.Itoa(parentCat.ID) + `) and parent_id =` + strconv.Itoa(parentCat.ID) + ` order by id desc`).Find(&childCategories).Error

		if err == nil {

			parentCat.ChildCategories = childCategories
		}

		final_categoriesList = append(final_categoriesList, parentCat)

	}

	return model.CategoriesList{Categories: final_categoriesList, Count: int(count)}, nil
}

