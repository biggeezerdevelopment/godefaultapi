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

type AssetGroupListResponse struct {
	//XMLName  xml.Name `xml:"SCAN_LIST_OUTPUT"`
	Response struct {
		AssetGroupList struct {
			AssetGroup []struct {
				ID                   string   `xml:"ID"`
				Title                string   `xml:"TITLE"`
				Appliance_id         []string `xml:"APPLIANCE_IDS"`
				Default_appliance_id string   `xml:"DEFAULT_APPLIANCE_ID"`
				Ip_set               []struct {
					Ip_range   []string `xml:"IP_RANGE"`
					Ip_address []string `xml:"IP"`
				} `xml:"IP_SET"`
			} `xml:"ASSET_GROUP"`
		} `xml:"ASSET_GROUP_LIST"`
	} `xml:"RESPONSE"`
}

func fetchVMAssetGroupList(client *godefaultapi.Client) (*AssetGroupListResponse, error) {
	var response AssetGroupListResponse
	err := client.Get(context.Background(), "/api/2.0/fo/asset/group/?action=list", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting asset group list: %w", err)
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

	assetGroupList, err := fetchVMAssetGroupList(client)
	if err != nil {
		log.Fatal(err)
	}

	// Convert to JSON
	jsonBytes, err := json.MarshalIndent(assetGroupList, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("JSON output:")
	//fmt.Println(string(jsonBytes))

	if len(jsonBytes) > 0 {
		err = os.WriteFile("vmassetgrouplist.json", jsonBytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Asset group list written to vmassetgrouplist.json")
	} else {
		fmt.Println("No asset group data to write")
	}
}
