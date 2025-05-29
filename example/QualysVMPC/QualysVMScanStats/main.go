package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ddfelts/godefaultapi"
)

type Response struct {
	TASK_PROCESSING struct {
		UNPROCESSED_SCANS       int `xml:"UNPROCESSED_SCANS"`
		VM_RECRYPT_BACKLOG      int `xml:"VM_RECRYPT_BACKLOG"`
		VM_DECRYPT_BACKLOG_SCAN struct {
			SCAN []struct {
				ID                  string `xml:"ID"`
				TITLE               string `xml:"TITLE"`
				STATUS              string `xml:"STATUS"`
				PROCESSING_PRIORITY string `xml:"PROCESSING_PRIORITY"`
				COUNT               int    `xml:"COUNT"`
			} `xml:"SCAN"`
		} `xml:"VM_DECRYPT_BACKLOG_SCAN"`
		VM_DECRYPT_BACKLOG_TASK struct {
			SCAN []struct {
				ID                  string `xml:"ID"`
				TITLE               string `xml:"TITLE"`
				STATUS              string `xml:"STATUS"`
				PROCESSING_PRIORITY string `xml:"PROCESSING_PRIORITY"`
				NBHOST              string `xml:"NBHOST"`
				PROCESSED           int    `xml:"PROCESSED"`
				TO_PROCESS          int    `xml:"TO_PROCESS"`
				SCAN_DATE           string `xml:"SCAN_DATE"`
				TASK_TYPE           string `xml:"TASK_TYPE"`
				TASK_STATUS         string `xml:"TASK_STATUS"`
				TASK_UPDATED_DATE   string `xml:"TASK_UPDATED_DATE"`
			} `xml:"SCAN"`
		} `xml:"VM_DECRYPT_BACKLOG_TASK"`
	} `xml:"TASK_PROCESSING"`
}

func fetchStats(client *godefaultapi.Client) (*Response, error) {
	var response Response
	err := client.Get(context.Background(), "/api/2.0/fo/scan/stats/?action=list", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting stats: %w", err)
	}
	return &response, nil
}

func main() {
	// Define command-line flags
	url := flag.String("url", "https://qualysapi.qg3.apps.qualys.com", "Qualys API URL")
	username := flag.String("username", "", "Qualys username")
	password := flag.String("password", "", "Qualys password")
	flag.Parse()

	// Validate required flags
	if *username == "" || *password == "" {
		log.Fatal("Username and password are required")
	}

	client := godefaultapi.NewClient(*url)
	client.SetBasicAuth(*username, *password)
	client.SetContentType(godefaultapi.ContentTypeXML)
	client.SetHeader("X-Requested-With", "GOQualysAPI")

	stats, err := fetchStats(client)
	if err != nil {
		log.Fatal(err)
	}

	// Print raw response for debugging
	jsonBytes, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("JSON output:")
	//fmt.Println(string(jsonBytes))

	if len(jsonBytes) > 0 {
		err = os.WriteFile("vmscanstats.json", jsonBytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Scan stats written to vmscanstats.json")
	} else {
		fmt.Println("No scan stats to write")
	}
}
