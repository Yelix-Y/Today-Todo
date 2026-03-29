package controllers

import (
	"Today-Todo/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type createTodoRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	DueAt       *time.Time `json:"due_at"`
}

type updateTodoRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Priority    *string    `json:"priority"`
	Completed   *bool      `json:"completed"`
	DueAt       *time.Time `json:"due_at"`
}

// CreateTodo 创建任务
func CreateTodo(c *gin.Context) {
	var req createTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo := models.Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    normalizePriority(req.Priority),
		DueAt:       req.DueAt,
	}

	if err := models.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建任务失败"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// GetTodos 获取列表
func GetTodos(c *gin.Context) {
	var todos []models.Todo

	query := models.DB.Order("created_at desc")
	if status := c.Query("completed"); status != "" {
		completed, err := strconv.ParseBool(status)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "completed 参数必须是 true 或 false"})
			return
		}
		query = query.Where("completed = ?", completed)
	}

	if err := query.Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询任务失败"})
		return
	}

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

	var req updateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Priority != nil {
		todo.Priority = normalizePriority(*req.Priority)
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	if req.DueAt != nil {
		todo.DueAt = req.DueAt
	}

	if err := models.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新任务失败"})
		return
	}

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

func normalizePriority(priority string) string {
	switch priority {
	case "high", "medium", "low":
		return priority
	default:
		return "medium"
	}
}
