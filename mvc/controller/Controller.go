package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ydsxiong/go-playground/mvc/modelview"
)

func SetupHTTPServerController() {
	http.HandleFunc("/", modelview.Index)
	http.HandleFunc("/show", modelview.Show)
	http.HandleFunc("/new", modelview.New)
	http.HandleFunc("/edit", modelview.Edit)
	http.HandleFunc("/insert", modelview.Insert)
	http.HandleFunc("/update", modelview.Update)
	http.HandleFunc("/delete", modelview.Delete)

	log.Println("http server started successfully, listening on port 8090!")
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("err occurred from http service" + err.Error())
	} else {
		fmt.Println("not supposed to reach here")
	}
}
