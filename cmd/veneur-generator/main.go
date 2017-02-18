package main

import (
	"flag"
	"fmt"
)

var (
	configFile = flag.String("f", "", "The config file to read for settings.")
)

func main() {
	flag.Parse()

	if configFile == nil || *configFile == "" {
		logrus.Fatal("You must specify a config file")
	}

	fmt.Println("hello world!")
}
