package main

import (
	//"context"
	//"encoding/csv"

	//"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	//"net/url"
	//"os"
	//"strconv"
	//"strings"
	//"time"

	"github.com/ddfelts/godefaultapi"
	//"github.com/schollz/progressbar/v3"
)

func main() {
	// Define command-line flags
	url := flag.String("url", "https://qualysapi.qg3.apps.qualys.com", "Qualys API URL")
	username := flag.String("username", "", "Qualys username")
	password := flag.String("password", "", "Qualys password")
	outputDir := flag.String("output", "output", "Output directory for results")
	flag.Parse()

	// Validate required flags
	if *username == "" || *password == "" {
		log.Fatal("Username and password are required")
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	client := godefaultapi.NewClient(*url)
	client.SetBasicAuth(*username, *password)
	client.SetRequestType(godefaultapi.ContentTypeJSON)
	client.SetResponseType(godefaultapi.ContentTypeXML)
	client.SetHeader("X-Requested-With", "GOQualysAPI")

	// Generate timestamp for filename
	timestamp := time.Now().Format("20060102_150405")

	// Fetch and write VM scan list
	scanList, err := fetchVMScanList(client)
	if err != nil {
		log.Fatalf("Error fetching VM scan list: %v", err)
	}

	// Write VM scan list with errors to CSV
	vmScanFilename := fmt.Sprintf("scan_list_error_%s.csv", timestamp)
	vmScanPath := filepath.Join(*outputDir, vmScanFilename)
	if err := writeVMScanListToCSV(scanList, vmScanPath); err != nil {
		log.Fatalf("Error writing VM scan list to CSV: %v", err)
	}

	// Write all VM scan list to CSV
	allVMScanFilename := fmt.Sprintf("scan_list_all_%s.csv", timestamp)
	allVMScanPath := filepath.Join(*outputDir, allVMScanFilename)
	if err := writeAllVMScanListToCSV(scanList, allVMScanPath); err != nil {
		log.Fatalf("Error writing all VM scan list to CSV: %v", err)
	}

	// Fetch and write compliance scan list
	comScanList, err := fetchComScanList(client)
	if err != nil {
		log.Fatalf("Error fetching compliance scan list: %v", err)
	}

	// Write compliance scan list with errors to CSV
	comScanFilename := fmt.Sprintf("compliance_scan_list_error_%s.csv", timestamp)
	comScanPath := filepath.Join(*outputDir, comScanFilename)
	if err := writeComplianceScanListToCSV(comScanList, comScanPath); err != nil {
		log.Fatalf("Error writing compliance scan list to CSV: %v", err)
	}

	// Write all compliance scan list to CSV
	allComScanFilename := fmt.Sprintf("compliance_scan_list_all_%s.csv", timestamp)
	allComScanPath := filepath.Join(*outputDir, allComScanFilename)
	if err := writeAllComplianceScanListToCSV(comScanList, allComScanPath); err != nil {
		log.Fatalf("Error writing all compliance scan list to CSV: %v", err)
	}

	// Fetch and write scanner list
	scannerList, err := fetchApplianceList(client)
	if err != nil {
		log.Fatalf("Error fetching scanner list: %v", err)
	}

	// Write scanner list to CSV
	scannerFilename := fmt.Sprintf("appliance_list_%s.csv", timestamp)
	scannerPath := filepath.Join(*outputDir, scannerFilename)
	if err := writeApplianceListToCSV(scannerList, scannerPath); err != nil {
		log.Fatalf("Error writing scanner list to CSV: %v", err)
	}

	// Write inactive appliances to CSV
	inactiveAppliancesFilename := fmt.Sprintf("inactive_appliances_%s.csv", timestamp)
	inactiveAppliancesPath := filepath.Join(*outputDir, inactiveAppliancesFilename)
	if err := writeInactiveAppliancesToCSV(scannerList, inactiveAppliancesPath); err != nil {
		log.Fatalf("Error writing inactive appliances to CSV: %v", err)
	}

	// Fetch and write user list
	userList, err := getUserList(client)
	if err != nil {
		log.Fatal(err)
	}

	// Write full user list to CSV
	fullListFilename := fmt.Sprintf("user_list_%s.csv", timestamp)
	fullListPath := filepath.Join(*outputDir, fullListFilename)
	if err := writeUserListToCSV(userList, fullListPath); err != nil {
		log.Fatalf("Error writing user list to CSV: %v", err)
	}

	// Write inactive users to CSV
	inactiveUsersFilename := fmt.Sprintf("90days_inactive_users_%s.csv", timestamp)
	inactiveUsersPath := filepath.Join(*outputDir, inactiveUsersFilename)
	if err := writeInactiveUsersToCSV(userList, inactiveUsersPath); err != nil {
		log.Fatalf("Error writing inactive users to CSV: %v", err)
	}

	// Fetch assets
	assets, err := fetchAssets(client, "1", "100")
	if err != nil {
		fmt.Printf("Error fetching assets: %v\n", err)
		return
	}

	// Write assets to CSV
	assetsFilename := filepath.Join(*outputDir, "agent_assets.csv")
	if err := writeAssetsToCSV(assets, assetsFilename); err != nil {
		fmt.Printf("Error writing assets to CSV: %v\n", err)
		return
	}

	// Fetch agent info
	agentInfo, err := fetchAgentInfo(client, "")
	if err != nil {
		fmt.Printf("Error fetching agent info: %v\n", err)
		return
	}

	// Write agent info to CSV
	agentInfoFilename := filepath.Join(*outputDir, "agent_binary_info.csv")
	if err := writeAgentInfoToCSV(agentInfo, agentInfoFilename); err != nil {
		fmt.Printf("Error writing agent info to CSV: %v\n", err)
		return
	}

	// Combine all CSV files into a single Excel spreadsheet
	if err := combineCSVToExcel(*outputDir, timestamp); err != nil {
		fmt.Printf("Error combining CSV files to Excel: %v\n", err)
		return
	}

	// Generate PDF report
	if err := generatePDFReport(*outputDir,
		timestamp,
		comScanList,
		scanList,
		scannerList,
		assets,
		userList,
	); err != nil {
		fmt.Printf("Error generating PDF report: %v\n", err)
		return
	}
}
