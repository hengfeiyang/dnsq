package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Configer 配置文件
type Configer struct {
	Server struct {
		TCPAddr string `yaml:"tcp_addr"`
		UDPAddr string `yaml:"udp_addr"`
	} `yaml:"server"`
	DNS struct {
		IN  []string `yaml:"in"`
		OUT []string `yaml:"out"`
	} `yaml:"dns"`
	Rule []string `yaml:"rule"`
}

// NewConfig 读取配置文件，实例化配置并返回
func NewConfig(filename string) (*Configer, error) {
	c := new(Configer)
	data, err := ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("config.parse error: %s", err.Error())
	}
	err = yaml.Unmarshal(data, c)
	return c, nil
}

// ReadFile read file content
func ReadFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
