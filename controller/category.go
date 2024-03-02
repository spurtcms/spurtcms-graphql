package controller

import (
	"context"
	"gqlserver/graph/model"
	"os"
	"strings"

	// "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CategoriesList(db *gorm.DB,ctx context.Context,limit, offset int)(model.CategoriesList,error){

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

	if err := db.Table("tbl_categories").Select("tbl_categories.*").Where("tbl_categories.is_deleted=0 and parent_id=0").Order("tbl_categories.id desc").Limit(limit).Offset(offset).Find(&categories).Error;err!=nil{

		return model.CategoriesList{},err
	}

	if err := db.Table("tbl_categories").Where("tbl_categories.is_deleted=0 and parent_id=0").Count(&count).Error;err!=nil{

		return model.CategoriesList{},err
	}

	var final_categoriesList []model.TblCategory

	for _,category := range categories{

		modified_path :=  pathUrl + strings.TrimPrefix(category.ImagePath,"/")

		category.ImagePath = modified_path

		var childCategories []model.TblCategory

		err := db.Table("tbl_categories").Select("tbl_categories.*").Where("tbl_categories.is_deleted=0 and parent_id=?",category.ID).Find(&childCategories).Error

		if err==nil{

			category.ChildCategories = childCategories
		}

		final_categoriesList = append(final_categoriesList, category)
	}

	return model.CategoriesList{Categories: final_categoriesList,Count: int(count)},nil
}


