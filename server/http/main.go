package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	t_conf "tvm-light/config"
	"tvm-light/proto/tm"
	t_server "tvm-light/trans"
)

func executeContract(writer http.ResponseWriter, request *http.Request) {
	var contractRequest *tm.ExecuteContractRequest;
	if err := json.NewDecoder(request.Body).Decode(&contractRequest); err != nil {
		fmt.Println(err)
		request.Body.Close()
	}
	server := t_server.NewContractService()
	response := server.ExecuteContract(nil, contractRequest)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		fmt.Println(err)
	}
}

func upLoadPackage(writer http.ResponseWriter, request *http.Request) {
	server := t_server.NewConsensusService()
	response := server.UploadData()
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		fmt.Println(err)
	}

}

func asyncTVM(writer http.ResponseWriter, request *http.Request) {
	var asyncRequest *t_server.AsyncTVMRequest
	if err := json.NewDecoder(request.Body).Decode(&asyncRequest);err != nil {
		fmt.Println(err)
		log.Fatal(err)
		request.Body.Close()
	}
	server := t_server.NewConsensusService()
	response := server.AsyncTVM(nil,asyncRequest)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		fmt.Println(err)
	}
}

func main() {
	http.HandleFunc("/executeContract", executeContract)
	http.HandleFunc("/upLoadPackage", upLoadPackage)
	http.HandleFunc("/asyncTVM", asyncTVM)
	fmt.Println(t_conf.TriasConfig.Port)
	err := http.ListenAndServe(":"+t_conf.GetPort(), nil)
	if err != nil {
		fmt.Println(err)
	}
}
