package trans

import (
	"fmt"
	"golang.org/x/net/context"
	enConf "tvm-light/config"
	"tvm-light/contract"
	"tvm-light/proto/tm"
	t_utils "tvm-light/utils"
	"tvm-light/validate"
)

func NewMWService() *server {
	return &server{}
}

const (
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
	var filePath = enConf.GetContractPath() + "/" + request.GetUser() + "/" + request.GetAddress() + "/"+ request.GetContractName() + "/";
	var fileName = request.GetContractName() + fileSuffix;
	isExists, err := t_utils.PathExists(filePath + fileName);
	if err != nil {
		fmt.Println("checkFilePathFails", err);
		return returnErrorResponse(err);
	}
	contract := contract.NewContract(enConf.GetOrderServer(), request.GetContractName(), request.GetContractType(), filePath, request.GetContractVersion(), enConf.GetChannelID(), enConf.GetOrdererOrgName(), request.GetCommand(), request.GetOperation());
	if !isExists {
		err := t_utils.FileDownLoad(filePath, fileName, enConf.GetIPFSAddress()+request.GetAddress());
		if err != nil {
			fmt.Println("Download contract happens a error", err);
			return returnErrorResponse(err);
		}
	}
	if _, err := t_utils.CheckFileMD5(filePath+fileName, request.GetCheckMD5()); err != nil {
		return returnErrorResponse(err);
	}
	cresp, err := contract.RunContract()
	if err != nil {
		fmt.Println("Failed to run contract", err);
		return returnErrorResponse(err);
	}

	resp := &tm.ExecuteContractResponse{
		Code:    1,
		Data:    cresp,
		Message: "success",
	}

	return resp;
}

func returnErrorResponse(err error) (*tm.ExecuteContractResponse) {
	resp := &tm.ExecuteContractResponse{
		Code:    -1,
		Data:    err.Error(),
		Message: "fail",
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
