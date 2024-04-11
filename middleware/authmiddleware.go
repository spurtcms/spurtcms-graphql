package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"spurtcms-graphql/controller"
	"time"

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

		err:= errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized,err)

		return "", err

	}

	if token == controller.SpecialToken {

		c.Set("token",token)

		return next(ctx)

	}

	currentTime := time.Now().In(controller.TimeZone).Unix()

	memberid,groupid,err :=  member.VerifyTokenWithExpiryTime(token,os.Getenv("JWT_SECRET"),currentTime)
	
	if err != nil {

		c.AbortWithError(http.StatusUnauthorized,err)
		
		return "", err 
	}

	c.Set("memberid",memberid)

	c.Set("groupid",groupid)

	c.Set("token",token)

	return next(ctx)
}

