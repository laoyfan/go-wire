package main

import (
	"fmt"
	"go-wire/wire"
)

func main() {
	app, err := wire.InitApp()
	if err != nil {
		fmt.Println("应用初始化异常:", err)
	}
	if err = app.Run(); err != nil {
		fmt.Println("应用运行失败:", err)
	}
}
