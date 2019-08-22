package main

import (
	"flag"
	"github.com/ydsxiong/gorestapiclient/app"
	"github.com/ydsxiong/gorestapiclient/app/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func main() {

	config := config.GetDefaultConfig()
	config.ServerPort = "8081"

	configPtr := flag.Bool("config", false, "to use external config or not")
	flag.Parse()
	if *configPtr {
		data, err := ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			log.Fatal(err)
		}
	}

	app := &app.App{}
	app.Initialize(config)
	app.Run(config.ServerPort)
}
