package qq

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthQq_GetUserInfo(t *testing.T) {
	authQq := NewAuthQq(&AuthConfig{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "",
		ApplyUnionID: "1",
	})
	token := ""
	userInfo, err := authQq.GetUserInfo(token)
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println(userInfo)
	assert.Equal(t, nil, err)
}
