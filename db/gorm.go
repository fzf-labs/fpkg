package db

import (
	"database/sql"

	"gorm.io/gorm"
)

// GetHealthStatus 检查链接是否健康
func GetHealthStatus(gorm *gorm.DB) string {
	sqlDB, err := gorm.DB()
	if err != nil {
		return "unhealth"
	}
	// 验证与数据库的连接是否仍然存在
	err = sqlDB.Ping()
	if err != nil {
		return "unhealth"
	}
	err = gorm.Raw(`select 1`).Error
	if err != nil {
		return "unhealth"
	}
	return "health"
}

// GetState 获取目前数据库状态参数
func GetState(gorm *gorm.DB) *sql.DBStats {
	db, err := gorm.DB()
	if err != nil {
		return nil
	}
	state := db.Stats()
	return &state
}
