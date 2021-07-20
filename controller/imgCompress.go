package controller

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiyouhpy/image/base"
	"github.com/xiyouhpy/image/config"
	"github.com/xiyouhpy/image/model/img"
	"github.com/xiyouhpy/image/service/compress"
)

// ImgCompress 图片压缩接口
func ImgCompress(c *gin.Context) {
	srcName := c.DefaultQuery("src_name", "")
	newName := c.DefaultQuery("new_name", "")
	strIndex := c.DefaultQuery("compress_idx", "news_logo")
	dstName := getCompressImgFilePath(strIndex)
	if dstName == "" || srcName == "" {
		JsonRet(c, base.ErrParamsError)
		return
	}
	if newName == "" {
		newName = filepath.Base(srcName)
	}

	// 处理网络图片情况
	if strings.Contains(srcName, "https://") || strings.Contains(srcName, "http://") {
		var err error
		srcName, err = img.Download(srcName)
		if err != nil {
			JsonRet(c, base.ErrDownloadError)
			return
		}
	}

	// 执行图片压缩逻辑
	err := compress.ImgCompress(srcName, dstName, config.ResultPath+newName)
	if err != nil {
		JsonRet(c, base.ErrServiceError)
		return
	}

	JsonRet(c, base.ErrSuccess, config.ResultPath+newName)
	return
}

// getCompressImgFilePath 查找索引文件路径
func getCompressImgFilePath(strIndex string) string {
	strFileName := config.WaterPath + strIndex + ".png"
	if _, err := os.Stat(strFileName); os.IsNotExist(err) {
		return ""
	}

	return strFileName
}
