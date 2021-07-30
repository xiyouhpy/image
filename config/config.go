package config

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xiyouhpy/image/base"
)

// GetLogo 根据 logo名判断该文件是否存在并获取 logo 文件详细路径
func GetLogo(logName string) string {
	strFileName := base.LogoDir + logName
	if !strings.Contains(strFileName, ".png") {
		strFileName = strFileName + ".png"
	}
	if _, err := os.Stat(strFileName); os.IsNotExist(err) {
		logrus.Errorf("getLogo Stat err, logo:%s, err:%s", strFileName, err.Error())
		return ""
	}

	return strFileName
}

// GetTtf 根据 ttf名判断该文件是否存在并获取 ttf 文件详细路径
func GetTtf(ttfName string) string {
	strFileName := base.TtfDir + ttfName
	if !strings.Contains(strFileName, ".ttf") {
		strFileName = strFileName + ".ttf"
	}
	if _, err := os.Stat(strFileName); os.IsNotExist(err) {
		logrus.Errorf("getTtf Stat err, logo:%s, err:%s", strFileName, err.Error())
		return ""
	}

	return strFileName
}
