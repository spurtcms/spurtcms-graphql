package main

import (
	ginhandler "gqlserver/ginHandler"
	"os"

	"github.com/gin-gonic/gin"
)

const defaultPort = "8080"

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		
		port = defaultPort
	}

	r :=gin.Default()

	r.POST("/query",ginhandler.GraphQLHandler())

	r.Run(":"+port)
}
