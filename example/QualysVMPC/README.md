# Qualys VM/PC Tools

This directory contains tools for interacting with the Qualys Vulnerability Management and Policy Compliance APIs.

## Programs

### QualysVMScanList

A tool for retrieving a list of VM scans.

#### Usage
```bash
go run QualysVMScanList/main.go -username <username> -password <password>
```

#### Flags
- `-url`: Qualys API URL (default: https://qualysapi.qg3.apps.qualys.com)
- `-username`: Qualys username (required)
- `-password`: Qualys password (required)
- `-daysBack`: Number of days to look back for scans (default: 1)

#### Features
- Retrieves VM scan list in XML format
- Converts XML response to JSON for output
- Writes results to vmscanlist.json
- Supports custom date range for scan retrieval

#### Example
```bash
go run QualysVMScanList/main.go -username user -password pass -daysBack 7
```

### QualysVMScanStats

A tool for retrieving VM scan statistics.

#### Usage
```bash
go run QualysVMScanStats/main.go -username <username> -password <password>
```

#### Flags
- `-url`: Qualys API URL (default: https://qualysapi.qg3.apps.qualys.com)
- `-username`: Qualys username (required)
- `-password`: Qualys password (required)
- `-daysBack`: Number of days to look back for scan statistics (default: 1)

#### Features
- Retrieves VM scan statistics in XML format
- Writes results to scan_stats.json
- Supports custom date range for statistics

#### Example
```bash
go run QualysVMScanStats/main.go -username user -password pass -daysBack 30
```

### QualysVMScanSummary

A tool for retrieving VM scan summaries.

#### Usage
```bash
go run QualysVMScanSummary/main.go -username <username> -password <password>
```

#### Flags
- `-url`: Qualys API URL (default: https://qualysapi.qg3.apps.qualys.com)
- `-username`: Qualys username (required)
- `-password`: Qualys password (required)
- `-daysBack`: Number of days to look back for scan summaries (default: 1)

#### Features
- Retrieves VM scan summaries in XML format
- Writes results to summary.xml
- Supports custom date range for summaries

#### Example
```bash
go run QualysVMScanSummary/main.go -username user -password pass -daysBack 14
```

### QualysAddUser
Automates the process of adding users to Qualys through their API.

Usage:
```bash
go run QualysAddUser/main.go -username user -password pass -input users.csv
```

Flags:
- `-url`: Qualys API URL (default: https://qualysapi.qg3.apps.qualys.com)
- `-username`: Qualys username (required)
- `-password`: Qualys password (required)
- `-input`: Input CSV file path (required)
- `-doemail`: Send welcome email to users (default: false)
- `-address1`: User's address line 1 (default: JCI)
- `-city`: User's city (default: Milwaukee)
- `-zipcode`: User's zip code (default: 53202)
- `-state`: User's state (default: Wisconsin)

Features:
- Reads user information from a CSV file
- Checks for existing users by email
- Creates users with customizable parameters
- Includes progress bar for status tracking
- Creates timestamped output file with results
- Enforces rate limiting with 1-second delay between requests

CSV Format:
```csv
user,email,support_group
John Doe,john.doe@example.com,group1
Jane Smith,jane.smith@example.com,group2
```

## Common Features

All programs in this directory:
- Support basic authentication
- Handle XML responses from the API
- Include proper error handling
- Support custom API URLs
- Use the Qualys VM/PC API v2.0
- Support date-based filtering
- Write results to files for persistence 