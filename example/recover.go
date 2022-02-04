package main

import (
	"log"
	"net/http"
	"zweb"
	"zweb/middlewares/logger"
	"zweb/middlewares/recover"
)

func main() {
	r := zweb.New()
	r.Use(logger.Logger(), recover.Recovery())

	r.GET("/", func(ctx *zweb.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	r.GET("/panic", func(ctx *zweb.Context) {
		names := []string{"1", "2"}
		ctx.String(http.StatusOK, names[10])
	})

	log.Fatal(r.Run(":9999"))
}
