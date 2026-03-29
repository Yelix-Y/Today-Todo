package main

import (
	"Today-Todo/controllers"
	"Today-Todo/models"
	"Today-Todo/routers"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. 初始化数据库
	models.ConnectDatabase()

	// 2. 初始化任务调度器（多线程处理健康提醒）
	controllers.InitScheduler()

	// 3. 初始化路由
	r := routers.SetupRouter()

	// 4. 启动服务
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	log.Println("🚀 服务已启动在 http://localhost:8080")
	log.Println("✅ 多线程任务调度器已启动")
	log.Println("💧 喝水提醒：每90分钟")
	log.Println("🚶 站立提醒：每60分钟")
	log.Println("📵 防沉迷提醒：每120分钟")

	// 5. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务...")
	controllers.StopScheduler()
}
