package main

import (
	redisService "awesomeProject/redis_"
	"awesomeProject/router"
	"awesomeProject/utils"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()
	// 使用中间件Logger日志功能
	r.Use(gin.Logger())
	// 使用中间件Recovery
	r.Use(gin.Recovery())
	// 配置不信任所有代理,
	r.SetTrustedProxies([]string{"192.168.2.56"})
	// 初始化路由
	router.InitRouter(r)

	// 初始化数据库
	utils.InitDB()
	// 初始化Redis
	redisService.InitRedis()

	//go timer.InitTimer()

	//r.GET()

	r.Run(":8080")

}
