package models

import (
	"time"

	"gorm.io/gorm"
)

// Todo 任务实体，支持基础待办和完成状态统计。
type Todo struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Completed   bool           `json:"completed"`
	Priority    string         `json:"priority"`
	DueAt       *time.Time     `json:"due_at"`
}
