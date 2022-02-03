package main

import (
	"log"
	"net/http"
	"time"
	"zweb"
	"zweb/middlewares/logger"
)

func onlyForV2() zweb.HandlerFunc {
	return func(context *zweb.Context) {
		// start time
		t := time.Now()
		// if a server error occurred
		//context.Fail(http.StatusInternalServerError, "Internal server Error")

		// Calculate the time
		log.Printf("[%d] %s in %v for group v2", context.StatusCode, context.Req.RequestURI, time.Since(t))
	}
}

// curl http://127.0.0.1:10004/
// curl http://127.0.0.1:10004/v2/hello/sharpe

func main() {
	r := zweb.New()
	r.Use(logger.Logger()) // global middleware

	r.GET("/", func(ctx *zweb.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Zweb</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2())

	{
		v2.GET("/hello/:name", func(ctx *zweb.Context) {
			// expect /hello/sharpe
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
		})
	}

	log.Fatal(r.Run(":10004"))
}
