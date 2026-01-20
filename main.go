package main

import (
	"Today-Todo/models"
	"Today-Todo/routers"
)

func main() {
	// 1. 初始化数据库
	models.ConnectDatabase()

	// 2. 初始化路由
	r := routers.SetupRouter()

	// 3. 启动服务
	r.Run(":8080")
}
