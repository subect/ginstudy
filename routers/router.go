package routers

import (
	"gindemo/middleware/jwt"
	"gindemo/pkg/setting"
	"gindemo/pkg/upload"
	"gindemo/routers/api"
	v1 "gindemo/routers/api/v1"
	_ "github.com/EDDYCJY/go-gin-example/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

type RoomTemplate struct {
	Pusher_common_client struct {
		EtcdUrls []string `yaml:"etcd_urls"`
	}
	Generatecenter_common_client struct {
		EtcdUrls []string `yaml:"etcd_urls"`
	}
	Common_client struct {
		EtcdUrls []string `yaml:"etcd_urls"`
	}
	Room struct {
		ListenPort  string   `yaml:"listen_port"`
		RpcBasePath string   `yaml:"rpc_base_path"`
		GameType    string   `yaml:"game_type"`
		IsTestEnv   string   `yaml:"is_test_env"`
		EtcdUrls    []string `yaml:"etcd_urls"`
	}
	Log struct {
		Level string `yaml:"level"`
	}
	Pprof struct {
		Enabled bool   `yaml:"enabled"`
		Address string `yaml:"address"`
	}
	Report struct {
		Address      string `yaml:"address"`
		Level        string `yaml:"level"`
		KafkaTopic   string `yaml:"kafka_topic"`
		KafkaBrokers string `yaml:"kafka_brokers"`
	}
}

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/auth", api.GetAuth)

	r.GET("/nacos", GetAddr)

	r.POST("/upload", api.UploadImage)

	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

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
