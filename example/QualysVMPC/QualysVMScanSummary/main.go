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
	ScanDateSince string `json:"scan_datetime_since"`
	Action        string `json:"action"`
}

func getDateTime(daysBack int) string {
	now := time.Now().UTC()
	past := now.AddDate(0, 0, -daysBack)
	return past.Format("2006-01-02T15:04:05Z")
}

func fetchSummary(client *godefaultapi.Client, daysBack int) (interface{}, error) {
	dateTime := getDateTime(daysBack)
	var response interface{}
	request := Request{ScanDateSince: dateTime, Action: "list"}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}
	fmt.Println(string(jsonBytes))
	err = client.Post(context.Background(), "/api/2.0/fo/scan/vm/summary/", []byte(jsonBytes), &response)
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
	daysBack := flag.Int("daysBack", 1, "Days back to get summary for")
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

	summary, err := fetchSummary(client, *daysBack)
	if err != nil {
		log.Fatal(err)
	}

	// Print raw response for debugging
	fmt.Printf("Raw response: %+v\n", summary)

	summaryBytes, err := xml.MarshalIndent(summary, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Print XML for debugging
	fmt.Println("XML output:")
	fmt.Println(string(summaryBytes))
	if len(summaryBytes) > 0 {
		err = os.WriteFile("summary.xml", summaryBytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("No summary to write")
	}
}
