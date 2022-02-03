package main

import (
	"log"
	"net/http"
	"zweb"
)

//  curl "http://127.0.0.1:10002/hello?name=sharpe"
// curl -i http://127.0.0.1:10002/
// curl "http://127.0.0.1:10002/login" -X POST -d 'username=sharpe&password=123456'
// curl "http://127.0.0.1:10002/xxx" -X POST -d 'username=sharpe&password=123456'
// curl "http://127.0.0.1:10002/hello/sharpe"
// curl "http://127.0.0.1:10002/assets/program/go/base/slice.md"

func main() {
	r := zweb.New()
	r.GET("/", func(ctx *zweb.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Zweb!</h1>")
	})

	r.GET("/hello", func(ctx *zweb.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})

	r.POST("/login", func(ctx *zweb.Context) {
		ctx.JSON(http.StatusOK, zweb.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	r.GET("/hello/:name", func(ctx *zweb.Context) {
		// expect /hello/sharpe
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
	})

	r.GET("/assets/*filepath", func(ctx *zweb.Context) {
		ctx.JSON(http.StatusOK, zweb.H{"filepath": ctx.Param("filepath")})
	})

	log.Fatal(r.Run(":10002"))
}
