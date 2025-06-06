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

type SCriteria struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Criteria struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type SFilters struct {
	Criteria SCriteria `json:"Criteria"`
}

type Filters struct {
	Criteria []Criteria `json:"Criteria"`
}

type Preferences struct {
	StartFromId  string `json:"startFromId"`
	LimitResults string `json:"limitResults"`
}

type ServiceRequest struct {
	Filters SFilters `json:"filters"`
	//Preferences Preferences `json:"preferences"`
}

type Request struct {
	ServiceRequest ServiceRequest `json:"ServiceRequest"`
}
type Response struct {
	ServiceResponse struct {
		ResponseCode    string                   `json:"responseCode"`
		ResponseMessage string                   `json:"responseMessage"`
		ResponseError   string                   `json:"responseError"`
		Count           int                      `json:"count"`
		Data            []map[string]interface{} `json:"data"`
		HasMoreRecords  string                   `json:"hasMoreRecords"`
		LastId          int                      `json:"lastId"`
	} `json:"ServiceResponse"`
}

var count int

func fetchAssets(client *godefaultapi.Client, host bool, tagName string, startFromId string, limitResults string) ([]map[string]interface{}, error) {
	count++
	fmt.Printf("Fetching page %d\n", count)

	request := Request{
		ServiceRequest: ServiceRequest{
			Filters: SFilters{
				Criteria: SCriteria{
					Field:    "tagName",
					Operator: "EQUALS",
					Value:    tagName,
				},
			},
		},
	}

	requestJSON, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	var response Response
	if host {
		err = client.Post(context.Background(), "/qps/rest/2.0/search/am/hostasset", []byte(requestJSON), &response)
	} else {
		err = client.Post(context.Background(), "/qps/rest/2.0/search/am/asset", []byte(requestJSON), &response)
	}
	if err != nil {
		return nil, err
	}

	assets := response.ServiceResponse.Data
	if response.ServiceResponse.HasMoreRecords == "true" {
		nextId := fmt.Sprintf("%d", response.ServiceResponse.LastId+1)
		nextAssets, err := fetchAssets(client, host, tagName, nextId, limitResults)
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
	limitResults := flag.String("limit", "1000", "Number of results per page max 1000")
	tagName := flag.String("tag", "", "Tag name to search for")
	host := flag.Bool("host", false, "Search for hosts")
	flag.Parse()

	// Validate required flags
	if *username == "" || *password == "" {
		log.Fatal("Username and password are required")
	}

	if *tagName == "" {
		log.Fatal("Tag name is required")
	}

	client := godefaultapi.NewClient(*url)
	client.SetBasicAuth(*username, *password)
	client.SetContentType(godefaultapi.ContentTypeJSON)
	client.SetHeader("X-Requested-With", "GOQualysAPI")

	assets, err := fetchAssets(client, *host, *tagName, "0", *limitResults)
	if err != nil {
		log.Fatal(err)
	}

	// Write assets to JSON file
	file, err := json.MarshalIndent(assets, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if *host {
		err = os.WriteFile("hostsbytags.json", file, 0644)
	} else {
		err = os.WriteFile("assetsbytags.json", file, 0644)
	}
	if err != nil {
		log.Fatal(err)
	}
	if *host {
		fmt.Printf("Total tagged hosts retrieved: %d\n", len(assets))
	} else {
		fmt.Printf("Total tagged assets retrieved: %d\n", len(assets))
	}

	fmt.Println("Tagged Assets written to bytags.json")
}
