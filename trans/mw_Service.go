package trans

import (
	"fmt"
	"golang.org/x/net/context"
	enConf "triasVM/config"
	t_utils "triasVM/utils"
	"tvm-light/contract"
	"tvm-light/proto/tm"
	"tvm-light/validate"
)

func NewMWService() *server {
	return &server{}
}

const (
	pathPrefix = "./source/contract/"
	fileSuffix = ".go"
)

type server struct {
}

func (serv *server) ExecuteContract(ctx context.Context, request *tm.ExecuteContractRequest) (*tm.ExecuteContractResponse) {
	// TODO validate
	isCorect, err := validate.RequestValidate(request);
	if !isCorect || err != nil {
		fmt.Println("Contract validate fails", err);
		return returnErrorResponse(err);
	}
	// TODO CheckContract is install
	var filePath = pathPrefix + request.GetContractName() + fileSuffix;
	isExists, err := t_utils.PathExists(filePath);
	if err != nil {
		fmt.Println("checkFilePathFails", err);
		return returnErrorResponse(err);
	}
	contract := contract.NewContract(enConf.PeerAddress, request.GetContractName(), request.GetContractType(), filePath, enConf.ContractVersion, enConf.ChannelID, enConf.OrdererOrgName, request.GetCommand(), request.GetOpration());
	if !isExists {
		err := t_utils.FileDownLoad(filePath, request.GetAddress());
		if err != nil {
			fmt.Println("Download contract happens a error", err);
			return returnErrorResponse(err);
		}
		installErr := contract.InstallContract();
		if installErr != nil {

			return returnErrorResponse(err);
		}
	}

	cresp, err := contract.RunContract()
	if err != nil {
		fmt.Println("Failed to run contract", err);
		return returnErrorResponse(err);
	}

	resp := &tm.ExecuteContractResponse{
		Code:    1,
		Data:    string(cresp),
		Message: "success",
	}

	return resp;
}

func returnErrorResponse(err error) (*tm.ExecuteContractResponse) {
	resp := &tm.ExecuteContractResponse{
		Code:    0,
		Data:    err.Error(),
		Message: "faile",
	}
	return resp;
}

func returnSuccessResponse(data string) (*tm.ExecuteContractResponse) {
	resp := &tm.ExecuteContractResponse{
		Code:    1,
		Data:    data,
		Message: "success",
	}
	return resp;
}
