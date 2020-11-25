package main

import (
	"flag"
	"fmt"
	"github.com/shybily/beansproxy/resources"
	"github.com/shybily/beansproxy/server"
	"os"
)

var (
	config string
)

func main() {
	flag.StringVar(&config, "config", "", "config file")
	flag.Parse()

	if len(config) <= 0 {
		flag.Usage()
		os.Exit(1)
	}

	opt, err := server.NewYamlProxyOptions(config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	resources.InitInstance(opt.Instances, opt.Strategy)
	_ = (server.NewProxyServer(opt)).Listen()
}
