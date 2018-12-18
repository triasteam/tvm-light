package main

import (
	"encoding/json"
	"log"
	"net/http"
	tm "tvm-light/proto/tm"
	t_server "tvm-light/trans"
)


type Contract struct {
	address string `json:"address"`
	checkMD5 string `json:"checkMD5"`
	command string `json:"command"`
	contractName string `json:"contractName"`
	contractType string `json:"contractType"`
	vmVersion string `json:"vmVersion"`
	sequence string `json:"sequence"`
	timestamp int32 `json:"timestamp"`
	user string `json:"user"`
	signature string `json:"signature"`
	opration string `json:"opration"`

}

type ExecuteResponse struct {
	code int32 `json:"code"`
	message string `json:"message"`
	data string `json:"data"`
}

func executeContract(writer http.ResponseWriter, request *http.Request) {
	var contractRequest *tm.ExecuteContractRequest;
	if err:= json.NewDecoder(request.Body).Decode(&contractRequest); err != nil {
		request.Body.Close()
		log.Fatal(err)
	}
	server := t_server.NewMWService();
	response := server.ExecuteContract(nil, contractRequest)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/executeContract", executeContract)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}