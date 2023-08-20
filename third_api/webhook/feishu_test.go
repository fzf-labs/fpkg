package webhook

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeiShu_SendMsg(t *testing.T) {
	feiShu := NewFeiShu(&FeiShuConfig{
		URL:  "",
		Sign: "",
	})
	err := feiShu.SendText("测试")
	if err != nil {
		return
	}
	fmt.Println(err)
	assert.Equal(t, nil, err)
}
