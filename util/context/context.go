package context

import (
	"github.com/gin-gonic/gin"
	"github.com/never-afk/Devloop/util"
	"net/http"
)

type wrap struct {
	*gin.Context
}

func (w *wrap) Success(data interface{}) {
	w.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
	})
	return
}

func (w *wrap) FailErr(err error) {

	w.JSON(http.StatusNotAcceptable, gin.H{
		"code":   10000,
		"errors": util.GetErrors(err),
	})
	return
}

func (w *wrap) FailMsg(msg string, code ...int) {
	statusCode := 10001
	if len(code) == 1 {
		statusCode = code[0]
	}
	w.JSON(http.StatusNotAcceptable, gin.H{
		"code": statusCode,
		"msg":  msg,
	})
	return
}

func New(c *gin.Context) *wrap {
	return &wrap{Context: c}
}
