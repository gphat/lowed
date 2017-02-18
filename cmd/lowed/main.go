package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"lowed"

	"github.com/Sirupsen/logrus"
	"github.com/go-yaml/yaml"
)

var (
	configFile = flag.String("f", "", "The config file to read for settings.")
)

func main() {
	flag.Parse()

	if configFile == nil || *configFile == "" {
		logrus.Fatal("You must specify a config file")
	}

	config, err := ReadConfig(*configFile)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to read config file")
	}
	fmt.Println(config)
}

// ReadConfig attempts to read a config file and unmarshal it's YAML
// in to a config struct.
func ReadConfig(path string) (c lowed.Config, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	bts, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(bts, &c)
	if err != nil {
		return
	}

	return c, nil
}
