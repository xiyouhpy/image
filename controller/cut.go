package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiyouhpy/image/base"
	"github.com/xiyouhpy/image/service/cut"
	"github.com/xiyouhpy/image/service/deliver"
	"github.com/xiyouhpy/image/util"
)

// ImgCut 图片裁剪接口
func ImgCut(c *gin.Context) {
	// 该参数由 x0,y0,x1,y1 四个值按照这个顺序组成，每个值之间使用英文逗号区分
	strXYSite := c.DefaultQuery("xy_site", "")
	srcName := c.DefaultQuery("src_name", "")
	strQuality := c.DefaultQuery("quality", "")

	var arrXYSite []int
	arrTempSite := strings.Split(strXYSite, ",")
	for _, strTmpSite := range arrTempSite {
		intTmpSite, err := strconv.Atoi(strTmpSite)
		if err != nil || intTmpSite < 0 {
			JsonRet(c, base.ErrParamsError)
			return
		}
		arrXYSite = append(arrXYSite, intTmpSite)
	}
	intQuality, err := strconv.Atoi(strQuality)
	if err != nil || len(arrXYSite) != 4 || arrXYSite[2] <= 0 || arrXYSite[3] <= 0 || srcName == "" {
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

	// 执行图片裁剪逻辑
	md5 := base.GetMd5(srcName)
	dstName := fmt.Sprintf("%scut_%d_%s", util.ImgCutDir, time.Now().Unix(), md5[len(md5)-20:]+".jpg")
	imgCut := cut.ClipInfo{
		SrcName: srcName,
		DstName: dstName,
		X0:      arrXYSite[0],
		Y0:      arrXYSite[1],
		X1:      arrXYSite[2],
		Y1:      arrXYSite[3],
		Quality: intQuality,
	}
	err = imgCut.ImgClip()
	if err != nil {
		JsonRet(c, base.ErrServiceError)
		return
	}

	JsonRet(c, base.ErrSuccess, imgCut.DstName)
	return
}
