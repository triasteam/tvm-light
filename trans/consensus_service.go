package trans

import (
	"encoding/json"
	"golang.org/x/net/context"
	"os"
	t_conf "tvm-light/config"
	t_utils "tvm-light/utils"
)

func NewConsensusService() *consensusServer {
	return &consensusServer{}
}

type consensusServer struct{}

type AsyncTVMRequest struct {
	IpfsHash string `json:"ipfsHash,omitempty"`
}
type uploadResponseBody struct {
	IpfsHash string `json:"ipfsHash"`
}

const (
	tar_suffix = ".tar.gz"
)

func (s *consensusServer) UploadData() *t_conf.CommonResponse {
	// package data
	var filePath string = t_conf.TriasConfig.PackagePath + "/data.tar.gz"
	err := t_utils.Compress(t_conf.TriasConfig.DataPath, filePath)
	if err != nil {
		return createErrorCommonResponse(err, -1)
	}
	// uploadIpfs
	hash, err := t_utils.AddFile(filePath)
	if err != nil {
		return createErrorCommonResponse(err, -1)
	}
	data := uploadResponseBody{
		IpfsHash: hash,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return createErrorCommonResponse(err, -1)
	}
	return createSuccessCommonResponse(string(jsonData))
}

func (s *consensusServer) AsyncTVM(ctx context.Context, request *AsyncTVMRequest) *t_conf.CommonResponse {
	// download package from ipfs
	var fileName string = string(t_conf.TriasConfig.PackagePath + request.IpfsHash + tar_suffix)
	if err := t_utils.GetFile(request.IpfsHash, fileName); err != nil {
		return createErrorCommonResponse(err, -1)
	}
	// stop docker-compose
	if err := t_utils.StopTVM(); err != nil {
		return createErrorCommonResponse(err, -1)
	}
	// clean data files
	if err := os.RemoveAll(t_conf.TriasConfig.DataPath); err != nil {
		return createErrorCommonResponse(err, -1)
	}
	// decomporess data file
	if err := t_utils.DeCompress(fileName, t_conf.TriasConfig.DataPath[:len(t_conf.TriasConfig.DataPath)-5]); err != nil {
		return createErrorCommonResponse(err, -1)
	}
	// chown 5984.5984
	if err := t_utils.ModifyPathUserGroup(t_conf.TriasConfig.CouchdbInfo.Path,t_conf.TriasConfig.CouchdbInfo.Port,t_conf.TriasConfig.CouchdbInfo.Port); err != nil {
		return createErrorCommonResponse(err, -1)
	}
	// start docker-compose
	if err := t_utils.StartTVM(); err != nil {
		return createErrorCommonResponse(err, -1)
	}
	return createSuccessCommonResponse("")
}

func createErrorCommonResponse(err error, code int32) *t_conf.CommonResponse {
	resp := &t_conf.CommonResponse{
		Code:    code,
		Data:    err.Error(),
		Message: "fail",
	}
	return resp;
}

func createSuccessCommonResponse(data string) *t_conf.CommonResponse {
	resp := &t_conf.CommonResponse{
		Code:    1,
		Data:    data,
		Message: "success",
	}
	return resp;
}

