package main

import (
	_ "sportApi/filter"
	_ "sportApi/routers"

	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.BConfig.WebConfig.Session.SessionOn = true

	httpport, _ := beego.AppConfig.String("httpport")
	fmt.Println("server port:%s", httpport)

	beego.Run()
}
