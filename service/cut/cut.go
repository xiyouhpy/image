package cut

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xiyouhpy/image/base"
)

// 图片按照原图宽高等比例缩放、图片按照给定宽高等比例缩放

// ClipInfo 定义裁剪尺寸信息
type ClipInfo struct {
	SrcName string // 裁剪图片名（包含路径信息）
	DstName string // 生成图片名（包含路径信息）
	X0      int    // 图片裁剪的 x0 位置
	Y0      int    // 图片裁剪的 y0 位置
	X1      int    // 图片裁剪的 x1 位置
	Y1      int    // 图片裁剪的 y1 位置
	Quality int    // 图片裁剪的质量
}

// ImgClip 图片裁剪接口，给定指定坐标值进行截取
func (c *ClipInfo) ImgClip() error {
	imgSrc, srcErr := os.Open(c.SrcName)
	if srcErr != nil {
		logrus.Warnf("os.Open err, src_name:%s, err:%s", c.SrcName, srcErr.Error())
		return srcErr
	}
	defer imgSrc.Close()

	imgDst, dstErr := os.Create(c.DstName)
	if dstErr != nil {
		logrus.Warnf("os.Create err, dst_name:%s, err:%s", c.DstName, dstErr.Error())
		return dstErr
	}
	defer imgDst.Close()

	var imgErr error
	fileType := strings.Replace(path.Ext(c.SrcName), ".", "", 1)
	switch fileType {
	case "png":
		imgErr = c.pngClip(imgSrc, imgDst)
		break
	case "jpg", "jpeg":
		imgErr = c.jpegClip(imgSrc, imgDst)
		break
	default:
		imgErr = base.ErrUnknownError
		break
	}

	return imgErr
}

// jpegClip jpeg 图片裁剪处理
func (c *ClipInfo) jpegClip(imgSrc io.Reader, dstImg io.Writer) error {
	if imgSrc == nil || dstImg == nil {
		logrus.Warnf("param err, src_name:%s, dst_name:%s", c.SrcName, c.DstName)
		return base.ErrParamsError
	}

	// decode 图片
	srcDec, srcErr := jpeg.Decode(imgSrc)
	if srcErr != nil {
		logrus.Warnf("jpeg.Decode err, dst_name:%s, err:%s", c.SrcName, srcErr.Error())
		return srcErr
	}

	// jpeg 裁剪处理图片
	img := srcDec.(*image.YCbCr)
	subImg := img.SubImage(image.Rect(c.X0, c.Y0, c.X1, c.Y1)).(*image.YCbCr)

	// encode 图片
	srcErr = jpeg.Encode(dstImg, subImg, &jpeg.Options{Quality: c.Quality})
	if srcErr != nil {
		logrus.Warnf("jpeg.Encode err, dst_name:%s, err:%s", c.SrcName, srcErr.Error())
		return srcErr
	}

	return nil
}

// pngClip png 图片裁剪处理
func (c *ClipInfo) pngClip(imgSrc io.Reader, dstImg io.Writer) error {
	if imgSrc == nil || dstImg == nil {
		logrus.Warnf("param err, src_name:%s, dst_name:%s", c.SrcName, c.DstName)
		return base.ErrParamsError
	}

	// decode 图片
	srcDec, srcErr := png.Decode(imgSrc)
	if srcErr != nil {
		logrus.Warnf("png.Decode err, src_name:%s, err:%s", c.SrcName, srcErr.Error())
		return srcErr
	}

	// png 裁剪处理图片
	img := srcDec.(*image.NRGBA)
	subImg := img.SubImage(image.Rect(c.X0, c.Y0, c.X1, c.Y1)).(*image.NRGBA)

	// encode 图片
	srcErr = png.Encode(dstImg, subImg)
	if srcErr != nil {
		logrus.Warnf("png.Encode err, src_name:%s, err:%s", c.SrcName, srcErr.Error())
		return srcErr
	}

	return nil
}
