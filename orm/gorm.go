package orm

import (
	"database/sql"

	"gorm.io/gorm"
)

const (
	unhealthy = "unhealthy"
	health    = "health"
)

// GetHealthStatus 检查链接是否健康
func GetHealthStatus(gormDB *gorm.DB) string {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return unhealthy
	}
	// 验证与数据库的连接是否仍然存在
	err = sqlDB.Ping()
	if err != nil {
		return unhealthy
	}
	err = gormDB.Raw(`select 1`).Error
	if err != nil {
		return unhealthy
	}
	return health
}

// GetState 获取目前数据库状态参数
func GetState(gormDB *gorm.DB) *sql.DBStats {
	db, err := gormDB.DB()
	if err != nil {
		return nil
	}
	state := db.Stats()
	return &state
}
