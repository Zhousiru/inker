package api

import (
	"net/http"

	"github.com/Zhousiru/inker/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var debugMode = true

// Init 初始化 API 服务
func Init() {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(cors.Default())

	normal := r.Group("/")
	{
		normal.GET("/login", login)
		normal.GET("/getArticle", getArticle)
		normal.GET("/paginateHome", paginateHome)
		normal.GET("/paginateSearch", paginateSearch)
		normal.GET("/getFile", getFile)
	}

	authorized := r.Group("/manage").Use(authMiddleware())
	{
		authorized.GET("/updateUser", updateUser)
		authorized.POST("/newArticle", newArticle)
		authorized.GET("/deleteArticle", deleteArticle)
		authorized.POST("/updateArticle", updateArticle)
		authorized.POST("/uploadFile", uploadFile)
		authorized.GET("/deleteFile", deleteFile)
		authorized.GET("/paginateFile", paginateFile)
	}

	r.NoRoute(func(c *gin.Context) {
		response(c, http.StatusNotFound, "invalid method", nil, nil)
	})

	r.Run(config.Conf.APIAddr)
}
