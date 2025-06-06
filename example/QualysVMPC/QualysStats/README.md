# Qualys Status Report Generator

This program generates comprehensive status reports for your Qualys environment, including user activity, appliance status, and scan information. The reports are generated in both PDF and Excel formats.

## Features

- **User Management**
  - Tracks user login activity
  - Identifies inactive users
  - Monitors users with no login activity for > 90 days
  - Visual representation of user status distribution

- **Appliance Management**
  - Monitors appliance update status
  - Tracks appliances updated in last 90 days
  - Identifies appliances not updated in > 90 days
  - Visual representation of appliance update status

- **Asset Management**
  - Tracks agent assets
  - Identifies inactive agent assets
  - Monitors agent versions and status

- **Scan Management**
  - Tracks VM scans
  - Monitors compliance scans
  - Identifies scan errors

## Output Formats

### PDF Report
The program generates a detailed PDF report containing:
- Summary statistics
- User status distribution pie chart
- Appliance update status pie chart
- Detailed tables for users, appliances, and assets
- Table of contents for easy navigation

### Excel Report
The program also generates an Excel spreadsheet with multiple sheets:
- VM Scan Errors
- All VM Scans
- Compliance Scan Errors
- All Compliance Scans
- Appliance List
- Inactive Appliances
- User List
- Inactive Users
- Agent Assets
- Agent Binary Info


## Usage

Run the program with the following command:
```bash
qualysstats -username <username> -password <password>
```

### Command Line Options

- `-output`: Specify the output directory for reports (default: current directory)
- `-username`: Username (required)
- `-password`: Password (required)
- `-url`: Qualys api url (default set to GuidePoints api url)


## Output Files

The program generates the following files:
- `qualys_report_[timestamp].pdf`: Comprehensive PDF report
- `qualys_report_[timestamp].xlsx`: Excel spreadsheet with detailed data
- `*_[timestamp].csv`: Csv files
- Individual CSV files for each data category

## Report Contents

### PDF Report Sections
1. Summary
   - Total Users
   - Inactive Users
   - Users with No Login > 90 Days
   - Total Agents
   - Total Appliances
   - Appliance Update Status
   - Total Scans

2. User Status Distribution
   - Pie chart showing active/inactive users
   - Detailed user tables

3. Appliance Update Status
   - Pie chart showing update status
   - Detailed appliance tables

4. Asset Information
   - Inactive agent assets
   - Agent status and versions

### Excel Report Sheets
- VM Scan Errors
- All VM Scans
- Compliance Scan Errors
- All Compliance Scans
- Appliance List
- Inactive Appliances
- User List
- Inactive Users
- Agent Assets
- Agent Binary Info

## Error Handling

The program includes error handling for:
- API connection issues
- Invalid credentials
- File generation errors
- Data processing errors
