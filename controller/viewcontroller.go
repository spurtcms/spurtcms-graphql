package controller

import (

	"github.com/gin-gonic/gin"
)

func GetPlayGroundView(c *gin.Context){

	c.HTML(200,"query.html",gin.H{"AuthToken":SpecialToken})

}

func GetDocumentationView(c *gin.Context){

	c.HTML(200,"index.html",nil)

}