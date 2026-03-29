package models

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	// 连接数据库，使用轻量 sqlite 方便本地和移动端联调。
	DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	// 自动迁移：统一维护待办与健康统计所需表结构。
	err = DB.AutoMigrate(&Todo{}, &WaterRecord{}, &StandRecord{}, &ShortVideoRecord{})
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}
}
