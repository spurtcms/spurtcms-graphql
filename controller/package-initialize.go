package controller

import (
	"os"
	authPkg "github.com/spurtcms/auth"
	ecomPkg "github.com/spurtcms/ecommerce"
	memberPkg "github.com/spurtcms/member"
	teamPkg "github.com/spurtcms/team"
    chanPkg  "github.com/spurtcms/channels"
)

var (
	AuthInstance                *authPkg.Auth
	TeamAuthInstance            *teamPkg.Teams
	MemberAuthInstance          *memberPkg.Member
	EcomAuthInstance            *ecomPkg.Ecommerce
	EcomInstance                *ecomPkg.Ecommerce
	MemberInstance              *memberPkg.Member
	TeamInstance                *teamPkg.Teams
	ChannelInstance             *chanPkg.Channel
	ChannelAuthInstance         *chanPkg.Channel
	// NewRole                  *role.PermissionConfig
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
func AuthConfig() *authPkg.Auth {

	AuthInstance = authPkg.AuthSetup(authPkg.Config{
		SecretKey: os.Getenv("JWT_SECRET"),
		DB:        DB,
	})

	return AuthInstance
}

func GetMemberInstance() *memberPkg.Member {

	MemberAuthInstance = memberPkg.MemberSetup(memberPkg.Config{
		DB:               DB,
		AuthEnable:       true,
		PermissionEnable: false,
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
		PermissionEnable: false,
		DB: DB,
		Auth: AuthInstance,
	})

	return EcomAuthInstance
}

func GetEcomInstanceWithoutAuth() *ecomPkg.Ecommerce{

	EcomInstance = ecomPkg.EcommerceSetup(ecomPkg.Config{
		DB: DB,
		Auth: AuthInstance,
	})

	return EcomInstance
}

func GetTeamInstance() *teamPkg.Teams{

	TeamAuthInstance = teamPkg.TeamSetup(teamPkg.Config{
		DB: DB,
		AuthEnable: true,
		PermissionEnable: false,
		Auth: AuthInstance,
	})

	return TeamAuthInstance
}

func GetTeamInstanceWithoutAuth() *teamPkg.Teams{

	TeamInstance = teamPkg.TeamSetup(teamPkg.Config{
		DB: DB,
		Auth: AuthInstance,
	})

	return TeamInstance
}

func GetChannelInstance() *chanPkg.Channel{

	ChannelAuthInstance =  chanPkg.ChannelSetup(chanPkg.Config{
		DB: DB,
		AuthEnable: true,
		PermissionEnable: false,
		Auth: AuthInstance,
	})

	return ChannelAuthInstance
}

func GetChannelInstanceWithoutAuth() *chanPkg.Channel{

	ChannelInstance =  chanPkg.ChannelSetup(chanPkg.Config{
		DB: DB,
		AuthEnable: true,
		PermissionEnable: false,
		Auth: AuthInstance,
	})

	return ChannelInstance
}
