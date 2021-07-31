package watermark

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golang/freetype"
	"github.com/sirupsen/logrus"
	"github.com/xiyouhpy/image/base"
)

// FontInfo 定义添加的文字信息
type FontInfo struct {
	Size     float64 // 文字大小
	Message  string  // 文字内容
	Position int     // 文字存放位置
	Dx       int     // 文字x轴留白距离
	Dy       int     // 文字y轴留白距离
	R        uint8   // 文字颜色值RGBA中的R值
	G        uint8   // 文字颜色值RGBA中的G值
	B        uint8   // 文字颜色值RGBA中的B值
	A        uint8   // 文字颜色值RGBA中的A值
}

// ImgWatermark 图片压缩图片
// 		srcName 表示需要被压缩合唱的图片文件名
//		logoName 表示需要用来压缩合成的logo图片文件名
func ImgWatermark(srcName string, logoName string) (string, error) {
	if srcName == "" || logoName == "" {
		logrus.Warnf("ImgWatermark params err, src_name:%s, dst_name:%s", srcName, logoName)
		return "", errors.New("params err")
	}

	srcDec, srcErr := imgDecode(srcName)
	if srcErr != nil {
		logrus.Warnf("imgDecode err, img_name:%s, err:%s", srcName, srcErr.Error())
		return "", srcErr
	}
	logoDec, dstErr := imgDecode(logoName)
	if dstErr != nil {
		logrus.Warnf("imgDecode err, img_name:%s, err:%s", logoName, dstErr.Error())
		return "", dstErr
	}

	// 获取源图片参数信息
	srcBounds := srcDec.Bounds()
	// 设置水印图片坐标信息
	location := image.Pt(srcBounds.Min.X+10, srcBounds.Min.Y+10)
	// 设置水印图片RGBA信息
	srcRGBA := image.NewNRGBA(srcBounds)
	// 设置源图片参数信息
	draw.Draw(srcRGBA, srcBounds, srcDec, image.ZP, draw.Src)
	// 设置水印图片参数信息
	draw.Draw(srcRGBA, logoDec.Bounds().Add(location), logoDec, image.ZP, draw.Over)

	// 生成合成图片，统一使用 jpeg 后缀（空间占用比较小）
	md5 := base.GetMd5(srcName + logoName)
	newName := fmt.Sprintf("%simage_%d_%s", base.ImageDir, time.Now().Unix(), md5[len(md5)-20:]+".jpg")
	if imgErr := imgEncode(newName, srcRGBA); imgErr != nil {
		logrus.Warnf("imgEncode err, img_name:%s, err:%s", newName, imgErr.Error())
		return "", imgErr
	}

	return newName, nil
}

// TextWatermark 文字压缩图片
// 		srcName 表示需要被压缩合唱的图片文件名
//		ttfName 表示需要用来压缩合成的logo图片文件名
func (font *FontInfo) TextWatermark(srcName string, ttfName string) (string, error) {
	if srcName == "" || ttfName == "" {
		logrus.Warnf("TextWatermark params err, src_name:%s, ttf_name:%s", srcName, ttfName)
		return "", errors.New("params err")
	}
	srcImgDec, srcErr := imgDecode(srcName)
	if srcErr != nil {
		logrus.Warnf("TextWatermark decode err, img_name:%s, err:%s", srcName, srcErr.Error())
		return "", srcErr
	}

	// 获取源图片参数信息
	srcBounds := srcImgDec.Bounds()
	// 设置水印图片RGBA信息
	srcRGBA := image.NewNRGBA(srcBounds)
	// 设置背景（使用原图背景）
	for y := 0; y < srcRGBA.Bounds().Dy(); y++ {
		for x := 0; x < srcRGBA.Bounds().Dx(); x++ {
			srcRGBA.Set(x, y, srcImgDec.At(x, y))
		}
	}
	srcRGBA, srcErr = font.setTextWaterMark(srcRGBA, ttfName)
	if srcErr != nil {
		logrus.Warnf("setWaterMark err, fontinfo:%v, err:%s", font, srcErr.Error())
		return "", srcErr
	}

	// 生成合成图片，统一使用 jpeg 后缀（空间占用比较小）
	md5 := base.GetMd5(srcName + ttfName)
	newName := fmt.Sprintf("%simage_%d_%s", base.ImageDir, time.Now().Unix(), md5[len(md5)-20:]+".jpg")
	if imgErr := imgEncode(newName, srcRGBA); imgErr != nil {
		logrus.Warnf("imgEncode err, img_name:%s, err:%s", newName, imgErr.Error())
		return "", imgErr
	}

	return newName, nil
}

// setTextWaterMark 添加文字水印
func (font *FontInfo) setTextWaterMark(srcRGBA *image.NRGBA, ttfName string) (*image.NRGBA, error) {
	if srcRGBA == nil || ttfName == "" {
		logrus.Warnf("setTextWaterMark params err, ttf_name:%s", ttfName)
		return nil, errors.New("params err")
	}

	fontBytes, fontBytesErr := ioutil.ReadFile(ttfName)
	if fontBytesErr != nil {
		logrus.Warnf("ioutil.ReadFile err, ttf_name:%s, err:%s", ttfName, fontBytesErr.Error())
		return nil, fontBytesErr
	}
	fontParse, fontParseErr := freetype.ParseFont(fontBytes)
	if fontParseErr != nil {
		logrus.Warnf("freetype.ParseFont err, err:%s", fontParseErr.Error())
		return nil, fontParseErr
	}

	f := freetype.NewContext()
	// 设置屏幕每英寸的分辨率，建议72
	f.SetDPI(72)
	// 设置用于绘制文本的字体
	f.SetFont(fontParse)
	// 以磅为单位设置字体大小
	f.SetFontSize(font.Size)
	// 设置剪裁矩形以进行绘制
	f.SetClip(srcRGBA.Bounds())
	// 设置目标图像
	f.SetDst(srcRGBA)
	// 设置绘制操作的源图像
	f.SetSrc(image.NewUniform(color.RGBA{R: font.R, G: font.G, B: font.B, A: font.A}))
	// 设置水印文字出现位置
	pt := freetype.Pt((srcRGBA.Bounds().Dx()-len(font.Message)*4)/2, (srcRGBA.Bounds().Dy()-font.Dy)/2)
	if _, err := f.DrawString(font.Message, pt); err != nil {
		logrus.Warnf("DrawString err, info:%s, err:%s", font.Message, err.Error())
		return nil, err
	}

	return srcRGBA, nil
}

// imgDecode 图片解码
func imgDecode(imgName string) (image.Image, error) {
	imgBin, imgErr := os.Open(imgName)
	if imgErr != nil {
		logrus.Warnf("os.Open err, img_name:%s, err:%s", imgName, imgErr.Error())
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

// imgEncode 图片编码
func imgEncode(imgName string, srcRGBA *image.NRGBA) error {
	if _, imgErr := os.Stat(imgName); os.IsExist(imgErr) {
		if imgErr = os.Remove(imgName); imgErr != nil {
			logrus.Warnf("os.Remove err, img_name:%s, err:%s", imgName, imgErr.Error())
			return imgErr
		}
	}
	imgNew, imgErr := os.Create(imgName)
	if imgErr != nil {
		logrus.Warnf("os.Create err, img_name:%s, err:%s", imgName, imgErr.Error())
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
