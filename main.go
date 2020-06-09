package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Cleaning folders before start...")
	//Clean chunks folder to always start from scratch
	cleanFolders()
	fmt.Println("Folders cleaned, parsing config...")
	//Create GCS client object
	client := gcsClient{}
	//Get merchantPortal from environment variables
	conf := getConfig()
	fmt.Println("Config parsed, downloading & parsing chunks...")
	//Create client with merchantPortal
	client.createClient(conf)
	//Close client after finished
	defer client.close()
	fmt.Println("Lines will be written only for valid chunks")
	//Download & process all chunks sequentially
	client.downloadAndProcessChunks()
}

func cleanFolders() {
	chunks, err := filepath.Glob("./chunks/*")

	if err != nil {
		log.Fatal(err)
	}

	for _, chunk := range chunks {
		if e := os.Remove(chunk); e != nil {
			log.Fatal(e)
		}
	}

	var logs []string
	logs, err = filepath.Glob("./logs/*")

	if err != nil {
		log.Fatal(err)
	}

	for _, logFile := range logs {
		if e := os.Remove(logFile); e != nil {
			log.Fatal(e)
		}
	}
}
