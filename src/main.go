package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	var left_time int
	if now_time < 12 {
		left_time = 12 -now_time
	} else {
		left_time = 24 + 8 - now_time
	}
	fmt.Println(left_time)
	t := time.NewTimer(time.Hour * time.Duration(left_time))
	defer t.Stop()
	for {
		<- t.C
		old_length := len(getEpidemicData())
		err := updateEpidemicData()
		if err != nil {
			panic(err)
		}
		if old_length != len(getEpidemicData()) {
			fmt.Println("updated")
			t.Reset(time.Hour * time.Duration(24 + 8 - time.Now().Hour()))
		} else {
			fmt.Println("update faile")
			t.Reset(time.Hour * 1)
		}
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
	if err := InitEpidemicData(); err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.Dir("../static")))
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
			value["confirm_rate"] = value["confirm"].(float64) - previous["confirm"].(float64)
			value["suspect_rate"] = value["suspect"].(float64) - previous["suspect"].(float64)
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
