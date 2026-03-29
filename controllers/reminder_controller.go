package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StreamReminders 通过 SSE 推送实时提醒。
func StreamReminders(c *gin.Context) {
	hub := GetReminderHub()
	ch := hub.Subscribe()
	defer hub.Unsubscribe(ch)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.Status(http.StatusOK)
	c.Writer.Flush()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case event, ok := <-ch:
			if !ok {
				return
			}
			payload, _ := json.Marshal(event)
			_, _ = c.Writer.Write([]byte("event: reminder\n"))
			_, _ = c.Writer.Write([]byte("data: " + string(payload) + "\n\n"))
			c.Writer.Flush()
		}
	}
}

// GetReminderConfig 返回提醒间隔，供客户端定时器兜底。
func GetReminderConfig(c *gin.Context) {
	s := GetScheduler()
	if s == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "调度器尚未初始化"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": s.ReminderConfig(),
	})
}
