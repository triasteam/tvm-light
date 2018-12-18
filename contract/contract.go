package contract

import (
	"io/ioutil"
	"log"
	"os/exec"
	tvm_conf "tvm-light/config"
)

type Contract struct {
	peerAddress     string
	contractName    string
	contractType    string
	contractPath    string
	contractVersion string
	channelID       string
	orgName         string
	args            string
	action          string
}

const (
	CMD_PEER = "peer"
)
var docker_command []string = []string{"docker","exec","cli"}

func NewContract(peerAddress string, contractName string, contractType string, contractPath string, contractVersion string, channelID string, orgName string, args string, action string) *Contract {
	return &Contract{peerAddress: peerAddress, contractName: contractName, contractType: contractType, contractPath: contractPath, contractVersion: contractVersion, channelID: channelID, orgName: orgName, args: args, action: action}
}

func (c *Contract) InstallContract() error {
	params := []string{"chaincode", "install", "-n", c.contractName, "-p", c.contractPath, "-v", c.contractVersion};
	cmd := exec.Command(CMD_PEER, append(docker_command,params...)...);
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		return err
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println(string(opBytes))
	return nil;
}

func (c *Contract) RunContract() ([]byte, error) {

	var resp []byte;
	switch c.action {
	case "instantiate":
		return c.instantiate()
		break;
	default:
		return c.execute()
		break;
	}
	return resp, nil;
}

func (c *Contract) instantiate() ([]byte, error) {
	params := []string{"chaincode", "instantiate", "-o", tvm_conf.OrderServer, "-C", c.channelID, "-n", c.contractName, "-v", c.contractVersion, "-c", c.args};
	cmd := exec.Command(CMD_PEER, append(docker_command,params...)...);
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return opBytes, err;
}

func (c *Contract) execute() ([]byte, error) {
	// peer chaincode $action -C $channelID -n $cname -c $args
	params := []string{"chaincode", c.action, "-C", c.channelID, "-n", c.contractName, "-c", c.args};
	cmd := exec.Command(CMD_PEER, append(docker_command,params...)...);
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return opBytes, err
}
