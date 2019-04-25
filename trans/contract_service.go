package trans

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"strings"
	t_conf "tvm-light/config"
	"tvm-light/contract"
	"tvm-light/proto/tm"
	t_utils "tvm-light/utils"
	"tvm-light/validate"
)

func NewContractService() *contractServer {
	return &contractServer{}
}

const (
	fileSuffix           = ".go"
	error_code           = -1
	error_not_install    = -100
	dbname_separator     = "_"
)

type contractServer struct {
}

func (serv *contractServer) ExecuteContract(ctx context.Context, request *tm.ExecuteContractRequest) (*tm.ExecuteContractResponse) {
	isCorect, err := validate.RequestValidate(request);
	if !isCorect || err != nil {
		fmt.Println("Contract validate fails", err);
		return returnErrorResponse(err, error_code);
	}
	var filePath = t_conf.TriasConfig.ContractPath + "/" + request.GetUser() + "/" + request.GetAddress() + "/" + request.GetContractName() + "/";
	var fileName = request.GetContractName() + fileSuffix;
	isExists, err := t_utils.PathExists(filePath + fileName);
	if err != nil {
		fmt.Println("checkFilePathFails", err);
		return returnErrorResponse(err, error_code);
	}
	cont := contract.NewContract(t_conf.TriasConfig.OrderServer, request.GetContractName(), request.GetContractType(), filePath, request.GetContractVersion(), t_conf.TriasConfig.ChannelID, t_conf.TriasConfig.OrdererOrgName, request.GetCommand(), request.GetOperation())
	if !isExists {
		if strings.EqualFold(request.GetOperation(), "install") {
			err := t_utils.FileDownLoad(filePath, fileName, t_conf.TriasConfig.IPFSAddress+request.GetAddress());
			if err != nil {
				fmt.Println("Download contract happens a error", err);
				return returnErrorResponse(err, error_code);
			}
		} else {
			dbName := t_utils.GetCouchDBName(t_conf.TriasConfig.ChannelID + dbname_separator + request.ContractName)
			dbExists := t_utils.CheckCouchDBExists(dbName)
			if !dbExists {
				return returnErrorResponse(errors.New("Contract hasn't been installed"), error_not_install);
			} else {
				err := t_utils.FileDownLoad(filePath, fileName, t_conf.TriasConfig.IPFSAddress+request.GetAddress());
				if err != nil {
					fmt.Println("Download contract happens a error", err);
					return returnErrorResponse(err, error_code);
				}
			}
		}
	}
	if _, err := t_utils.CheckFileMD5(filePath+fileName, request.GetCheckMD5()); err != nil {
		return returnErrorResponse(err, error_code);
	}
	cresp, err := cont.RunContract()
	if err != nil {
		fmt.Println("Failed to run contract", err);
		return returnErrorResponse(err, error_code);
	} else {
		c_hash := calculateHash(request)
		contract.UpdateCurrentHash(t_conf.BasicHashKey, c_hash)
	}

	resp := &tm.ExecuteContractResponse{
		Code:    1,
		Data:    cresp,
		Message: "success",
	}

	return resp;
}

func returnErrorResponse(err error, code int32) (*tm.ExecuteContractResponse) {
	resp := &tm.ExecuteContractResponse{
		Code:    code,
		Data:    err.Error(),
		Message: "fail",
	}
	return resp;
}

func calculateHash(request *tm.ExecuteContractRequest) string {
	var message string = string(request.GetContractName() + request.GetUser() + request.GetOperation() + string(request.GetTimestamp()) + request.GetCheckMD5());
	return message;
}

