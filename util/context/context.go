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

func (w *wrap) FailErr(code int, err error) {
	w.JSON(http.StatusOK, gin.H{
		"code":   code,
		"errors": util.GetErrors(err),
	})
	return
}

func (w *wrap) FailMsg(code int, msg string) {
	w.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
	return
}

func New(c *gin.Context) *wrap {
	return &wrap{Context: c}
}
