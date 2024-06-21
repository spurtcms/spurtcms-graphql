package controller

import (
	"os"

	newauth "github.com/spurtcms/auth"
	memberPkg "github.com/spurtcms/member"
)

var (
	NewAuth *newauth.Auth
	// NewTeam            *team.Teams
	// NewRole            *role.PermissionConfig
	// ChannelConfig      *chn.Channel
	// CategoryConfig     *cat.Categories
	MemberConfig *memberPkg.Member
	// MemberaccessConfig *memaccess.AccessControl
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
