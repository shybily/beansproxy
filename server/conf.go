package server

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func NewYamlProxyOptions(filename string) (*ProxyOptions, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	opt := &ProxyOptions{}
	if err := yaml.Unmarshal(data, opt); err != nil {
		return nil, err
	}
	return opt, nil
}

func NewJsonProxyOptions(filename string) (*ProxyOptions, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	opt := &ProxyOptions{}
	if err := json.Unmarshal(data, opt); err != nil {
		return nil, err
	}
	return opt, nil
}

type ProxyOptions struct {
	Instances      []string `json:"instances" yaml:"instances"`
	Listen         string   `json:"listen" yaml:"listen"`
	MaxConnections int      `json:"max_connections" yaml:"max_connections"`
	Strategy       string   `json:"strategy" yaml:"strategy"`
}
