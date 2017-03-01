package main

import (
	"fmt"

	"git.vodjk.com/totoro/totoro/common/util"
	"gopkg.in/yaml.v2"
)

// Configer 配置文件
type Configer struct {
	Server struct {
		TCPAddr string `yaml:"tcp_addr"`
		UDPAddr string `yaml:"udp_addr"`
	} `yaml:"server"`
	DNS struct {
		Mode int      `yaml:"mode"`
		IN   []string `yaml:"in"`
		OUT  []string `yaml:"out"`
	} `yaml:"dns"`
	Rule []string `yaml:"rule"`
}

// NewConfig 读取配置文件，实例化配置并返回
func NewConfig(filename string) (*Configer, error) {
	c := new(Configer)
	data, err := util.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("config.parse error: %s", err.Error())
	}
	err = yaml.Unmarshal(data, c)
	return c, nil
}
