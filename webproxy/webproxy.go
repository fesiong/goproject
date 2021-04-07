package webproxy

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fesiong/goproject/convert"
	"net/url"
	"path/filepath"
	"strings"
)

func WebProxy(link string, selfLink string) (string, error) {
	self, err := url.Parse(selfLink)
	if err != nil {
		return "", err
	}
	selfHost := self.Host

	link, _ = url.QueryUnescape(link)
	if strings.HasPrefix(link, "//") {
		link = fmt.Sprintf("https:%s", link)
	}
	if strings.HasPrefix(link, "http") {
		resp, err := convert.Request(link)
		if err != nil {
			return "", err
		}

		htmlR := strings.NewReader(resp.Body)
		doc, err := goquery.NewDocumentFromReader(htmlR)
		if err != nil {
			return "", err
		}
		aLinks := doc.Find("a")
		for i := range aLinks.Nodes {
			href, exists := aLinks.Eq(i).Attr("href")
			if exists {
				newHref := ParseLink(href, link)
				if newHref != "" {
					if !strings.Contains(href, selfHost) {
						newHref = fmt.Sprintf("%s?link=%s", selfLink, url.QueryEscape(href))
						aLinks.Eq(i).SetAttr("href", newHref)
					}
				}
			}
		}
		imgLinks := doc.Find("img")
		for i := range imgLinks.Nodes {
			src, exists := imgLinks.Eq(i).Attr("src")
			if exists {
				newHref := ParseLink(src, link)
				if newHref != "" {
					if !strings.Contains(src, selfHost) {
						imgLinks.Eq(i).SetAttr("src", newHref)
					}
				}
			}
		}
		styleLinks := doc.Find("link")
		for i := range styleLinks.Nodes {
			src, exists := styleLinks.Eq(i).Attr("href")
			if exists {
				newHref := ParseLink(src, link)
				if newHref != "" {
					if !strings.Contains(src, selfHost) {
						styleLinks.Eq(i).SetAttr("href", newHref)
					}
				}
			}
		}
		scriptLinks := doc.Find("script")
		for i := range scriptLinks.Nodes {
			src, exists := scriptLinks.Eq(i).Attr("src")
			if exists {
				newHref := ParseLink(src, link)
				if newHref != "" {
					if !strings.Contains(src, selfHost) {
						scriptLinks.Eq(i).SetAttr("src", newHref)
					}
				}
			}
		}

		htmlStr, _ := doc.Html()
		return htmlStr, nil
	}

	return "", errors.New("无法读取页面")
}

func ParseLink(link string, baseUrl string) string {
	//过滤不同源url
	if strings.Contains(link, "javascript") || strings.Contains(link, "void") || link == "#" || link == "./" || link == "../" || link == "../../" {
		return ""
	}

	link = replaceDot(link, baseUrl)

	return link
}

func replaceDot(currUrl string, baseUrl string) string {
	if strings.HasPrefix(currUrl, "//") {
		currUrl = fmt.Sprintf("https:%s", currUrl)
	}
	urlInfo, err := url.Parse(currUrl)
	if err != nil {
		return ""
	}

	if urlInfo.Scheme != "" {
		return currUrl
	}
	baseInfo, err := url.Parse(baseUrl)
	if err != nil {
		return ""
	}

	u := baseInfo.Scheme + "://" + baseInfo.Host
	var path string
	if strings.Index(urlInfo.Path, "/") == 0 {
		path = urlInfo.Path
	} else {
		path = filepath.Dir(baseInfo.Path) + "/" + urlInfo.Path
	}

	rst := make([]string, 0)
	pathArr := strings.Split(path, "/")

	// 如果path是已/开头，那在rst加入一个空元素
	if pathArr[0] == "" {
		rst = append(rst, "")
	}
	for _, p := range pathArr {
		if p == ".." {
			if len(rst) > 0 {
				if rst[len(rst)-1] == ".." {
					rst = append(rst, "..")
				} else {
					rst = rst[:len(rst)-1]
				}
			}
		} else if p != "" && p != "." {
			rst = append(rst, p)
		}
	}
	return u + strings.Join(rst, "/")
}
