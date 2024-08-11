package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	engine := gin.Default()
	engine.GET("/", func(context *gin.Context) {
		result := gin.H{
			"message": "你好",
		}
		context.JSON(http.StatusOK, result)

	})

	engine.Run(":8888")

}
