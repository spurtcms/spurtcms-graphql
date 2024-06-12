package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"spurtcms-graphql/controller"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
)

// Implement the AuthMiddleware function
func AuthMiddleware(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {

	c, ok := ctx.Value(controller.ContextKey).(*gin.Context)

	if !ok {

		controller.ErrorLog.Printf("Auth context error: %v", ok)
	}

	token := c.GetHeader("Authorization")

	// tokenString := strings.Replace(token, "Bearer ", "", 1)

	if token == "" {

		err := errors.New("unauthorized access")

		controller.ErrorLog.Printf("invalid token error: %s", err)

		c.AbortWithError(http.StatusUnauthorized, err)

		return "", err

	}

	if token == controller.SpecialToken {

		c.Set("token", token)

		return next(ctx)

	}

	memberid, groupid, tokenType, err := controller.AuthInstance.MemberVerifyToken(token, os.Getenv("JWT_SECRET"))

	if err != nil {

		controller.ErrorLog.Printf("Verify token error: %s", err)

		c.AbortWithError(http.StatusUnauthorized, err)

		return "", err
	}

	c.Set("memberid", memberid)

	c.Set("groupid", groupid)

	c.Set("token", token)

	c.Set("tokenType", tokenType)

	return next(ctx)
}
