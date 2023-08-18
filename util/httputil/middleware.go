package httputil

import (
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
)

// ResponseMiddleware  响应中间件
func ResponseMiddleware(_ *req.Client, resp *req.Response) error {
	if resp.Err != nil { // 有一个潜在的错误，例如网络错误或解封错误(SetSuccessResult或SetError之前被调用)。
		if dump := resp.Dump(); dump != "" { // 将转储内容附加到原始基础错误以帮助排除故障。
			resp.Err = errors.New(fmt.Sprintf("%s\nraw content:\n%s", resp.Err.Error(), resp.Dump()))
		}
		return nil // 如果存在底层错误，则跳过以下逻辑。
	}
	// 极端情况:既没有错误响应，也没有成功响应;
	// 转储内容以帮助排除故障。
	if !resp.IsSuccessState() {
		resp.Err = errors.New(fmt.Sprintf("bad response, raw content:\n%s", resp.Dump()))
	}
	return nil
}
