package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
)

type Endpoint []string

type Conf struct {
	Endpoints []Endpoint
	API       string
}

type QueryResult struct {
	Count     float64
	Duration  time.Duration
	QueryType string
}

var conf Conf

func readConf() {
	if _, err := toml.DecodeFile("conf.toml", &conf); err != nil {
		log.Fatal(err)
	}
}

func main() {
	readConf()
	if len(conf.Endpoints) == 0 {
		log.Fatal("missing conf")
	}
	for {
		// re-read on each loop
		readConf()
		for _, ep := range conf.Endpoints {
			fmt.Println("requesting", ep[1])
			makeRequest(ep)
		}
		time.Sleep(5 * time.Second)
	}

}

func makeRequest(ep Endpoint) {
	qres := QueryResult{QueryType: ep[0]}

	start := time.Now()

	res, err := http.Get(ep[1])

	if err != nil {
		fmt.Println("err", err)
		return
	}

	defer res.Body.Close()

	qres.Duration = time.Since(start)

	data := map[string]interface{}{}

	json.NewDecoder(res.Body).Decode(&data)

	fmt.Println(data)

	count, ok := data["count"].(float64)
	if !ok {
		fmt.Printf("unexpected count type: %T\n", data["count"])
	}
	qres.Count = count
	saveResult(qres)
}

func saveResult(qres QueryResult) {
	results := map[string]interface{}{
		"query_type": qres.QueryType,
		"count":      qres.Count,
		"resp_time":  qres.Duration / 1000000, //ms
		"created":    time.Now().UnixNano() / 1000000,
	}

	data, err := json.Marshal(results)

	res, err := http.Post(conf.API, "application/json", bytes.NewBuffer(data))

	if err != nil {
		fmt.Println("err", err)
		return
	}

	defer res.Body.Close()
}
