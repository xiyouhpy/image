package controller

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/xiyouhpy/image/base"
	"github.com/xiyouhpy/image/config"
	"github.com/xiyouhpy/image/service/watermark"
)

// ImgWatermark 图片压缩接口
func ImgWatermark(c *gin.Context) {
	srcName := c.DefaultQuery("src_name", "")
	newName := c.DefaultQuery("new_name", "")
	strIndex := c.DefaultQuery("watermark_idx", "news_logo")
	dstName := getWmFilePath(strIndex, "pic")
	if dstName == "" || srcName == "" {
		JsonRet(c, base.ErrParamsError)
		return
	}
	if newName == "" {
		newName = base.GetMd5(srcName) + filepath.Base(srcName)
	}
	if !strings.Contains(newName, ".") {
		newName = newName + ".jpg"
	}

	// 处理网络图片情况
	if strings.Contains(srcName, "https://") || strings.Contains(srcName, "http://") {
		var err error
		srcName, err = base.Download(srcName, newName)
		if err != nil {
			JsonRet(c, base.ErrDownloadError)
			return
		}
	}

	// 执行图片压缩逻辑
	err := watermark.ImgWatermark(srcName, dstName, config.ResultDir+newName)
	if err != nil {
		JsonRet(c, base.ErrServiceError)
		return
	}

	JsonRet(c, base.ErrSuccess, config.ResultDir+newName)
	return
}

// TextWatermark 文字压缩接口
func TextWatermark(c *gin.Context) {
	srcName := c.DefaultQuery("src_name", "")
	newName := c.DefaultQuery("new_name", "")
	strMessage := c.DefaultQuery("message", "")
	strIndex := c.DefaultQuery("ttf_name", "")
	ttfName := getWmFilePath(strIndex, "text")
	if strMessage == "" || srcName == "" {
		JsonRet(c, base.ErrParamsError)
		return
	}
	if newName == "" {
		newName = base.GetMd5(srcName) + filepath.Base(srcName)
	}
	if !strings.Contains(newName, ".") {
		newName = newName + ".jpg"
	}

	// 处理网络图片情况
	if strings.Contains(srcName, "https://") || strings.Contains(srcName, "http://") {
		var err error
		srcName, err = base.Download(srcName, newName)
		if err != nil {
			JsonRet(c, base.ErrDownloadError)
			return
		}
	}

	// 执行图片压缩逻辑
	fontInfo := &watermark.FontInfo{Size: 10, Message: strMessage, Position: 0, Dx: 0, Dy: 0, R: 255, G: 140, B: 0, A: 250}
	err := fontInfo.TextWatermark(srcName, ttfName, config.ResultDir+newName)
	if err != nil {
		JsonRet(c, base.ErrServiceError)
		return
	}

	JsonRet(c, base.ErrSuccess, "file://"+config.ResultDir+newName)
	return
}

// getWmFilePath 获取压缩的文件路径
func getWmFilePath(strIndex string, WatermarkType string) string {
	var strFileName string

	switch WatermarkType {
	case "pic":
		strFileName = config.WaterPicDir + strIndex + ".png"
		break
	case "text":
		strFileName = config.WaterTtfDir + strIndex + ".ttf"
		break
	default:
		break
	}

	if _, err := os.Stat(strFileName); os.IsNotExist(err) {
		logrus.Warnf("os.Stat err, file:%s, err:%s", strFileName, err.Error())
		return ""
	}

	return strFileName
}
