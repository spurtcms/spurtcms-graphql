package controller

import(

	"github.com/gin-gonic/gin"
)

func GetPlayGroundView(c *gin.Context){

	c.HTML(200,"query.html",nil)
	
}