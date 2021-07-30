package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiyouhpy/image/base"
	"github.com/xiyouhpy/image/config"
	"github.com/xiyouhpy/image/service/deliver"
	"github.com/xiyouhpy/image/service/watermark"
)

// ImgWatermark 图片压缩接口
func ImgWatermark(c *gin.Context) {
	srcName := c.DefaultQuery("src_name", "")
	logoName := c.DefaultQuery("logo_name", "news_logo")
	logoName = config.GetLogo(logoName)
	if logoName == "" || srcName == "" {
		JsonRet(c, base.ErrParamsError)
		return
	}

	// 处理网络图片情况
	if strings.Contains(srcName, "https://") || strings.Contains(srcName, "http://") {
		tmpName, err := deliver.Download(srcName)
		if err != nil {
			JsonRet(c, base.ErrDownloadError)
			return
		}
		srcName = tmpName
	}

	// 执行图片压缩逻辑
	newName, err := watermark.ImgWatermark(srcName, logoName)
	if err != nil {
		JsonRet(c, base.ErrServiceError)
		return
	}

	JsonRet(c, base.ErrSuccess, newName)
	return
}

// TextWatermark 文字压缩接口
func TextWatermark(c *gin.Context) {
	srcName := c.DefaultQuery("src_name", "")
	message := c.DefaultQuery("message", "")
	ttfName := c.DefaultQuery("ttf_name", "msyh")
	ttfName = config.GetTtf(ttfName)
	if message == "" || srcName == "" {
		JsonRet(c, base.ErrParamsError)
		return
	}

	// 处理网络图片情况
	if strings.Contains(srcName, "https://") || strings.Contains(srcName, "http://") {
		var err error
		tmpName, err := deliver.Download(srcName)
		if err != nil {
			JsonRet(c, base.ErrDownloadError)
			return
		}
		srcName = tmpName
	}

	// 执行图片压缩逻辑
	fontInfo := &watermark.FontInfo{Size: 20, Message: message, Position: 0, Dx: 0, Dy: 0, R: 255, G: 140, B: 0, A: 250}
	newName, err := fontInfo.TextWatermark(srcName, ttfName)
	if err != nil {
		JsonRet(c, base.ErrServiceError)
		return
	}

	JsonRet(c, base.ErrSuccess, newName)
	return
}
