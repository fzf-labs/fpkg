// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package user_model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameSysJob = "sys_job"

// SysJob mapped from table <sys_job>
type SysJob struct {
	ID        string         `gorm:"column:id;primaryKey;comment:编号" json:"id"`                // 编号
	Name      string         `gorm:"column:name;not null;comment:岗位名称" json:"name"`            // 岗位名称
	Code      string         `gorm:"column:code;comment:岗位编码" json:"code"`                     // 岗位编码
	Remark    string         `gorm:"column:remark;comment:备注" json:"remark"`                   // 备注
	Sort      int64          `gorm:"column:sort;not null;comment:排序值" json:"sort"`             // 排序值
	Status    int16          `gorm:"column:status;not null;comment:0=禁用 1=开启" json:"status"`   // 0=禁用 1=开启
	CreatedAt time.Time      `gorm:"column:created_at;not null;comment:创建时间" json:"createdAt"` // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;comment:更新时间" json:"updatedAt"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;comment:删除时间" json:"deletedAt"`          // 删除时间
}

// TableName SysJob's table name
func (*SysJob) TableName() string {
	return TableNameSysJob
}
