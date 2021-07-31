package router

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/xiyouhpy/image/base"

	"github.com/gin-gonic/gin"
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
	r.GET("/image/imgWatermark", controller.ImgWatermark)
	// 文字压缩图片接口
	r.GET("/image/textWatermark", controller.TextWatermark)
}

// initPath ...
func initPath() bool {
	// 检查并创建 result 目录
	if _, err := os.Stat(base.ImageDir); os.IsNotExist(err) {
		err = os.MkdirAll(base.ImageDir, os.ModePerm)
		if err != nil {
			logrus.Warnf("initPath create err, dir:%s, err:%s", base.ImageDir, err.Error())
			return false
		}
	}

	// 检查并创建 download 目录
	if _, err := os.Stat(base.TmpDir); os.IsNotExist(err) {
		err = os.MkdirAll(base.TmpDir, os.ModePerm)
		if err != nil {
			logrus.Warnf("initPath create err, dir:%s, err:%s", base.TmpDir, err.Error())
			return false
		}
	}

	return true
}
