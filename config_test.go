package testcase

import (
	"testing"

	"bitbucket.org/jellynian/contrib/config"
)

func TestConfig(t *testing.T) {
	conf := config.Default()
	t.Log("read result is :", conf.Get("user").String())

	t.Log("read result is :", conf.Get("global.addr").String())
}

func TestConfigOnly(t *testing.T) {
	conf := config.Default()
	t.Log("read result is :", conf.Get("user").String())
}
