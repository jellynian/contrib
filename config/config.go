package config

import (
	fl "flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	yaml "gopkg.in/yaml.v2"
)

var (
	conf   = make(map[interface{}]interface{})
	locker sync.Mutex
	flag   uint32 // ==1 已经读取过配置 ==0 需要重新加载配置
)

type Configure interface {
	Get(key string) *Configure
	Interface() interface{}
	Int() int
	Int64() int64
	String() string
	Float64() float64
	Map() map[interface{}]interface{}
}

type configure struct {
	keys []string
}

func getConfigFileName() string {
	fl.Parse()
	configFlag := fl.String("config", "", "config file path")
	if *configFlag != "" {
		return *configFlag
	}
	config := os.Getenv("CONFIG")
	if config != "" {
		return config
	}
	return "./config.yaml"
}

func TryReload() {
	if atomic.LoadUint32(&flag) == 0 {
		return
	} else {
		locker.Lock()
		defer locker.Unlock()
		if flag == 1 {
			atomic.StoreUint32(&flag, 0)
		}
	}
}

func MustLoad() {
	if atomic.LoadUint32(&flag) == 1 {
		return
	}

	locker.Lock()
	defer locker.Unlock()

	if flag == 0 {
		configPath := getConfigFileName()
		if _, err := os.Stat(configPath); err != nil {
			if os.IsNotExist(err) {
				log.Panic("Config file not fund", err.Error())
			}
		}
		if configPath == "" {
			log.Panic("config path is nil")
		}
		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Panic("can't read config file due to:", err.Error())
		}
		err = yaml.Unmarshal(data, &conf)
		if err != nil {
			log.Panic(err.Error())
		}

		atomic.StoreUint32(&flag, 1)
	}

}

func Default() *configure {
	MustLoad()
	return new(configure)
}

func (c *configure) Get(key string) *configure {
	MustLoad()
	c.keys = strings.Split(key, ".")
	return c
}

func (c *configure) Interface() interface{} {
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

func (c *configure) String() string {

	v, ok := c.Interface().(string)
	if !ok {
		panic(fmt.Sprintf("value %v is not string", c.Interface()))
	}
	return v
}

func (c *configure) Int64() int64 {
	v, ok := c.Interface().(int64)
	if !ok {
		panic(fmt.Sprintf("value %v is not int64", c.Interface()))
	}
	return v
}

func (c *configure) Int() int {
	v, ok := c.Interface().(int)
	if !ok {
		panic(fmt.Sprintf("value %v is not int", c.Interface()))
	}
	return v
}

func (c *configure) Float64() float64 {
	v, ok := c.Interface().(float64)
	if !ok {
		panic(fmt.Sprintf("value %v is not float64", c.Interface()))
	}
	return v
}

func (c *configure) Map() map[interface{}]interface{} {
	v, ok := c.Interface().(map[interface{}]interface{})
	if !ok {
		panic(fmt.Sprintf("value %v is not map[interface{}]interface{}", c.Interface()))
	}
	return v
}
