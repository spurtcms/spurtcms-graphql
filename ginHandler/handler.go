package ginhandler

import (
	"context"
	"gqlserver/controller"
	"gqlserver/graph"
	"gqlserver/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)


func GraphQLHandler() gin.HandlerFunc{

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(),Directives: graph.DirectiveRoot{Auth: middleware.AuthMiddleware}}))

	return func(c *gin.Context){

		gincontext := c

		ctx := context.WithValue(c.Request.Context(),controller.ContextKey,gincontext)

		srv.ServeHTTP(c.Writer,c.Request.WithContext(ctx))
	}

}

func PlaygroundHandler() gin.HandlerFunc{

	h := playground.Handler("GraphQL playground", "/query")

	return func(c *gin.Context){
		
		h.ServeHTTP(c.Writer,c.Request)
	}

}

