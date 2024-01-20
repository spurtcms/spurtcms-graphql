package main

import (
	"gqlserver/controller"
	ginhandler "gqlserver/ginHandler"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const defaultPort string = "8080"

func main() {

	er := godotenv.Load()

	if er != nil {

		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {

		port = defaultPort
	}

	r := gin.Default()

	r.LoadHTMLGlob("view/*")

	r.POST("/query", ginhandler.GraphQLHandler())

	r.GET("/play", controller.GetPlayGroundView)

	r.GET("/", ginhandler.PlaygroundHandler())

	r.Run(":" + port)
}
