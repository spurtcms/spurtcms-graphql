package controller

import (
	"context"
	"errors"
	"gqlserver/graph/model"
	"github.com/gin-gonic/gin"
	"github.com/spurtcms/spurtcms-content/lms"
	"gorm.io/gorm"
)

func Pagelist(db *gorm.DB, spaceid int) (model.PageAndPagegroups, error) {

	Pg.MemAuth = GetAuthorizationWithoutToken(db)

	pagegroups, pages, subpages, err := Pg.MemberPageList(spaceid)

	if err != nil {

		return model.PageAndPagegroups{}, err
	}

	var pggz []model.Pagegroups

	for _, pgg := range pagegroups {

		pggObj := model.Pagegroups{
			GroupID:    pgg.GroupId,
			NewGroupID: pgg.NewGroupId,
			Name:       pgg.Name,
			OrderIndex: pgg.OrderIndex,
		}

		pggz = append(pggz, pggObj)
	}

	var pgz []model.Pages

	for _, page := range pages {

		pgObj := model.Pages{
			PgID:       page.PgId,
			Name:       page.Name,
			NewPgID:    page.NewPgId,
			Content:    page.Content,
			Pgroupid:   page.Pgroupid,
			NewGrpID:   page.NewGrpId,
			OrderIndex: page.OrderIndex,
			ParentID:   page.ParentId,
		}

		pgz = append(pgz, pgObj)
	}

	var spgz []model.Subpages

	for _, spg := range subpages {

		spgObj := model.Subpages{
			SpgID:       spg.SpgId,
			NewSpID:     spg.NewSpId,
			Name:        spg.Name,
			Content:     spg.Content,
			ParentID:    spg.ParentId,
			NewParentID: spg.NewParentId,
			PgroupID:    spg.NewParentId,
			NewPgroupID: spg.NewParentId,
			OrderIndex:  spg.OrderIndex,
		}

		spgz = append(spgz, spgObj)

	}

	pages_and_pgg := model.PageAndPagegroups{Pages: pgz, SubPages: spgz, Pagegroups: pggz}

	return pages_and_pgg,nil
}

func PageContent(db *gorm.DB,ctx context.Context, pageid int) (string, error) {

	c,_ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetHeader("Authorization")

	if token == ""{

		return "",errors.New("Login Required!")

	}

	Pg.MemAuth = GetAuthorization(token,db)

	pageDetails, err := Pg.GetPageContent(pageid)

	if err!=nil{

		return "",err
	}

	pageContent := pageDetails.PageDescription

	return pageContent,nil
}

func UpdateHighlights(db *gorm.DB,ctx context.Context,highlights model.Highlights)(bool,error){

	c,_ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	Pg.MemAuth = GetAuthorization(token,db)

	highlightRequest := lms.HighlightsReq{
                        Pageid: highlights.Pageid,
						Content: highlights.Content,
						Start: highlights.StartOffset,
						Offset: highlights.EndOffset,
						SelectPara: highlights.SelectPara,
						ContentColor: highlights.ContentColor,

	}

	isUpdated,err := Pg.UpdateHighlights(highlightRequest)

	if err!=nil || !isUpdated{

		return isUpdated,err
	}

	return isUpdated,nil
}

func UpdateNotes(db *gorm.DB,ctx context.Context, pageid int, notes string) (bool, error) {

	c,_ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	Pg.MemAuth = GetAuthorization(token,db)

	isUpdated,err := Pg.UpdateNotes(pageid,notes)

	if err!=nil || !isUpdated{

		return isUpdated,err
	}

	return isUpdated,nil
}

func GetNotesOrHighlights(db *gorm.DB, ctx context.Context, pageid int,contentType string) ([]model.MemberNotesHighlight, error) {

	if contentType == "" || pageid == 0{

		var Error error 

		if contentType == ""{

			Error = errors.New("Content-type Required!")

		}else if(pageid == 0){

			Error = errors.New("Page-id Required!")

		}

		return []model.MemberNotesHighlight{},Error
	}

	c,_ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	Pg.MemAuth = GetAuthorization(token,db)

	var notesOrHilitz []lms.TblMemberNotesHighlight

	var pckg_err error

	if contentType == "highlights"{

		notesOrHilitz,pckg_err =  Pg.GetHighlights(pageid)

	}else if contentType == "notes"{

		notesOrHilitz,pckg_err =  Pg.GetNotes(pageid)
	}

	if pckg_err!=nil{

		return []model.MemberNotesHighlight{},pckg_err
	}

	var conv_notesOrHilitz []model.MemberNotesHighlight

	for _,noteHilit := range notesOrHilitz{

		hilitObj := model.MemberNotesHighlight{
			        ID: noteHilit.Id,
					MemberID: noteHilit.MemberId,
					PageID: noteHilit.PageId,
					NotesHighlightsContent: noteHilit.NotesHighlightsContent,
					NotesHighlightsType: noteHilit.NotesHighlightsType,
					HighlightsConfiguration: noteHilit.HighlightsConfiguration,
					CreatedBy: noteHilit.CreatedBy,
					CreatedOn: noteHilit.CreatedOn,
					ModifiedOn: &noteHilit.ModifiedOn,
					ModifiedBy: &noteHilit.ModifiedBy,
					DeletedOn: &noteHilit.DeletedOn,
					DeletedBy: &noteHilit.DeletedBy,
					IsDeleted: noteHilit.IsDeleted,

		}

		conv_notesOrHilitz = append(conv_notesOrHilitz,hilitObj)
	}

	return conv_notesOrHilitz,nil
}

func DeleteNotesOrHighlights(db *gorm.DB,ctx context.Context, contentID int) (bool, error) {

	c,_ := ctx.Value(ContextKey).(*gin.Context)

	token := c.GetString("token")

	Pg.MemAuth = GetAuthorization(token,db)

	isDeleted ,err := Pg.RemoveHighlightsandNotes(contentID)

	if err!=nil || !isDeleted{

		return isDeleted,err
	}

	return isDeleted,nil
}