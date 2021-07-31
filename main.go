package main

import (
	"time"

	"github.com/xiyouhpy/image/base"
	"github.com/xiyouhpy/image/config"
	"github.com/xiyouhpy/image/router"
)

func init() {
	// 获取日志信息
	base.Logger()
}

func main() {
	// 清理过期文件
	go func() {
		for {
			intNow := time.Now().Unix() - 3600
			config.CleanFile(base.ImageDir, intNow)
			config.CleanFile(base.TmpDir, intNow)
			time.Sleep(time.Minute * 10)
		}
	}()

	router.Server()
	return
}
