package shard

import (
	"gorm.io/sharding"
	"time"
)

func NewShardingPlugin(table string, shardingKey string, num uint) *sharding.Sharding {
	return sharding.Register(sharding.Config{
		ShardingKey:         shardingKey,
		NumberOfShards:      num,
		PrimaryKeyGenerator: sharding.PKSnowflake,
	}, table)
}

// NewMonthShardingPlugin 按月份分表
func NewMonthShardingPlugin(table string, shardingKey string) *sharding.Sharding {
	return sharding.Register(sharding.Config{
		ShardingKey:         shardingKey,
		PrimaryKeyGenerator: sharding.PKCustom,
		ShardingAlgorithm: func(columnValue interface{}) (suffix string, err error) {
			t := time.Unix(columnValue.(int64), 00)
			return "_" + t.Format("2006") + t.Format("01"), nil
		},
		PrimaryKeyGeneratorFn: func(tableIdx int64) int64 {
			return tableIdx
		},
	}, table)
}
