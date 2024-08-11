package main

import (
	"fmt"
	"github.com/MikeShekera/L0/services"
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
	err := getJSONs()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second)
}

func getJSONs() error {

	files, err := os.ReadDir(jsonPath)
	if err != nil {
		log.Fatal(err)
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
		err = services.PublishNATS(sc, "foo", json)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
