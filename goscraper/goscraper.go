package goscraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

var (
	EscapedFragment string = "_escaped_fragment_="
	fragmentRegexp         = regexp.MustCompile("#!(.*)")
)

type Scraper struct {
	URL                *url.URL
	EscapedFragmentURL *url.URL
	MaxRedirect        int
}

type Document struct {
	Body    bytes.Buffer
	Preview DocumentPreview
}

type DocumentPreview struct {
	Icon        string
	Name        string
	Title       string
	Description string
	Images      []string
	Link        string
}

func Scrape(uri string, maxRedirect int) (*Document, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	return (&Scraper{URL: u, MaxRedirect: maxRedirect}).Scrape()
}

func (scraper *Scraper) Scrape() (*Document, error) {
	doc, err := scraper.getDocument()
	if err != nil {
		return nil, err
	}
	err = scraper.parseDocument(doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (scraper *Scraper) getURL() string {
	if scraper.EscapedFragmentURL != nil {
		return scraper.EscapedFragmentURL.String()
	}
	return scraper.URL.String()
}

func (scraper *Scraper) toFragmentURL() error {
	unescapedurl, err := url.QueryUnescape(scraper.URL.String())
	if err != nil {
		return err
	}
	matches := fragmentRegexp.FindStringSubmatch(unescapedurl)
	if len(matches) > 1 {
		escapedFragment := EscapedFragment
		for _, r := range matches[1] {
			b := byte(r)
			if avoidByte(b) {
				continue
			}
			if escapeByte(b) {
				escapedFragment += url.QueryEscape(string(r))
			} else {
				escapedFragment += string(r)
			}
		}

		p := "?"
		if len(scraper.URL.Query()) > 0 {
			p = "&"
		}
		fragmentURL, err := url.Parse(strings.Replace(unescapedurl, matches[0], p+escapedFragment, 1))
		if err != nil {
			return err
		}
		scraper.EscapedFragmentURL = fragmentURL
	} else {
		p := "?"
		if len(scraper.URL.Query()) > 0 {
			p = "&"
		}
		fragmentURL, err := url.Parse(unescapedurl + p + EscapedFragment)
		if err != nil {
			return err
		}
		scraper.EscapedFragmentURL = fragmentURL
	}
	return nil
}

func (scraper *Scraper) getDocument() (*Document, error) {
	scraper.MaxRedirect -= 1
	if strings.Contains(scraper.URL.String(), "#!") {
		err := scraper.toFragmentURL()
		if err != nil {
			return nil, err
		}
	}
	if strings.Contains(scraper.URL.String(), EscapedFragment) {
		scraper.EscapedFragmentURL = scraper.URL
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, scraper.getURL(), http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "GoScraper")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.Request.URL.String() != scraper.getURL() {
		scraper.EscapedFragmentURL = nil
		scraper.URL = resp.Request.URL
	}
	b, err := convertUTF8(resp.Body, resp.Header.Get("content-type"))
	if err != nil {
		return nil, err
	}
	doc := &Document{Body: b, Preview: DocumentPreview{Link: scraper.URL.String()}}

	return doc, nil
}

func convertUTF8(content io.Reader, contentType string) (bytes.Buffer, error) {
	buff := bytes.Buffer{}
	content, err := charset.NewReader(content, contentType)
	if err != nil {
		return buff, err
	}
	_, err = io.Copy(&buff, content)
	if err != nil {
		return buff, err
	}
	return buff, nil
}

func (scraper *Scraper) parseDocument(doc *Document) error {
	t := html.NewTokenizer(&doc.Body)
	var ogImage bool
	var headPassed bool
	var hasFragment bool
	var hasCanonical bool
	var canonicalURL *url.URL
	doc.Preview.Images = []string{}
	// saves previews' link in case that <link rel="canonical"> is found after <meta property="og:url">
	link := doc.Preview.Link
	// set default value to site name if <meta property="og:site_name"> not found
	doc.Preview.Name = scraper.URL.Host
	// set default icon to web root if <link rel="icon" href="/favicon.ico"> not found
	doc.Preview.Icon = fmt.Sprintf("%s://%s%s", scraper.URL.Scheme, scraper.URL.Host, "/favicon.ico")
	for {
		tokenType := t.Next()
		if tokenType == html.ErrorToken {
			return nil
		}
		if tokenType != html.SelfClosingTagToken && tokenType != html.StartTagToken && tokenType != html.EndTagToken {
			continue
		}
		token := t.Token()

		switch token.Data {
		case "head":
			if tokenType == html.EndTagToken {
				headPassed = true
			}
		case "body":
			headPassed = true

		case "link":
			var canonical bool
			var hasIcon bool
			var href string
			for _, attr := range token.Attr {
				if cleanStr(attr.Key) == "rel" && cleanStr(attr.Val) == "canonical" {
					canonical = true
				}
				if cleanStr(attr.Key) == "rel" && strings.Contains(cleanStr(attr.Val), "icon") {
					hasIcon = true
				}
				if cleanStr(attr.Key) == "href" {
					href = attr.Val
				}
				if href != "" && canonical && link != href {
					hasCanonical = true
					var err error
					canonicalURL, err = url.Parse(href)
					if err != nil {
						return err
					}
				}
				if href != "" && hasIcon {
					doc.Preview.Icon = href
				}
			}

		case "meta":
			if len(token.Attr) != 2 {
				break
			}
			if metaFragment(token) && scraper.EscapedFragmentURL == nil {
				hasFragment = true
			}
			var property string
			var content string
			for _, attr := range token.Attr {
				if cleanStr(attr.Key) == "property" || cleanStr(attr.Key) == "name" {
					property = attr.Val
				}
				if cleanStr(attr.Key) == "content" {
					content = attr.Val
				}
			}
			switch cleanStr(property) {
			case "og:site_name":
				doc.Preview.Name = content
			case "og:title":
				doc.Preview.Title = content
			case "og:description":
				doc.Preview.Description = content
			case "description":
				if doc.Preview.Description == "" {
					doc.Preview.Description = content
				}
			case "og:url":
				doc.Preview.Link = content
			case "og:image":
				ogImage = true
				ogImgURL, err := url.Parse(content)
				if err != nil {
					return err
				}
				if !ogImgURL.IsAbs() {
					ogImgURL, err = url.Parse(fmt.Sprintf("%s://%s%s", scraper.URL.Scheme, scraper.URL.Host, ogImgURL.Path))
					if err != nil {
						return err
					}
				}

				doc.Preview.Images = []string{ogImgURL.String()}
			}
		case "title":
			if tokenType == html.StartTagToken {
				t.Next()
				token = t.Token()
				if doc.Preview.Title == "" {
					doc.Preview.Title = token.Data
				}
			}
		case "img":
			for _, attr := range token.Attr {
				if cleanStr(attr.Key) == "src" {
					imgURL, err := url.Parse(attr.Val)
					if err != nil {
						return err
					}
					if !imgURL.IsAbs() {
						doc.Preview.Images = append(doc.Preview.Images, fmt.Sprintf("%s://%s%s", scraper.URL.Scheme, scraper.URL.Host, imgURL.Path))
					} else {
						doc.Preview.Images = append(doc.Preview.Images, attr.Val)
					}
				}
			}
		}
		if hasCanonical && headPassed && scraper.MaxRedirect > 0 {
			if !canonicalURL.IsAbs() {
				absCanonical, err := url.Parse(fmt.Sprintf("%s://%s%s", scraper.URL.Scheme, scraper.URL.Host, canonicalURL.Path))
				if err != nil {
					return err
				}
				canonicalURL = absCanonical
			}
			scraper.URL = canonicalURL
			scraper.EscapedFragmentURL = nil
			fdoc, err := scraper.getDocument()
			if err != nil {
				return err
			}
			*doc = *fdoc
			return scraper.parseDocument(doc)
		}
		if hasFragment && headPassed && scraper.MaxRedirect > 0 {
			err := scraper.toFragmentURL()
			if err != nil {
				return err
			}
			fdoc, err := scraper.getDocument()
			if err != nil {
				return err
			}
			*doc = *fdoc
			return scraper.parseDocument(doc)
		}
		if doc.Preview.Title != "" && doc.Preview.Description != "" && ogImage && headPassed {
			return nil
		}
	}
}

func avoidByte(b byte) bool {
	i := int(b)
	if i == 127 || (i >= 0 && i <= 31) {
		return true
	}
	return false
}

func escapeByte(b byte) bool {
	i := int(b)
	if i == 32 || i == 35 || i == 37 || i == 38 || i == 43 || (i >= 127 && i <= 255) {
		return true
	}
	return false
}

func metaFragment(token html.Token) bool {
	var name string
	var content string

	for _, attr := range token.Attr {
		if cleanStr(attr.Key) == "name" {
			name = attr.Val
		}
		if cleanStr(attr.Key) == "content" {
			content = attr.Val
		}
	}
	if name == "fragment" && content == "!" {
		return true
	}
	return false
}

func cleanStr(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}
