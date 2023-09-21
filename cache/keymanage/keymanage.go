package keymanage

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type KeyManage struct {
	ServerName string
	List       map[string]KeyPrefix
}

// New 实例化key管理器
func New(serverName string) *KeyManage {
	return &KeyManage{
		ServerName: serverName,
		List:       make(map[string]KeyPrefix),
	}
}

// AddKey 添加一个key prefix
func (p *KeyManage) AddKey(prefix string, expirationTime time.Duration, remark string) *KeyPrefix {
	if _, ok := p.List[prefix]; ok {
		panic(fmt.Sprintf("key %s is exsit, please change one", prefix))
	}
	key := KeyPrefix{
		ServerName:     p.ServerName,
		PrefixName:     prefix,
		Remark:         remark,
		ExpirationTime: expirationTime,
	}
	p.List[prefix] = key
	return &key
}

// Document 导出MD文档
func (p *KeyManage) Document() string {
	str := `|ServerName|PrefixName|ttl(s)|Remark` + "\n" + `|--|--|--|--|` + "\n"

	if len(p.List) > 0 {
		for _, m := range p.List {
			str += `|` + p.ServerName + `|` + m.PrefixName + `|` + strconv.FormatFloat(m.ExpirationTime.Seconds(), 'f', -1, 64) + `|` + m.Remark + `|` + "\n"
		}
	}
	return str
}

// KeyPrefix key前缀
type KeyPrefix struct {
	ServerName     string
	PrefixName     string
	Remark         string
	ExpirationTime time.Duration
}

// Key 获取key
func (p *KeyPrefix) Key(keys ...string) string {
	return strings.Join(append([]string{p.ServerName, p.PrefixName}, keys...), ":")
}

// Keys 获取keys
func (p *KeyPrefix) Keys(keys []string) []string {
	result := make([]string, 0)
	if len(keys) > 0 {
		for _, key := range keys {
			result = append(result, strings.Join([]string{p.ServerName, p.PrefixName, key}, ":"))
		}
	}
	return result
}

// TTL 获取key的过期时间time.Duration
func (p *KeyPrefix) TTL() time.Duration {
	return p.ExpirationTime
}

// TTLSecond 获取key的过期时间 Second
func (p *KeyPrefix) TTLSecond() int {
	return int(p.ExpirationTime / time.Second)
}
