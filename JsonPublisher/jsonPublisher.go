package main

import (
	"02.08.2024-L0/ServerPart"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	jsonPath = "E:\\Go\\02.08.2024-L0\\JsonPublisher\\jsons"
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
	err, sc := ServerPart.ConnectNats(clientID)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for _, json := range jsonList {
		ServerPart.PublishNATS(sc, "foo", json)
	}
}
