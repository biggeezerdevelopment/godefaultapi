package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	//"time"

	"github.com/ddfelts/godefaultapi"
)

func getUserList(client *godefaultapi.Client) (*USERLIST, error) {
	var response USERLIST
	err := client.Get(context.Background(), "/msp/user_list.php", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting user list: %w", err)
	}
	return &response, nil
}

func fetchApplianceList(client *godefaultapi.Client) (*ApplianceListResponse, error) {
	var response ApplianceListResponse
	// Calculate date 30 days ago
	//sinceDate := time.Now().AddDate(0, 0, -90).Format("2006-01-02")
	params := url.Values{}
	params.Add("action", "list")
	params.Add("output_mode", "full")
	//params.Add("include_cloud_info", "1")

	err := client.Get(context.Background(), "/api/2.0/fo/appliance/?"+params.Encode(), nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting scanner list: %w", err)
	}
	return &response, nil
}

func fetchVMScanList(client *godefaultapi.Client) (*ScanListResponse, error) {
	var response ScanListResponse
	err := client.Get(context.Background(), "/api/2.0/fo/scan/?action=list", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting scan list: %w", err)
	}
	return &response, nil
}

func fetchComScanList(client *godefaultapi.Client) (*ScanListResponse, error) {
	var response ScanListResponse
	err := client.Get(context.Background(), "/api/2.0/fo/scan/compliance/?action=list", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting compliance scan list: %w", err)
	}
	return &response, nil
}

var count int

func fetchAssets(client *godefaultapi.Client, startFromId string, limitResults string) ([]HostAsset, error) {
	count++
	fmt.Printf("Fetching page %d\n", count)

	request := map[string]interface{}{
		"ServiceRequest": map[string]interface{}{
			"preferences": map[string]string{
				"startFromId":  startFromId,
				"limitResults": limitResults,
			},
			"filters": map[string]interface{}{
				"Criteria": map[string]string{
					"field":    "tagName",
					"operator": "EQUALS",
					"value":    "Cloud Agent",
				},
			},
		},
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Request: %+v\n", string(requestJSON))

	var response ServiceResponse
	err = client.Post(context.Background(), "/qps/rest/2.0/search/am/hostasset", requestJSON, &response)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Response: %+v\n", response)

	var assets []HostAsset
	assets = response.Data.HostAsset
	if response.HasMoreRecords == "true" {
		lastId, _ := strconv.Atoi(response.LastId)
		nextId := fmt.Sprintf("%d", lastId+1)
		nextAssets, err := fetchAssets(client, nextId, limitResults)
		if err != nil {
			return nil, err
		}
		assets = append(assets, nextAssets...)
	}

	return assets, nil
}

func fetchAgentInfo(client *godefaultapi.Client, agentId string) (*BinaryInfoResponse, error) {
	// Create a context with a timeout of 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := map[string]interface{}{
		"ServiceRequest": map[string]interface{}{
			"data": map[string]interface{}{
				"BinaryInfo": map[string]interface{}{
					"platform":     "ALL",
					"architecture": "ALL",
				},
			},
		},
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Request: %+v\n", string(requestJSON))
	var response BinaryInfoResponse
	err = client.Post(ctx, "/qps/rest/1.0/process/ca/binaryinfo", requestJSON, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting agent info: %w", err)
	}

	return &response, nil
}
