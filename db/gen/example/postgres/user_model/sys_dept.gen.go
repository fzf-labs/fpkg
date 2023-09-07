// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package user_model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameSysDept = "sys_dept"

// SysDept mapped from table <sys_dept>
type SysDept struct {
	ID          string         `gorm:"column:id;primaryKey;comment:编号" json:"id"`                // 编号
	Pid         string         `gorm:"column:pid;not null;comment:父级id" json:"pid"`              // 父级id
	Name        string         `gorm:"column:name;not null;comment:部门简称" json:"name"`            // 部门简称
	FullName    string         `gorm:"column:full_name;not null;comment:部门全称" json:"fullName"`   // 部门全称
	Responsible string         `gorm:"column:responsible;comment:负责人" json:"responsible"`        // 负责人
	Phone       string         `gorm:"column:phone;comment:负责人电话" json:"phone"`                  // 负责人电话
	Email       string         `gorm:"column:email;comment:负责人邮箱" json:"email"`                  // 负责人邮箱
	Type        int16          `gorm:"column:type;not null;comment:1=公司 2=子公司 3=部门" json:"type"` // 1=公司 2=子公司 3=部门
	Status      int16          `gorm:"column:status;not null;comment:0=禁用 1=开启" json:"status"`   // 0=禁用 1=开启
	Sort        int64          `gorm:"column:sort;not null;comment:排序值" json:"sort"`             // 排序值
	CreatedAt   time.Time      `gorm:"column:created_at;not null;comment:创建时间" json:"createdAt"` // 创建时间
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null;comment:更新时间" json:"updatedAt"` // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;comment:删除时间" json:"deletedAt"`          // 删除时间
}

// TableName SysDept's table name
func (*SysDept) TableName() string {
	return TableNameSysDept
}
