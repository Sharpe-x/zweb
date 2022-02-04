package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"zweb"
	"zweb/middlewares/logger"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := zweb.New()
	r.Use(logger.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("../templates/*")
	r.Static("/assets", "../static")

	stu1 := &student{Name: "sharpe", Age: 20}
	stu2 := &student{Name: "neo", Age: 22}
	r.GET("/", func(c *zweb.Context) {
		c.HTMLRender(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *zweb.Context) {
		c.HTMLRender(http.StatusOK, "arr.tmpl", zweb.H{
			"title":  "zweb",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *zweb.Context) {
		c.HTMLRender(http.StatusOK, "custom_func.tmpl", zweb.H{
			"title": "zweb",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	log.Fatal(r.Run(":9999"))

}
