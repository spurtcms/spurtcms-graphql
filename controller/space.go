package controller

import (
	"context"
	"os"
	"spurtcms-graphql/graph/model"

	"github.com/gin-gonic/gin"
	spaces "github.com/spurtcms/pkgcontent/spaces"
	"gorm.io/gorm"
)

func SpaceList(db *gorm.DB, ctx context.Context, limit, offset int, categoriesID *int) (*model.SpaceDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	var pathUrl string

	if os.Getenv("DOMAIN_URL") != "" {

		pathUrl = os.Getenv("DOMAIN_URL")

	} else {

		pathUrl = os.Getenv("LOCAL_URL")
	}

	spaceAuth := spaces.Space{Authority: GetAuthorization(token.(string), db)}

	spacelist, count, err := spaceAuth.GetGraphqlSpacelist(limit, offset, pathUrl, categoriesID)

	if err != nil {

		return &model.SpaceDetails{}, err
	}

	var final_spacelist []model.Space

	for _, space := range spacelist {

		var conv_categories []model.Category

		for _, category := range space.CategoryNames {

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

			conv_categories = append(conv_categories, conv_category)
		}

		conv_space := model.Space{
			ID:               space.Id,
			SpaceName:        space.SpacesName,
			SpaceSlug:        space.SpacesSlug,
			SpaceDescription: space.SpacesDescription,
			ImagePath:        space.ImagePath,
			LanguageID:       space.LanguageId,
			CreatedOn:        space.CreatedOn,
			CreatedBy:        space.CreatedBy,
			ModifiedOn:       &space.ModifiedOn,
			ModifiedBy:       &space.ModifiedBy,
			CategoryID:       space.PageCategoryId,
			Categories:       conv_categories,
		}

		final_spacelist = append(final_spacelist, conv_space)

	}

	return &model.SpaceDetails{Spacelist: final_spacelist, Count: int(count)}, nil
}

func SpaceDetails(db *gorm.DB, ctx context.Context, spaceId int) (*model.Space, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	var pathUrl string

	if os.Getenv("DOMAIN_URL") != "" {

		pathUrl = os.Getenv("DOMAIN_URL")

	} else {

		pathUrl = os.Getenv("LOCAL_URL")
	}

	spaceAuth := spaces.Space{Authority: GetAuthorization(token.(string), db)}

	space, err := spaceAuth.GetGraphqlSpaceDetails(spaceId, pathUrl)

	if err != nil {

		return &model.Space{}, err
	}

	var conv_categories []model.Category

	for _, category := range space.CategoryNames {

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

		conv_categories = append(conv_categories, conv_category)
	}

	conv_space := model.Space{
		ID:               space.Id,
		SpaceName:        space.SpacesName,
		SpaceSlug:        space.SpacesSlug,
		SpaceDescription: space.SpacesDescription,
		ImagePath:        space.ImagePath,
		LanguageID:       space.LanguageId,
		CreatedOn:        space.CreatedOn,
		CreatedBy:        space.CreatedBy,
		ModifiedOn:       &space.ModifiedOn,
		ModifiedBy:       &space.ModifiedBy,
		CategoryID:       space.PageCategoryId,
		Categories:       conv_categories,
	}

	return &conv_space, nil
}

func PagesAndPageGroupsBySpaceId(db *gorm.DB, ctx context.Context, spaceId int) (*model.PageAndPageGroups, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	token, _ := c.Get("token")

	spaceAuth := spaces.Space{Authority: GetAuthorization(token.(string), db)}

	pagez, subpagez, pagegroupz, err := spaceAuth.GetPagesAndPagegroupsUnderSpace(spaceId)

	if err != nil {

		return &model.PageAndPageGroups{}, err
	}

	var conv_pages []model.Page

	var conv_subpages []model.SubPage

	var conv_pagegroups []model.PageGroup

	for _, page := range pagez {

		conv_page := model.Page{
			ID:          page.Id,
			PageName:    page.PageTitle,
			Content:     page.PageDescription,
			PagegroupID: page.PageGroupId,
			OrderIndex:  page.OrderIndex,
			ParentID:    page.ParentId,
			Status:      page.Status,
			CreatedOn:   page.CreatedOn,
			CreatedBy:   page.CreatedBy,
			ModifiedOn:  &page.ModifiedOn,
			ModifiedBy:  &page.ModifiedBy,
		}

		conv_pages = append(conv_pages, conv_page)
	}

	for _, subpage := range subpagez {

		conv_subpage := model.SubPage{
			ID:          subpage.Id,
			SubpageName: subpage.PageTitle,
			Conent:      subpage.PageDescription,
			ParentID:    subpage.ParentId,
			PageGroupID: subpage.PageGroupId,
			OrderIndex:  subpage.PageSuborder,
			Status:      subpage.Status,
			CreatedOn:   subpage.CreatedOn,
			CreatedBy:   subpage.CreatedBy,
			ModifiedOn:  &subpage.ModifiedOn,
			ModifiedBy:  &subpage.ModifiedBy,
		}

		conv_subpages = append(conv_subpages, conv_subpage)
	}

	for _, pagegroup := range pagegroupz {

		conv_pagegroup := model.PageGroup{
			ID:            pagegroup.Id,
			PagegroupName: pagegroup.GroupName,
			OrderIndex:    pagegroup.OrderIndex,
			CreatedOn:     pagegroup.CreatedOn,
			CreatedBy:     pagegroup.CreatedBy,
			ModifiedOn:    &pagegroup.ModifiedOn,
			ModifiedBy:    &pagegroup.ModifiedBy,
		}

		conv_pagegroups = append(conv_pagegroups, conv_pagegroup)
	}

	return &model.PageAndPageGroups{Pages: conv_pages, Subpages: conv_subpages, Pagegroups: conv_pagegroups}, nil

}
