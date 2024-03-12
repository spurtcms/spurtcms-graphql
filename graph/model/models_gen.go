// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type AdditionalFields struct {
	Sections []Section `json:"sections,omitempty"`
	Fields   []Field   `json:"fields,omitempty"`
}

type Author struct {
	AuthorID         int       `json:"AuthorId"`
	FirstName        string    `json:"FirstName"`
	LastName         string    `json:"LastName"`
	Email            string    `json:"Email"`
	MobileNo         *string   `json:"MobileNo,omitempty"`
	IsActive         *int      `json:"IsActive,omitempty"`
	ProfileImage     *string   `json:"ProfileImage,omitempty"`
	ProfileImagePath *string   `json:"ProfileImagePath,omitempty"`
	CreatedOn        time.Time `json:"CreatedOn"`
	CreatedBy        int       `json:"CreatedBy"`
}

type CategoriesList struct {
	Categories []Category `json:"categories"`
	Count      int        `json:"count"`
}

type Category struct {
	ID           int        `json:"id"`
	CategoryName string     `json:"categoryName"`
	CategorySlug string     `json:"categorySlug"`
	Description  string     `json:"description"`
	ImagePath    string     `json:"imagePath"`
	CreatedOn    time.Time  `json:"createdOn"`
	CreatedBy    int        `json:"createdBy"`
	ModifiedOn   *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy   *int       `json:"modifiedBy,omitempty"`
	ParentID     int        `json:"parentId"`
}

type Channel struct {
	ID                 int        `json:"id"`
	ChannelName        string     `json:"channelName"`
	ChannelDescription string     `json:"channelDescription"`
	SlugName           string     `json:"slugName"`
	FieldGroupID       int        `json:"fieldGroupId"`
	IsActive           int        `json:"isActive"`
	CreatedOn          time.Time  `json:"createdOn"`
	CreatedBy          int        `json:"createdBy"`
	ModifiedOn         *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy         *int       `json:"modifiedBy,omitempty"`
}

type ChannelDetails struct {
	Channellist []Channel `json:"channellist"`
	Count       int       `json:"count"`
}

type ChannelEntries struct {
	ID               int               `json:"id"`
	Title            string            `json:"title"`
	Slug             string            `json:"slug"`
	Description      string            `json:"description"`
	UserID           int               `json:"userId"`
	ChannelID        int               `json:"channelId"`
	Status           int               `json:"status"`
	IsActive         int               `json:"isActive"`
	CreatedOn        time.Time         `json:"createdOn"`
	CreatedBy        int               `json:"createdBy"`
	ModifiedBy       *int              `json:"modifiedBy,omitempty"`
	ModifiedOn       *time.Time        `json:"modifiedOn,omitempty"`
	CoverImage       string            `json:"coverImage"`
	ThumbnailImage   string            `json:"thumbnailImage"`
	MetaTitle        string            `json:"metaTitle"`
	MetaDescription  string            `json:"metaDescription"`
	Keyword          string            `json:"keyword"`
	CategoriesID     string            `json:"categoriesId"`
	RelatedArticles  string            `json:"relatedArticles"`
	FeaturedEntry    *int              `json:"featuredEntry,omitempty"`
	Categories       [][]Category      `json:"categories"`
	AdditionalFields *AdditionalFields `json:"additionalFields,omitempty"`
	AuthorDetails    *Author           `json:"authorDetails"`
	MemberProfile    []MemberProfile   `json:"memberProfile,omitempty"`
}

type ChannelEntriesDetails struct {
	ChannelEntriesList []ChannelEntries `json:"channelEntriesList"`
	Count              int              `json:"count"`
}

type Field struct {
	FieldID          int            `json:"fieldId"`
	FieldName        string         `json:"fieldName"`
	FieldTypeID      int            `json:"fieldTypeId"`
	MandatoryField   int            `json:"mandatoryField"`
	OptionExist      int            `json:"optionExist"`
	CreatedOn        time.Time      `json:"createdOn"`
	CreatedBy        int            `json:"createdBy"`
	ModifiedOn       *time.Time     `json:"modifiedOn,omitempty"`
	ModifiedBy       *int           `json:"modifiedBY,omitempty"`
	FieldDesc        string         `json:"fieldDesc"`
	OrderIndex       int            `json:"orderIndex"`
	ImagePath        string         `json:"imagePath"`
	DatetimeFormat   *string        `json:"datetimeFormat,omitempty"`
	TimeFormat       *string        `json:"timeFormat,omitempty"`
	SectionParentID  *int           `json:"sectionParentId,omitempty"`
	CharacterAllowed *int           `json:"characterAllowed,omitempty"`
	FieldTypeName    string         `json:"fieldTypeName"`
	FieldValue       *FieldValue    `json:"fieldValue,omitempty"`
	FieldOptions     []FieldOptions `json:"fieldOptions,omitempty"`
}

type FieldOptions struct {
	ID          int        `json:"id"`
	OptionName  string     `json:"optionName"`
	OptionValue string     `json:"optionValue"`
	CreatedOn   time.Time  `json:"createdOn"`
	CreatedBy   int        `json:"createdBy"`
	ModifiedOn  *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy  *int       `json:"modifiedBY,omitempty"`
}

type FieldValue struct {
	ID         int        `json:"id"`
	FieldValue string     `json:"fieldValue"`
	CreatedOn  time.Time  `json:"createdOn"`
	CreatedBy  int        `json:"createdBy"`
	ModifiedOn *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy *int       `json:"modifiedBY,omitempty"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Member struct {
	ID               int           `json:"id"`
	FirstName        string        `json:"firstName"`
	LastName         string        `json:"lastName"`
	Email            string        `json:"email"`
	MobileNo         string        `json:"mobileNo"`
	IsActive         int           `json:"isActive"`
	ProfileImage     string        `json:"profileImage"`
	ProfileImagePath string        `json:"profileImagePath"`
	CreatedOn        time.Time     `json:"createdOn"`
	CreatedBy        int           `json:"createdBy"`
	ModifiedOn       *time.Time    `json:"modifiedOn,omitempty"`
	ModifiedBy       *int          `json:"modifiedBy,omitempty"`
	MemberGroupID    int           `json:"memberGroupId"`
	Group            []MemberGroup `json:"group,omitempty"`
}

type MemberDetails struct {
	FirstName        string  `json:"firstName"`
	LastName         string  `json:"lastName"`
	Mobile           string  `json:"mobile"`
	Email            string  `json:"email"`
	Password         string  `json:"password"`
	IsActive         *int    `json:"isActive,omitempty"`
	ProfileImage     *string `json:"profileImage,omitempty"`
	ProfileImagePath *string `json:"profileImagePath,omitempty"`
	Username         *string `json:"username,omitempty"`
	GroupID          *int    `json:"groupId,omitempty"`
}

type MemberGroup struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	IsActive    int        `json:"isActive"`
	CreatedOn   time.Time  `json:"createdOn"`
	CreatedBy   int        `json:"createdBy"`
	ModifiedOn  *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy  *int       `json:"modifiedBy,omitempty"`
}

type MemberProfile struct {
	MemberID        *int        `json:"memberId,omitempty"`
	ProfileName     *string     `json:"profileName,omitempty"`
	ProfileSlug     *string     `json:"profileSlug,omitempty"`
	ProfilePage     *string     `json:"profilePage,omitempty"`
	MemberDetails   interface{} `json:"memberDetails,omitempty"`
	CompanyName     *string     `json:"companyName,omitempty"`
	CompanyLocation *string     `json:"companyLocation,omitempty"`
	CompanyLogo     *string     `json:"companyLogo,omitempty"`
	About           *string     `json:"about,omitempty"`
	SeoTitle        *string     `json:"seoTitle,omitempty"`
	SeoDescription  *string     `json:"seoDescription,omitempty"`
	SeoKeyword      *string     `json:"seoKeyword,omitempty"`
}

type Page struct {
	ID          int        `json:"id"`
	PageName    string     `json:"pageName"`
	Content     string     `json:"content"`
	PagegroupID int        `json:"pagegroupId"`
	OrderIndex  int        `json:"orderIndex"`
	ParentID    int        `json:"parentId"`
	Status      string     `json:"status"`
	CreatedOn   time.Time  `json:"createdOn"`
	CreatedBy   int        `json:"created_by"`
	ModifiedOn  *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy  *int       `json:"modifiedBy,omitempty"`
}

type PageAndPageGroups struct {
	Pages      []Page      `json:"pages"`
	Subpages   []SubPage   `json:"subpages"`
	Pagegroups []PageGroup `json:"pagegroups"`
}

type PageGroup struct {
	ID            int        `json:"id"`
	PagegroupName string     `json:"pagegroupName"`
	OrderIndex    int        `json:"orderIndex"`
	CreatedOn     time.Time  `json:"createdOn"`
	CreatedBy     int        `json:"created_by"`
	ModifiedOn    *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy    *int       `json:"modifiedBy,omitempty"`
}

type Section struct {
	SectionID     *int       `json:"sectionId,omitempty"`
	SectionName   string     `json:"sectionName"`
	SectionTypeID int        `json:"sectionTypeId"`
	CreatedOn     time.Time  `json:"createdOn"`
	CreatedBy     int        `json:"createdBy"`
	ModifiedOn    *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy    *int       `json:"modifiedBY,omitempty"`
	OrderIndex    int        `json:"orderIndex"`
}

type Space struct {
	ID               int        `json:"id"`
	SpaceName        string     `json:"spaceName"`
	SpaceSlug        string     `json:"spaceSlug"`
	SpaceDescription string     `json:"spaceDescription"`
	ImagePath        string     `json:"imagePath"`
	LanguageID       int        `json:"languageId"`
	CreatedOn        time.Time  `json:"createdOn"`
	CreatedBy        int        `json:"createdBy"`
	ModifiedOn       *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy       *int       `json:"modifiedBy,omitempty"`
	CategoryID       int        `json:"categoryId"`
	Categories       []Category `json:"categories"`
}

type SpaceDetails struct {
	Spacelist []Space `json:"spacelist"`
	Count     int     `json:"count"`
}

type SubPage struct {
	ID          int        `json:"id"`
	SubpageName string     `json:"subpageName"`
	Conent      string     `json:"conent"`
	ParentID    int        `json:"parentId"`
	PageGroupID int        `json:"pageGroupId"`
	OrderIndex  int        `json:"orderIndex"`
	Status      string     `json:"status"`
	CreatedOn   time.Time  `json:"createdOn"`
	CreatedBy   int        `json:"created_by"`
	ModifiedOn  *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy  *int       `json:"modifiedBy,omitempty"`
}
