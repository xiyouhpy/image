package img

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xiyouhpy/image/config"
)

// Download 下载url的图片，返回下载文件名
func Download(strURL string) (string, error) {
	client := new(http.Client)
	client.Timeout = time.Second * 600
	rsp, err := client.Get(strURL)
	if err != nil {
		logrus.Warnf("client.Get err, url:%s, err:%s", strURL, err.Error())
		return "", err
	}
	defer rsp.Body.Close()

	dstFile := config.DownloadDir + filepath.Base(strURL)
	fileSize, _ := strconv.ParseInt(rsp.Header.Get("Content-Length"), 10, 32)
	if !isDownload(dstFile, fileSize) {
		file, fileErr := os.Create(dstFile)
		if fileErr != nil {
			logrus.Warnf("os.Create err, file:%s, err:%s", dstFile, fileErr.Error())
			return "", fileErr
		}
		defer file.Close()

		if _, err = io.Copy(file, rsp.Body); err != nil {
			logrus.Warnf("os.Copy err, err:%s", err.Error())
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
