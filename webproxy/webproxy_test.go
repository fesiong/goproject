package webproxy

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"testing"
)

func TestWebProxy(t *testing.T) {
	app := iris.New()

	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.HTML(fmt.Sprintf("<h1>请访问：http://%s/proxy?link=https://www.baidu.com/</h1>", ctx.Host()))
	})
	app.Handle("GET", "/proxy", func(ctx iris.Context) {
		body, err := WebProxy(ctx.URLParam("link"), fmt.Sprintf("http://%s/proxy", ctx.Host()))
		if err != nil {
			t.Skip(err.Error())
			return
		}

		ctx.HTML(body)
	})

	app.Run(iris.Addr(":8088"))

	t.Skip()
}
