package controllers

import (
	"Today-Todo/models" // 引入刚才定义的 models 包
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateTodo 创建任务
func CreateTodo(c *gin.Context) {
	var todo models.Todo // 使用 models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用 models.DB
	models.DB.Create(&todo)
	c.JSON(http.StatusOK, todo)
}

// GetTodos 获取列表
func GetTodos(c *gin.Context) {
	var todos []models.Todo
	models.DB.Find(&todos)
	c.JSON(http.StatusOK, todos)
}

// UpdateTodo 修改任务
func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	if err := models.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到任务"})
		return
	}

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Save(&todo)
	c.JSON(http.StatusOK, todo)
}

// DeleteTodo 删除任务
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	if err := models.DB.Unscoped().Delete(&models.Todo{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
