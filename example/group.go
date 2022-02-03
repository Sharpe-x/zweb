package main

import (
	"log"
	"net/http"
	"zweb"
)

//  curl "http://127.0.0.1:10003/hello?name=sharpe"
// curl "http://127.0.0.1:10003/v1"
//  curl "http://127.0.0.1:10003/v2/"

func main() {
	r := zweb.New()
	r.GET("/index", func(ctx *zweb.Context) {
		ctx.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(ctx *zweb.Context) {
			ctx.HTML(http.StatusOK, "<h1>Hello Zweb!</h1>")
		})

		v1.GET("/hello", func(ctx *zweb.Context) {
			// expect /hello?name=sharpe
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(ctx *zweb.Context) {
			// expect /hello/sharpe
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
		})
		v2.POST("/login", func(ctx *zweb.Context) {
			ctx.JSON(http.StatusOK, zweb.H{
				"username": ctx.PostForm("username"),
				"password": ctx.PostForm("password"),
			})
		})
	}

	log.Fatal(r.Run(":10003"))
}
