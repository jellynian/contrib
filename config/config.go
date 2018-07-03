package config

import (
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
	DefaultConfigure Configure
)

type Configure interface {
	Get(key string) *ValueInterface
	TryReload()
}

type configure struct {
	conf     map[interface{}]interface{}
	locker   sync.Mutex
	flag     uint32
	filePath string
}

type ValueInterface struct {
	value interface{}
}

func init() {
	DefaultConfigure = New("./config.yaml")
}

func New(confPath string) Configure {
	return &configure{
		filePath: confPath,
	}
}

func Default() Configure {
	return DefaultConfigure
}

func TryReload() {
	DefaultConfigure.TryReload()
}

func (c *configure) TryReload() {
	if atomic.LoadUint32(&c.flag) == 0 {
		return
	} else {
		c.locker.Lock()
		defer c.locker.Unlock()
		if c.flag == 1 {
			atomic.StoreUint32(&c.flag, 0)
		}
	}
}

func (c *configure) MustLoaded() {
	if atomic.LoadUint32(&c.flag) == 1 {
		return
	}
	c.locker.Lock()
	defer c.locker.Unlock()
	if c.flag == 0 {
		if c.filePath == "" {
			panic("config path is nil")
		}
		if _, err := os.Stat(c.filePath); err != nil {
			if os.IsNotExist(err) {
				panic(err.Error())
			}
		}

		data, err := ioutil.ReadFile(c.filePath)
		if err != nil {
			panic(err.Error())
		}
		err = yaml.Unmarshal(data, &c.conf)
		if err != nil {
			log.Panic(err.Error())
		}
		atomic.StoreUint32(&c.flag, 1)
	}
}

func (c *configure) Get(key string) *ValueInterface {
	c.MustLoaded()
	var (
		ok    bool
		value = c.conf
		cKeys []string
	)
	keys := strings.Split(key, ".")
	cKeys = keys[:len(keys)-1]
	for _, v := range cKeys {
		value, ok = value[v].(map[interface{}]interface{})
		if !ok {
			panic(fmt.Sprintf("key %s is not map[interface{}]interface{} ", v))
		}
	}
	return &ValueInterface{
		value: value[keys[len(keys)-1]],
	}
}

func (c *ValueInterface) String() string {

	v, ok := c.value.(string)
	if !ok {
		panic(fmt.Sprintf("Value %v is not string", c.value))
	}
	return v
}

func (c *ValueInterface) Int64() int64 {
	v, ok := c.value.(int64)
	if !ok {
		panic(fmt.Sprintf("Value %v is not int64", c.value))
	}
	return v
}

func (c *ValueInterface) Int() int {
	v, ok := c.value.(int)
	if !ok {
		panic(fmt.Sprintf("Value %v is not int", c.value))
	}
	return v
}

func (c *ValueInterface) Float64() float64 {
	v, ok := c.value.(float64)
	if !ok {
		panic(fmt.Sprintf("Value %v is not float64", c.value))
	}
	return v
}

func (c *ValueInterface) Map() map[interface{}]interface{} {
	v, ok := c.value.(map[interface{}]interface{})
	if !ok {
		panic(fmt.Sprintf("Value %v is not map[interface{}]interface{}", c.value))
	}
	return v
}
