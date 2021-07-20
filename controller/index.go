package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/xiyouhpy/image/base"
)

var ctx = context.Background()

func JsonRet(c *gin.Context, err error, data ...interface{}) {
	cause := errors.Cause(err)
	e, ok := cause.(base.Error)
	if !ok {
		e = base.ErrUnknownError
	}

	response := map[string]interface{}{
		"no":  e.Errno(),
		"msg": e.Error(),
	}
	if data != nil {
		response["data"] = data
	}
	c.JSON(http.StatusOK, response)
	return
}
