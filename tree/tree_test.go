package tree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义我们自己的菜单对象
type SystemMenu struct {
	ID       int    `json:"id"`        // id
	FatherID int    `json:"father_id"` // 上级菜单id
	Name     string `json:"name"`      // 菜单名
	Route    string `json:"route"`     // 页面路径
	Icon     string `json:"icon"`      // 图标路径
}

func (s SystemMenu) GetID() int {
	return s.ID
}

func (s SystemMenu) GetPid() int {
	return s.FatherID
}

func (s SystemMenu) GetData() any {
	return s
}

func (s SystemMenu) IsRoot() bool {
	// 这里通过FatherID等于0 或者 FatherID等于自身ID表示顶层根节点
	return s.FatherID == 0 || s.FatherID == s.ID
}

// endregion

type SystemMenus []SystemMenu

// ConvertToINodeArray 将当前数组转换成父类 INode 接口 数组
func (s SystemMenus) ConvertToINodeArray() (nodes []INode) {
	for _, v := range s {
		nodes = append(nodes, v)
	}
	return
}

func TestGenerateTree(t *testing.T) {
	// 模拟获取数据库中所有菜单，在其它所有的查询中，也是首先将数据库中所有数据查询出来放到数组中，
	// 后面的遍历递归，都在这个 allMenu中进行，而不是在数据库中进行递归查询，减小数据库压力。
	allMenu := SystemMenus{
		{ID: 1, FatherID: 0, Name: "系统总览", Route: "/systemOverview", Icon: "icon-system"},
		{ID: 2, FatherID: 0, Name: "系统配置", Route: "/systemConfig", Icon: "icon-config"},

		{ID: 3, FatherID: 1, Name: "资产", Route: "/asset", Icon: "icon-asset"},
		{ID: 4, FatherID: 1, Name: "动环", Route: "/pe", Icon: "icon-pe"},

		{ID: 5, FatherID: 2, Name: "菜单配置", Route: "/menuConfig", Icon: "icon-menu-config"},
		{ID: 6, FatherID: 3, Name: "设备", Route: "/device", Icon: "icon-device"},
		{ID: 7, FatherID: 3, Name: "机柜", Route: "/device", Icon: "icon-device"},
	}

	// 模拟从数据库中查询出 '设备'
	device := SystemMenu{ID: 1, FatherID: 0, Name: "系统总览", Route: "/systemOverview", Icon: "icon-system"}
	// 查询 设备 的所有父节点
	respNodes := FindRelationSubNode(allMenu.ConvertToINodeArray(), device)

	fmt.Println(respNodes)
	assert.True(t, len(respNodes) > 0)
}
