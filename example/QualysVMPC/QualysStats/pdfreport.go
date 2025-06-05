package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/wcharczuk/go-chart/v2"
)

func writeUserTableToPDF(pdf *gofpdf.Fpdf, userList *USERLIST) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Inactive Users (No Login for > 90 Days)")
	pdf.Ln(10)
	pdf.Cell(190, 10, fmt.Sprintf("Total Users: %d", len(userList.USER_LIST.USER)))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	headers := []string{"UserLogin", "UserId", "Email", "Status", "LastLogin", "Role"}
	colWidths := []float64{40, 20, 50, 20, 30, 30}
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)
	for _, user := range userList.USER_LIST.USER {
		// Skip if last login date is empty or N/A
		if user.UserLastLogin == "" || user.UserLastLogin == "N/A" {
			continue
		}

		// Parse the last login date
		lastLogin, err := time.Parse(time.RFC3339, user.UserLastLogin)
		if err != nil {
			continue
		}

		// Check if user hasn't logged in for 90 days
		if lastLogin.Before(ninetyDaysAgo) {
			// Truncate email to 20 characters
			email := user.CONTACT_INFO.EMAIL
			if len(email) > 20 {
				email = email[:20] + "..."
			}
			// Truncate LastLogin to 11 characters
			lastLoginStr := user.UserLastLogin
			if len(lastLoginStr) > 11 {
				lastLoginStr = lastLoginStr[:11] + "..."
			}
			pdf.CellFormat(colWidths[0], 10, user.UserLogin, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[1], 10, user.UserId, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[2], 10, email, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[3], 10, user.UserStatus, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[4], 10, lastLoginStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[5], 10, user.UserRole, "1", 0, "", false, 0, "")
			pdf.Ln(10)
		}
	}
}

func writeInactiveUserTableToPDF(pdf *gofpdf.Fpdf, userList *USERLIST) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Inactive Users or Hasn't Logged in")
	pdf.Ln(10)
	pdf.Cell(190, 10, fmt.Sprintf("Total Users: %d", len(userList.USER_LIST.USER)))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	headers := []string{"UserLogin", "UserId", "Email", "Status", "LastLogin", "Role"}
	colWidths := []float64{40, 20, 50, 20, 30, 30}
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	//ninetyDaysAgo := time.Now().AddDate(0, 0, -90)
	for _, user := range userList.USER_LIST.USER {
		// Skip if last login date is empty or N/A
		//if user.UserLastLogin == "" || user.UserLastLogin == "N/A" {
		//	continue
		//}

		// Parse the last login date
		//lastLogin, err := time.Parse(time.RFC3339, user.UserLastLogin)
		//if err != nil {
		//	continue
		//}

		// Check if user hasn't logged in for 90 days
		if user.UserStatus == "Inactive" || user.UserLastLogin == "" || user.UserLastLogin == "N/A" {
			// Truncate email to 20 characters
			email := user.CONTACT_INFO.EMAIL
			if len(email) > 20 {
				email = email[:20] + "..."
			}
			// Truncate LastLogin to 11 characters
			lastLoginStr := user.UserLastLogin
			if len(lastLoginStr) > 11 {
				lastLoginStr = lastLoginStr[:11] + "..."
			}
			pdf.CellFormat(colWidths[0], 10, user.UserLogin, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[1], 10, user.UserId, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[2], 10, email, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[3], 10, user.UserStatus, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[4], 10, lastLoginStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[5], 10, user.UserRole, "1", 0, "", false, 0, "")
			pdf.Ln(10)
		}
	}
}

func writeRecentApplianceTableToPDF(pdf *gofpdf.Fpdf, applianceList *ApplianceListResponse) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Appliances Updated in Last 90 Days")
	pdf.Ln(10)
	pdf.Cell(190, 10, fmt.Sprintf("Total Appliances: %d", len(applianceList.Response.ApplianceList.Appliance)))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	headers := []string{"ID", "Name", "SoftVersion", "Status", "LastUpdate", "Type"}
	colWidths := []float64{40, 40, 30, 20, 30, 30}
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)
	for _, appliance := range applianceList.Response.ApplianceList.Appliance {
		// Parse the last update date
		lastUpdate, err := time.Parse("2006-01-02T15:04:05Z", appliance.LastUpdate)
		if err != nil {
			continue
		}

		// Check if appliance was updated in the last 90 days
		if lastUpdate.After(ninetyDaysAgo) {
			// Truncate LastUpdate to 10 characters
			lastUpdateStr := appliance.LastUpdate
			if len(lastUpdateStr) > 10 {
				lastUpdateStr = lastUpdateStr[:10] + "..."
			}
			pdf.CellFormat(colWidths[0], 10, appliance.ID, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[1], 10, appliance.Name, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[2], 10, appliance.SoftVersion, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[3], 10, appliance.Status, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[4], 10, lastUpdateStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[5], 10, appliance.Type, "1", 0, "", false, 0, "")
			//pdf.CellFormat(colWidths[6], 10, appliance.Model, "1", 0, "", false, 0, "")
			//pdf.CellFormat(colWidths[7], 10, appliance.Serial, "1", 0, "", false, 0, "")
			pdf.Ln(10)
		}
	}
}

func writeInactiveApplianceTableToPDF(pdf *gofpdf.Fpdf, applianceList *ApplianceListResponse) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Appliances Not Updated in Last 90 Days")
	pdf.Ln(10)
	pdf.Cell(190, 10, fmt.Sprintf("Total Appliances: %d", len(applianceList.Response.ApplianceList.Appliance)))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	headers := []string{"ID", "Name", "SoftVersion", "Status", "LastUpdate", "Type"}
	colWidths := []float64{40, 40, 30, 20, 30, 30}
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)
	for _, appliance := range applianceList.Response.ApplianceList.Appliance {
		// Parse the last update date
		lastUpdate, err := time.Parse("2006-01-02T15:04:05Z", appliance.LastUpdate)
		if err != nil {
			continue
		}

		// Check if appliance was not updated in the last 90 days
		if lastUpdate.Before(ninetyDaysAgo) {
			// Truncate LastUpdate to 10 characters
			lastUpdateStr := appliance.LastUpdate
			if len(lastUpdateStr) > 10 {
				lastUpdateStr = lastUpdateStr[:10] + "..."
			}
			pdf.CellFormat(colWidths[0], 10, appliance.ID, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[1], 10, appliance.Name, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[2], 10, appliance.SoftVersion, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[3], 10, appliance.Status, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[4], 10, lastUpdateStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[5], 10, appliance.Type, "1", 0, "", false, 0, "")
			//pdf.CellFormat(colWidths[6], 10, appliance.Model, "1", 0, "", false, 0, "")
			//pdf.CellFormat(colWidths[7], 10, appliance.Serial, "1", 0, "", false, 0, "")
			pdf.Ln(10)
		}
	}
}

func writeInactiveAgentAssetsTableToPDF(pdf *gofpdf.Fpdf, assets []HostAsset) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Inactive Agent Assets")
	pdf.Ln(10)
	pdf.Cell(190, 10, fmt.Sprintf("Total Agent Assets: %d", len(assets)))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	headers := []string{"ID", "Name", "LastScan", "AgentStatus", "LastCheckedIn", "AgentVersion"}
	colWidths := []float64{20, 40, 30, 40, 30, 30}
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	for _, asset := range assets {
		if asset.AgentInfo.Status == "STATUS_INACTIVE" {
			// Truncate LastScan to 10 characters
			lastScanStr := asset.LastVulnScan
			if len(lastScanStr) > 10 {
				lastScanStr = lastScanStr[:10] + "..."
			}
			// Truncate LastCheckedIn to 10 characters
			lastCheckedInStr := asset.AgentInfo.LastCheckedIn
			if len(lastCheckedInStr) > 10 {
				lastCheckedInStr = lastCheckedInStr[:10] + "..."
			}
			pdf.CellFormat(colWidths[0], 10, asset.ID, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[1], 10, asset.Name, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[2], 10, lastScanStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[3], 10, asset.AgentInfo.Status, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[4], 10, lastCheckedInStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(colWidths[5], 10, asset.AgentInfo.AgentVersion, "1", 0, "", false, 0, "")
			pdf.Ln(10)
		}
	}
}

func addTableOfContents(pdf *gofpdf.Fpdf) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Table of Contents")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 12)
	sections := []string{
		"Summary",
		"User Status Distribution",
		"Inactive Users (No Login for > 90 Days)",
		"Inactive Users or Hasn't Logged in",
		"Appliance Update Status",
		"Appliances Updated in Last 90 Days",
		"Appliances Not Updated in Last 90 Days",
		"Inactive Agent Assets",
	}

	for _, section := range sections {
		pdf.Cell(190, 10, section)
		pdf.Ln(10)
	}
}

func addSummaryPage(pdf *gofpdf.Fpdf, userCount, assetCount, applianceCount int, comCount int, vmCount int, userList *USERLIST, scannerList *ApplianceListResponse) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Summary")
	pdf.Ln(20)

	// Calculate inactive users and users with no login for 90+ days
	inactiveUsers := 0
	noLogin90Days := 0
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)

	for _, user := range userList.USER_LIST.USER {
		if user.UserStatus == "Inactive" {
			inactiveUsers++
		}

		if user.UserLastLogin == "" || user.UserLastLogin == "N/A" {
			noLogin90Days++
			continue
		}

		lastLogin, err := time.Parse(time.RFC3339, user.UserLastLogin)
		if err != nil {
			continue
		}

		if lastLogin.Before(ninetyDaysAgo) {
			noLogin90Days++
		}
	}

	// Calculate appliance update statistics
	updatedAppliances := 0
	notUpdatedAppliances := 0

	for _, appliance := range scannerList.Response.ApplianceList.Appliance {
		lastUpdate, err := time.Parse("2006-01-02T15:04:05Z", appliance.LastUpdate)
		if err != nil {
			continue
		}

		if lastUpdate.After(ninetyDaysAgo) {
			updatedAppliances++
		} else {
			notUpdatedAppliances++
		}
	}

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(130, 10, "Total Users:")
	pdf.Cell(140, 10, fmt.Sprintf("%d", userCount))
	pdf.Ln(10)

	pdf.Cell(130, 10, "Inactive Users:")
	pdf.Cell(140, 10, fmt.Sprintf("%d", inactiveUsers))
	pdf.Ln(10)

	pdf.Cell(130, 10, "No Login > 90 Days:")
	pdf.Cell(140, 10, fmt.Sprintf("%d", noLogin90Days))
	pdf.Ln(10)

	pdf.Cell(130, 10, "Total Agents:")
	pdf.Cell(140, 10, fmt.Sprintf("%d", assetCount))
	pdf.Ln(10)

	pdf.Cell(130, 10, "Total Appliances:")
	pdf.Cell(140, 10, fmt.Sprintf("%d", applianceCount))
	pdf.Ln(10)

	pdf.Cell(130, 10, "Appliances Updated in Last 90 Days:")
	pdf.Cell(140, 10, fmt.Sprintf("%d", updatedAppliances))
	pdf.Ln(10)

	pdf.Cell(130, 10, "Appliances Not Updated > 90 Days:")
	pdf.Cell(140, 10, fmt.Sprintf("%d", notUpdatedAppliances))
	pdf.Ln(10)

	pdf.Cell(130, 10, "Total Compliance Scans:")
	pdf.Cell(140, 10, fmt.Sprintf("%d", comCount))
	pdf.Ln(10)

	pdf.Cell(130, 10, "Total VM Scans:")
	pdf.Cell(140, 10, fmt.Sprintf("%d", vmCount))
	pdf.Ln(10)
}

func createUserStatsPieChart(userList *USERLIST) (string, error) {
	// Calculate statistics
	totalUsers := len(userList.USER_LIST.USER)
	inactiveUsers := 0
	noLogin90Days := 0
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)

	for _, user := range userList.USER_LIST.USER {
		if user.UserStatus == "Inactive" {
			inactiveUsers++
		}

		if user.UserLastLogin == "" || user.UserLastLogin == "N/A" {
			noLogin90Days++
			continue
		}

		lastLogin, err := time.Parse(time.RFC3339, user.UserLastLogin)
		if err != nil {
			continue
		}

		if lastLogin.Before(ninetyDaysAgo) {
			noLogin90Days++
		}
	}

	activeUsers := totalUsers - inactiveUsers - noLogin90Days

	// Create pie chart
	pie := chart.PieChart{
		Title: "User Status Distribution",
		Values: []chart.Value{
			{
				Label: fmt.Sprintf("Active Users (%d)", activeUsers),
				Value: float64(activeUsers),
			},
			{
				Label: fmt.Sprintf("Inactive Users (%d)", inactiveUsers),
				Value: float64(inactiveUsers),
			},
			{
				Label: fmt.Sprintf("No Login > 90 Days (%d)", noLogin90Days),
				Value: float64(noLogin90Days),
			},
		},
	}

	// Create temporary file for the chart
	tempFile, err := os.CreateTemp("", "user_stats_*.png")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()

	// Create a new file for writing
	f, err := os.Create(tempPath)
	if err != nil {
		return "", fmt.Errorf("error creating chart file: %w", err)
	}

	// Render chart to file
	if err := pie.Render(chart.PNG, f); err != nil {
		f.Close()
		os.Remove(tempPath)
		return "", fmt.Errorf("error rendering chart: %w", err)
	}
	f.Close()

	return tempPath, nil
}

func createApplianceStatsPieChart(scannerList *ApplianceListResponse) (string, error) {
	// Calculate statistics
	updatedAppliances := 0
	notUpdatedAppliances := 0
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)

	for _, appliance := range scannerList.Response.ApplianceList.Appliance {
		lastUpdate, err := time.Parse("2006-01-02T15:04:05Z", appliance.LastUpdate)
		if err != nil {
			continue
		}

		if lastUpdate.After(ninetyDaysAgo) {
			updatedAppliances++
		} else {
			notUpdatedAppliances++
		}
	}

	// Create pie chart
	pie := chart.PieChart{
		Title: "Appliance Update Status",
		Values: []chart.Value{
			{
				Label: fmt.Sprintf("Updated in Last 90 Days (%d)", updatedAppliances),
				Value: float64(updatedAppliances),
			},
			{
				Label: fmt.Sprintf("Not Updated in 90 Days (%d)", notUpdatedAppliances),
				Value: float64(notUpdatedAppliances),
			},
		},
	}

	// Create temporary file for the chart
	tempFile, err := os.CreateTemp("", "appliance_stats_*.png")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()

	// Create a new file for writing
	f, err := os.Create(tempPath)
	if err != nil {
		return "", fmt.Errorf("error creating chart file: %w", err)
	}

	// Render chart to file
	if err := pie.Render(chart.PNG, f); err != nil {
		f.Close()
		os.Remove(tempPath)
		return "", fmt.Errorf("error rendering chart: %w", err)
	}
	f.Close()

	return tempPath, nil
}

func generatePDFReport(outputDir string,
	timestamp string,
	comScanList *ScanListResponse,
	scanList *ScanListResponse,
	scannerList *ApplianceListResponse,
	assets []HostAsset,
	userList *USERLIST) error {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Title
	pdf.Cell(190, 10, "Qualys Status Report")
	pdf.Ln(20)

	// Add timestamp
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 10, fmt.Sprintf("Generated on: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(20)

	// Add table of contents
	addTableOfContents(pdf)

	// Add summary page
	addSummaryPage(pdf,
		len(userList.USER_LIST.USER),
		len(assets),
		len(scannerList.Response.ApplianceList.Appliance),
		len(comScanList.Response.ScanList.Scan),
		len(scanList.Response.ScanList.Scan),
		userList,
		scannerList,
	)

	// Create and add user statistics pie chart
	userChartPath, err := createUserStatsPieChart(userList)
	if err != nil {
		return fmt.Errorf("error creating user stats chart: %w", err)
	}
	defer os.Remove(userChartPath)

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "User Status Distribution")
	pdf.Ln(10)

	// Add the user chart to the PDF
	pdf.Image(userChartPath, 10, 30, 190, 0, false, "", 0, "")
	pdf.Ln(150)

	// Add user table to PDF
	writeUserTableToPDF(pdf, userList)
	writeInactiveUserTableToPDF(pdf, userList)

	// Create and add appliance statistics pie chart
	applianceChartPath, err := createApplianceStatsPieChart(scannerList)
	if err != nil {
		return fmt.Errorf("error creating appliance stats chart: %w", err)
	}
	defer os.Remove(applianceChartPath)

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Appliance Update Status")
	pdf.Ln(10)

	//Add the appliance chart to the PDF
	pdf.Image(applianceChartPath, 10, 30, 190, 0, false, "", 0, "")
	pdf.Ln(150)

	// Add appliance table to PDF
	writeRecentApplianceTableToPDF(pdf, scannerList)
	writeInactiveApplianceTableToPDF(pdf, scannerList)

	// Add agent assets table to PDF
	writeInactiveAgentAssetsTableToPDF(pdf, assets)

	// Save PDF
	pdfPath := filepath.Join(outputDir, "qualys_report_"+timestamp+".pdf")
	return pdf.OutputFileAndClose(pdfPath)
}
