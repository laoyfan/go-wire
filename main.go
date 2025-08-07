package main

import (
	"fmt"
	"go-wire/bootstrap"
)

func main() {
	app, err := bootstrap.InitApp()
	if err != nil {
		fmt.Println("应用初始化异常:", err)
	}
	if err = app.Run(); err != nil {
		fmt.Println("应用启动失败:", err)
	}
}
