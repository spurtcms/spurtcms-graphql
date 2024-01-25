package controller

import(

	"github.com/gin-gonic/gin"
)

func GetPlayGroundView(c *gin.Context){

	if AuthToken==""{

		c.HTML(200,"query.html",gin.H{"AuthToken":SpecialToken})

	}else{

		c.HTML(200,"query.html",gin.H{"AuthToken":AuthToken})
	}
	
}