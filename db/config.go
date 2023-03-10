package db

import "time"

// MysqlConfig
// @Description: mysql配置
type MysqlConfig struct {
	DataSourceName  string        `json:"DataSourceName"`
	ShowLog         bool          `json:"ShowLog"`
	MaxIdleConn     int           `json:"MaxIdleConn"`
	MaxOpenConn     int           `json:"MaxOpenConn"`
	ConnMaxLifeTime time.Duration `json:"ConnMaxLifeTime"`
}
