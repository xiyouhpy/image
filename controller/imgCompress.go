package controller

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiyouhpy/image/base"
	"github.com/xiyouhpy/image/model/img"
	"github.com/xiyouhpy/image/service/compress"
)

// DstPath 压缩的水印图片目录
const DstPath = "./config/water_pic/"

// NewPath 压缩的结果图片目录
const NewPath = "./data/result/"

// ImgCompress 图片压缩接口
func ImgCompress(c *gin.Context) {
	srcName := c.DefaultQuery("src_name", "")
	newName := c.DefaultQuery("new_name", "")
	strIndex := c.DefaultQuery("compress_idx", "news_log")

	// 判断并获取压缩图片
	dstName := getCompressImgFilePath(strIndex)
	if dstName == "" {
		JsonRet(c, base.ErrParamsError)
		return
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
	err := compress.ImgCompress(srcName, dstName, NewPath+newName)
	if err != nil {
		JsonRet(c, base.ErrServiceError)
		return
	}

	JsonRet(c, base.ErrSuccess)
	return
}

// getCompressImgFilePath 查找索引文件路径
func getCompressImgFilePath(strIndex string) string {
	strFileName := strIndex + ".png"
	if _, err := os.Stat(strFileName); os.IsNotExist(err) {
		return ""
	}

	return DstPath + strFileName
}
