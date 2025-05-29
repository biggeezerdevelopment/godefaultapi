package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	//"time"

	"github.com/ddfelts/godefaultapi"
)

/*
	func getDateTime(daysBack int) string {
		now := time.Now().UTC()
		past := now.AddDate(0, 0, -daysBack)
		return past.Format("2006-01-02T15:04:05Z")
	}
*/

type Request struct {
	Action string `json:"action"`
}

type ScanListResponse struct {
	//XMLName  xml.Name `xml:"SCAN_LIST_OUTPUT"`
	Response struct {
		ScanList struct {
			Scan []struct {
				Ref    string `xml:"REF"`
				Title  string `xml:"TITLE"`
				Status struct {
					State    string `xml:"STATE"`
					SubState string `xml:"SUB_STATE"`
				} `xml:"STATUS"`
				LaunchDate         string `xml:"LAUNCH_DATETIME"`
				EndDate            string `xml:"END_DATETIME"`
				Duration           string `xml:"DURATION"`
				ProcessingPriority string `xml:"PROCESSING_PRIORITY"`
				Target             string `xml:"TARGET"`
			} `xml:"SCAN"`
		} `xml:"SCAN_LIST"`
	} `xml:"RESPONSE"`
}

func fetchVMScanList(client *godefaultapi.Client) (*ScanListResponse, error) {
	var response ScanListResponse
	err := client.Get(context.Background(), "/api/2.0/fo/scan/?action=list", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting scan list: %w", err)
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
	client.SetRequestType(godefaultapi.ContentTypeJSON)
	client.SetResponseType(godefaultapi.ContentTypeXML)
	client.SetHeader("X-Requested-With", "GOQualysAPI")

	scanList, err := fetchVMScanList(client)
	if err != nil {
		log.Fatal(err)
	}

	// Convert to JSON
	jsonBytes, err := json.MarshalIndent(scanList, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("JSON output:")
	//fmt.Println(string(jsonBytes))

	if len(jsonBytes) > 0 {
		err = os.WriteFile("vmscanlist.json", jsonBytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Scan list written to vmscanlist.json")
	} else {
		fmt.Println("No scan data to write")
	}
}
