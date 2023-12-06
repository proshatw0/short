package response

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"status_servis/src/structs"
)

type JsonResponseNode struct {
	Id           int     `json:"id"`
	Pid          int     `json:"pid"`
	Link         *string `json:"link,omitempty"`
	IP           *string `json:"ip,omitempty"`
	TimeInterval *string `json:"time_interval,omitempty"`
	Cout         int     `json:"cout"`
}

type JsonResponse struct {
	Table []*JsonResponseNode
	Size  int
}

type Processor struct {
	First  *int
	Second *int
	Third  *int
}

func NewProcessor(link, ip, time *int) *Processor {
	value := []int{0, 1, 2}

	var p Processor

	switch *link {
	case 1:
		p.First = &value[1]
	case 2:
		p.Second = &value[1]
	case 3:
		p.Third = &value[1]
	}

	switch *ip {
	case 1:
		p.First = &value[0]
	case 2:
		p.Second = &value[0]
	case 3:
		p.Third = &value[0]
	}

	switch *time {
	case 1:
		p.First = &value[2]
	case 2:
		p.Second = &value[2]
	case 3:
		p.Third = &value[2]
	}

	return &p
}

func (js *JsonResponse) Append(id int, pid int, link, ip, time *string, cout int) {
	node := &JsonResponseNode{
		Id:           id,
		Pid:          pid,
		Link:         link,
		IP:           ip,
		TimeInterval: time,
		Cout:         cout,
	}

	js.Table = append(js.Table, node)
	js.Size++
}

func (js *JsonResponse) parentalElement(value *string, status int, pids *structs.Hash_Table) {
	switch status {
	case 0:
		js.Append(js.Size+1, 0, nil, value, nil, 1)
	case 1:
		js.Append(js.Size+1, 0, value, nil, nil, 1)
	case 2:
		js.Append(js.Size+1, 0, nil, nil, value, 1)
	}
	pids.Hset(*value, strconv.Itoa(js.Size))
}

func (js *JsonResponse) subElement(value *string, status int, id int, pids *structs.Hash_Table) {
	fmt.Println(*value, status)
	if id == 0 {
		switch status {
		case 0:
			js.Append(js.Size+1, js.Size, nil, value, nil, 1)
		case 1:
			js.Append(js.Size+1, js.Size, value, nil, nil, 1)
		case 2:
			js.Append(js.Size+1, js.Size, nil, nil, value, 1)
		}
		pids.Hset(strconv.Itoa(js.Size-1), *value)
	} else {
		switch status {
		case 0:
			js.Append(js.Size+1, id, nil, value, nil, 1)
		case 1:
			js.Append(js.Size+1, id, value, nil, nil, 1)
		case 2:
			js.Append(js.Size+1, id, nil, nil, value, 1)
		}
	}
	err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
	if err != nil {
		val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
		val += "\n" + strconv.Itoa(js.Size)
		pids.Hset(strconv.Itoa(js.Size-1), val)
	}
}

func (js *JsonResponse) processingSubElement(pi *string, status *int, value *string) int {
	pidarr := strings.Split(*pi, "\n")
	for _, val := range pidarr {
		pidu, _ := strconv.Atoi(val)
		switch *status {
		case 0:
			if js.Table[pidu-1].IP != nil && *js.Table[pidu-1].IP != *value {
				continue
			}
		case 1:
			if js.Table[pidu-1].Link != nil && *js.Table[pidu-1].Link != *value {
				continue
			}
		case 2:
			if js.Table[pidu-1].TimeInterval != nil && *js.Table[pidu-1].TimeInterval != *value {
				continue
			}

		}
		js.Table[pidu-1].Cout++
		return pidu
	}
	return 0
}

func (js *JsonResponse) CreateReport(statistics *structs.Queue, status Processor) {
	links := structs.NewHashTable(10)
	pids := structs.NewHashTable(10)
	for statistics.Head != nil {
		value := strings.Split(statistics.Head.Data, "\n")
		i, _ := links.Hget(value[*status.First])
		if i == "" {
			js.parentalElement(&value[*status.First], *status.First, &links)
			if status.Second == nil {
				statistics.Qpop()
				continue
			}
			js.subElement(&value[*status.Second], *status.Second, 0, &pids)
			if status.Third == nil {
				statistics.Qpop()
				continue
			}
			js.subElement(&value[*status.Third], *status.Third, 0, &pids)
			statistics.Qpop()
			continue
		}
		id, _ := strconv.Atoi(i)
		js.Table[id-1].Cout++
		if status.Second == nil {
			statistics.Qpop()
			continue
		}
		pi, _ := pids.Hget(i)
		sub_index := js.processingSubElement(&pi, status.Second, &value[*status.Second])
		if sub_index == 0 {
			js.subElement(&value[*status.Second], *status.Second, id, &pids)
			if status.Third == nil {
				statistics.Qpop()
				continue
			}
			js.subElement(&value[*status.Third], *status.Third, 0, &pids)
			statistics.Qpop()
			continue
		}
		if status.Third == nil {
			statistics.Qpop()
			continue
		}
		datatime, _ := pids.Hget(strconv.Itoa(sub_index))
		key := js.processingSubElement(&datatime, status.Third, &value[*status.Third])
		if key != 0 {
			statistics.Qpop()
			continue
		} else {
			js.subElement(&value[*status.Third], *status.Third, sub_index, &pids)
			statistics.Qpop()
		}
	}
	jsonData, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		fmt.Println("JОшибка при маршалинге в JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}
