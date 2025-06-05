package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/xuri/excelize/v2"
)

func writeUserListToCSV(userList *USERLIST, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"UserLogin",
		"UserId",
		"Email",
		"Status",
		"LastLogin",
		"Role",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write data
	for _, user := range userList.USER_LIST.USER {
		record := []string{
			user.UserLogin,
			user.UserId,
			user.CONTACT_INFO.EMAIL,
			user.UserStatus,
			user.UserLastLogin,
			user.UserRole,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %w", err)
		}
	}

	fmt.Printf("User list written to %s\n", filename)
	return nil
}

func writeInactiveUsersToCSV(userList *USERLIST, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"UserLogin",
		"UserID",
		"NoLogin90",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Calculate the date 90 days ago
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)

	// Write data for users matching criteria
	for _, user := range userList.USER_LIST.USER {
		// Skip if last login date is empty or N/A
		if user.UserLastLogin == "" || user.UserLastLogin == "N/A" {
			// Write users with N/A login date to CSV
			record := []string{
				user.UserLogin,
				user.UserId,
				"N/A",
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing CSV record: %w", err)
			}
			continue
		}

		// Parse the last login date
		lastLogin, err := time.Parse(time.RFC3339, user.UserLastLogin)
		if err != nil {
			fmt.Printf("Failed to parse date for user %s: %s\n", user.UserLogin, err)
			continue
		}

		// Check if the date is in the future (indicating no login)
		if lastLogin.After(time.Now()) {
			record := []string{
				user.UserLogin,
				user.UserId,
				user.UserLastLogin,
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing CSV record: %w", err)
			}
			continue
		}

		// Check if user hasn't logged in for 90 days
		if lastLogin.Before(ninetyDaysAgo) {
			record := []string{
				user.UserLogin,
				user.UserId,
				user.UserLastLogin,
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing CSV record: %w", err)
			}
		}
	}

	fmt.Printf("Inactive users list written to %s\n", filename)
	return nil
}

func writeVMScanListToCSV(scanList *ScanListResponse, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Ref",
		"Title",
		"Status",
		"Duration",
		"LaunchDate",
		"EndDate",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write data for scans with Error status
	for _, scan := range scanList.Response.ScanList.Scan {
		if scan.Status.State == "Error" {
			record := []string{
				scan.Ref,
				scan.Title,
				scan.Status.State,
				scan.Duration,
				scan.LaunchDate,
				scan.EndDate,
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing CSV record: %w", err)
			}
		}
	}

	fmt.Printf("VM scan list written to %s\n", filename)
	return nil
}

func writeApplianceListToCSV(applianceList *ApplianceListResponse, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()
	//fmt.Printf("Scanner list: %+v\n", scannerList)
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID",
		"Name",
		"SoftVersion",
		"Status",
		"MLVersion",
		"VulnSigVersion",
		"LastUpdate",
		"Type",
		"Model",
		"Serial",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write data
	for _, appliance := range applianceList.Response.ApplianceList.Appliance {
		record := []string{
			appliance.ID,
			appliance.Name,
			appliance.SoftVersion,
			appliance.Status,
			appliance.MLVersion,
			appliance.VulnSigVersion,
			appliance.LastUpdate,
			appliance.Type,
			appliance.Model,
			appliance.Serial,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %w", err)
		}
	}

	fmt.Printf("Appliance list written to %s\n", filename)
	return nil
}

func writeAssetsToCSV(assets []HostAsset, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID",
		"Name",
		"AgentVersion",
		"ManifestVersion",
		"AgentId",
		"OS",
		"LastScan",
		"AgentStatus",
		"LastCheckedIn",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write data
	for _, asset := range assets {
		record := []string{
			asset.ID,
			asset.Name,
			asset.AgentInfo.AgentVersion,
			asset.AgentInfo.ManifestVersion.Vm,
			asset.AgentInfo.AgentId,
			asset.Os,
			asset.LastVulnScan,
			asset.AgentInfo.Status,
			asset.AgentInfo.LastCheckedIn,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %w", err)
		}
	}

	fmt.Printf("Asset list written to %s\n", filename)
	return nil
}

func writeAgentInfoToCSV(agentInfo *BinaryInfoResponse, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Platform",
		"Version",
		"Extension",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write data
	for _, platform := range agentInfo.Data.AllBinaryInfo.Platforms.Platform {
		// Clean version string using regex to remove all non-version characters
		re := regexp.MustCompile(`[^0-9.]`)
		version := re.ReplaceAllString(platform.Version, "")

		record := []string{
			platform.Name,
			version,
			platform.Extension,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %w", err)
		}
	}

	fmt.Printf("Agent info written to %s\n", filename)
	return nil
}

func writeComplianceScanListToCSV(scanList *ScanListResponse, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Ref",
		"Title",
		"Status",
		"Duration",
		"LaunchDate",
		"EndDate",
		"Target",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write data for scans with Error status
	for _, scan := range scanList.Response.ScanList.Scan {
		if scan.Status.State == "Error" {
			record := []string{
				scan.Ref,
				scan.Title,
				scan.Status.State,
				scan.Duration,
				scan.LaunchDate,
				scan.EndDate,
				scan.Target,
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing CSV record: %w", err)
			}
		}
	}

	fmt.Printf("Compliance scan list written to %s\n", filename)
	return nil
}

func writeAllVMScanListToCSV(scanList *ScanListResponse, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Ref",
		"Title",
		"Status",
		"Duration",
		"LaunchDate",
		"EndDate",
		"Target",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write all scan data
	for _, scan := range scanList.Response.ScanList.Scan {
		record := []string{
			scan.Ref,
			scan.Title,
			scan.Status.State,
			scan.Duration,
			scan.LaunchDate,
			scan.EndDate,
			scan.Target,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %w", err)
		}
	}

	fmt.Printf("All VM scan list written to %s\n", filename)
	return nil
}

func writeAllComplianceScanListToCSV(scanList *ScanListResponse, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Ref",
		"Title",
		"Status",
		"Duration",
		"LaunchDate",
		"EndDate",
		"Target",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write all scan data
	for _, scan := range scanList.Response.ScanList.Scan {
		record := []string{
			scan.Ref,
			scan.Title,
			scan.Status.State,
			scan.Duration,
			scan.LaunchDate,
			scan.EndDate,
			scan.Target,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %w", err)
		}
	}

	fmt.Printf("All compliance scan list written to %s\n", filename)
	return nil
}

func writeInactiveAppliancesToCSV(applianceList *ApplianceListResponse, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID",
		"Name",
		"SoftVersion",
		"Status",
		"LastUpdate",
		"Type",
		"Model",
		"Serial",
		"DaysSinceLastUpdate",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Calculate the date 7 days ago
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	// Write data for appliances that haven't checked in for 7 days
	for _, appliance := range applianceList.Response.ApplianceList.Appliance {
		// Parse the last update date
		lastUpdate, err := time.Parse("2006-01-02T15:04:05Z", appliance.LastUpdate)
		if err != nil {
			fmt.Printf("Failed to parse date for appliance %s: %s\n", appliance.Name, err)
			continue
		}

		// Check if appliance hasn't updated in 7 days
		if lastUpdate.Before(sevenDaysAgo) {
			daysSinceUpdate := int(time.Since(lastUpdate).Hours() / 24)
			record := []string{
				appliance.ID,
				appliance.Name,
				appliance.SoftVersion,
				appliance.Status,
				appliance.LastUpdate,
				appliance.Type,
				appliance.Model,
				appliance.Serial,
				fmt.Sprintf("%d", daysSinceUpdate),
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing CSV record: %w", err)
			}
		}
	}

	fmt.Printf("Inactive appliances list written to %s\n", filename)
	return nil
}

func combineCSVToExcel(outputDir string, timestamp string) error {
	// Create a new Excel file
	f := excelize.NewFile()
	defer f.Close()

	// Define the CSV files and their corresponding sheet names
	files := []struct {
		filename  string
		sheetName string
	}{
		{"scan_list_error_" + timestamp + ".csv", "VM Scan Errors"},
		{"scan_list_all_" + timestamp + ".csv", "All VM Scans"},
		{"compliance_scan_list_error_" + timestamp + ".csv", "Compliance Scan Errors"},
		{"compliance_scan_list_all_" + timestamp + ".csv", "All Compliance Scans"},
		{"appliance_list_" + timestamp + ".csv", "Appliance List"},
		{"inactive_appliances_" + timestamp + ".csv", "Inactive Appliances"},
		{"user_list_" + timestamp + ".csv", "User List"},
		{"90days_inactive_users_" + timestamp + ".csv", "Inactive Users"},
		{"agent_assets.csv", "Agent Assets"},
		{"agent_binary_info.csv", "Agent Binary Info"},
	}

	// Process each CSV file
	for _, file := range files {
		// Open the CSV file
		csvFile, err := os.Open(filepath.Join(outputDir, file.filename))
		if err != nil {
			fmt.Printf("Warning: Could not open file %s: %v\n", file.filename, err)
			continue
		}
		defer csvFile.Close()

		// Create a new sheet
		sheetName := file.sheetName
		index, err := f.NewSheet(sheetName)
		if err != nil {
			return fmt.Errorf("error creating sheet %s: %w", sheetName, err)
		}
		f.SetActiveSheet(index)

		// Read CSV data
		reader := csv.NewReader(csvFile)
		records, err := reader.ReadAll()
		if err != nil {
			return fmt.Errorf("error reading CSV %s: %w", file.filename, err)
		}

		// Write data to Excel sheet
		for i, record := range records {
			for j, cell := range record {
				col := string(rune('A' + j))
				cellName := fmt.Sprintf("%s%d", col, i+1)
				f.SetCellValue(sheetName, cellName, cell)
			}
		}

		// Auto-fit columns
		cols, err := f.GetCols(sheetName)
		if err != nil {
			return fmt.Errorf("error getting columns for sheet %s: %w", sheetName, err)
		}
		for idx, col := range cols {
			largestWidth := 0
			for _, cellValue := range col {
				cellWidth := len(fmt.Sprintf("%v", cellValue))
				if cellWidth > largestWidth {
					largestWidth = cellWidth
				}
			}
			colName := string(rune('A' + idx))
			f.SetColWidth(sheetName, colName, colName, float64(largestWidth+2))
		}
	}

	// Delete the default Sheet1
	f.DeleteSheet("Sheet1")

	// Save the Excel file
	excelPath := filepath.Join(outputDir, "qualys_report_"+timestamp+".xlsx")
	if err := f.SaveAs(excelPath); err != nil {
		return fmt.Errorf("error saving Excel file: %w", err)
	}

	fmt.Printf("Combined CSV files written to %s\n", excelPath)
	return nil
}
