package config

import (
	"github.com/xiyouhpy/image/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// GetLogo 根据 logo名判断该文件是否存在并获取 logo 文件详细路径
func GetLogo(logName string) string {
	strFileName := util.LogoDir + logName
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
	strFileName := util.TtfDir + ttfName
	if !strings.Contains(strFileName, ".ttf") {
		strFileName = strFileName + ".ttf"
	}
	if _, err := os.Stat(strFileName); os.IsNotExist(err) {
		logrus.Errorf("getTtf Stat err, logo:%s, err:%s", strFileName, err.Error())
		return ""
	}

	return strFileName
}

// CleanFile 清理过期文件
func CleanFile(cleanDir string, expireTime int64) error {
	fileList, listErr := ioutil.ReadDir(cleanDir)
	if listErr != nil {
		logrus.Errorf("CleanFile ReadDir err, cleanDir:%s, err:%s", cleanDir, listErr.Error())
		return listErr
	}

	for _, file := range fileList {
		fileName := file.Name()
		if fileName != "." && fileName != ".." && file.ModTime().Unix() < expireTime {
			err := os.RemoveAll(filepath.Join(cleanDir, "/", fileName))
			logrus.Infof("CleanFile RemoveAll, cleanName:%s, err:%+v", fileName, err)
		}
	}

	return nil
}
