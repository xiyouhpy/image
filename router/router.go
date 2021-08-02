package router

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xiyouhpy/image/base"
	"github.com/xiyouhpy/image/config"
	"github.com/xiyouhpy/image/controller"
)

// tmpFileExpireTime 临时文件过期时间设置
const tmpFileExpireTime = 3600

// Server ...
func Server() {
	r := gin.Default()

	// 创建必要的目录
	initPath()

	// 清理过期文件
	go cleanPath()

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

// cleanPath ...
func cleanPath() {
	ticker := time.NewTicker(time.Minute * time.Duration(10))
	for range ticker.C {
		intNow := time.Now().Unix() - tmpFileExpireTime
		config.CleanFile(base.ImageDir, intNow)
		config.CleanFile(base.TmpDir, intNow)
	}
}
