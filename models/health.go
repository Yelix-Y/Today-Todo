package models

import (
	"time"
)

// WaterRecord 喝水记录
type WaterRecord struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `json:"user_id"`
	Amount    int       `json:"amount"` // 毫升
	CreatedAt time.Time `json:"created_at"`
	Date      string    `json:"date"` // YYYY-MM-DD
}

// StandRecord 站立记录
type StandRecord struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `json:"user_id"`
	Duration  int       `json:"duration"` // 秒
	CreatedAt time.Time `json:"created_at"`
	Date      string    `json:"date"` // YYYY-MM-DD
}

// ShortVideoRecord 短视频刷屏记录，便于反沉迷统计。
type ShortVideoRecord struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `json:"user_id"`
	Count     int       `json:"count"` // 本次计数，默认 1
	CreatedAt time.Time `json:"created_at"`
	Date      string    `json:"date"` // YYYY-MM-DD
}

// DailyProgress 每日进度统计
type DailyProgress struct {
	Date              string  `json:"date"`
	CompletedTodos    int64   `json:"completed_todos"`     // 已完成任务数
	TotalTodos        int64   `json:"total_todos"`         // 总任务数
	WaterTotal        int     `json:"water_total"`         // 总喝水量(ml)
	WaterTarget       int     `json:"water_target"`        // 目标喝水量(ml)
	WaterProgress     float64 `json:"water_progress"`      // 喝水进度百分比
	StandTotalMinutes int     `json:"stand_total_minutes"` // 总站立时长(分钟)
	StandTarget       int     `json:"stand_target"`        // 目标站立时长(分钟)
	StandProgress     float64 `json:"stand_progress"`      // 站立进度百分比
	ShortVideoCount   int     `json:"short_video_count"`   // 短视频次数
	FocusScore        int     `json:"focus_score"`         // 专注得分(0-100)
	StandingCount     int64   `json:"standing_count"`      // 站立次数
	WaterCheckins     int64   `json:"water_checkins"`      // 喝水打卡次数
}
