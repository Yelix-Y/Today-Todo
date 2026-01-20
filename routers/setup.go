package routers

import (
	"Today-Todo/controllers" // 引入 controllers
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		// 注意：这里调用的函数名都变大写了
		v1.POST("/todos", controllers.CreateTodo)
		v1.GET("/todos", controllers.GetTodos)
		v1.PUT("/todos/:id", controllers.UpdateTodo)
		v1.DELETE("/todos/:id", controllers.DeleteTodo)
	}

	return r
}
