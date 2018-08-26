package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
)

type Request struct {
	RequestType string   `json:"RequestType"`
	RequestBody KeyValue `json:"RequestBody"`
}

type Server struct {
	Host       string
	Port       string
	StoredData map[string]string
}

type Response struct {
	ResponseStatus string   `json:"ResponseStatus"`
	ResponseBody   KeyValue `json:"ResponseBody"`
}

type KeyValue struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func (s *Server) GetConnectionString() string {
	return s.Host + ":" + s.Port
}

func (s *Server) HandleRequest(request []byte) (KeyValue, error) {

	fmt.Println(string(request))
	request = bytes.Trim(request, "\x00")

	var req Request
	err := json.Unmarshal(request, &req)

	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println("req: ", req)
	fmt.Println("RequestType: ", req.RequestType)
	switch req.RequestType {
	case "GET":
		return s.get(req.RequestBody), nil
	case "PUT":
		return s.store(req.RequestBody), nil
	default:
		fmt.Println("Type not found: ", req.RequestType)
	}
	return KeyValue{
		Key:   "",
		Value: "",
	}, errors.New("Type not Found")
}

func (s *Server) store(kv KeyValue) KeyValue {
	fmt.Println("Store Called")
	s.StoredData[kv.Key] = kv.Value
	fmt.Println("Data Stored.")
	return kv
}

func (s *Server) get(kv KeyValue) KeyValue {
	fmt.Println("Get Called")
	result := s.StoredData[kv.Key]

	return KeyValue{
		Key:   kv.Key,
		Value: result,
	}
}

func main() {

	server := Server{
		Host:       "localhost",
		Port:       "3333",
		StoredData: make(map[string]string),
	}

	s, err := net.Listen("tcp", server.GetConnectionString())

	if err != nil {
		fmt.Println("Failed to open socket.")
		fmt.Println("Connection string is ", server.GetConnectionString())
		os.Exit(1)
	}

	fmt.Println("Server listens on ", server.GetConnectionString())

	for {
		connection, err := s.Accept()

		if err != nil {
			fmt.Println("Connection accept is failure.")
			os.Exit(1)
		}
		buf := make([]byte, 1024)

		reqLen, err := connection.Read(buf)
		fmt.Println("Buffer read size is ", reqLen)

		kv, err := server.HandleRequest(buf)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(kv)

		resp := Response{
			ResponseStatus: "success",
			ResponseBody: KeyValue{
				Key:   kv.Key,
				Value: kv.Value,
			},
		}

		response, err := json.Marshal(resp)

		fmt.Println(string(response))

		connection.Write([]byte(response))

		connection.Close()
	}

	// When the process stops, close the socket.
	defer s.Close()
}
