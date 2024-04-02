package middleware

import (
	"context"
	"fmt"
	"spurtcms-graphql/controller"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"

	"github.com/spurtcms/pkgcore/member"
)

// Implement the AuthMiddleware function
func AuthMiddleware(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {

	c, _ := ctx.Value(controller.ContextKey).(*gin.Context)

	token := c.GetHeader("Authorization")

	// tokenString := strings.Replace(token, "Bearer ", "", 1)

	if token == "" {

		return "", fmt.Errorf("Unauthorized")

	}

	if token == controller.SpecialToken {

		c.Set("token",token)

		return next(ctx)

	}

	memberid,groupid,err := member.VerifyToken(token,os.Getenv("JWT_SECRET"))

	if err != nil {

		return nil, fmt.Errorf("Unauthorized: %v", err)
	}

	c.Set("memberid",memberid)

	c.Set("groupid",groupid)

	c.Set("token",token)

	return next(ctx)
}

