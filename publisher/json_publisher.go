package main

import (
	"02.08.2024-L0/services"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	jsonPath = "publisher/jsons"
	clientID = "jsonPublisher"
)

func main() {
	getJSONs()
	time.Sleep(2 * time.Second)
}

func getJSONs() error {

	f, err := os.Open(jsonPath)
	if err != nil {
		return err
	}

	files, err := f.Readdir(0)
	if err != nil {
		return err
	}

	var jsonList [][]byte
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			os.Chdir(jsonPath)
			json, err := os.ReadFile(file.Name())
			fmt.Printf("%s: %s \n", file.Name(), json)
			if err == nil {
				jsonList = append(jsonList, json)
			}
		}
	}

	sendJSONs(jsonList)
	return nil
}

func sendJSONs(jsonList [][]byte) {
	err, sc := services.ConnectNats(clientID)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for _, json := range jsonList {
		services.PublishNATS(sc, "foo", json)
		fmt.Println(json)
	}
}
