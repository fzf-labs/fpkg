package jsonutil

import (
	"bytes"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/encoder"
	"regexp"
	"strings"
	"text/scanner"
)

// EncodeSortMapKeys 编码,map排序
func EncodeSortMapKeys(v interface{}) ([]byte, error) {
	return encoder.Encode(v, encoder.SortMapKeys)
}

// SortMapKeys 编码,json key 排序
func SortMapKeys(v []byte) ([]byte, error) {
	root, err := sonic.Get(v)
	if err != nil {
		return nil, err
	}
	err = root.SortKeys(false)
	if err != nil {
		return nil, err
	}
	return root.MarshalJSON()
}

// EncodeEscapeHTML 编码,去除html转义
func EncodeEscapeHTML(v interface{}) ([]byte, error) {
	return encoder.Encode(v, encoder.EscapeHTML)
}

// `(?s:` enable match multi line
var jsonMLComments = regexp.MustCompile(`(?s:/\*.*?\*/\s*)`)

// StripComments 去除 JSON 字符串的注释
func StripComments(src string) string {
	// multi line comments
	if strings.Contains(src, "/*") {
		src = jsonMLComments.ReplaceAllString(src, "")
	}

	// single line comments
	if !strings.Contains(src, "//") {
		return strings.TrimSpace(src)
	}

	// strip inline comments
	var s scanner.Scanner

	s.Init(strings.NewReader(src))
	s.Filename = "comments"
	s.Mode ^= scanner.SkipComments // don't skip comments

	buf := new(bytes.Buffer)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		txt := s.TokenText()
		if !strings.HasPrefix(txt, "//") && !strings.HasPrefix(txt, "/*") {
			buf.WriteString(txt)
			// } else {
			// fmt.Printf("%s: %s\n", s.Position, txt)
		}
	}

	return buf.String()
}
