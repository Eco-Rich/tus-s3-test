package route

import (
	"gin/tus"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	tusHandler := tus.InitTusHandler()
	r.Any("/files/*path", gin.WrapH(tusHandler))

	return r
}
