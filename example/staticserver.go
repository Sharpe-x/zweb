package main

import (
	"log"
	"zweb"
)

// http://127.0.0.1:10002/assets/example/group.go
func main() {
	r := zweb.New()
	r.Static("/assets", "../") // root 可以是相对或者绝对路径
	log.Fatal(r.Run(":10002"))
}
