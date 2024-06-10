package controller

import (
	"os"

	newauth "github.com/spurtcms/auth"
	ecomPkg "github.com/spurtcms/ecommerce"
	memberPkg "github.com/spurtcms/member"
	"github.com/spurtcms/team"
	teampkg "github.com/spurtcms/team"
)

var (
	NewAuth *newauth.Auth
	TeamConfig             *team.Teams
	// NewRole             *role.PermissionConfig
	// ChannelConfig       *chn.Channel
	// CategoryConfig      *cat.Categories
	MemberConfig           *memberPkg.Member
	// MemberaccessConfig  *memaccess.AccessControl
	EcomConfig             *ecomPkg.Ecommerce
)

// AuthCofing
func AuthConfig() *newauth.Auth {

	NewAuth = newauth.AuthSetup(newauth.Config{
		SecretKey: os.Getenv("JWT_SECRET"),
		DB:        DB,
	})

	return NewAuth
}

func GetMemberInstance() *memberPkg.Member {

	MemberConfig = memberPkg.MemberSetup(memberPkg.Config{
		DB:               DB,
		AuthEnable:       true,
		PermissionEnable: true,
		Auth:             NewAuth,
	})

	return MemberConfig
}

func GetMemberInstanceWithoutAuth() *memberPkg.Member {

	MemberConfig = memberPkg.MemberSetup(memberPkg.Config{
		DB:               DB,
		AuthEnable:       false,
		PermissionEnable: true,
		Auth:             NewAuth,
	})

	return MemberConfig
}

func GetEcomInstance() *ecomPkg.Ecommerce{

	EcomConfig = ecomPkg.EcommerceSetup(ecomPkg.Config{
		AuthEnable: true,
		PermissionEnable: true,
		DB: DB,
		Auth: NewAuth,
	})

	return EcomConfig
}

func GetEcomInstanceWithoutAuth() *ecomPkg.Ecommerce{

	EcomConfig = ecomPkg.EcommerceSetup(ecomPkg.Config{
		AuthEnable: false,
		PermissionEnable: true,
		DB: DB,
		Auth: NewAuth,
	})

	return EcomConfig
}

func GetTeamInstance() *team.Teams{

	TeamConfig = teampkg.TeamSetup(teampkg.Config{
		DB: DB,
		AuthEnable: true,
		PermissionEnable: true,
		Auth: NewAuth,
	})

	return TeamConfig
}
