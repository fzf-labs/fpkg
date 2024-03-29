package qq

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type UserInfo struct {
	Ret            int    `json:"ret"`            // 返回码
	Msg            string `json:"msg"`            // 如果ret<0，会有相应的错误信息提示，返回数据全部用UTF-8编码。
	Nickname       string `json:"nickname"`       // 用户在QQ空间的昵称。
	Figureurl      string `json:"figureurl"`      // 大小为30×30像素的QQ空间头像URL。
	Figureurl1     string `json:"figureurl_1"`    // 大小为50×50像素的QQ空间头像URL。
	Figureurl2     string `json:"figureurl_2"`    // 大小为100×100像素的QQ空间头像URL。
	FigureurlQQ1   string `json:"figureurl_qq_1"` // 大小为40×40像素的QQ头像URL。
	FigureurlQQ2   string `json:"figureurl_qq_2"` // 大小为100×100像素的QQ头像URL。需要注意，不是所有的用户都拥有QQ的100x100的头像，但40x40像素则是一定会有。
	Gender         string `json:"gender"`         // 性别。 如果获取不到则默认返回"男"
	Openid         string `json:"openid"`
	Unionid        string `json:"unionid"`
	IDentifierCode string `json:"identifier_code"`
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// AuthConfig 基本配置
type AuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	ApplyUnionID string
}

// AuthQq QQ授权登录
type AuthQq struct {
	config *AuthConfig // 配置信息
}

// NewAuthQq QQ授权登录实例化
func NewAuthQq(conf *AuthConfig) *AuthQq {
	return &AuthQq{config: conf}
}

func (a *AuthQq) GetState(state string) string {
	if state == "" {
		return uuid.NewString()
	}
	return state
}

// AuthorizationCode 获取accessToken
// code -> accessToken
// 文档：https://wiki.connect.qq.com/%E4%BD%BF%E7%94%A8authorization_code%E8%8E%B7%E5%8F%96access_token
// 接口：https://graph.qq.com/oauth2.0/token
func (a *AuthQq) AuthorizationCode(code string) (*AccessToken, error) {
	var accessToken AccessToken
	resp, err := req.R().
		SetQueryParam("grant_type", "authorization_code").
		SetQueryParam("client_id", a.config.ClientID).
		SetQueryParam("client_secret", a.config.ClientSecret).
		SetQueryParam("code", code).
		SetQueryParam("redirect_uri", a.config.RedirectURL).
		SetSuccessResult(&accessToken).
		Get("https://graph.qq.com/oauth2.0/token")
	if err != nil {
		return nil, err
	}
	if resp.IsSuccessState() {
		return nil, errors.New("http response err")
	}
	return &accessToken, nil
}

// GetOpenid 获取openid
// 文档：https://wiki.connect.qq.com/%E8%8E%B7%E5%8F%96%E7%94%A8%E6%88%B7openid_oauth2-0
// 接口：https://graph.qq.com/oauth2.0/me
func (a *AuthQq) GetOpenid(accessToken string) (clientID, openid, unionID string, err error) {
	resp, err := req.R().
		SetQueryParam("access_token", accessToken).
		SetQueryParam("unionID", a.config.ApplyUnionID). // 申请unionID，0：不申请，1：申请
		Get("https://graph.qq.com/oauth2.0/me")
	if err != nil {
		return "", "", "", err
	}
	content := resp.String()
	content = a.removeCallback(content)
	clientID = gjson.Get(content, "client_id").String()
	openid = gjson.Get(content, "openid").String()
	unionID = gjson.Get(content, "unionID").String()
	return
}

// GetUserInfo 获取用户信息
// 文档：https://wiki.connect.qq.com/get_user_info
// 接口：https://graph.qq.com/user/get_user_info
func (a *AuthQq) GetUserInfo(accessToken string) (*UserInfo, error) {
	clientID, openid, unionID, err := a.GetOpenid(accessToken)
	if err != nil {
		return nil, err
	}
	if clientID != a.config.ClientID {
		return nil, errors.New("qq login err clientID")
	}
	resp, err := req.R().
		SetQueryParam("access_token", accessToken).
		SetQueryParam("oauth_consumer_key", clientID).
		SetQueryParam("openid", openid).
		Get("https://graph.qq.com/user/get_user_info")
	if err != nil {
		return nil, err
	}

	content := resp.String()
	ret := gjson.Get(content, "ret").Int()
	msg := gjson.Get(content, "msg").String()

	if ret != 0 {
		return nil, errors.New("get_user_info:ret=" + strconv.FormatInt(ret, 10) + ",msg=" + msg)
	}

	userInfo := &UserInfo{}
	err = json.Unmarshal([]byte(content), userInfo)
	if err != nil {
		return nil, err
	}
	userInfo.Openid = openid
	userInfo.Unionid = unionID
	if a.config.ApplyUnionID == "0" {
		userInfo.IDentifierCode = openid
	} else {
		userInfo.IDentifierCode = unionID
	}
	return userInfo, nil
}

// GetUserInfoByCode
// 根据code获取用户信息
// 流程为先使用code换取accessToken，然后根据accessToken获取用户信息
func (a *AuthQq) GetUserInfoByCode(code string) (*UserInfo, error) {
	token, err := a.AuthorizationCode(code)
	if err != nil {
		return nil, err
	}
	return a.GetUserInfo(token.AccessToken)
}

// qq有些接口返回的数据带了callback，例如：callback( {"error":100020,"error_description":"code is reused error"} );
// 这里将callback去掉
func (a *AuthQq) removeCallback(content string) string {
	prefix := "callback("
	suffix := ");"
	content = strings.TrimSpace(content)
	if strings.Index(content, "callback(") == 0 {
		content = content[len(prefix) : len(content)-len(suffix)]
		content = strings.TrimSpace(content)
	}
	return content
}
