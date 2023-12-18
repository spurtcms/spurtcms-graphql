package controller

import (
	"gqlserver/graph/model"
	"time"

	"github.com/spurtcms/spurtcms-content/lms"
	"gorm.io/gorm"
)

func Spacelist(db *gorm.DB,filter model.Filter)(model.SpacesDetails,error){

	Sp.MemAuth = GetAuthorizationWithoutToken(db)

	spacelist,space_count,err := Sp.MemberSpaceList(filter.Limit,filter.Offset,lms.Filter{Keyword: *filter.Keyword,CategoryId: *filter.CategoryID})

	if err!=nil{

		return model.SpacesDetails{},err 
	}

	var conv_spacelist []model.TblSpacesAliases

	for _, space := range spacelist {

		var spaceobj model.TblSpacesAliases

		var catNames []model.TblCategory

		for _, catname := range space.CategoryNames {

			var parentWith_child []model.Result

			for _, pwc := range catname.ParentWithChild {

				parentWith_child = append(parentWith_child,model.Result(pwc))
			}

			catObj := model.TblCategory{
				ID:                 catname.Id,
				CategoryName:       catname.CategoryName,
				CategorySlug:       catname.CategorySlug,
				Description:        catname.Description,
				ImagePath:          catname.ImagePath,
				CreatedOn:          catname.CreatedOn,
				CreatedBy:          catname.CreatedBy,
				ModifiedOn:         catname.ModifiedOn,
				ModifiedBy:         catname.ModifiedBy,
				IsDeleted:          catname.IsDeleted,
				DeletedOn:          catname.DeletedOn,
				DeletedBy:          catname.DeletedBy,
				ParentID:           catname.ParentId,
				CreatedDate:        catname.CreatedDate,
				ModifiedDate:       catname.ModifiedDate,
				DateString:         catname.DateString,
				ParentCategoryName: catname.ParentCategoryName,
				Parent:             catname.Parent,
				ParentWithChild:    parentWith_child,
			}

			catNames = append(catNames, catObj)
		}

		spaceobj = model.TblSpacesAliases{
			ID:                space.Id,
			SpacesID:          space.SpacesId,
			SpacesName:        space.SpacesName,
			LanguageID:        space.LanguageId,
			SpacesSlug:        space.SpacesSlug,
			SpacesDescription: space.SpacesDescription,
			ImagePath:         space.ImagePath,
			CreatedOn:         space.CreatedOn,
			CreatedBy:         space.CreatedBy,
			ModifiedOn:        (*time.Time)(&space.ModifiedOn),
			ModifiedBy:        space.ModifiedBy,
			DeletedOn:         (*time.Time)(&space.DeletedOn),
			DeletedBy:         space.DeletedBy,
			IsDeleted:         space.IsDeleted,
			PageCategoryID:    space.PageCategoryId,
			ParentID:          space.ParentId,
			CreatedDate:       space.CreatedDate,
			ModifiedDate:      space.ModifiedDate,
			CategoryNames:     catNames,
			CategoryID:        space.CategoryId,
			FullSpaceAccess:   space.FullSpaceAccess,
		}

		conv_spacelist = append(conv_spacelist, spaceobj)
	}

	spaceDetails := model.SpacesDetails{Spaces: conv_spacelist, Count: int(space_count)}

	return spaceDetails, nil
}