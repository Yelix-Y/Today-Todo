package models

import (
	"github.com/glebarez/sqlite" // <---  glebarez
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	// 连接数据库
	// 这里的逻辑没变，只是底层的驱动变了，它会自动创建 gorm.db 文件
	DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		// 我修改了这里，万一报错，把具体的错误原因打印出来，而不是只说“失败”
		panic("连接数据库失败: " + err.Error())
	}

	// 自动迁移
	err = DB.AutoMigrate(&Todo{})
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}
}
