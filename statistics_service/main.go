package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"

	"status_servis/src/response"
	"status_servis/src/structs"
)

type JsonRequest struct {
	Method       string  `json:"method"`
	Table        string  `json:"table"`
	Link         string  `json:"link"`
	IP           *string `json:"ip,omitempty"`
	TimeInterval *string `json:"time_interval,omitempty"`
}

func (request *JsonRequest) ProcessRequest(conn net.Conn) error {
	switch request.Method {
	case "POST":
		return request.processPostRequest()
	case "GET":
		return request.processGetRequest(conn)
	default:
		return errors.New("unsupported method")
	}
}

func (request *JsonRequest) processPostRequest() error {
	address := "127.0.0.1:6379"

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return err
	}
	defer conn.Close()

	requestData := JsonRequest{
		Method:       "POST",
		Table:        "statistics",
		Link:         request.Link,
		IP:           request.IP,
		TimeInterval: request.TimeInterval,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Error sending JSON data:", err)
		return err
	}

	return nil
}

func (request *JsonRequest) processGetRequest(conn net.Conn) error {
	address := "127.0.0.1:6379"

	conn1, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return err
	}
	defer conn.Close()
	requestData := JsonRequest{
		Method:       "GET",
		Table:        "statistics",
		Link:         request.Link,
		IP:           request.IP,
		TimeInterval: request.TimeInterval,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	_, err = conn1.Write(jsonData)
	if err != nil {
		fmt.Println("Error sending JSON data:", err)
		return err
	}

	var gobResponse structs.Queue
	decoder := gob.NewDecoder(conn1)
	err = decoder.Decode(&gobResponse)
	if err != nil {
		fmt.Println("Error decoding gob:", err)
		return err
	}

	link, _ := strconv.Atoi(request.Link)
	ip, _ := strconv.Atoi(*request.IP)
	time, _ := strconv.Atoi(*request.TimeInterval)
	rep := response.NewProcessor(&link, &ip, &time)
	js := response.JsonResponse{}
	js.CreateReport(&gobResponse, *rep)

	jsonData, err = json.Marshal(js)
	if err != nil {
		fmt.Println("Ошибка при кодировании в JSON:", err)
		return err
	}

	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Ошибка при отправке данных:", err)
		return err
	}
	return nil
}

func main() {
	address := "127.0.0.1:1333"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error when starting the server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("The server is listening on:", address)

	var mutex sync.Mutex

	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("Error accepting connection:", err)
					continue
				}
				go handleConnection(conn, &mutex)
			}
		}()
	}
	wg.Wait()
}

func handleConnection(conn net.Conn, mutex *sync.Mutex) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)

	var request JsonRequest
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	mutex.Lock()
	err = request.ProcessRequest(conn)
	fmt.Println(err)
	mutex.Unlock()
}
