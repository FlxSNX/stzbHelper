package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"stzbHelper/http/route/api"
	"stzbHelper/web"
)

func RegisterRoute(r *gin.Engine) {
	staticRoute(r)
	api.Register(r.Group("/v1"))
}

func staticRoute(r *gin.Engine) {
	assetsFS, err := fs.Sub(web.PublicAssets, "dist")
	if err != nil {
		fmt.Println("初始化静态资源出错")
		return
	}
	staticServer := http.FileServer(http.FS(assetsFS))

	r.NoRoute(func(c *gin.Context) {
		// 获取请求路径
		reqpath := c.Request.URL.Path

		// 处理根路径默认指向index.html
		if reqpath == "/" {
			reqpath = "/index.html"

		} else if reqpath[len(reqpath)-1:] == "/" { // 最后一位字符是斜杠时去除 否则会触发下面的打不开静态资源 然后500状态码
			reqpath = reqpath[:len(reqpath)-1]
		}

		// 尝试打开静态文件
		file, err := assetsFS.Open(reqpath[1:])

		if errors.Is(err, fs.ErrNotExist) {
			// 文件不存在，返回自定义404
			c.JSON(404, gin.H{
				"message": "404 - Page Not Found",
			})
			return
		} else if err != nil {
			// 其他类型的错误
			log.Println("静态资源访问错误: " + err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer file.Close()

		// 检查是否是目录
		fileInfo, err := file.Stat()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if fileInfo.IsDir() {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		staticServer.ServeHTTP(c.Writer, c.Request)
	})
}
