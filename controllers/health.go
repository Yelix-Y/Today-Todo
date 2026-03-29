package controllers

import (
	"Today-Todo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type recordWaterRequest struct {
	UserID uint `json:"user_id"`
	Amount int  `json:"amount" binding:"required"`
}

type recordStandRequest struct {
	UserID   uint `json:"user_id"`
	Duration int  `json:"duration"` // 秒
}

type recordShortVideoRequest struct {
	UserID uint `json:"user_id"`
	Count  int  `json:"count"`
}

// RecordWater 记录喝水，默认一杯 250ml。
func RecordWater(c *gin.Context) {
	var req recordWaterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount <= 0 {
		req.Amount = 250
	}

	record := models.WaterRecord{
		UserID:    req.UserID,
		Amount:    req.Amount,
		Date:      todayDate(),
		CreatedAt: time.Now(),
	}

	if err := models.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "喝水记录成功", "data": record})
}

// RecordStand 记录站立，默认 5 分钟。
func RecordStand(c *gin.Context) {
	var req recordStandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Duration <= 0 {
		req.Duration = 300
	}

	record := models.StandRecord{
		UserID:    req.UserID,
		Duration:  req.Duration,
		Date:      todayDate(),
		CreatedAt: time.Now(),
	}

	if err := models.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "站立记录成功", "data": record})
}

// RecordShortVideo 记录刷短视频次数。
func RecordShortVideo(c *gin.Context) {
	var req recordShortVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Count <= 0 {
		req.Count = 1
	}

	record := models.ShortVideoRecord{
		UserID:    req.UserID,
		Count:     req.Count,
		Date:      todayDate(),
		CreatedAt: time.Now(),
	}

	if err := models.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "短视频记录成功", "data": record})
}

// GetDailyProgress 获取每日进度。
func GetDailyProgress(c *gin.Context) {
	userID := c.Query("user_id")
	date := c.DefaultQuery("date", todayDate())

	// 查询喝水记录
	var waterRecords []models.WaterRecord
	models.DB.Where("user_id = ? AND date = ?", userID, date).Find(&waterRecords)

	waterTotal := 0
	for _, record := range waterRecords {
		waterTotal += record.Amount
	}

	// 查询站立记录
	var standRecords []models.StandRecord
	models.DB.Where("user_id = ? AND date = ?", userID, date).Find(&standRecords)

	standTotalSeconds := 0
	for _, record := range standRecords {
		standTotalSeconds += record.Duration
	}

	// 查询短视频次数
	var shortRecords []models.ShortVideoRecord
	models.DB.Where("user_id = ? AND date = ?", userID, date).Find(&shortRecords)

	shortVideoCount := 0
	for _, record := range shortRecords {
		shortVideoCount += record.Count
	}

	// 查询任务统计：按创建日期归档今日任务数量。
	var totalTodos int64
	var completedTodos int64
	models.DB.Model(&models.Todo{}).
		Where("date(created_at) = date(?)", date).
		Count(&totalTodos)
	models.DB.Model(&models.Todo{}).
		Where("date(created_at) = date(?) AND completed = ?", date, true).
		Count(&completedTodos)

	// 计算进度
	waterTarget := 2000 // 目标喝水 2000ml
	standTarget := 30   // 目标站立 30 分钟
	standMinutes := standTotalSeconds / 60

	progress := models.DailyProgress{
		Date:              date,
		CompletedTodos:    completedTodos,
		TotalTodos:        totalTodos,
		WaterTotal:        waterTotal,
		WaterTarget:       waterTarget,
		WaterProgress:     clampProgress(float64(waterTotal) / float64(waterTarget) * 100),
		StandTotalMinutes: standMinutes,
		StandTarget:       standTarget,
		StandProgress:     clampProgress(float64(standMinutes) / float64(standTarget) * 100),
		ShortVideoCount:   shortVideoCount,
		FocusScore:        calcFocusScore(completedTodos, totalTodos, shortVideoCount),
		StandingCount:     int64(len(standRecords)),
		WaterCheckins:     int64(len(waterRecords)),
	}

	c.JSON(http.StatusOK, gin.H{"data": progress})
}

func todayDate() string {
	return time.Now().Format("2006-01-02")
}

func clampProgress(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

func calcFocusScore(completed int64, total int64, shortVideoCount int) int {
	score := 70
	if total > 0 {
		score += int(float64(completed) / float64(total) * 25)
	}
	score -= shortVideoCount * 3
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}
