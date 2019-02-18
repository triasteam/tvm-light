package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type TVMconf struct {
	OrderServer     string `yaml:"orderServer"`
	ContractPath    string `yaml:"contractPath"`
	ChannelID       string `yaml:"channelID"`
	OrdererOrgName  string `yaml:"ordererOrgName"`
	IPFSAddress     string `yaml:"IPFSAddress"`
	DockerPath      string `yaml:"dockerPath"`
	Port            string `yaml:"port"`
	IpfsAPIAddress  string `yaml:"ipfsAPIAddress"`
	GOPATH          string `yaml:"goPath"`
	ComposeFilePath string `yaml:"conposeFilePath"`
	DataPath        string `yaml:"dataPath"`
	PackagePath     string `yaml:"packagePath"`
	CouchdbPort     int `yaml:"couchdbPort"`
	CouchdbPath     string `yaml:"couchdbPath"`
}

var TriasConfig = TVMconf{}

func init() {
	var filePath = "/opt/gopath/src/tvm-light/config.yml"
	data, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(data, &TriasConfig)
}

func GetOrderServer() string {
	return TriasConfig.OrderServer;
}
func GetContractPath() string {
	return TriasConfig.ContractPath;
}

func GetChannelID() string {
	return TriasConfig.ChannelID;
}

func GetOrdererOrgName() string {
	return TriasConfig.OrdererOrgName
}

func GetIPFSAddress() string {
	return TriasConfig.IPFSAddress
}

func GetDockerPath() string {
	return TriasConfig.DockerPath
}

func GetPort() string {
	return TriasConfig.Port
}
