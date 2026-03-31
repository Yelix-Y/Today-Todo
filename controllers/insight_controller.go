package controllers

import (
	"Today-Todo/models"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type insightTask struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Priority string `json:"priority"`
}

type todayInsight struct {
	Date              string        `json:"date"`
	RiskLevel         string        `json:"risk_level"`
	Momentum          string        `json:"momentum"`
	SuggestedAction   string        `json:"suggested_action"`
	SuggestedNudge    string        `json:"suggested_nudge"`
	FocusScore        int           `json:"focus_score"`
	CompletionRate    float64       `json:"completion_rate"`
	HydrationProgress float64       `json:"hydration_progress"`
	StandingProgress  float64       `json:"standing_progress"`
	TopTasks          []insightTask `json:"top_tasks"`
}

// GetTodayInsights 返回今日行为洞察，作为差异化能力给前端直接消费。
func GetTodayInsights(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "1")
	date := c.DefaultQuery("date", todayDate())

	var waterRecords []models.WaterRecord
	var standRecords []models.StandRecord
	var shortRecords []models.ShortVideoRecord
	var todos []models.Todo

	models.DB.Where("user_id = ? AND date = ?", userID, date).Find(&waterRecords)
	models.DB.Where("user_id = ? AND date = ?", userID, date).Find(&standRecords)
	models.DB.Where("user_id = ? AND date = ?", userID, date).Find(&shortRecords)
	models.DB.Where("date(created_at) = date(?)", date).Find(&todos)

	waterTotal := 0
	for _, record := range waterRecords {
		waterTotal += record.Amount
	}

	standTotalSeconds := 0
	for _, record := range standRecords {
		standTotalSeconds += record.Duration
	}

	shortVideoCount := 0
	for _, record := range shortRecords {
		shortVideoCount += record.Count
	}

	var totalTodos int64
	var completedTodos int64
	models.DB.Model(&models.Todo{}).
		Where("date(created_at) = date(?)", date).
		Count(&totalTodos)
	models.DB.Model(&models.Todo{}).
		Where("date(created_at) = date(?) AND completed = ?", date, true).
		Count(&completedTodos)

	completionRate := 0.0
	if totalTodos > 0 {
		completionRate = float64(completedTodos) / float64(totalTodos) * 100
	}

	hydrationProgress := clampProgress(float64(waterTotal) / 2000 * 100)
	standingProgress := clampProgress(float64(standTotalSeconds/60) / 30 * 100)
	focusScore := calcFocusScore(completedTodos, totalTodos, shortVideoCount)

	riskLevel := "low"
	momentum := "稳定推进"
	action := "保持当前节奏，优先完成 1 个中高优任务。"
	nudge := "建议下一轮 25 分钟专注后站立 5 分钟。"

	if shortVideoCount >= 5 || focusScore < 50 {
		riskLevel = "high"
		momentum = "注意力波动"
		action = "先停止信息流干扰 20 分钟，只做 1 个最小任务。"
		nudge = "立刻离开短视频，补水并开始 10 分钟冲刺。"
	} else if shortVideoCount >= 3 || completionRate < 40 || hydrationProgress < 35 {
		riskLevel = "medium"
		momentum = "节奏偏慢"
		action = "缩小任务粒度，优先完成 1 个高优先级任务。"
		nudge = "先喝水，再做一个 15 分钟的可完成小任务。"
	}

	topTasks := pickTopTasks(todos, 3)

	c.JSON(http.StatusOK, gin.H{
		"data": todayInsight{
			Date:              date,
			RiskLevel:         riskLevel,
			Momentum:          momentum,
			SuggestedAction:   action,
			SuggestedNudge:    nudge,
			FocusScore:        focusScore,
			CompletionRate:    completionRate,
			HydrationProgress: hydrationProgress,
			StandingProgress:  standingProgress,
			TopTasks:          topTasks,
		},
	})
}

func pickTopTasks(todos []models.Todo, limit int) []insightTask {
	openTodos := make([]models.Todo, 0, len(todos))
	for _, todo := range todos {
		if !todo.Completed {
			openTodos = append(openTodos, todo)
		}
	}

	sort.SliceStable(openTodos, func(i, j int) bool {
		pi := priorityWeight(openTodos[i].Priority)
		pj := priorityWeight(openTodos[j].Priority)
		if pi != pj {
			return pi > pj
		}
		return openTodos[i].CreatedAt.Before(openTodos[j].CreatedAt)
	})

	if len(openTodos) > limit {
		openTodos = openTodos[:limit]
	}

	result := make([]insightTask, 0, len(openTodos))
	for _, todo := range openTodos {
		result = append(result, insightTask{
			ID:       todo.ID,
			Title:    todo.Title,
			Priority: todo.Priority,
		})
	}

	return result
}

func priorityWeight(priority string) int {
	switch priority {
	case "high":
		return 3
	case "medium":
		return 2
	default:
		return 1
	}
}
