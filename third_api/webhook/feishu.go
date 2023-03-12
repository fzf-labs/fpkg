package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/imroc/req/v3"
	"strconv"
	"time"
)

type FeiShuConfig struct {
	Url  string
	Sign string
}

type FeiShu struct {
	cfg *FeiShuConfig
}

func NewFeiShu(cfg *FeiShuConfig) *FeiShu {
	return &FeiShu{cfg: cfg}
}

func NewFeiShuByCfg(cfg *FeiShuConfig) *FeiShu {
	return &FeiShu{cfg: cfg}
}

type SendMsg struct {
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
	MsgType   string `json:"msg_type"`
	Content   struct {
		Text string `json:"text"`
	} `json:"content"`
}

// SendMsg 发送消息
func (f *FeiShu) SendMsg(msg string) error {
	timestamp := time.Now().Unix()
	sign, err := f.GenSign(f.cfg.Sign, timestamp)
	if err != nil {
		return err
	}
	param := SendMsg{
		Timestamp: strconv.FormatInt(timestamp, 10),
		Sign:      sign,
		MsgType:   "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: msg,
		},
	}
	resp, err := req.R().SetBody(param).Post(f.cfg.Url)
	if err != nil {
		return err
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("bad response status: %s", resp.Status)
	}
	return nil
}

func (f *FeiShu) GenSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
