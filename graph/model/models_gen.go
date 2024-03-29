// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type ChannelDetails struct {
	Channellist []TblChannel `json:"channellist"`
	Count       int          `json:"count"`
}

type ChannelEntries struct {
	ChannelEntryList []TblChannelEntries `json:"channelEntryList"`
	Count            int                 `json:"count"`
}

type ChannelEntryDetails struct {
	ChannelEntryList *ChannelEntries    `json:"channelEntryList,omitempty"`
	ChannelEntry     *TblChannelEntries `json:"channelEntry,omitempty"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Member struct {
	ID               int           `json:"id"`
	UUID             string        `json:"uuid"`
	FirstName        string        `json:"firstName"`
	LastName         string        `json:"lastName"`
	Email            string        `json:"email"`
	MobileNo         string        `json:"mobileNo"`
	IsActive         int           `json:"isActive"`
	ProfileImage     string        `json:"profileImage"`
	ProfileImagePath string        `json:"profileImagePath"`
	LastLogin        int           `json:"lastLogin"`
	IsDeleted        int           `json:"isDeleted"`
	DeletedOn        *time.Time    `json:"deletedOn,omitempty"`
	DeletedBy        *int          `json:"deletedBy,omitempty"`
	CreatedOn        time.Time     `json:"createdOn"`
	CreatedDate      *string       `json:"createdDate,omitempty"`
	CreatedBy        int           `json:"createdBy"`
	ModifiedOn       *time.Time    `json:"modifiedOn,omitempty"`
	ModifiedBy       *int          `json:"modifiedBy,omitempty"`
	MemberGroupID    int           `json:"memberGroupId"`
	Group            []MemberGroup `json:"group,omitempty"`
	Password         string        `json:"password"`
	Username         string        `json:"username"`
	Otp              *int          `json:"otp,omitempty"`
	OtpExpiry        *time.Time    `json:"otpExpiry,omitempty"`
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
	IsDeleted   int        `json:"isDeleted"`
	CreatedOn   time.Time  `json:"createdOn"`
	CreatedBy   int        `json:"createdBy"`
	ModifiedOn  *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy  *int       `json:"modifiedBy,omitempty"`
	DateString  *string    `json:"dateString,omitempty"`
}

type TblCategory struct {
	ID           int        `json:"id"`
	CategoryName string     `json:"categoryName"`
	CategorySlug string     `json:"categorySlug"`
	Description  string     `json:"description"`
	ImagePath    string     `json:"imagePath"`
	CreatedOn    time.Time  `json:"createdOn"`
	CreatedBy    int        `json:"createdBy"`
	ModifiedOn   *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy   *int       `json:"modifiedBy,omitempty"`
	IsDeleted    int        `json:"isDeleted"`
	DeletedOn    *time.Time `json:"deletedOn,omitempty"`
	DeletedBy    *int       `json:"deletedBy,omitempty"`
	ParentID     int        `json:"parentId"`
}

type TblChannel struct {
	ID                 int        `json:"id"`
	ChannelName        string     `json:"channelName"`
	ChannelDescription string     `json:"channelDescription"`
	SlugName           string     `json:"slugName"`
	FieldGroupID       int        `json:"fieldGroupId"`
	IsActive           int        `json:"isActive"`
	IsDeleted          int        `json:"isDeleted"`
	CreatedOn          time.Time  `json:"createdOn"`
	CreatedBy          int        `json:"createdBy"`
	ModifiedOn         *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy         *int       `json:"modifiedBy,omitempty"`
}

type TblChannelEntries struct {
	ID              int             `json:"id"`
	Title           string          `json:"title"`
	Slug            string          `json:"slug"`
	Description     string          `json:"description"`
	UserID          int             `json:"userId"`
	ChannelID       int             `json:"channelId"`
	Status          int             `json:"status"`
	IsActive        int             `json:"isActive"`
	IsDeleted       int             `json:"isDeleted"`
	DeletedBy       *int            `json:"deletedBy,omitempty"`
	DeletedOn       *time.Time      `json:"deletedOn,omitempty"`
	CreatedOn       time.Time       `json:"createdOn"`
	CreatedBy       int             `json:"createdBy"`
	ModifiedBy      *int            `json:"modifiedBy,omitempty"`
	ModifiedOn      *time.Time      `json:"modifiedOn,omitempty"`
	CoverImage      string          `json:"coverImage"`
	ThumbnailImage  string          `json:"thumbnailImage"`
	MetaTitle       string          `json:"metaTitle"`
	MetaDescription string          `json:"metaDescription"`
	Keyword         string          `json:"keyword"`
	CategoriesID    string          `json:"categoriesId"`
	RelatedArticles string          `json:"relatedArticles"`
	Categories      [][]TblCategory `json:"categories" gorm:"-"`
}
