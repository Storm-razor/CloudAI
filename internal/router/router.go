package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wwwzy/CloudAI/internal/controller"
	"github.com/wwwzy/CloudAI/internal/middlerware"
)

func SetUpRouters(r *gin.Engine,
	uc *controller.UserController,
	fc *controller.FileController) {
	api := r.Group("/api")
	{
		publicUser := api.Group("/user")
		{
			publicUser.POST("/register", uc.Register)
			publicUser.POST("/login", uc.Login)
		}

		file := api.Group("/files")
		file.Use(middlerware.JWTAuth())
		{
			file.POST("/upload", fc.Upload)
			file.GET("/page", fc.PageList)
			file.GET("/download", fc.Download)
			file.DELETE("/delete", fc.Delete)
			file.POST("/folder", fc.CreateFolder)
			file.POST("/move", fc.BatchMove)
			file.GET("/search", fc.Search)
			file.PUT("/rename", fc.Rename)
			file.GET("/path", fc.GetPath)
			file.GET("/idPath", fc.GetIDPath)
		}
	}
}
