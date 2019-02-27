package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type TVMconf struct {
	OrderServer     string  `yaml:"orderServer"`
	ContractPath    string  `yaml:"contractPath"`
	ChannelID       string  `yaml:"channelID"`
	OrdererOrgName  string  `yaml:"ordererOrgName"`
	IPFSAddress     string  `yaml:"IPFSAddress"`
	DockerPath      string  `yaml:"dockerPath"`
	Port            string  `yaml:"port"`
	IpfsAPIAddress  string  `yaml:"ipfsAPIAddress"`
	GOPATH          string  `yaml:"goPath"`
	ComposeFilePath string  `yaml:"conposeFilePath"`
	DataPath        string  `yaml:"dataPath"`
	PackagePath     string  `yaml:"packagePath"`
	CouchdbInfo     CouchDB `yaml:"couchdbInfo"`
}

type CouchDB struct {
	Port     int `yaml:"port"`
	Path     string `yaml:"path"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var TriasConfig = TVMconf{}

func init() {
	var filePath = "config.yml"
	data, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(data, &TriasConfig)
}



func GetPort() string {
	return TriasConfig.Port
}
