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

type Preferences struct {
	StartFromId  string `json:"startFromId"`
	LimitResults string `json:"limitResults"`
}

type ServiceRequest struct {
	Preferences Preferences `json:"preferences"`
}

type Request struct {
	ServiceRequest ServiceRequest `json:"ServiceRequest"`
}
type Response struct {
	ServiceResponse struct {
		ResponseCode   string                   `json:"responseCode"`
		Count          int                      `json:"count"`
		Data           []map[string]interface{} `json:"data"`
		HasMoreRecords string                   `json:"hasMoreRecords"`
		LastId         int                      `json:"lastId"`
	} `json:"ServiceResponse"`
}

var count int

func fetchAssets(client *godefaultapi.Client, startFromId string, limitResults string) ([]map[string]interface{}, error) {
	count++
	fmt.Printf("Fetching page %d\n", count)

	request := Request{
		ServiceRequest: ServiceRequest{
			Preferences: Preferences{
				StartFromId:  startFromId,
				LimitResults: limitResults,
			},
		},
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var response Response
	err = client.Post(context.Background(), "/qps/rest/2.0/search/am/asset", []byte(requestJSON), &response)
	if err != nil {
		return nil, err
	}

	assets := response.ServiceResponse.Data
	if response.ServiceResponse.HasMoreRecords == "true" {
		nextId := fmt.Sprintf("%d", response.ServiceResponse.LastId+1)
		nextAssets, err := fetchAssets(client, nextId, limitResults)
		if err != nil {
			return nil, err
		}
		assets = append(assets, nextAssets...)
	}

	return assets, nil
}

func main() {
	// Define command-line flags
	url := flag.String("url", "https://qualysapi.qg3.apps.qualys.com", "Qualys API URL")
	username := flag.String("username", "", "Qualys username")
	password := flag.String("password", "", "Qualys password")
	limitResults := flag.String("limit", "100", "Number of results per page max 1000")
	flag.Parse()

	// Validate required flags
	if *username == "" || *password == "" {
		log.Fatal("Username and password are required")
	}

	client := godefaultapi.NewClient(*url)
	client.SetBasicAuth(*username, *password)
	client.SetRequestType(godefaultapi.ContentTypeJSON)
	client.SetResponseType(godefaultapi.ContentTypeJSON)
	client.SetHeader("X-Requested-With", "GOQualysAPI")

	assets, err := fetchAssets(client, "0", *limitResults)
	if err != nil {
		log.Fatal(err)
	}

	// Write assets to JSON file
	file, err := json.MarshalIndent(assets, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("assets.json", file, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total assets retrieved: %d\n", len(assets))
	fmt.Println("Assets written to assets.json")
}
