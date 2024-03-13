package main

import (
	"gqlserver/controller"
	ginhandler "gqlserver/ginHandler"
	"gqlserver/middleware"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	r.Use(middleware.CorsMiddleware())

	r.Static("/public","./public")

	var htmlfiles []string

	filepath.Walk("./", func(path string, info os.FileInfo, err error) error {

		if strings.HasSuffix(path, ".html") {

			htmlfiles = append(htmlfiles, path)

		}

		return nil
	})

	r.LoadHTMLFiles(htmlfiles...)

	r.GET("/apidocs",controller.GetDocumentationView)

	r.POST("/query", ginhandler.GraphQLHandler())

	r.GET("/play", controller.GetPlayGroundView)

	r.GET("/", ginhandler.PlaygroundHandler())

	r.Run(":" + port)
}
