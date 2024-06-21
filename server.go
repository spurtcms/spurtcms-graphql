package main

import (
	"spurtcms-graphql/controller"
	ginhandler "spurtcms-graphql/ginHandler"
	"spurtcms-graphql/middleware"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const defaultPort string = "8080"

func main() {

	port := os.Getenv("PORT")

	if port == "" {

		port = defaultPort
	}

	r := gin.Default()

	r.Use(middleware.CorsMiddleware())

	r.Static("/public","./public")

	r.Static("/view","./view")

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

	r.GET("image-resize", controller.ImageResize)

	r.Run(":" + port)
}
