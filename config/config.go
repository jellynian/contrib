package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var (
	conf         = make(map[interface{}]interface{})
	onceLoadConf sync.Once
)

type Configure struct {
	keys []string
}

func mustReadConf() {
	var config_path = "./config.yaml"
	if _, err := os.Stat(config_path); err != nil {
		if os.IsNotExist(err) {
			logrus.Error("Config file not fund", err.Error())
			os.Exit(-1)
		}
	}
	if config_path == "" {
		os.Exit(-1)
	}
	data, err := ioutil.ReadFile(config_path)
	if err != nil {
		logrus.Panic(err.Error())
	}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		logrus.Panic(err.Error())
	}
}

func Default() *Configure {
	onceLoadConf.Do(mustReadConf)
	return new(Configure)
}

func (c *Configure) Get(key string) *Configure {
	c.keys = strings.Split(key, ".")
	return c
}

func (c *Configure) Interface() interface{} {
	var (
		ok    bool
		value = conf
		cKeys []string
	)

	lenKeys := len(c.keys)

	cKeys = c.keys[:lenKeys-1]

	for _, v := range cKeys {
		value, ok = value[v].(map[interface{}]interface{})
		if !ok {
			panic(fmt.Sprintf("key %s not fund from Configure file ", v))
		}
	}
	v := value[c.keys[lenKeys-1]]
	return v
}

func (c *Configure) String() string {

	v, ok := c.Interface().(string)
	if !ok {
		panic(fmt.Sprintf("value %s is not string", v))
	}
	return v
}

func (c *Configure) Int64() int64 {
	v, ok := c.Interface().(int64)
	if !ok {
		panic(fmt.Sprintf("value %v is not int64", v))
	}
	return v
}

func (c *Configure) Int() int {
	v, ok := c.Interface().(int)
	if !ok {
		panic(fmt.Sprintf("value %v is not int64", v))
	}
	return v
}

func (c *Configure) Float64() float64 {
	v, ok := c.Interface().(float64)
	if !ok {
		panic(fmt.Sprintf("value %v is not float64", v))
	}
	return v
}

func (c *Configure) Map() map[interface{}]interface{} {
	v, ok := c.Interface().(map[interface{}]interface{})
	if !ok {
		panic(fmt.Sprintf("value %v is not map[interface{}]interface{}", v))
	}
	return v
}
