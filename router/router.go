package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xiyouhpy/image/config"
	"github.com/xiyouhpy/image/controller"
)

// Server ...
func Server() {
	r := gin.Default()

	// 创建必要的目录
	initPath()

	// 路由注册和跳转
	registerService(r)

	// 服务监听端口
	r.Run(":8000")
}

// registerService ...
func registerService(r *gin.Engine) {
	// 图片压缩图片接口
	r.GET("/image/imgCompress", controller.ImgCompress)
}

// initPath ...
func initPath() bool {
	// 检查并创建 result 目录
	if _, err := os.Stat(config.ResultPath); os.IsNotExist(err) {
		err = os.MkdirAll(config.ResultPath, os.ModePerm)
		if err != nil {
			return false
		}
	}

	// 检查并创建 download 目录
	if _, err := os.Stat(config.DownloadDir); os.IsNotExist(err) {
		err = os.MkdirAll(config.DownloadDir, os.ModePerm)
		if err != nil {
			return false
		}
	}

	return true
}
