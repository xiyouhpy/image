package deliver

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xiyouhpy/image/base"

	"github.com/sirupsen/logrus"
)

// downloadTimeOut 下载超时设置
const downloadTimeOut = time.Second * 60

// Download 下载url的图片，返回下载文件名
func Download(strURL string) (string, error) {
	client := new(http.Client)
	client.Timeout = downloadTimeOut
	rsp, err := client.Get(strURL)
	if err != nil {
		logrus.Warnf("client.Get err, url:%s, err:%s", strURL, err.Error())
		return "", err
	}
	defer rsp.Body.Close()

	dstFile := base.TmpDir + filepath.Base(strURL)
	if !strings.Contains(dstFile, ".") {
		dstFile = base.TmpDir + filepath.Base(strURL) + ".jpg"
	}
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

	return dstFile, nil
}
