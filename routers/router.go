package routers

import (
	"gindemo/middleware/jwt"
	"gindemo/pkg/setting"
	"gindemo/routers/api"
	v1 "gindemo/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/auth", api.GetAuth)

	apiv1 := r.Group("/api/v1")
	{
		apiv1.Use(jwt.JWT())

		apiv1.GET("/tag/:id", v1.GetArticle)

		apiv1.GET("/tags", v1.GetArticles)

		apiv1.POST("/tags", v1.AddArticle)

		apiv1.PUT("/tags/:id", v1.EditArticle)

		apiv1.DELETE("/tags/:id", v1.DeleteArticle)

	}
	return r
}
