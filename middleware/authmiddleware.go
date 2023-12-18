package middleware

import (
	"context"
	"fmt"
	"gqlserver/controller"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	
	"github.com/spurtcms/spurtcms-core/member"
)

// Implement the AuthMiddleware function
func AuthMiddleware(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {

	c, _ := ctx.Value(controller.ContextKey).(*gin.Context)

	token := c.GetHeader("Authorization")

	// tokenString := strings.Replace(token, "Bearer ", "", 1)

	if token == "" {

		return "", fmt.Errorf("Unauthorized")

	}

	userid,roleid,err := member.VerifyToken(token,os.Getenv("JWT_SECRET"))

	if err != nil {

		return nil, fmt.Errorf("unauthorized: %v", err)
	}

	c.Set("userid",userid)

	c.Set("roleid",roleid)

	c.Set("token",token)

	return next(ctx)
}

