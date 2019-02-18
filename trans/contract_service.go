package trans

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"strings"
	enConf "tvm-light/config"
	"tvm-light/contract"
	"tvm-light/proto/tm"
	t_utils "tvm-light/utils"
	"tvm-light/validate"
)

func NewContractService() *contractServer {
	return &contractServer{}
}

const (
	fileSuffix = ".go"
	error_code = -1
	error_not_install = -100
)

type contractServer struct {
}

func (serv *contractServer) ExecuteContract(ctx context.Context, request *tm.ExecuteContractRequest) (*tm.ExecuteContractResponse) {
	// TODO validate
	isCorect, err := validate.RequestValidate(request);
	if !isCorect || err != nil {
		fmt.Println("Contract validate fails", err);
		return returnErrorResponse(err,error_code);
	}
	// TODO CheckContract is install
	var filePath = enConf.GetContractPath() + "/" + request.GetUser() + "/" + request.GetAddress() + "/"+ request.GetContractName() + "/";
	var fileName = request.GetContractName() + fileSuffix;
	isExists, err := t_utils.PathExists(filePath + fileName);
	if err != nil {
		fmt.Println("checkFilePathFails", err);
		return returnErrorResponse(err,error_code);
	}
	contract := contract.NewContract(enConf.GetOrderServer(), request.GetContractName(), request.GetContractType(), filePath, request.GetContractVersion(), enConf.GetChannelID(), enConf.GetOrdererOrgName(), request.GetCommand(), request.GetOperation());
	if !isExists {
		if strings.EqualFold(request.GetOperation(),"install"){
			err := t_utils.FileDownLoad(filePath, fileName, enConf.GetIPFSAddress()+request.GetAddress());
			if err != nil {
				fmt.Println("Download contract happens a error", err);
				return returnErrorResponse(err,error_code);
			}
		} else {
			return returnErrorResponse(errors.New("Contract hasn't been installed"),error_not_install);
		}
	}
	if _, err := t_utils.CheckFileMD5(filePath+fileName, request.GetCheckMD5()); err != nil {
		return returnErrorResponse(err,error_code);
	}
	cresp, err := contract.RunContract()
	if err != nil {
		fmt.Println("Failed to run contract", err);
		return returnErrorResponse(err,error_code);
	}

	if strings.HasSuffix(cresp,"\n") {
		cresp = cresp[:len(cresp)-1]
	}

	resp := &tm.ExecuteContractResponse{
		Code:    1,
		Data:    cresp,
		Message: "success",
	}

	return resp;
}

func returnErrorResponse(err error,code int32) (*tm.ExecuteContractResponse) {
	resp := &tm.ExecuteContractResponse{
		Code:    code,
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
