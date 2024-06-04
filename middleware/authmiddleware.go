package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"spurtcms-graphql/controller"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"

	"github.com/spurtcms/pkgcore/auth"
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

	currentTime := time.Now().In(controller.TimeZone).Unix()

	currentTime1 := time.Now().Unix()

	fmt.Println("log", controller.TimeZone, currentTime, currentTime1)

	memberid, groupid, tokenType, err := auth.VerifyTokenWithExpiryTime(token, os.Getenv("JWT_SECRET"), currentTime)

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
