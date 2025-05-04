package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"stzbHelper/http"
	"sync"
)

func StartHttpService(wait *sync.WaitGroup) {
	log.Println("HTTP服务启动")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	http.RegisterRoute(r)

	log.Println("http://127.0.0.1:9527 浏览器打开此地址控制软件")

	err := r.Run(":9527")

	if err != nil {
		log.Fatal("http服务启动失败:" + err.Error())
		wait.Done()
		return
	}
}
