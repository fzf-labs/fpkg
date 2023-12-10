// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package gorm_gen_model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const TableNameAdminDemo = "admin_demo"

// AdminDemo mapped from table <admin_demo>
type AdminDemo struct {
	ID            string           `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid();comment:编号" json:"id"`          // 编号
	Username      string           `gorm:"column:username;type:character varying(50);not null;comment:用户名" json:"username"`        // 用户名
	Password      string           `gorm:"column:password;type:character varying(128);not null;comment:密码" json:"password"`        // 密码
	Nickname      string           `gorm:"column:nickname;type:character varying(50);not null;comment:昵称" json:"nickname"`         // 昵称
	Avatar        string           `gorm:"column:avatar;type:character varying(255);comment:头像" json:"avatar"`                     // 头像
	Gender        int16            `gorm:"column:gender;type:smallint;not null;comment:0=保密 1=女 2=男" json:"gender"`                // 0=保密 1=女 2=男
	Email         string           `gorm:"column:email;type:character varying(50);comment:邮件" json:"email"`                        // 邮件
	Mobile        string           `gorm:"column:mobile;type:character varying(15);comment:手机号" json:"mobile"`                     // 手机号
	JobID         string           `gorm:"column:job_id;type:character varying(50);comment:岗位" json:"jobId"`                       // 岗位
	DeptID        string           `gorm:"column:dept_id;type:character varying(50);comment:部门" json:"deptId"`                     // 部门
	RoleIds       datatypes.JSON   `gorm:"column:role_ids;type:json;comment:角色集" json:"roleIds"`                                   // 角色集
	Salt          string           `gorm:"column:salt;type:character varying(32);not null;comment:盐值" json:"salt"`                 // 盐值
	Status        int16            `gorm:"column:status;type:smallint;not null;comment:0=禁用 1=开启" json:"status"`                   // 0=禁用 1=开启
	Motto         string           `gorm:"column:motto;type:character varying(255);comment:个性签名" json:"motto"`                     // 个性签名
	CreatedAt     time.Time        `gorm:"column:created_at;type:timestamp with time zone;not null;comment:创建时间" json:"createdAt"` // 创建时间
	UpdatedAt     time.Time        `gorm:"column:updated_at;type:timestamp with time zone;not null;comment:更新时间" json:"updatedAt"` // 更新时间
	DeletedAt     gorm.DeletedAt   `gorm:"column:deleted_at;type:timestamp with time zone;comment:删除时间" json:"deletedAt"`          // 删除时间
	AdminLogDemos []*AdminLogDemo  `gorm:"foreignKey:admin_id" json:"adminLogDemos"`
	AdminRoles    []*AdminRoleDemo `gorm:"joinForeignKey:admin_id;joinReferences:role_id;many2many:admin_to_role_demo" json:"adminRoles"`
}

// TableName AdminDemo's table name
func (*AdminDemo) TableName() string {
	return TableNameAdminDemo
}
