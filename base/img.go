package base

import (
	"errors"
	"github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"
)

// ImgDecode 图片解码
func ImgDecode(imgPath string) (image.Image, error) {
	imgBin, imgErr := os.Open(imgPath)
	if imgErr != nil {
		logrus.Warnf("os.Open err, img_name:%s, err:%s", imgPath, imgErr.Error())
		return nil, imgErr
	}
	defer imgBin.Close()

	var imgDec image.Image
	fileType := strings.Replace(path.Ext(imgPath), ".", "", 1)
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

// ImgEncode 图片编码
func ImgEncode(imgPath string, srcRGBA *image.NRGBA) error {
	if _, imgErr := os.Stat(imgPath); os.IsExist(imgErr) {
		if imgErr = os.Remove(imgPath); imgErr != nil {
			logrus.Warnf("os.Remove err, img_name:%s, err:%s", imgPath, imgErr.Error())
			return imgErr
		}
	}
	imgNew, imgErr := os.Create(imgPath)
	if imgErr != nil {
		logrus.Warnf("os.Create err, img_name:%s, err:%s", imgPath, imgErr.Error())
		return imgErr
	}
	defer imgNew.Close()

	fileType := strings.Replace(path.Ext(imgPath), ".", "", 1)
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
