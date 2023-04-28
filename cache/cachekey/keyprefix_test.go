package cachekey

import (
	"fmt"
	"testing"
	"time"
)

func TestKeyPrefixes_Document(t *testing.T) {
	prefixes := NewKeyPrefixes("user")
	prefixes.AddKeyPrefix("uuid", time.Hour, "uuid")
	prefixes.AddKeyPrefix("dl", time.Second*5, "分布式锁")
	fmt.Println(prefixes.Document())
}

func TestKeyPrefixes_AddKeyPrefix(t *testing.T) {
	prefixes := NewKeyPrefixes("user")
	prefixes.AddKeyPrefix("uuid", time.Hour, "uuid")
	prefixes.AddKeyPrefix("uuid", time.Hour, "uuid")
}
