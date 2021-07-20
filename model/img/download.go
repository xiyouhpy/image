package img

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// dataDir 下载文件目录
const dataDir = "./data/"

// Download 下载url的图片，返回下载文件名
func Download(strURL string) (string, error) {
	client := new(http.Client)
	client.Timeout = time.Second * 600
	rsp, err := client.Get(strURL)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	dstFile := dataDir + filepath.Base(strURL)
	fileSize, _ := strconv.ParseInt(rsp.Header.Get("Content-Length"), 10, 32)
	if !isDownload(dstFile, fileSize) {
		file, err := os.Create(dstFile)
		if err != nil {
			logrus.Warn("os.Create err, file:%s, err:%s", dstFile, err.Error())
			return "", err
		}
		defer file.Close()

		if _, err := io.Copy(file, rsp.Body); err != nil {
			logrus.Warn("os.Copy err, err:%s", err.Error())
			return "", err
		}
	}

	return dstFile, nil
}

// isDownload 判断下载的文件是否存在
func isDownload(fileName string, fileSize int64) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	if info.Size() != fileSize {
		os.Remove(fileName)
		return false
	}

	return true
}