package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ddfelts/godefaultapi"
)

type Request struct {
	Action        string `json:"action"`
	Ref           string `json:"scan_reference"`
	ScanDateSince string `json:"scan_date_since"`
}

func getDateTime(daysBack int) string {
	now := time.Now().UTC()
	past := now.AddDate(0, 0, -daysBack)
	return past.Format("2006-01-02T15:04:05Z")
}

func fetchStats(client *godefaultapi.Client, reference string) (interface{}, error) {
	var response interface{}
	request := Request{Action: "list", Ref: reference, ScanDateSince: getDateTime(1)}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}
	err = client.Post(context.Background(), "/api/2.0/fo/scan/vm/summary/", jsonBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting stats: %w", err)
	}
	return response, nil
}

func main() {
	// Define command-line flags
	url := flag.String("url", "https://qualysapi.qg3.apps.qualys.com", "Qualys API URL")
	username := flag.String("username", "", "Qualys username")
	password := flag.String("password", "", "Qualys password")
	reference := flag.String("reference", "", "Scan reference")
	flag.Parse()

	// Validate required flags
	if *username == "" || *password == "" {
		log.Fatal("Username and password are required")
	}

	client := godefaultapi.NewClient(*url)
	client.SetBasicAuth(*username, *password)
	client.SetRequestType(godefaultapi.ContentTypeJSON)
	client.SetResponseType(godefaultapi.ContentTypeXML)
	client.SetHeader("X-Requested-With", "GOQualysAPI")

	stats, err := fetchStats(client, *reference)
	if err != nil {
		log.Fatal(err)
	}

	// Print raw response for debugging
	fmt.Printf("Raw response: %+v\n", stats)

	statsBytes, err := xml.MarshalIndent(stats, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Print XML for debugging
	fmt.Println("XML output:")
	fmt.Println(string(statsBytes))
	if len(statsBytes) > 0 {
		err = os.WriteFile("stats.xml", statsBytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("No stats to write")
	}
}
