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
	AuthInstance                *newauth.Auth
	TeamAuthInstance            *team.Teams
	MemberAuthInstance          *memberPkg.Member
	EcomAuthInstance            *ecomPkg.Ecommerce
	EcomInstance                *ecomPkg.Ecommerce
	MemberInstance              *memberPkg.Member
	TeamInstance                *team.Teams
	// NewRole                  *role.PermissionConfig
	// ChannelConfig            *chn.Channel
	// CategoryConfig           *cat.Categories
	// MemberaccessConfig       *memaccess.AccessControl
)

func init(){

	AuthConfig()

	GetMemberInstance()

	GetMemberInstanceWithoutAuth()
	
	GetEcomInstance()

	GetEcomInstanceWithoutAuth()

	GetTeamInstance()

	GetTeamInstanceWithoutAuth()
}

// AuthCofing
func AuthConfig() *newauth.Auth {

	AuthInstance = newauth.AuthSetup(newauth.Config{
		SecretKey: os.Getenv("JWT_SECRET"),
		DB:        DB,
	})

	return AuthInstance
}

func GetMemberInstance() *memberPkg.Member {

	MemberAuthInstance = memberPkg.MemberSetup(memberPkg.Config{
		DB:               DB,
		AuthEnable:       true,
		PermissionEnable: true,
		Auth:             AuthInstance,
	})

	return MemberAuthInstance
}

func GetMemberInstanceWithoutAuth() *memberPkg.Member {

	MemberInstance = memberPkg.MemberSetup(memberPkg.Config{
		DB:               DB,
		Auth:             AuthInstance,
	})

	return MemberInstance
}

func GetEcomInstance() *ecomPkg.Ecommerce{

	EcomAuthInstance = ecomPkg.EcommerceSetup(ecomPkg.Config{
		AuthEnable: true,
		PermissionEnable: true,
		DB: DB,
		Auth: AuthInstance,
	})

	return EcomAuthInstance
}

func GetEcomInstanceWithoutAuth() *ecomPkg.Ecommerce{

	EcomInstance = ecomPkg.EcommerceSetup(ecomPkg.Config{
		AuthEnable: false,
		PermissionEnable: true,
		DB: DB,
		Auth: AuthInstance,
	})

	return EcomInstance
}

func GetTeamInstance() *team.Teams{

	TeamAuthInstance = teampkg.TeamSetup(teampkg.Config{
		DB: DB,
		AuthEnable: true,
		PermissionEnable: true,
		Auth: AuthInstance,
	})

	return TeamAuthInstance
}

func GetTeamInstanceWithoutAuth() *team.Teams{

	TeamInstance = teampkg.TeamSetup(teampkg.Config{
		DB: DB,
		Auth: AuthInstance,
	})

	return TeamInstance
}
