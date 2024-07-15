package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"spurtcms-graphql/graph/model"
)

type MailConfig struct {
	Emails       []string
	MailUsername string
	MailPassword string
	SmtpPort     string
	SmtpHost     string
	Subject      string
	TimeOut      time.Duration
}

type MailImages struct {
	Owndesk   string
	Twitter   string
	Facebook  string
	LinkedIn  string
	Youtube   string
	Instagram string
}

type SocialMedias struct {
	Linkedin  string
	Twitter   string
	Facebook  string
	Instagram string
	Youtube   string
}

type StorageType struct {
	Id           int
	Local        string
	Aws          datatypes.JSONMap `gorm:"type:jsonb"`
	Azure        datatypes.JSONMap `gorm:"type:jsonb"`
	Drive        datatypes.JSONMap `gorm:"type:jsonb"`
	SelectedType string
}

type EmailConfiguration struct {
	Id           int
	SmtpConfig   datatypes.JSONMap `gorm:"type:jsonb"`
	SelectedType string
}

type TblGraphqlSettings struct {
	Id          int       `gorm:"primaryKey;auto_increment;"`
	TokenName   string    `gorm:"type:character varying"`
	Description string    `gorm:"type:character varying"`
	Duration    string    `gorm:"type:character varying"`
	CreatedBy   int       `gorm:"type:integer"`
	CreatedOn   time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy  int       `gorm:"type:integer;DEFAULT:NULL"`
	ModifiedOn  time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy   int       `gorm:"type:integer;DEFAULT:NULL"`
	DeletedOn   time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	IsDeleted   int       `gorm:"type:integer;DEFAULT:0"`
	Token       string    `gorm:"type:character varying"`
	ExpiryTime  time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
}

type ModelConfig struct {
	DB *gorm.DB
}

func (model ModelConfig) GetApiSettings(apikey string, graphqlsetting *TblGraphqlSettings) error {

	if err := model.DB.Debug().Model(TblGraphqlSettings{}).Where("is_deleted = 0 and token = ?", apikey).First(&graphqlsetting).Error; err != nil {

		return err
	}

	return nil
}

func (model ModelConfig) GetStorageType(storageType *StorageType) error {

	if err := model.DB.Debug().Table("tbl_storage_types").First(&storageType).Error; err != nil {

		return err
	}

	return nil
}

func (model ModelConfig) GetEmailConfig(emailConfig *EmailConfiguration) error {

	if err := model.DB.Debug().Table("tbl_email_configurations").First(&emailConfig).Error; err != nil {

		return err
	}

	return nil
}

func (model ModelConfig) GetEmailTemplate(templateSlug string,emailTemplate *model.EmailTemplate) error {

	if err := model.DB.Debug().Table("tbl_email_templates").Where("is_deleted=0 and template_slug = ?", templateSlug).First(&emailTemplate).Error; err != nil {

		return err
	}

	return nil
}
