package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ddfelts/godefaultapi"
)

type Criteria struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Filters struct {
	Criteria []Criteria `json:"Criteria"`
}

type Preferences struct {
	StartFromId  string `json:"startFromId"`
	LimitResults string `json:"limitResults"`
}

type ServiceRequest struct {
	Filters     Filters     `json:"filters"`
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

func fetchAssets(client *godefaultapi.Client, criteria []Criteria, startFromId string, limitResults string, host bool) ([]map[string]interface{}, error) {
	count++
	fmt.Printf("Fetching page %d\n", count)

	request := Request{
		ServiceRequest: ServiceRequest{
			Filters: Filters{
				Criteria: criteria,
			},
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
	//fmt.Println(string(requestJSON))
	var response Response
	if host {
		err = client.Post(context.Background(), "/qps/rest/2.0/search/am/hostasset", []byte(requestJSON), &response)
	} else {
		err = client.Post(context.Background(), "/qps/rest/2.0/search/am/asset", []byte(requestJSON), &response)
	}
	if err != nil {
		return nil, err
	}
	//fmt.Println(response)
	assets := response.ServiceResponse.Data
	if response.ServiceResponse.HasMoreRecords == "true" {
		nextId := fmt.Sprintf("%d", response.ServiceResponse.LastId+1)
		nextAssets, err := fetchAssets(client, criteria, nextId, limitResults, host)
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
	criteriaStr := flag.String("criteria", "", "Search criteria in format 'field1:operator1:value1,field2:operator2:value2'")
	host := flag.Bool("host", false, "Search for hosts")
	flag.Parse()

	// Validate required flags
	if *username == "" || *password == "" {
		log.Fatal("Username and password are required")
	}

	if *criteriaStr == "" {
		log.Fatal("At least one search criterion is required")
	}

	// Parse multiple criteria
	var criteria []Criteria
	validFields := map[string]bool{
		"customAttributeKey": true, "tagId": true, "created": true, "name": true,
		"id": true, "type": true, "tagName": true, "activationKey": true,
		"agentUuid": true, "updated": true, "customAttributeValue": true,
	}

	for _, criterion := range strings.Split(*criteriaStr, ",") {
		parts := strings.Split(criterion, ":")
		if len(parts) != 3 {
			log.Fatal("Invalid criterion format. Use 'field:operator:value'")
		}
		field := strings.TrimSpace(parts[0])

		if !*host && !validFields[field] {
			log.Fatal("Invalid field. When host flag is not set, valid fields are: customAttributeKey, tagId, created, name, id, type, tagName, activationKey, agentUuid, updated, customAttributeValue")
		}

		criteria = append(criteria, Criteria{
			Field:    field,
			Operator: strings.TrimSpace(parts[1]),
			Value:    strings.TrimSpace(parts[2]),
		})
	}

	client := godefaultapi.NewClient(*url)
	client.SetBasicAuth(*username, *password)
	client.SetRequestType(godefaultapi.ContentTypeJSON)
	client.SetResponseType(godefaultapi.ContentTypeJSON)
	client.SetHeader("X-Requested-With", "GOQualysAPI")

	assets, err := fetchAssets(client, criteria, "0", *limitResults, *host)
	if err != nil {
		log.Fatal(err)
	}

	// Write assets to JSON file
	file, err := json.MarshalIndent(assets, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("search_results.json", file, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total assets found: %d\n", len(assets))
	fmt.Println("Search results written to search_results.json")
}
