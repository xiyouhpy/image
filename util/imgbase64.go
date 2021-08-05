package util

import (
	"encoding/base64"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/xiyouhpy/image/base"
)

// ImageToBase64 对图片进行 base64 编码
func ImageToBase64(imgName string) (string, error) {
	if imgName == "" {
		logrus.Warnf("param err, src_name:%s", imgName)
		return "", base.ErrParamsError
	}
	if _, err := os.Stat(imgName); os.IsNotExist(err) {
		logrus.Warnf("file not exist, src_name:%s", imgName)
		return "", base.ErrParamsError
	}

	imgFile, fileErr := os.Open(imgName)
	if fileErr != nil {
		logrus.Warnf("os.Open err, src_name:%s, err:%s", imgName, fileErr.Error())
		return "", fileErr
	}
	imgData, dataErr := ioutil.ReadAll(imgFile)
	if dataErr != nil {
		logrus.Warnf("ioutil.ReadAll err, src_name:%s, err:%s", imgName, dataErr.Error())
		return "", dataErr
	}

	strBase64 := base64.StdEncoding.EncodeToString(imgData)

	return strBase64, nil
}
