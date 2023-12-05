package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"dms/src/checkURL"
	"dms/src/structs"
	"dms/src/workfile"
)

type JsonRequest struct {
	Method       string  `json:"method"`
	Table        string  `json:"table"`
	Link         string  `json:"link"`
	IP           *string `json:"ip,omitempty"`
	TimeInterval *string `json:"time_interval,omitempty"`
}

func (request *JsonRequest) ProcessRequest(conn net.Conn, links *structs.Link_source, statistics *structs.Queue) (string, error) {
	switch request.Method {
	case "POST":
		return request.processPostRequest(links, statistics)
	case "GET":
		return request.processGetRequest(conn, links, statistics)
	default:
		return "", errors.New("unsupported method")
	}
}

func (request *JsonRequest) processPostRequest(links *structs.Link_source, statistics *structs.Queue) (string, error) {
	switch request.Table {
	case "links":
		if !checkURL.CheckURL(request.Link) {
			return "", errors.New("link is not available")
		}
		var link string
		var err error
		base := strconv.Itoa(links.Original_link.Cout + 1)
		if links.Original_link.Cout != 0 {
			link, err = links.Original_link.Hget(request.Link)
			if err == nil {
				return link, nil
			} else {
				err = links.Original_link.Hset(request.Link, base)
				if err != nil {
					return "", err
				}
				links.Short_link.Hset(base, request.Link)
				if err != nil {
					return "", err
				}
				return base, nil
			}
		} else {
			err = links.Original_link.Hset(request.Link, base)
			if err != nil {
				return "", err
			}
			links.Short_link.Hset(base, request.Link)
			if err != nil {
				return "", err
			}
			return base, nil
		}
	case "statistics":
		err := statistics.Qpush(*request.IP + "\n" + request.Link + "\n" + *request.TimeInterval)
		return "", err
	default:
		return "", errors.New("table not found")
	}
}

func (request *JsonRequest) processGetRequest(conn net.Conn, links *structs.Link_source, statistics *structs.Queue) (string, error) {
	switch request.Table {
	case "links":
		if links.Short_link.Cout != 0 {
			link, err := links.Short_link.Hget(request.Link)
			if err != nil {
				return "", err
			}
			return link, nil
		} else {
			return "", errors.New("table is clear")
		}

	case "statistics":
		encoder := gob.NewEncoder(conn)
		err := encoder.Encode(statistics)
		if err != nil {
			return "", err
		}
		return "", nil
	default:
		return "", errors.New("table not found")
	}
}

func main() {
	var links structs.Link_source
	err := workfile.LoadFromGobFile(&links, "data/links.gob")
	if err != nil {
		fmt.Println(err)
		links = structs.NewLinkSource(1024)
	}

	statistics := structs.Queue{}
	err = workfile.LoadFromGobFile(&statistics, "data/statistics.gob")
	if err != nil {
		fmt.Println(err)
	}

	address := "127.0.0.1:6379"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error when starting the server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("The server is listening on:", address)

	var mutex sync.Mutex

	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(links *structs.Link_source, statistics *structs.Queue) {
			defer wg.Done()
			for {
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("Error accepting connection:", err)
					continue
				}
				go handleConnection(conn, &mutex, links, statistics)
			}
		}(&links, &statistics)
	}

	go func(links *structs.Link_source, statistics *structs.Queue) {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			mutex.Lock()
			workfile.SaveToGobFile(links, "data/links.gob")
			workfile.SaveToGobFile(statistics, "data/statistics.gob")
			mutex.Unlock()
		}
	}(&links, &statistics)

	wg.Wait()
}

func handleConnection(conn net.Conn, mutex *sync.Mutex, links *structs.Link_source, statistics *structs.Queue) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)

	var request JsonRequest
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	fmt.Println(request.Method, request.Table, request.Link)

	mutex.Lock()
	value, err := request.ProcessRequest(conn, links, statistics)
	mutex.Unlock()

	if err != nil {
		response := fmt.Sprintln(err)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing response:", err)
			return
		}
	} else {
		response := fmt.Sprintln(value)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing response:", err)
			return
		}
	}
}
