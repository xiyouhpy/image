package main

import (
	"github.com/xiyouhpy/image/base"
	"github.com/xiyouhpy/image/router"
)

func init() {
	// 获取日志信息
	base.Logger()
}

func main() {
	router.Server()
	return
}
