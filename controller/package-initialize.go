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
	MemberAuthInstance *memberPkg.Member
	MemberInstance     *memberPkg.Member
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

	MemberAuthInstance = memberPkg.MemberSetup(memberPkg.Config{
		DB:               DB,
		AuthEnable:       true,
		PermissionEnable: false,
		Auth:             NewAuth,
	})

	return MemberAuthInstance
}

func GetMemberInstanceWithoutAuth() *memberPkg.Member {

	MemberInstance = memberPkg.MemberSetup(memberPkg.Config{
		DB:   DB,
		Auth: NewAuth,
	})

	return MemberInstance
}
