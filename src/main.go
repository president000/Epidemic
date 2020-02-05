package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type EpidemicData struct {
	data []byte
	lock sync.Mutex
}
var epidemic_data EpidemicData
func getEpidemicData() []byte {
	epidemic_data.lock.Lock()
	ret := epidemic_data.data
	epidemic_data.lock.Unlock()
	return ret
}

func setEpidemicData(data []byte) {
	epidemic_data.lock.Lock()
	epidemic_data.data = data
	epidemic_data.lock.Unlock()
}

func InitEpidemicData() error {
	err := updateEpidemicData()
	if err != nil {
		return err
	}
	go updateEpidemicDataTimer()
	return err
}

func updateEpidemicDataTimer() {
	now_time := time.Now().Hour()
	left_time := 12
	if now_time >= 12 {
		left_time += 24
	}
	t := time.NewTimer(time.Hour * time.Duration(left_time))
	defer t.Stop()
	for {
		<- t.C
		err := updateEpidemicData()
		if err != nil {
			panic(err)
		}
		t.Reset(time.Hour * 24)
	}
}

func updateEpidemicData() error {
	day_add_list, err := getDayAddList()
	if err != nil {
		return err
	}
	err = calcInflectionPoint(day_add_list)
	if err != nil {
		return err
	}
	data, err := json.Marshal(day_add_list)
	if err != nil {
		return err
	}
	setEpidemicData(data)
	return err
}

func main() {
	if err := updateEpidemicData(); err != nil {
		panic(err)
	}
	fmt.Println(string(getEpidemicData()))
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/api", ApiRequest)
	http.ListenAndServe(":8080", nil)
}

func ApiRequest(w http.ResponseWriter, r *http.Request) {
	w.Write(getEpidemicData())
}

func calcInflectionPoint(data []interface{}) error {
	for i := 0; i < len(data); i++ {
		value := data[i].(map[string]interface{})
		if i == 0 {
			value["confirm_rate"] = 0
			value["suspect_rate"] = 0
		} else {
			previous := data[i - 1].(map[string]interface{})
			temp_1, err := strconv.Atoi(value["confirm"].(string))
			if err != nil {
				return err
			}
			temp_2, err := strconv.Atoi(previous["confirm"].(string))
			if err != nil {
				return err
			}
			temp_3, err := strconv.Atoi(value["suspect"].(string))
			if err != nil {
				return err
			}
			temp_4, err := strconv.Atoi(previous["suspect"].(string))
			if err != nil {
				return err
			}
			value["confirm_rate"] = temp_1 - temp_2
			value["suspect_rate"] = temp_3 - temp_4
		}
	}
	return nil
}

func getDayAddList() ([]interface {}, error) {
	resp, err := http.Get("https://view.inews.qq.com/g2/getOnsInfo?name=disease_h5")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(dat["data"].(string)), &dat); err != nil {
		return nil, err
	}
	add_list := dat["chinaDayAddList"].([]interface {})
	return add_list, nil
}
