package main

import (
	"encoding/json"
	"log"
	"net/http"
	"tvm-light/proto/tm"
	t_server "tvm-light/trans"
)

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
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}