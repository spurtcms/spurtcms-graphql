// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type AdditionalFields struct {
	Sections []Section `json:"sections,omitempty"`
	Fields   []Field   `json:"fields,omitempty"`
}

type Author struct {
	AuthorID         int       `json:"AuthorId" gorm:"column:id"`
	FirstName        string    `json:"FirstName"`
	LastName         string    `json:"LastName"`
	Email            string    `json:"Email"`
	MobileNo         *string   `json:"MobileNo,omitempty"`
	IsActive         *int      `json:"IsActive,omitempty"`
	ProfileImagePath *string   `json:"ProfileImagePath,omitempty"`
	CreatedOn        time.Time `json:"CreatedOn"`
	CreatedBy        int       `json:"CreatedBy"`
}

type CartSummary struct {
	SubTotal       string `json:"subTotal"`
	ShippingAmount int    `json:"shippingAmount"`
	TotalTax       string `json:"totalTax"`
	TotalCost      string `json:"totalCost"`
	TotalQuantity  int    `json:"totalQuantity"`
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
	FeaturedEntry    int               `json:"featuredEntry"`
	ViewCount        int               `json:"viewCount"`
	Categories       [][]Category      `json:"categories" gorm:"-"`
	AdditionalFields *AdditionalFields `json:"additionalFields,omitempty" gorm:"-"`
	AuthorDetails    Author            `json:"authorDetails" gorm:"-"`
	MemberProfile    MemberProfile     `json:"memberProfile" gorm:"-"`
	Author           *string           `json:"author,omitempty"`
	SortOrder        *int              `json:"sortOrder,omitempty"`
	CreateTime       *time.Time        `json:"createTime,omitempty"`
	PublishedTime    *time.Time        `json:"publishedTime,omitempty"`
	ReadingTime      *int              `json:"readingTime,omitempty"`
	Tags             *string           `json:"tags,omitempty"`
	Excerpt          *string           `json:"excerpt,omitempty"`
	ImageAltTag      *string           `json:"imageAltTag,omitempty"`
}

type ChannelEntriesDetails struct {
	ChannelEntriesList []ChannelEntries `json:"channelEntriesList"`
	Count              int              `json:"count"`
}

type ClaimData struct {
	ProfileName   string `json:"profileName"`
	ProfileSlug   string `json:"profileSlug"`
	WorkMail      string `json:"workMail"`
	CompanyNumber string `json:"companyNumber"`
	PersonName    string `json:"personName"`
}

type EcomOrderedProductDetails struct {
	EcommerceProduct EcommerceProduct `json:"EcommerceProduct"`
	OrderStatuses    []OrderStatus    `json:"OrderStatuses"`
}

type EcommerceCartDetails struct {
	CartList    []EcommerceProduct `json:"cartList"`
	CartSummary CartSummary        `json:"cartSummary"`
	Count       int                `json:"Count"`
}

type EcommerceOrder struct {
	ID              int        `json:"id"`
	UUID            string     `json:"uuid"`
	CustomerID      int        `json:"customerId"`
	Status          string     `json:"status"`
	ShippingAddress string     `json:"shippingAddress"`
	IsDeleted       int        `json:"isDeleted"`
	CreatedOn       time.Time  `json:"createdOn"`
	ModifiedOn      *time.Time `json:"modifiedOn,omitempty"`
	Price           int        `json:"price"`
	Tax             int        `json:"tax"`
	TotalCost       int        `json:"totalCost"`
}

type EcommerceProduct struct {
	ID                 int        `json:"id"`
	CategoriesID       int        `json:"categoriesId"`
	ProductName        string     `json:"productName"`
	ProductSlug        string     `json:"productSlug"`
	ProductDescription string     `json:"productDescription"`
	ProductImagePath   string     `json:"productImagePath"`
	ProductYoutubePath *string    `json:"productYoutubePath,omitempty"`
	ProductVimeoPath   *string    `json:"productVimeoPath,omitempty"`
	Sku                string     `json:"sku"`
	Tax                int        `json:"tax"`
	Totalcost          int        `json:"totalcost"`
	IsActive           int        `json:"isActive"`
	CreatedOn          time.Time  `json:"createdOn"`
	CreatedBy          int        `json:"createdBy"`
	ModifiedOn         *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy         *int       `json:"modifiedBy,omitempty"`
	IsDeleted          int        `json:"isDeleted"`
	DeletedBy          *int       `json:"deletedBy,omitempty"`
	DeletedOn          *time.Time `json:"deletedOn,omitempty"`
	ViewCount          *int       `json:"viewCount,omitempty"`
	ProductPrice       int        `json:"productPrice"`
	DiscountPrice      *int       `json:"discountPrice,omitempty"`
	SpecialPrice       *int       `json:"specialPrice,omitempty"`
	ProductImageArray  []string   `json:"productImageArray,omitempty"`
	CartID             int        `json:"cartId"`
	ProductID          int        `json:"productId"`
	CustomerID         int        `json:"customerId"`
	Quantity           int        `json:"quantity"`
	CartCreatedOn      time.Time  `json:"cartCreatedOn"`
	CartModifiedOn     *time.Time `json:"cartModifiedOn,omitempty"`
	CartIsDeleted      int        `json:"cartIsDeleted"`
	CartDeletedOn      *time.Time `json:"cartDeletedOn,omitempty"`
	OrderID            *int       `json:"orderId,omitempty"`
	OrderUniqueID      *string    `json:"orderUniqueId,omitempty"`
	OrderQuantity      *int       `json:"orderQuantity,omitempty"`
	OrderPrice         *int       `json:"orderPrice,omitempty"`
	OrderTax           *int       `json:"orderTax,omitempty"`
	OrderStatus        *string    `json:"orderStatus,omitempty"`
	OrderCustomer      *int       `json:"orderCustomer,omitempty"`
	OrderTime          *time.Time `json:"orderTime,omitempty"`
	PaymentMode        *string    `json:"paymentMode,omitempty"`
	ShippingDetails    *string    `json:"shippingDetails,omitempty"`
}

type EcommerceProducts struct {
	ProductList []EcommerceProduct `json:"productList"`
	Count       int                `json:"count"`
}

type EmailTemplate struct {
	ID              int        `json:"id"`
	TemplateName    string     `json:"templateName"`
	TemplateSlug    string     `json:"templateSlug"`
	TemplateSubject string     `json:"templateSubject"`
	TemplateMessage string     `json:"templateMessage"`
	IsActive        int        `json:"IsActive"`
	CreatedOn       time.Time  `json:"createdOn"`
	CreatedBy       int        `json:"createdBy"`
	ModifiedOn      *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy      *int       `json:"modifiedBy,omitempty"`
	IsDeleted       int        `json:"isDeleted"`
	DeletedOn       *time.Time `json:"deletedOn,omitempty"`
	DeletedBy       *int       `json:"deletedBy,omitempty"`
	IsDefault       *int       `json:"isDefault,omitempty"`
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

type Job struct {
	ID             int        `json:"id"`
	CategoriesID   int        `json:"categoriesId"`
	Category       *Category `json:"category,omitempty" gorm:"foreignKey:CategoriesID;references:ID"`
	JobTitle       string     `json:"jobTitle"`
	JobSlug        string     `json:"jobSlug"`
	JobDescription string     `json:"jobDescription"`
	JobLocation    string     `json:"jobLocation"`
	JobType        string     `json:"jobType"`
	Education      string     `json:"education"`
	Department     *string    `json:"department,omitempty"`
	Experience     *string    `json:"experience,omitempty"`
	Salary         string     `json:"salary"`
	CreatedOn      time.Time  `json:"createdOn"`
	CreatedBy      int        `json:"createdBy"`
	IsDeleted      *int       `json:"isDeleted,omitempty"`
	DeletedOn      *time.Time `json:"deletedOn,omitempty"`
	DeletedBy      *int       `json:"deletedBy,omitempty"`
	Keyword        *string    `json:"keyword,omitempty"`
	Skill          string     `json:"skill"`
	MinimumYears   int        `json:"minimumYears"`
	MaximumYears   int        `json:"maximumYears"`
	PostedDate     time.Time  `json:"postedDate"`
	ValidThrough   time.Time  `json:"validThrough"`
	Status         int        `json:"status"`
	ModifiedOn     *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy     *int       `json:"modifiedBy,omitempty"`
}

type JobFilter struct {
	JobTitle     graphql.Omittable[*string] `json:"jobTitle,omitempty"`
	JobLocation  graphql.Omittable[*string] `json:"jobLocation,omitempty"`
	CategoryID   graphql.Omittable[*int]    `json:"categoryId,omitempty"`
	CategorySlug graphql.Omittable[*string] `json:"categorySlug,omitempty"`
	KeyWord      graphql.Omittable[*string] `json:"keyWord,omitempty"`
	MinimumYears graphql.Omittable[*int]    `json:"minimumYears,omitempty"`
	MaximumYears graphql.Omittable[*int]    `json:"maximumYears,omitempty"`
	DatePosted   graphql.Omittable[*string] `json:"datePosted,omitempty"`
}

type JobsList struct {
	Jobs  []Job `json:"jobs"`
	Count int   `json:"count"`
}

type LoginDetails struct {
	MemberProfileData MemberProfile `json:"memberProfileData"`
	Token             string        `json:"token"`
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
	Group            []MemberGroup `json:"group,omitempty" gorm:"-"`
	Password         *string       `json:"password,omitempty"`
	Username         *string       `json:"username,omitempty"`
	StorageType      *string       `json:"storageType,omitempty"`
}

type MemberDetails struct {
	FirstName        string                     `json:"firstName"`
	LastName         graphql.Omittable[*string] `json:"lastName,omitempty"`
	Mobile           graphql.Omittable[*string] `json:"mobile,omitempty"`
	Email            string                     `json:"email"`
	Password         graphql.Omittable[*string] `json:"password,omitempty"`
	IsActive         graphql.Omittable[*int]    `json:"isActive,omitempty"`
	ProfileImage     graphql.Omittable[*string] `json:"profileImage,omitempty"`
	ProfileImagePath graphql.Omittable[*string] `json:"profileImagePath,omitempty"`
	Username         graphql.Omittable[*string] `json:"username,omitempty"`
	GroupID          graphql.Omittable[*int]    `json:"groupId,omitempty"`
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
	ID              *int        `json:"id,omitempty"`
	MemberID        *int        `json:"memberId,omitempty"`
	ProfileName     *string     `json:"profileName,omitempty"`
	ProfileSlug     *string     `json:"profileSlug,omitempty"`
	ProfilePage     *string     `json:"profilePage,omitempty"`
	MemberDetails   interface{} `json:"memberDetails,omitempty" gorm:"column:member_details;type:jsonb"`
	CompanyName     *string     `json:"companyName,omitempty"`
	CompanyLocation *string     `json:"companyLocation,omitempty"`
	CompanyLogo     *string     `json:"companyLogo,omitempty"`
	About           *string     `json:"about,omitempty"`
	SeoTitle        *string     `json:"seoTitle,omitempty"`
	SeoDescription  *string     `json:"seoDescription,omitempty"`
	SeoKeyword      *string     `json:"seoKeyword,omitempty"`
	Linkedin        *string     `json:"linkedin,omitempty"`
	Twitter         *string     `json:"twitter,omitempty"`
	Website         *string     `json:"website,omitempty"`
	CreatedBy       *int        `json:"createdBy,omitempty"`
	CreatedOn       *time.Time  `json:"createdOn,omitempty"`
	ModifiedOn      *time.Time  `json:"modifiedOn,omitempty"`
	ModifiedBy      *int        `json:"modifiedBy,omitempty"`
	ClaimStatus     *int        `json:"claimStatus,omitempty"`
	StorageType     *string     `json:"storageType,omitempty"`
	IsActive        *int        `json:"isActive,omitempty"`
}

type MemberSettings struct {
	ID                int        `json:"id"`
	AllowRegistration int        `json:"allowRegistration"`
	MemberLogin       string     `json:"memberLogin"`
	ModifiedBy        *int       `json:"modifiedBy,omitempty"`
	ModifiedOn        *time.Time `json:"modifiedOn,omitempty"`
	NotificationUsers string     `json:"notificationUsers"`
}

type Mutation struct {
}

type OrderProductDetails struct {
	ID        int `json:"id"`
	OrderID   int `json:"orderId"`
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
	Price     int `json:"price"`
	Tax       int `json:"tax"`
}

type OrderStatus struct {
	ID          int       `json:"id"`
	OrderID     int       `json:"orderId"`
	OrderStatus string    `json:"orderStatus"`
	CreatedBy   int       `json:"createdBy"`
	CreatedOn   time.Time `json:"createdOn"`
}

type OrderStatusNames struct {
	ID          int        `json:"id"`
	Status      string     `json:"status"`
	Description *string    `json:"description,omitempty"`
	IsActive    int        `json:"isActive"`
	CreatedBy   int        `json:"createdBy"`
	CreatedOn   time.Time  `json:"createdOn"`
	ModifiedBy  *int       `json:"modifiedBy,omitempty"`
	ModifiedOn  *time.Time `json:"modifiedOn,omitempty"`
	IsDeleted   int        `json:"isDeleted"`
}

type OrderSummary struct {
	SubTotal       string                  `json:"subTotal"`
	ShippingAmount graphql.Omittable[*int] `json:"shippingAmount,omitempty"`
	TotalTax       string                  `json:"totalTax"`
	TotalCost      string                  `json:"totalCost"`
	TotalQuantity  int                     `json:"totalQuantity"`
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

type ProductFilter struct {
	ReleaseDate   graphql.Omittable[*string]  `json:"releaseDate,omitempty"`
	StartingPrice graphql.Omittable[*int]     `json:"startingPrice,omitempty"`
	EndingPrice   graphql.Omittable[*int]     `json:"endingPrice,omitempty"`
	CategoryName  graphql.Omittable[*string]  `json:"categoryName,omitempty"`
	CategoryID    graphql.Omittable[*int]     `json:"categoryId,omitempty"`
	StarRatings   graphql.Omittable[*float64] `json:"starRatings,omitempty"`
	SearchKeyword graphql.Omittable[*string]  `json:"searchKeyword,omitempty"`
}

type ProductPricing struct {
	ID        int       `json:"id"`
	PriceID   int       `json:"priceId"`
	Sku       string    `json:"sku"`
	Priority  int       `json:"priority"`
	Price     int       `json:"price"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Type      string    `json:"type"`
}

type ProductSort struct {
	Price     graphql.Omittable[*int] `json:"price,omitempty"`
	Date      graphql.Omittable[*int] `json:"date,omitempty"`
	ViewCount graphql.Omittable[*int] `json:"viewCount,omitempty"`
}

type ProfileData struct {
	CompanyName     string                     `json:"companyName"`
	ProfileName     string                     `json:"profileName"`
	ProfileSlug     string                     `json:"profileSlug"`
	CompanyLocation graphql.Omittable[*string] `json:"companyLocation,omitempty"`
	CompanyLogo     graphql.Omittable[*string] `json:"companyLogo,omitempty"`
	About           graphql.Omittable[*string] `json:"about,omitempty"`
	Website         graphql.Omittable[*string] `json:"website,omitempty"`
	Twitter         graphql.Omittable[*string] `json:"twitter,omitempty"`
	Linkedin        graphql.Omittable[*string] `json:"linkedin,omitempty"`
	CompanyProfile  graphql.Omittable[*string] `json:"companyProfile,omitempty"`
	SeoTitle        graphql.Omittable[*string] `json:"seoTitle,omitempty"`
	SeoDescription  graphql.Omittable[*string] `json:"seoDescription,omitempty"`
	SeoKeyword      graphql.Omittable[*string] `json:"seoKeyword,omitempty"`
}

type Query struct {
}

type RequireData struct {
	AuthorDetails    graphql.Omittable[*bool] `json:"authorDetails,omitempty"`
	Categories       graphql.Omittable[*bool] `json:"categories,omitempty"`
	MemberProfile    graphql.Omittable[*bool] `json:"memberProfile,omitempty"`
	AdditionalFields graphql.Omittable[*bool] `json:"additionalFields,omitempty"`
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

type TblEcommerceCart struct {
	ID         int        `json:"id"`
	ProductID  int        `json:"productId"`
	CustomerID int        `json:"customerId"`
	Quantity   int        `json:"quantity"`
	CreatedOn  time.Time  `json:"createdOn"`
	ModifiedOn *time.Time `json:"modifiedOn,omitempty"`
	IsDeleted  int        `json:"isDeleted"`
	DeletedOn  *time.Time `json:"deletedOn,omitempty"`
}

type ApplicationInput struct {
	Name           string                     `json:"name"`
	EmailID        string                     `json:"emailId"`
	MobileNo       int                        `json:"mobileNo"`
	JobType        string                     `json:"jobType"`
	Gender         string                     `json:"gender"`
	Location       string                     `json:"location"`
	Education      string                     `json:"education"`
	Graduation     int                        `json:"graduation"`
	CompanyName    graphql.Omittable[*string] `json:"companyName,omitempty"`
	Experience     int                        `json:"experience"`
	Skills         string                     `json:"skills"`
	ApplicantImage graphql.Upload             `json:"applicantImage"`
	CurrentSalary  graphql.Omittable[*int]    `json:"currentSalary,omitempty"`
	ExpectedSalary int                        `json:"expectedSalary"`
	Resume         graphql.Upload             `json:"resume"`
}

type CustomerDetails struct {
	ID               int        `json:"id"`
	FirstName        string     `json:"firstName"`
	LastName         *string    `json:"lastName,omitempty"`
	MobileNo         string     `json:"mobileNo"`
	Email            string     `json:"email"`
	Username         string     `json:"username"`
	Password         string     `json:"password"`
	IsActive         int        `json:"isActive"`
	ProfileImage     *string    `json:"profileImage,omitempty"`
	ProfileImagePath *string    `json:"profileImagePath,omitempty"`
	CreatedOn        time.Time  `json:"createdOn"`
	CreatedBy        int        `json:"createdBy"`
	ModifiedOn       *time.Time `json:"modifiedOn,omitempty"`
	ModifiedBy       *int       `json:"modifiedBy,omitempty"`
	IsDeleted        *int       `json:"IsDeleted,omitempty"`
	DeletedOn        *time.Time `json:"DeletedOn,omitempty"`
	HouseNo          *string    `json:"houseNo,omitempty" gorm:"-"`
	Area             *string    `json:"Area,omitempty" gorm:"-"`
	City             *string    `json:"city,omitempty"`
	State            *string    `json:"state,omitempty"`
	Country          *string    `json:"country,omitempty"`
	ZipCode          *string    `json:"zipCode,omitempty"`
	StreetAddress    *string    `json:"streetAddress,omitempty"`
	MemberID         *int       `json:"memberId,omitempty"`
}

type CustomerInput struct {
	FirstName     string                     `json:"firstName"`
	LastName      graphql.Omittable[*string] `json:"lastName,omitempty"`
	MobileNo      graphql.Omittable[*string] `json:"mobileNo,omitempty"`
	Email         string                     `json:"email"`
	Username      graphql.Omittable[*string] `json:"username,omitempty"`
	Password      graphql.Omittable[*string] `json:"password,omitempty"`
	IsActive      graphql.Omittable[*int]    `json:"isActive,omitempty"`
	ProfileImage  graphql.Omittable[*string] `json:"profileImage,omitempty"`
	City          graphql.Omittable[*string] `json:"city,omitempty"`
	State         graphql.Omittable[*string] `json:"state,omitempty"`
	Country       graphql.Omittable[*string] `json:"country,omitempty"`
	ZipCode       graphql.Omittable[*string] `json:"zipCode,omitempty"`
	StreetAddress graphql.Omittable[*string] `json:"streetAddress,omitempty"`
}

type OrderFilter struct {
	Status         graphql.Omittable[*string]  `json:"status,omitempty"`
	StartingPrice  graphql.Omittable[*int]     `json:"startingPrice,omitempty"`
	EndingPrice    graphql.Omittable[*int]     `json:"endingPrice,omitempty"`
	StartingDate   graphql.Omittable[*string]  `json:"startingDate,omitempty"`
	EndingDate     graphql.Omittable[*string]  `json:"endingDate,omitempty"`
	CategoryName   graphql.Omittable[*string]  `json:"categoryName,omitempty"`
	CategoryID     graphql.Omittable[*int]     `json:"categoryId,omitempty"`
	StarRatings    graphql.Omittable[*float64] `json:"starRatings,omitempty"`
	SearchKeyword  graphql.Omittable[*string]  `json:"searchKeyword,omitempty"`
	OrderID        graphql.Omittable[*string]  `json:"orderId,omitempty"`
	UpcomingOrders graphql.Omittable[*int]     `json:"upcomingOrders,omitempty"`
	OrderHistory   graphql.Omittable[*int]     `json:"orderHistory,omitempty"`
}

type OrderPayment struct {
	ID          int    `json:"id"`
	OrderID     int    `json:"orderId"`
	PaymentMode string `json:"paymentMode"`
}

type OrderProduct struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
	Price     int `json:"price"`
	Tax       int `json:"tax"`
	TotalCost int `json:"totalCost"`
}

type OrderSort struct {
	Price graphql.Omittable[*int] `json:"price,omitempty"`
	Date  graphql.Omittable[*int] `json:"date,omitempty"`
}

func(Category) TableName() string{

    return "tbl_categories"
}