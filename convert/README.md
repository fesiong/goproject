# golang 抓取网页并将其他编码(gbk,gb2312,big5等)中文文字转换成uft8编码的字符串处理函数

最近用golang采集网页中遇到了各种不能识别的的乱码字符串，他们大多编码是gbk、gb2312、big5、windows-1252 等编码。有时候，网页上并没有声明编码，却使用上面这种编码的网页也有，也有网页声明的编码和实际使用的编码不同的网页，导致网页编码转换工作带来诸多不便，更多的是根据提示的编码转换出来依然还是乱码的问题，着实让人头疼。于是乎，为了得到一个通用可行的中文字符串编码转换方法，本人通过网络上上百万个网站测试，采集数据回来进行编码转换，终于总结出来了一套绝大部分都能顺利将网页中文字符串编码都转换成utf-8编码的方法。

## golang项目直接引用
安装依赖包
```shell script
go get github.com/fesiong/goproject/convert
```
使用说明  
对外公开有3个函数，Request函数支持请求网络页面，并自动检测页面内容的编码，转换成utf-8，ToUtf8函数支持传入的字符串会自动检测编码，并转换成utf-8，Convert函数需要传入原始编码和输出编码，如果原始编码传入出错，则转换出来的文本会乱码

请求网络页面，并自动检测页面内容的编码，转换成utf-8
```go
    link := "http://www.youth.cn/"
    resp, err := Request(link)
    if err != nil {
        t.Error(err.Error())
    }
```
传入的字符串会自动检测编码，并转换成utf-8
```go
    content := "中国青年网"
    content = ToUtf8(content)
```
传入原始编码和输出编码
```go
    content := "中国青年网"
    content = Convert(content, "utf-8", "utf-8")
```
## 源码地址
[github.com/fesiong/goproject](https://github.com/fesiong/goproject)


## 核心的编码转换判断函数
```go
package convert

import (
	"github.com/axgle/mahonia"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/net/html/charset"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type RequestData struct {
	Header     http.Header
	Request    *http.Request
	Body       string
	Status     string
	StatusCode int
}

/**
 * 请求网络页面，并自动检测页面内容的编码，转换成utf-8
 */
func Request(urlPath string) (*RequestData, error) {
	resp, body, errs := gorequest.New().Timeout(90 * time.Second).Get(urlPath).End()
	if len(errs) > 0 {
		//如果是https,则尝试退回http请求
		if strings.HasPrefix(urlPath, "https") {
			urlPath = strings.Replace(urlPath, "https://", "http://", 1)
			return Request(urlPath)
		}
		return nil, errs[0]
	}
	defer resp.Body.Close()
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	body = toUtf8(body, contentType)

	requestData := RequestData{
		Header:     resp.Header,
		Request:    resp.Request,
		Body:       body,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}

	return &requestData, nil
}

/**
 * 对外公开的编码转换接口，传入的字符串会自动检测编码，并转换成utf-8
 */
func ToUtf8(content string) string {
	return toUtf8(content, "")
}

/**
 * 内部编码判断和转换，会自动判断传入的字符串编码，并将它转换成utf-8
 */
func toUtf8(content string, contentType string) string {
	var htmlEncode string

	if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
		htmlEncode = "gb18030"
	} else if strings.Contains(contentType, "big5") {
		htmlEncode = "big5"
	} else if strings.Contains(contentType, "utf-8") {
		htmlEncode = "utf-8"
	}
	if htmlEncode == "" {
		//先尝试读取charset
		reg := regexp.MustCompile(`(?is)<meta[^>]*charset\s*=["']?\s*([A-Za-z0-9\-]+)`)
		match := reg.FindStringSubmatch(content)
		if len(match) > 1 {
			contentType = strings.ToLower(match[1])
			if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
				htmlEncode = "gb18030"
			} else if strings.Contains(contentType, "big5") {
				htmlEncode = "big5"
			} else if strings.Contains(contentType, "utf-8") {
				htmlEncode = "utf-8"
			}
		}
		if htmlEncode == "" {
			reg = regexp.MustCompile(`(?is)<title[^>]*>(.*?)<\/title>`)
			match = reg.FindStringSubmatch(content)
			if len(match) > 1 {
				aa := match[1]
				_, contentType, _ = charset.DetermineEncoding([]byte(aa), "")
				htmlEncode = strings.ToLower(htmlEncode)
				if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
					htmlEncode = "gb18030"
				} else if strings.Contains(contentType, "big5") {
					htmlEncode = "big5"
				} else if strings.Contains(contentType, "utf-8") {
					htmlEncode = "utf-8"
				}
			}
		}
	}
	if htmlEncode != "" && htmlEncode != "utf-8" {
		content = Convert(content, htmlEncode, "utf-8")
	}

	return content
}

/**
 * 编码转换
 * 需要传入原始编码和输出编码，如果原始编码传入出错，则转换出来的文本会乱码
 */
func Convert(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

```