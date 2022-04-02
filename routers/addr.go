package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAddr(c *gin.Context) {
	RoomConfitTemplate
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  e.GetMsg(code),
	})
}
