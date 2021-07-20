package compress

import (
	"errors"
	"github.com/sirupsen/logrus"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"
)

// ImgCompress 图片压缩图片
// srcImgName 表示需要被压缩合唱的图片文件名，dstImgName 表示需要用来压缩合成的logo图片文件名，newImgName 表示生成压缩合成的文件名
func ImgCompress(srcImgName string, dstImgName string, newImgName string) error {
	srcImgDec, srcErr := imgDecode(srcImgName)
	if srcErr != nil {
		logrus.Warn("imgDecode err, img_name:%s, err:%s", srcImgName, srcErr.Error())
		return srcErr
	}
	dstImgDec, dstErr := imgDecode(dstImgName)
	if dstErr != nil {
		logrus.Warn("imgDecode err, img_name:%s, err:%s", dstImgName, dstErr.Error())
		return dstErr
	}

	// 获取源图片参数信息
	srcBounds := srcImgDec.Bounds()
	// 设置水印图片坐标信息
	location := image.Pt(srcBounds.Min.X, srcBounds.Min.Y)
	// 设置水印图片RGBA信息
	srcRGBA := image.NewNRGBA(srcBounds)
	// 设置源图片参数信息
	draw.Draw(srcRGBA, srcBounds, srcImgDec, image.ZP, draw.Src)
	// 设置水印图片参数信息
	draw.Draw(srcRGBA, dstImgDec.Bounds().Add(location), dstImgDec, image.ZP, draw.Over)

	// 生成合成图片，统一使用 jpeg 后缀（空间占用比较小）
	if imgErr := imgEncode(newImgName, srcRGBA); imgErr != nil {
		logrus.Warn("imgEncode err, img_name:%s, err:%s", newImgName, imgErr.Error())
		return imgErr
	}
	return nil
}

// imgDecode 图片解码
// imgName 表示需要解码的图片文件名
// imgDec 表示返回该文件解码的数据信息
func imgDecode(imgName string) (image.Image, error) {
	imgBin, imgErr := os.Open(imgName)
	if imgErr != nil {
		logrus.Warn("os.Open err, img_name:%s, err:%s", imgName, imgErr.Error())
		return nil, imgErr
	}
	defer imgBin.Close()

	var imgDec image.Image
	fileType := strings.Replace(path.Ext(imgName), ".", "", 1)
	switch fileType {
	case "png":
		imgDec, imgErr = png.Decode(imgBin)
		break
	case "jpg", "jpeg":
		imgDec, imgErr = jpeg.Decode(imgBin)
		break
	default:
		imgErr = errors.New("img decode err")
		break
	}

	return imgDec, imgErr
}

// imgEncode 图片编码；
// imgName 表示编码的图片文件名，srcRGBA 表示编码的数据内容信息
func imgEncode(imgName string, srcRGBA *image.NRGBA) error {
	if _, imgErr := os.Stat(imgName); os.IsExist(imgErr) {
		if imgErr = os.Remove(imgName); imgErr != nil {
			logrus.Warn("os.Remove err, img_name:%s, err:%s", imgName, imgErr.Error())
			return imgErr
		}
	}
	imgNew, imgErr := os.Create(imgName)
	if imgErr != nil {
		logrus.Warn("os.Create err, img_name:%s, err:%s", imgName, imgErr.Error())
		return imgErr
	}
	defer imgNew.Close()

	fileType := strings.Replace(path.Ext(imgName), ".", "", 1)
	switch fileType {
	case "png":
		imgErr = png.Encode(imgNew, srcRGBA)
		break
	case "jpg", "jpeg":
		imgErr = jpeg.Encode(imgNew, srcRGBA, &jpeg.Options{Quality: 100})
		break
	default:
		imgErr = errors.New("img encode err")
		break
	}

	return imgErr
}
