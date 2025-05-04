package common

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    any    `json:"data"`
}

// Success 成功的返回
func (r Response) Success(c *gin.Context) {
	r.Code = 200
	if r.Message == "" {
		r.Message = "ok"
	}
	c.JSON(200, r)
}

// Error 发生错误的返回
func (r Response) Error(c *gin.Context) {
	if r.Message == "" {
		r.Message = "error"
	}

	if r.Code == 0 {
		r.Code = 500
	}

	c.JSON(200, r)
}
