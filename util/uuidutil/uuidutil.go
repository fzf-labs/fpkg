package uuidutil

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
	"github.com/teris-io/shortid"
)

// GenUUID 生成随机字符串，eg: 76d27e8c-a80e-48c8-ad20-e5562e0f67e4
func GenUUID() string {
	return uuid.NewString()
}

// GenShortID 生成一个id
func GenShortID() (string, error) {
	return shortid.Generate()
}

func KSUID() string { //2HhWuYvDuhvsOZWcVTujThVHPWf
	return ksuid.New().String()
}

func KSUIDByTime() string {
	s, _ := ksuid.NewRandomWithTime(time.Now())
	return strings.ToLower(s.String())
}
