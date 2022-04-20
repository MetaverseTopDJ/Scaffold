package main

import (
	"log"

	"github.com/MetaverseTopDJ/Scaffold/app"
)

/* 示例代码 */
func main() {
	if err := app.InitModule("./conf/dev/", []string{"base", "swagger", "postgres", "redis"}); err != nil {
		log.Fatal(err)
	}
	defer app.Destroy()
}
