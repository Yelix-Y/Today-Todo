package controllers

import (
	"Today-Todo/models"
	"Today-Todo/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var scheduler *services.Scheduler
var reminderHub = services.NewReminderHub()

// InitScheduler 初始化调度器
func InitScheduler() {
	scheduler = services.NewScheduler(reminderHub)
	scheduler.Start()
}

// StopScheduler 停止后台提醒任务。
func StopScheduler() {
	if scheduler != nil {
		scheduler.Stop()
	}
}

// GetScheduler 获取调度器，供其他控制器读取配置。
func GetScheduler() *services.Scheduler {
	return scheduler
}

// GetReminderHub 获取提醒中心，供 SSE 推流使用。
func GetReminderHub() *services.ReminderHub {
	return reminderHub
}

// GetUserState 获取用户状态
func GetUserState(c *gin.Context) {
	if scheduler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "调度器尚未初始化"})
		return
	}

	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 32)

	sm := scheduler.GetStateMachine(uint(userID))
	if sm == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户状态机未找到"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"state":   sm.GetCurrentState(),
	})
}

// ChangeUserState 切换用户状态
func ChangeUserState(c *gin.Context) {
	if scheduler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "调度器尚未初始化"})
		return
	}

	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 32)

	var req struct {
		State models.UserState `json:"state" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sm := scheduler.GetStateMachine(uint(userID))
	if sm == nil {
		// 如果状态机不存在，先注册用户
		scheduler.RegisterUser(uint(userID))
		sm = scheduler.GetStateMachine(uint(userID))
	}

	if !sm.TransitionTo(req.State) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "无效的状态转换",
			"current_state": sm.GetCurrentState(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "状态切换成功",
		"state":   sm.GetCurrentState(),
	})
}
