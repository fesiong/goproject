package convert

import (
	"testing"
)

// Test for Request.
func TestRequest(t *testing.T) {
	link := "http://www.fengying.com.tw/"
	resp, err := Request(link, nil)
	if err != nil {
		t.Error(err.Error())
	}
	t.Skip(resp)
}

// Test for ToUtf8.
func TestToUtf8(t *testing.T) {
	content := "中国青年网"
	content = ToUtf8(content)

	t.Skip(content)
}

// Test for Convert.
func TestConvert(t *testing.T) {
	content := "中国青年网"
	content = Convert(content, "utf-8", "utf-8")

	t.Skip(content)
}
