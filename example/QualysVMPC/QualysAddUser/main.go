package main

import (
	"context"
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ddfelts/godefaultapi"
	"github.com/schollz/progressbar/v3"
)

type AddUserRequest struct {
	Action       string `json:"action,omitempty"`
	AssetGroup   string `json:"asset_group,omitempty"`
	UserRole     string `json:"user_role,omitempty"`
	BusinessUnit string `json:"business_unit,omitempty"`
	Email        string `json:"email,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Title        string `json:"title,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Fax          string `json:"fax,omitempty"`
	Address1     string `json:"address1,omitempty"`
	Address2     string `json:"address2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	ZipCode      string `json:"zip_code,omitempty"`
	Country      string `json:"country,omitempty"`
	ExternalID   string `json:"external_id,omitempty"`
	SendEmail    int    `json:"send_email,omitempty"`
}

type USER struct {
	CONTACT_INFO struct {
		EMAIL string `xml:"EMAIL"`
	} `xml:"CONTACT_INFO"`
}

type USERLIST struct {
	XMLName   xml.Name `xml:"USER_LIST_OUTPUT"`
	Text      string   `xml:",chardata"`
	USER_LIST struct {
		USER []USER `xml:"USER"`
	} `xml:"USER_LIST"`
}

type USEROUTPUT struct {
	XMLName xml.Name `xml:"USER_OUTPUT"`
	Text    string   `xml:",chardata"`
	API     struct {
		Name     string `xml:"name,attr"`
		Username string `xml:"username,attr"`
		At       string `xml:"at,attr"`
	} `xml:"API"`
	RETURN struct {
		Text    string `xml:",chardata"`
		Status  string `xml:"status,attr"`
		MESSAGE string `xml:"MESSAGE"`
	} `xml:"RETURN"`
}

type UserRecord struct {
	Record    []string
	RowData   map[string]string
	FirstName string
	LastName  string
}

func AddUser(client *godefaultapi.Client, request AddUserRequest) (*USEROUTPUT, error) {
	// Build URL parameters
	params := url.Values{}
	params.Add("action", request.Action)
	if request.AssetGroup != "" {
		params.Add("asset_group", request.AssetGroup)
	}
	if request.UserRole != "" {
		params.Add("user_role", request.UserRole)
	}
	if request.BusinessUnit != "" {
		params.Add("business_unit", request.BusinessUnit)
	}
	if request.Email != "" {
		params.Add("email", request.Email)
	}
	if request.FirstName != "" {
		params.Add("first_name", request.FirstName)
	}
	if request.LastName != "" {
		params.Add("last_name", request.LastName)
	}
	if request.Title != "" {
		params.Add("title", request.Title)
	}
	if request.Phone != "" {
		params.Add("phone", request.Phone)
	}
	if request.Fax != "" {
		params.Add("fax", request.Fax)
	}
	if request.Address1 != "" {
		params.Add("address1", request.Address1)
	}
	if request.Address2 != "" {
		params.Add("address2", request.Address2)
	}
	if request.City != "" {
		params.Add("city", request.City)
	}
	if request.State != "" {
		params.Add("state", request.State)
	}
	if request.ZipCode != "" {
		params.Add("zip_code", request.ZipCode)
	}
	if request.Country != "" {
		params.Add("country", request.Country)
	}
	if request.ExternalID != "" {
		params.Add("external_id", request.ExternalID)
	}
	if request.SendEmail != 0 {
		params.Add("send_email", strconv.Itoa(request.SendEmail))
	}

	var response USEROUTPUT
	err := client.Get(context.Background(), "/msp/user.php?"+params.Encode(), nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error adding user: %w", err)
	}
	return &response, nil
}

func getUserList(client *godefaultapi.Client) (*USERLIST, error) {
	var response USERLIST
	err := client.Get(context.Background(), "/msp/user_list.php", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting user list: %w", err)
	}
	return &response, nil
}

func processCSV(client *godefaultapi.Client, inputFile string, doemail string, address1 string, city string, zipcode string, state string) error {
	// Open the CSV file
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error reading CSV header: %w", err)
	}

	// Count total lines for progress bar
	allRecords, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV records: %w", err)
	}
	totalRecords := len(allRecords)

	// Reset file pointer to after header
	file.Seek(0, 0)
	reader = csv.NewReader(file)
	reader.Read() // Skip header again

	// Create progress bar
	bar := progressbar.Default(int64(totalRecords))

	// Create output file
	outputFile := fmt.Sprintf("user_add_results_%s.csv", time.Now().Format("20060102_150405"))
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outFile.Close()

	// Create CSV writer
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Write header with additional column for status
	outputHeader := append(header, "username")
	if err := writer.Write(outputHeader); err != nil {
		return fmt.Errorf("error writing output header: %w", err)
	}

	// Get user list
	userlist, err := getUserList(client)
	if err != nil {
		log.Fatal(err)
	}

	// Process each row
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error reading CSV record: %w", err)
		}

		// Create map of column names to values
		rowData := make(map[string]string)
		for i, col := range header {
			rowData[col] = record[i]
		}

		// Split name and handle potential errors
		nameParts := strings.Split(rowData["user"], " ")
		if len(nameParts) < 2 {
			record = append(record, "Error: User name must contain both first and last name")
			writer.Write(record)
			bar.Add(1)
			continue
		}

		// Check if email already exists
		emailExists := false
		for _, existingUser := range userlist.USER_LIST.USER {
			if existingUser.CONTACT_INFO.EMAIL == rowData["email"] {
				emailExists = true
				record = append(record, "Skipped: Email already exists")
				writer.Write(record)
				bar.Add(1)
				break
			}
		}

		if emailExists {
			bar.Add(1)
			continue
		}

		sg := fmt.Sprintf("SG: %s", rowData["support_group"])
		var sendemail int
		if doemail == "true" {
			sendemail = 1
		} else {
			sendemail = 0
		}

		// Create AddUserRequest
		request := AddUserRequest{
			Action:       "add",
			Email:        rowData["email"],
			FirstName:    nameParts[0],
			LastName:     nameParts[1],
			UserRole:     "reader",
			BusinessUnit: "Unassigned",
			Title:        "TesterJCLTitle",
			Phone:        "1234567890",
			Address1:     address1,
			Address2:     "",
			City:         city,
			State:        state,
			ZipCode:      zipcode,
			Country:      "United States of America",
			ExternalID:   sg,
			SendEmail:    sendemail,
		}

		// Add user
		response, err := AddUser(client, request)
		if err != nil {
			record = append(record, fmt.Sprintf("Error: %v", err))
		} else {
			message := response.RETURN.MESSAGE
			if response.RETURN.Status == "SUCCESS" {
				if idx := strings.Index(message, " "); idx > 0 {
					message = message[:idx]
				}
			}
			record = append(record, message)
		}

		writer.Write(record)
		bar.Add(1)

		// Add delay between requests
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("\nResults written to %s\n", outputFile)
	return nil
}

func main() {
	// Define command-line flags
	url := flag.String("url", "https://qualysapi.qg3.apps.qualys.com", "Qualys API URL")
	username := flag.String("username", "", "Qualys username")
	password := flag.String("password", "", "Qualys password")
	inputFile := flag.String("input", "", "Input CSV file with user data")
	doemail := flag.String("doemail", "false", "Send email to users")
	address1 := flag.String("address1", "JCI", "Address1")
	city := flag.String("city", "Milwaukee", "City")
	zipcode := flag.String("zipcode", "53202", "Zipcode")
	state := flag.String("state", "Wisconsin", "State")

	flag.Parse()

	// Validate required flags
	if *username == "" || *password == "" {
		log.Fatal("Username and password are required")
	}

	if *inputFile == "" {
		log.Fatal("Input CSV file is required")
	}

	client := godefaultapi.NewClient(*url)
	client.SetBasicAuth(*username, *password)
	client.SetRequestType(godefaultapi.ContentTypeJSON)
	client.SetResponseType(godefaultapi.ContentTypeXML)
	client.SetHeader("X-Requested-With", "GOQualysAPI")

	if err := processCSV(client, *inputFile, *doemail, *address1, *city, *zipcode, *state); err != nil {
		log.Fatal(err)
	}

}
