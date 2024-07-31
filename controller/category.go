package controller

import (
	"context"
	"spurtcms-graphql/graph/model"

	"gorm.io/gorm"
)

func CategoriesList(db *gorm.DB, ctx context.Context, limit, offset, categoryGroupId *int, categoryGroupSlug *string, hierarchyLevel, excludeGroup, excludeParent, checkEntriesPresence *int) (*model.CategoriesList, error) {

	// fmt.Println("dbchk",db.Config.Dialector.Name())

	// c, _ := ctx.Value(ContextKey).(*gin.Context)

	// var (
	// 	FinalCategoryList                                                                 []model.Category
	// 	limitVal, offsetVal, categoryGrpIdVal ,excludeParentVal                           int
	// 	hierarchyLevelVal, checkEntriesPresenceVal, excludeGroupVal                       int
	// 	categoryGroupSlugVal                                                              string
	// )

	// if categoryGroupId != nil {

	// 	categoryGrpIdVal = *categoryGroupId
	// }

	// if limit != nil{

	// 	limitVal = *limit
	// }

	// if offset != nil{

	// 	offsetVal = *offset
	// }

	// if categoryGroupSlug != nil {

	// 	categoryGroupSlugVal = *categoryGroupSlug
	// }

	// if hierarchyLevel != nil {

	// 	hierarchyLevelVal = *hierarchyLevel
	// }

	// if checkEntriesPresence != nil {

	// 	checkEntriesPresenceVal = *checkEntriesPresence
	// }

	// if excludeGroup != nil {

	// 	excludeGroupVal = *excludeGroup
	// }

	// if excludeParent != nil {

	// 	excludeParentVal = *excludeParent
	// }

	// // memberid := c.GetInt("memberid")

	// categories,count, err := CategoryInstance.CategoryList(limitVal, offsetVal, categoryGrpIdVal, hierarchyLevelVal, checkEntriesPresenceVal,excludeGroupVal,excludeParentVal,categoryGroupSlugVal)

	// if err != nil {

	// 	ErrorLog.Printf("category list retrieval error: %s", err)

	// 	c.AbortWithError(http.StatusInternalServerError, err)

	// 	return &model.CategoriesList{}, err

	// }

	// for _, category := range categories {

	// 	localCategory := model.Category{
	// 		ID: category.Id,
	// 		CategoryName: category.CategoryName,
	// 		CategorySlug: category.CategorySlug,
	// 		Description: category.Description,
	// 		ImagePath: category.ImagePath,
	// 		CreatedOn: category.CreatedOn,
	// 		CreatedBy: category.CreatedBy,
	// 		ModifiedOn: &category.ModifiedOn,
	// 		ModifiedBy: &category.ModifiedBy,
	// 		ParentID: category.ParentId,
	// 	}

	// 	FinalCategoryList = append(FinalCategoryList, localCategory)
	// }

	// return &model.CategoriesList{Categories: FinalCategoryList, Count: count}, nil

	return &model.CategoriesList{}, nil
}
