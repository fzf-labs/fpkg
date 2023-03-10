package openai

import (
	"fmt"
	"github.com/imroc/req/v3"
	"log"
)

type ChatGPT struct {
	ApiKey string
}

func NewChatGPT(apiKey string) *ChatGPT {
	return &ChatGPT{ApiKey: apiKey}
}

type CompletionsReq struct {
	Model       string      `json:"model"`
	Prompt      string      `json:"prompt"`
	MaxTokens   int         `json:"max_tokens"`
	Temperature float32     `json:"temperature"`
	TopP        int         `json:"top_p"`
	N           int         `json:"n"`
	Stream      bool        `json:"stream"`
	Logprobs    interface{} `json:"logprobs"`
	Stop        string      `json:"stop"`
}

type CompletionsResp struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (c *ChatGPT) Completions(msg string) (string, error) {
	var result string
	url := "https://api.openai.com/v1/completions"
	param := &CompletionsReq{
		Model:       "text-davinci-003",
		Prompt:      msg,
		MaxTokens:   4000,
		Temperature: 0.7,
		TopP:        1,
		N:           1,
	}
	var response CompletionsResp
	client := req.C()
	resp, err := client.R().SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey)).SetBody(param).SetSuccessResult(&response).Post(url)
	if err != nil {
		return "", err
	}
	fmt.Println(resp)
	if !resp.IsSuccessState() {
		return "", fmt.Errorf("bad response status: %s", resp.Status)
	}
	if len(response.Choices) > 0 {
		result = response.Choices[0].Text
	}
	return result, nil
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ChatCompletionsReq struct {
	Model            string        `json:"model"`
	Messages         []ChatMessage `json:"messages"`
	MaxTokens        int           `json:"max_tokens"`
	Temperature      float32       `json:"temperature"`
	TopP             int           `json:"top_p"`
	FrequencyPenalty int           `json:"frequency_penalty"`
	PresencePenalty  int           `json:"presence_penalty"`
}

type ChatCompletionsResp struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (c *ChatGPT) ChatCompletions(messages []ChatMessage) (*ChatMessage, error) {
	log.Println(messages)
	var result *ChatMessage
	url := "https://api.openai.com/v1/chat/completions"
	param := &ChatCompletionsReq{
		Model:            "gpt-3.5-turbo",
		Messages:         messages,
		MaxTokens:        2000,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	var response ChatCompletionsResp
	client := req.C()
	resp, err := client.R().SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey)).SetBody(param).SetSuccessResult(&response).Post(url)
	if err != nil {
		return nil, err
	}
	log.Println(resp)
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("bad response status: %s", resp.Status)
	}
	if len(response.Choices) > 0 {
		result = &response.Choices[0].Message
	}
	return result, nil
}
