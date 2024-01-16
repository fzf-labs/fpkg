package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
)

type FeiShuConfig struct {
	URL  string
	Sign string
}

type FeiShu struct {
	cfg *FeiShuConfig
}

func NewFeiShu(cfg *FeiShuConfig) *FeiShu {
	return &FeiShu{cfg: cfg}
}

func (f *FeiShu) GenSign(secret string, timestamp int64) (string, error) {
	// timestamp + key 做sha256, 再进行base64 encode
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

type SendText struct {
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
	MsgType   string `json:"msg_type"`
	Content   struct {
		Text string `json:"text"`
	} `json:"content"`
}

// SendText 发送文本消息
func (f *FeiShu) SendText(msg string) error {
	timestamp := time.Now().Unix()
	sign, err := f.GenSign(f.cfg.Sign, timestamp)
	if err != nil {
		return err
	}
	param := SendText{
		Timestamp: strconv.FormatInt(timestamp, 10),
		Sign:      sign,
		MsgType:   "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: msg,
		},
	}
	resp, err := req.R().SetBody(param).Post(f.cfg.URL)
	if err != nil {
		return err
	}
	if !resp.IsSuccessState() {
		return errors.New(fmt.Sprintf("bad response status: %s", resp.Status))
	}
	return nil
}

type SendInteractive struct {
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
	MsgType   string `json:"msg_type"`
	Card      Card   `json:"card"`
}

type Card struct {
	Elements []struct {
		Tag  string `json:"tag"`
		Text struct {
			Content string `json:"content"`
			Tag     string `json:"tag"`
		} `json:"text,omitempty"`
		Actions []struct {
			Tag  string `json:"tag"`
			Text struct {
				Content string `json:"content"`
				Tag     string `json:"tag"`
			} `json:"text"`
			URL   string `json:"url"`
			Type  string `json:"type"`
			Value struct {
			} `json:"value"`
		} `json:"actions,omitempty"`
	} `json:"elements"`
	Header struct {
		Title struct {
			Content string `json:"content"`
			Tag     string `json:"tag"`
		} `json:"title"`
	} `json:"header"`
}

// SendInteractive 发送消息卡片
func (f *FeiShu) SendInteractive(card Card) error {
	timestamp := time.Now().Unix()
	sign, err := f.GenSign(f.cfg.Sign, timestamp)
	if err != nil {
		return err
	}
	param := SendInteractive{
		Timestamp: strconv.FormatInt(timestamp, 10),
		Sign:      sign,
		MsgType:   "text",
		Card:      card,
	}
	resp, err := req.R().SetBody(param).Post(f.cfg.URL)
	if err != nil {
		return err
	}
	if !resp.IsSuccessState() {
		return errors.New(fmt.Sprintf("bad response status: %s", resp.Status))
	}
	return nil
}
