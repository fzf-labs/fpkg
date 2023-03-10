package webhook

import (
	"fmt"
	"testing"
)

func TestFeiShu_SendMsg(t *testing.T) {
	feiShu := NewFeiShu(&FeiShuConfig{
		Url:  "",
		Sign: "",
	})
	err := feiShu.SendMsg("测试")
	if err != nil {
		return
	}
	fmt.Println(err)
}
