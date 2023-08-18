package wechat

import (
	"sync"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/urllink"
	"github.com/silenceper/wechat/v2/miniprogram/urlscheme"
)

var (
	miniProgramOnce sync.Once
	wxMiniProgram   WxMiniProgram
)

type MiniProgramConfig struct {
	AppID     string
	AppSecret string
}

type WxMiniProgram struct {
	config        *MiniProgramConfig
	wxMiniProgram *miniprogram.MiniProgram
}

// NewMiniProgram 实例化小程序
func NewMiniProgram(cfg *MiniProgramConfig, cache *RedisCache) *WxMiniProgram {
	miniProgramOnce.Do(func() {
		wc := wechat.NewWechat()
		config := &miniConfig.Config{
			AppID:     cfg.AppID,
			AppSecret: cfg.AppSecret,
		}
		if cache != nil {
			config.Cache = cache
		}
		MiniProgram := wc.GetMiniProgram(config)
		wxMiniProgram.config = cfg
		wxMiniProgram.wxMiniProgram = MiniProgram
	})
	return &wxMiniProgram
}

// GetMiniURLScheme 获取小程序URLScheme
func (w *WxMiniProgram) GetMiniURLScheme(envVersion urlscheme.EnvVersion) (string, error) {
	return w.wxMiniProgram.GetSURLScheme().Generate(&urlscheme.USParams{
		JumpWxa: &urlscheme.JumpWxa{
			Path:       "",
			Query:      "",
			EnvVersion: envVersion,
		},
		ExpireType:     1,
		ExpireTime:     0,
		ExpireInterval: 30,
	})
}

// GetMiniURLLink 获取小程序URLLink
func (w *WxMiniProgram) GetMiniURLLink(envVersion string) (string, error) {
	return w.wxMiniProgram.GetURLLink().Generate(&urllink.ULParams{
		Path:           "",
		Query:          "",
		EnvVersion:     envVersion,
		IsExpire:       false,
		ExpireType:     1,
		ExpireTime:     0,
		ExpireInterval: 30,
	})
}
