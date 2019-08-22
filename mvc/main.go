package main

import (
	"fmt"

	"github.com/ydsxiong/go-playground/mvc/controller"
)

func main() {
	controller.SetupHTTPServerController()
	fmt.Println("program exit")
}
