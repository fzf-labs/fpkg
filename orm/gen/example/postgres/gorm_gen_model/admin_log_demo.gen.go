// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package gorm_gen_model

import (
	"time"

	"gorm.io/datatypes"
)

const TableNameAdminLogDemo = "admin_log_demo"

// AdminLogDemo mapped from table <admin_log_demo>
type AdminLogDemo struct {
	ID        string         `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid();comment:编号" json:"id"`          // 编号
	AdminID   string         `gorm:"column:admin_id;type:uuid;not null;comment:管理员ID" json:"adminId"`                        // 管理员ID
	IP        string         `gorm:"column:ip;type:character varying(32);not null;comment:ip" json:"ip"`                     // ip
	URI       string         `gorm:"column:uri;type:character varying(200);not null;comment:请求路径" json:"uri"`                // 请求路径
	Useragent string         `gorm:"column:useragent;type:character varying(255);comment:浏览器标识" json:"useragent"`            // 浏览器标识
	Header    datatypes.JSON `gorm:"column:header;type:json;comment:header" json:"header"`                                   // header
	Req       datatypes.JSON `gorm:"column:req;type:json;comment:请求数据" json:"req"`                                           // 请求数据
	Resp      datatypes.JSON `gorm:"column:resp;type:json;comment:响应数据" json:"resp"`                                         // 响应数据
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp with time zone;not null;comment:创建时间" json:"createdAt"` // 创建时间
}

// TableName AdminLogDemo's table name
func (*AdminLogDemo) TableName() string {
	return TableNameAdminLogDemo
}
