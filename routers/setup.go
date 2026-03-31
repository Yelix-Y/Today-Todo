package routers

import (
	"Today-Todo/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 提供 Web 客户端静态资源，便于一键本地体验。
	r.Static("/web", "./web")
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		// TODO CRUD
		v1.POST("/todos", controllers.CreateTodo)
		v1.GET("/todos", controllers.GetTodos)
		v1.PUT("/todos/:id", controllers.UpdateTodo)
		v1.DELETE("/todos/:id", controllers.DeleteTodo)

		// 健康与防沉迷记录
		v1.POST("/health/water", controllers.RecordWater)
		v1.POST("/health/stand", controllers.RecordStand)
		v1.POST("/health/short-video", controllers.RecordShortVideo)
		v1.GET("/health/daily-progress", controllers.GetDailyProgress)
		v1.GET("/insights/today", controllers.GetTodayInsights)

		// 实时提醒事件流 + 配置
		v1.GET("/reminders/stream", controllers.StreamReminders)
		v1.GET("/reminders/config", controllers.GetReminderConfig)

		// 状态机接口（保留给扩展场景）
		v1.GET("/state/:user_id", controllers.GetUserState)
		v1.POST("/state/:user_id", controllers.ChangeUserState)
	}

	return r
}
