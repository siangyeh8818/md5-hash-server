package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/siangyeh8818/md5-hash-server/internal/middleware/jwt"
	v1 "github.com/siangyeh8818/md5-hash-server/internal/routers/api/v1"
	"github.com/siangyeh8818/gin.project.template/pkg/export"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	/*

		r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
		r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

		r.POST("/auth", api.GetAuth)

		r.POST("/upload", api.UploadImage)
	*/
	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//r.NoRoute(NoRouteHandler)

	//apiv1 := r.Group("/api/v1")
	r.POST("/hash", v1.UploadUrlfile)
	//apiv1.Use(jwt.JWT())
	//{

	//	apiv1.POST("/upload", v1.UploadUrlfile)

	//}

	return r
}
