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

func (js *JsonResponse) LinkIpTime(statistics *structs.Queue) {
	links := structs.NewHashTable(10)
	pids := structs.NewHashTable(10)
	for statistics.Head != nil {
		value := strings.Split(statistics.Head.Data, "\n")
		if js.Size == 0 {
			js.Append(js.Size+1, 0, &value[1], nil, nil, 1)
			links.Hset(value[1], strconv.Itoa(js.Size))
			js.Append(js.Size+1, js.Size, nil, &value[0], nil, 1)
			pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
			js.Append(js.Size+1, js.Size, nil, nil, &value[2], 1)
			pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
		} else {
			if i, _ := links.Hget(value[1]); i != "" {
				id, _ := strconv.Atoi(i)
				js.Table[id-1].Cout++
				if pi, _ := pids.Hget(i); pi != "" {
					pidarr := strings.Split(pi, "\n")
					key := 0
					sub_index := 0
					for _, val := range pidarr {
						pidu, _ := strconv.Atoi(val)
						if js.Table[pidu-1].IP != nil && *js.Table[pidu-1].IP == value[0] {
							js.Table[pidu-1].Cout++
							key = 1
							sub_index = pidu
							break
						}
					}
					if key == 0 {
						js.Append(js.Size+1, id, nil, &value[0], nil, 1)
						err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
						if err != nil {
							val, _ := pids.Hdel(strconv.Itoa(id))
							val += "\n" + strconv.Itoa(js.Size)
							pids.Hset(strconv.Itoa(id), val)
						}
						js.Append(js.Size+1, js.Size, nil, nil, &value[2], 1)
						err = pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
						if err != nil {
							val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
							val += "\n" + strconv.Itoa(js.Size)
							pids.Hset(strconv.Itoa(js.Size-1), val)
						}
					}
					if sub_index != 0 {
						datatime, _ := pids.Hget(strconv.Itoa(sub_index))
						keydata := 0
						if len(datatime) >= 1 {
							dataarr := strings.Split(datatime, "\n")
							for _, dat := range dataarr {
								timein, _ := strconv.Atoi(dat)
								if js.Table[timein-1].TimeInterval != nil && *js.Table[timein-1].TimeInterval == value[2] {
									js.Table[timein-1].Cout++
									keydata = 1
									break
								}
							}
						} else {
							js.Append(js.Size+1, sub_index, nil, nil, &value[2], 1)
							err := pids.Hset(strconv.Itoa(sub_index), strconv.Itoa(js.Size))
							if err != nil {
								val, _ := pids.Hdel(strconv.Itoa(sub_index))
								val += "\n" + strconv.Itoa(js.Size)
								pids.Hset(strconv.Itoa(sub_index), val)
							}
						}
						if keydata == 0 {
							js.Append(js.Size+1, sub_index, nil, nil, &value[2], 1)
							err := pids.Hset(strconv.Itoa(sub_index), strconv.Itoa(js.Size))
							if err != nil {
								val, _ := pids.Hdel(strconv.Itoa(sub_index))
								val += "\n" + strconv.Itoa(js.Size)
								pids.Hset(strconv.Itoa(sub_index), val)
							}
						}
					}
				} else {
					js.Table[id].Cout++
					js.Append(js.Size+1, id, nil, &value[0], nil, 1)
					err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
					if err != nil {
						val, _ := pids.Hdel(strconv.Itoa(id))
						val += "\n" + strconv.Itoa(js.Size)
						pids.Hset(strconv.Itoa(id), val)
					}
					js.Append(js.Size+1, js.Size, nil, nil, &value[2], 1)
					err = pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
					if err != nil {
						val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
						val += "\n" + strconv.Itoa(js.Size)
						pids.Hset(strconv.Itoa(js.Size-1), val)
					}
				}
			} else {
				js.Append(js.Size+1, 0, &value[1], nil, nil, 1)
				links.Hset(value[1], strconv.Itoa(js.Size))
				pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
				js.Append(js.Size+1, js.Size, nil, &value[0], nil, 1)
				pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
				js.Append(js.Size+1, js.Size, nil, nil, &value[2], 1)
			}
		}
		statistics.Qpop()
	}
	jsonData, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при маршалинге в JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}

func (js *JsonResponse) LinkTimeIp(statistics *structs.Queue) {
	links := structs.NewHashTable(10)
	pids := structs.NewHashTable(10)
	for statistics.Head != nil {
		value := strings.Split(statistics.Head.Data, "\n")
		if js.Size == 0 {
			js.Append(js.Size+1, 0, &value[1], nil, nil, 1)
			links.Hset(value[1], strconv.Itoa(js.Size))
			js.Append(js.Size+1, js.Size, nil, nil, &value[2], 1)
			pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
			js.Append(js.Size+1, js.Size, nil, &value[0], nil, 1)
			pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
		} else {
			if i, _ := links.Hget(value[1]); i != "" {
				id, _ := strconv.Atoi(i)
				js.Table[id-1].Cout++
				if pi, _ := pids.Hget(i); pi != "" {
					pidarr := strings.Split(pi, "\n")
					key := 0
					sub_index := 0
					for _, val := range pidarr {
						pidu, _ := strconv.Atoi(val)
						if js.Table[pidu-1].TimeInterval != nil && *js.Table[pidu-1].TimeInterval == value[2] {
							js.Table[pidu-1].Cout++
							key = 1
							sub_index = pidu
							break
						}
					}
					if key == 0 {
						js.Append(js.Size+1, id, nil, nil, &value[2], 1)
						err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
						if err != nil {
							val, _ := pids.Hdel(strconv.Itoa(id))
							val += "\n" + strconv.Itoa(js.Size)
							pids.Hset(strconv.Itoa(id), val)
						}
						js.Append(js.Size+1, js.Size, nil, &value[0], nil, 1)
						err = pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
						if err != nil {
							val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
							val += "\n" + strconv.Itoa(js.Size)
							pids.Hset(strconv.Itoa(js.Size-1), val)
						}
					}
					if sub_index != 0 {
						datatime, _ := pids.Hget(strconv.Itoa(sub_index))
						keydata := 0
						if len(datatime) >= 1 {
							dataarr := strings.Split(datatime, "\n")
							for _, dat := range dataarr {
								timein, _ := strconv.Atoi(dat)
								if js.Table[timein-1].IP != nil && *js.Table[timein-1].IP == value[0] {
									js.Table[timein-1].Cout++
									keydata = 1
									break
								}
							}
						} else {
							js.Append(js.Size+1, sub_index, nil, &value[0], nil, 1)
							err := pids.Hset(strconv.Itoa(sub_index), strconv.Itoa(js.Size))
							if err != nil {
								val, _ := pids.Hdel(strconv.Itoa(sub_index))
								val += "\n" + strconv.Itoa(js.Size)
								pids.Hset(strconv.Itoa(sub_index), val)
							}
						}
						if keydata == 0 {
							js.Append(js.Size+1, sub_index, nil, &value[0], nil, 1)
							err := pids.Hset(strconv.Itoa(sub_index), strconv.Itoa(js.Size))
							if err != nil {
								val, _ := pids.Hdel(strconv.Itoa(sub_index))
								val += "\n" + strconv.Itoa(js.Size)
								pids.Hset(strconv.Itoa(sub_index), val)
							}
						}
					}
				} else {
					js.Table[id].Cout++
					js.Append(js.Size+1, id, nil, nil, &value[2], 1)
					err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
					if err != nil {
						val, _ := pids.Hdel(strconv.Itoa(id))
						val += "\n" + strconv.Itoa(js.Size)
						pids.Hset(strconv.Itoa(id), val)
					}
					js.Append(js.Size+1, js.Size, nil, &value[0], nil, 1)
					err = pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
					if err != nil {
						val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
						val += "\n" + strconv.Itoa(js.Size)
						pids.Hset(strconv.Itoa(js.Size-1), val)
					}
				}
			} else {
				js.Append(js.Size+1, 0, &value[1], nil, nil, 1)
				links.Hset(value[1], strconv.Itoa(js.Size))
				pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
				js.Append(js.Size+1, js.Size, nil, nil, &value[2], 1)
				pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
				js.Append(js.Size+1, js.Size, nil, &value[0], nil, 1)
			}
		}
		statistics.Qpop()
	}
	jsonData, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при маршалинге в JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}

func (js *JsonResponse) IpLinkTime(statistics *structs.Queue) {
	ip := structs.NewHashTable(10)
	pids := structs.NewHashTable(10)
	for statistics.Head != nil {
		value := strings.Split(statistics.Head.Data, "\n")
		if js.Size == 0 {
			js.Append(js.Size+1, 0, nil, &value[0], nil, 1)
			ip.Hset(value[0], strconv.Itoa(js.Size))
			pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
			js.Append(js.Size+1, js.Size, &value[1], nil, nil, 1)
			pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
			js.Append(js.Size+1, js.Size, nil, nil, &value[2], 1)
		} else {
			if i, _ := ip.Hget(value[0]); i != "" {
				id, _ := strconv.Atoi(i)
				js.Table[id-1].Cout++
				if pi, _ := pids.Hget(i); pi != "" {
					pidarr := strings.Split(pi, "\n")
					key := 0
					sub_index := 0
					for _, val := range pidarr {
						pidu, _ := strconv.Atoi(val)
						if js.Table[pidu-1].Link != nil && *js.Table[pidu-1].Link == value[1] {
							js.Table[pidu-1].Cout++
							key = 1
							sub_index = pidu
							break
						}
					}
					if key == 0 {
						js.Append(js.Size+1, id, &value[1], nil, nil, 1)
						err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
						if err != nil {
							val, _ := pids.Hdel(strconv.Itoa(id))
							val += "\n" + strconv.Itoa(js.Size)
							pids.Hset(strconv.Itoa(id), val)
						}
						js.Append(js.Size+1, js.Size, nil, nil, &value[2], 1)
						err = pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
						if err != nil {
							val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
							val += "\n" + strconv.Itoa(js.Size)
							pids.Hset(strconv.Itoa(js.Size-1), val)
						}
					}
					if sub_index != 0 {
						datatime, _ := pids.Hget(strconv.Itoa(sub_index))
						keydata := 0
						if len(datatime) >= 1 {
							dataarr := strings.Split(datatime, "\n")
							for _, dat := range dataarr {
								timein, _ := strconv.Atoi(dat)
								if js.Table[timein-1].TimeInterval != nil && *js.Table[timein-1].TimeInterval == value[2] {
									js.Table[timein-1].Cout++
									keydata = 1
									break
								}
							}
						} else {
							js.Append(js.Size+1, sub_index, nil, nil, &value[2], 1)
							err := pids.Hset(strconv.Itoa(sub_index), strconv.Itoa(js.Size))
							if err != nil {
								val, _ := pids.Hdel(strconv.Itoa(sub_index))
								val += "\n" + strconv.Itoa(js.Size)
								pids.Hset(strconv.Itoa(sub_index), val)
							}
						}
						if keydata == 0 {
							js.Append(js.Size+1, sub_index, nil, nil, &value[2], 1)
							err := pids.Hset(strconv.Itoa(sub_index), strconv.Itoa(js.Size))
							if err != nil {
								val, _ := pids.Hdel(strconv.Itoa(sub_index))
								val += "\n" + strconv.Itoa(js.Size)
								pids.Hset(strconv.Itoa(sub_index), val)
							}
						}
					}
				} else {
					js.Table[id].Cout++
					js.Append(js.Size+1, id, &value[1], nil, nil, 1)
					err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
					if err != nil {
						val, _ := pids.Hdel(strconv.Itoa(id))
						val += "\n" + strconv.Itoa(js.Size)
						pids.Hset(strconv.Itoa(id), val)
					}
					js.Append(js.Size+1, js.Size, nil, nil, &value[2], 1)
					err = pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
					if err != nil {
						val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
						val += "\n" + strconv.Itoa(js.Size)
						pids.Hset(strconv.Itoa(js.Size-1), val)
					}
				}
			} else {
				js.Append(js.Size+1, 0, nil, &value[0], nil, 1)
				ip.Hset(value[0], strconv.Itoa(js.Size))
				pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
				js.Append(js.Size+1, js.Size, &value[1], nil, nil, 1)
				pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
				js.Append(js.Size+1, js.Size, nil, nil, &value[2], 1)
			}
		}
		statistics.Qpop()
	}
	jsonData, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при маршалинге в JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}

func (js *JsonResponse) IpTimeLink(statistics *structs.Queue) {
	ip := structs.NewHashTable(10)
	pids := structs.NewHashTable(10)
	for statistics.Head != nil {
		value := strings.Split(statistics.Head.Data, "\n")
		if js.Size == 0 {
			js.Append(js.Size+1, 0, nil, &value[0], nil, 1)
			ip.Hset(value[0], strconv.Itoa(js.Size))
			pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
			js.Append(js.Size+1, js.Size, nil, nil, &value[2], 1)
			pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
			js.Append(js.Size+1, js.Size, &value[1], nil, nil, 1)
		} else {
			if i, _ := ip.Hget(value[0]); i != "" {
				id, _ := strconv.Atoi(i)
				js.Table[id-1].Cout++
				if pi, _ := pids.Hget(i); pi != "" {
					pidarr := strings.Split(pi, "\n")
					key := 0
					sub_index := 0
					for _, val := range pidarr {
						pidu, _ := strconv.Atoi(val)
						if js.Table[pidu-1].TimeInterval != nil && *js.Table[pidu-1].TimeInterval == value[2] {
							js.Table[pidu-1].Cout++
							key = 1
							sub_index = pidu
							break
						}
					}
					if key == 0 {
						js.Append(js.Size+1, id, nil, nil, &value[2], 1)
						err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
						if err != nil {
							val, _ := pids.Hdel(strconv.Itoa(id))
							val += "\n" + strconv.Itoa(js.Size)
							pids.Hset(strconv.Itoa(id), val)
						}
						js.Append(js.Size+1, js.Size, &value[1], nil, nil, 1)
						err = pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
						if err != nil {
							val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
							val += "\n" + strconv.Itoa(js.Size)
							pids.Hset(strconv.Itoa(js.Size-1), val)
						}
					}
					if sub_index != 0 {
						datatime, _ := pids.Hget(strconv.Itoa(sub_index))
						keydata := 0
						if len(datatime) >= 1 {
							dataarr := strings.Split(datatime, "\n")
							for _, dat := range dataarr {
								timein, _ := strconv.Atoi(dat)
								if js.Table[timein-1].Link != nil && *js.Table[timein-1].Link == value[1] {
									js.Table[timein-1].Cout++
									keydata = 1
									break
								}
							}
						} else {
							js.Append(js.Size+1, sub_index, &value[1], nil, nil, 1)
							err := pids.Hset(strconv.Itoa(sub_index), strconv.Itoa(js.Size))
							if err != nil {
								val, _ := pids.Hdel(strconv.Itoa(sub_index))
								val += "\n" + strconv.Itoa(js.Size)
								pids.Hset(strconv.Itoa(sub_index), val)
							}
						}
						if keydata == 0 {
							js.Append(js.Size+1, sub_index, &value[1], nil, nil, 1)
							err := pids.Hset(strconv.Itoa(sub_index), strconv.Itoa(js.Size))
							if err != nil {
								val, _ := pids.Hdel(strconv.Itoa(sub_index))
								val += "\n" + strconv.Itoa(js.Size)
								pids.Hset(strconv.Itoa(sub_index), val)
							}
						}
					}
				} else {
					js.Table[id].Cout++
					js.Append(js.Size+1, id, nil, nil, &value[2], 1)
					err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
					if err != nil {
						val, _ := pids.Hdel(strconv.Itoa(id))
						val += "\n" + strconv.Itoa(js.Size)
						pids.Hset(strconv.Itoa(id), val)
					}
					js.Append(js.Size+1, js.Size, &value[1], nil, nil, 1)
					err = pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
					if err != nil {
						val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
						val += "\n" + strconv.Itoa(js.Size)
						pids.Hset(strconv.Itoa(js.Size-1), val)
					}
				}
			} else {
				js.Append(js.Size+1, 0, nil, &value[0], nil, 1)
				ip.Hset(value[0], strconv.Itoa(js.Size))
				pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
				js.Append(js.Size+1, js.Size, nil, nil, &value[2], 1)
				pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
				js.Append(js.Size+1, js.Size, &value[1], nil, nil, 1)
			}
		}
		statistics.Qpop()
	}
	jsonData, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при маршалинге в JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}

func (js *JsonResponse) TimeIpLink(statistics *structs.Queue) {
	time := structs.NewHashTable(10)
	pids := structs.NewHashTable(10)
	for statistics.Head != nil {
		value := strings.Split(statistics.Head.Data, "\n")
		if js.Size == 0 {
			js.Append(js.Size+1, 0, nil, nil, &value[2], 1)
			time.Hset(value[2], strconv.Itoa(js.Size))
			pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
			js.Append(js.Size+1, js.Size, nil, &value[0], nil, 1)
			pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
			js.Append(js.Size+1, js.Size, &value[1], nil, nil, 1)
		} else {
			if i, _ := time.Hget(value[2]); i != "" {
				id, _ := strconv.Atoi(i)
				js.Table[id-1].Cout++
				if pi, _ := pids.Hget(i); pi != "" {
					pidarr := strings.Split(pi, "\n")
					key := 0
					sub_index := 0
					for _, val := range pidarr {
						pidu, _ := strconv.Atoi(val)
						if js.Table[pidu-1].IP != nil && *js.Table[pidu-1].IP == value[0] {
							js.Table[pidu-1].Cout++
							key = 1
							sub_index = pidu
							break
						}
					}
					if key == 0 {
						js.Append(js.Size+1, id, nil, &value[0], nil, 1)
						err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
						if err != nil {
							val, _ := pids.Hdel(strconv.Itoa(id))
							val += "\n" + strconv.Itoa(js.Size)
							pids.Hset(strconv.Itoa(id), val)
						}
						js.Append(js.Size+1, js.Size, &value[1], nil, nil, 1)
						err = pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
						if err != nil {
							val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
							val += "\n" + strconv.Itoa(js.Size)
							pids.Hset(strconv.Itoa(js.Size-1), val)
						}
					}
					if sub_index != 0 {
						datatime, _ := pids.Hget(strconv.Itoa(sub_index))
						keydata := 0
						if len(datatime) >= 1 {
							dataarr := strings.Split(datatime, "\n")
							for _, dat := range dataarr {
								timein, _ := strconv.Atoi(dat)
								if js.Table[timein-1].Link != nil && *js.Table[timein-1].Link == value[1] {
									js.Table[timein-1].Cout++
									keydata = 1
									break
								}
							}
						} else {
							js.Append(js.Size+1, sub_index, &value[1], nil, nil, 1)
							err := pids.Hset(strconv.Itoa(sub_index), strconv.Itoa(js.Size))
							if err != nil {
								val, _ := pids.Hdel(strconv.Itoa(sub_index))
								val += "\n" + strconv.Itoa(js.Size)
								pids.Hset(strconv.Itoa(sub_index), val)
							}
						}
						if keydata == 0 {
							js.Append(js.Size+1, sub_index, &value[1], nil, nil, 1)
							err := pids.Hset(strconv.Itoa(sub_index), strconv.Itoa(js.Size))
							if err != nil {
								val, _ := pids.Hdel(strconv.Itoa(sub_index))
								val += "\n" + strconv.Itoa(js.Size)
								pids.Hset(strconv.Itoa(sub_index), val)
							}
						}
					}
				} else {
					js.Table[id].Cout++
					js.Append(js.Size+1, id, nil, &value[0], nil, 1)
					err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
					if err != nil {
						val, _ := pids.Hdel(strconv.Itoa(id))
						val += "\n" + strconv.Itoa(js.Size)
						pids.Hset(strconv.Itoa(id), val)
					}
					js.Append(js.Size+1, js.Size, &value[1], nil, nil, 1)
					err = pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
					if err != nil {
						val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
						val += "\n" + strconv.Itoa(js.Size)
						pids.Hset(strconv.Itoa(js.Size-1), val)
					}
				}
			} else {
				js.Append(js.Size+1, 0, nil, nil, &value[2], 1)
				time.Hset(value[2], strconv.Itoa(js.Size))
				pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
				js.Append(js.Size+1, js.Size, nil, &value[0], nil, 1)
				pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
				js.Append(js.Size+1, js.Size, &value[1], nil, nil, 1)
			}
		}
		statistics.Qpop()
	}
	jsonData, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при маршалинге в JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}

func (js *JsonResponse) TimeLinkIp(statistics *structs.Queue) {
	time := structs.NewHashTable(10)
	pids := structs.NewHashTable(10)
	for statistics.Head != nil {
		value := strings.Split(statistics.Head.Data, "\n")
		if js.Size == 0 {
			js.Append(js.Size+1, 0, nil, nil, &value[2], 1)
			time.Hset(value[2], strconv.Itoa(js.Size))
			pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
			js.Append(js.Size+1, js.Size, &value[1], nil, nil, 1)
			pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
			js.Append(js.Size+1, js.Size, nil, &value[0], nil, 1)
		} else {
			if i, _ := time.Hget(value[2]); i != "" {
				id, _ := strconv.Atoi(i)
				js.Table[id-1].Cout++
				if pi, _ := pids.Hget(i); pi != "" {
					pidarr := strings.Split(pi, "\n")
					key := 0
					sub_index := 0
					for _, val := range pidarr {
						pidu, _ := strconv.Atoi(val)
						if js.Table[pidu-1].Link != nil && *js.Table[pidu-1].Link == value[1] {
							js.Table[pidu-1].Cout++
							key = 1
							sub_index = pidu
							break
						}
					}
					if key == 0 {
						js.Append(js.Size+1, id, &value[1], nil, nil, 1)
						err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
						if err != nil {
							val, _ := pids.Hdel(strconv.Itoa(id))
							val += "\n" + strconv.Itoa(js.Size)
							pids.Hset(strconv.Itoa(id), val)
						}
						js.Append(js.Size+1, js.Size, nil, &value[0], nil, 1)
						err = pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
						if err != nil {
							val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
							val += "\n" + strconv.Itoa(js.Size)
							pids.Hset(strconv.Itoa(js.Size-1), val)
						}
					}
					if sub_index != 0 {
						datatime, _ := pids.Hget(strconv.Itoa(sub_index))
						keydata := 0
						if len(datatime) >= 1 {
							dataarr := strings.Split(datatime, "\n")
							for _, dat := range dataarr {
								timein, _ := strconv.Atoi(dat)
								if js.Table[timein-1].IP != nil && *js.Table[timein-1].IP == value[0] {
									js.Table[timein-1].Cout++
									keydata = 1
									break
								}
							}
						} else {
							js.Append(js.Size+1, sub_index, nil, &value[0], nil, 1)
							err := pids.Hset(strconv.Itoa(sub_index), strconv.Itoa(js.Size))
							if err != nil {
								val, _ := pids.Hdel(strconv.Itoa(sub_index))
								val += "\n" + strconv.Itoa(js.Size)
								pids.Hset(strconv.Itoa(sub_index), val)
							}
						}
						if keydata == 0 {
							js.Append(js.Size+1, sub_index, nil, &value[0], nil, 1)
							err := pids.Hset(strconv.Itoa(sub_index), strconv.Itoa(js.Size))
							if err != nil {
								val, _ := pids.Hdel(strconv.Itoa(sub_index))
								val += "\n" + strconv.Itoa(js.Size)
								pids.Hset(strconv.Itoa(sub_index), val)
							}
						}
					}
				} else {
					js.Table[id].Cout++
					js.Append(js.Size+1, id, &value[1], nil, nil, 1)
					err := pids.Hset(strconv.Itoa(id), strconv.Itoa(js.Size))
					if err != nil {
						val, _ := pids.Hdel(strconv.Itoa(id))
						val += "\n" + strconv.Itoa(js.Size)
						pids.Hset(strconv.Itoa(id), val)
					}
					js.Append(js.Size+1, js.Size, nil, &value[0], nil, 1)
					err = pids.Hset(strconv.Itoa(js.Size-1), strconv.Itoa(js.Size))
					if err != nil {
						val, _ := pids.Hdel(strconv.Itoa(js.Size - 1))
						val += "\n" + strconv.Itoa(js.Size)
						pids.Hset(strconv.Itoa(js.Size-1), val)
					}
				}
			} else {
				js.Append(js.Size+1, 0, nil, nil, &value[2], 1)
				time.Hset(value[2], strconv.Itoa(js.Size))
				pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
				js.Append(js.Size+1, js.Size, &value[1], nil, nil, 1)
				pids.Hset(strconv.Itoa(js.Size), strconv.Itoa(js.Size+1))
				js.Append(js.Size+1, js.Size, nil, &value[0], nil, 1)
			}
		}
		statistics.Qpop()
	}
	jsonData, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при маршалинге в JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}

func (js *JsonResponse) Priority(link int, ip int, time int, statistics *structs.Queue) {
	switch link {
	case 1:
		switch ip {
		case 2:
			js.LinkIpTime(statistics)
			fmt.Println("link=1, ip=2", "time=3")
		case 3:
			js.LinkTimeIp(statistics)
			fmt.Println("link=1, ip=3", "time=2")
		}
	case 2:
		switch ip {
		case 1:
			js.IpLinkTime(statistics)
			fmt.Println("link=2, ip=1", "time=3")
		case 3:
			js.TimeLinkIp(statistics)
			fmt.Println("link=2, ip=3", "time=1")
		}
	case 3:
		switch ip {
		case 1:
			js.IpTimeLink(statistics)
			fmt.Println("link=3, ip=1", "time=2")
		case 2:
			js.TimeIpLink(statistics)
			fmt.Println("link=3, ip=2", "time=1")
		}
	}
}
