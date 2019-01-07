package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type TVMconf struct {
	OrderServer    string `yaml:"orderServer"`
	ContractPath   string `yaml:"contractPath"`
	ChannelID      string `yaml:"channelID"`
	OrdererOrgName string `yaml:"ordererOrgName"`
	IPFSAddress    string `yaml:"IPFSAddress"`
	DockerPath     string `yaml:"dockerPath"`
	Port           string `yml:"port"`
}

var triasConfig = TVMconf{}

func init() {
	var filePath = "/opt/gopath/src/tvm-light/config.yml"
	data, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(data, &triasConfig)
}

func GetOrderServer() string {
	return triasConfig.OrderServer;
}
func GetContractPath() string {
	return triasConfig.ContractPath;
}

func GetChannelID() string {
	return triasConfig.ChannelID;
}

func GetOrdererOrgName() string {
	return triasConfig.OrdererOrgName
}

func GetIPFSAddress() string {
	return triasConfig.IPFSAddress
}

func GetDockerPath() string {
	return triasConfig.DockerPath
}

func GetPort() string {
	return triasConfig.Port
}
