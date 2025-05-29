# Qualys Asset Management Tools

This directory contains tools for interacting with the Qualys Asset Management API.

## Programs

### QualysAMSearch

A tool for searching assets in Qualys Asset Management.

#### Usage
```bash
go run QualysAMSearch/main.go -username <username> -password <password> -criteria "field:operator:value"
```

#### Flags
- `-url`: Qualys API URL (default: https://qualysapi.qg3.apps.qualys.com)
- `-username`: Qualys username (required)
- `-password`: Qualys password (required)
- `-limit`: Number of results per page (default: 100, max: 1000)
- `-criteria`: Search criteria in format 'field1:operator1:value1,field2:operator2:value2'
- `-host`: Search for hosts (default: false)

#### Field Restrictions
When `-host` flag is not set, valid fields are limited to:
- customAttributeKey
- tagId
- created
- name
- id
- type
- tagName
- activationKey
- agentUuid
- updated
- customAttributeValue

#### Example
```bash
# Search for assets with specific tag
go run QualysAMSearch/main.go -username user -password pass -criteria "tagName:CONTAINS:production"

# Search for hosts with specific name
go run QualysAMSearch/main.go -username user -password pass -host -criteria "name:CONTAINS:server"
```

### QualysByTag

A tool for retrieving assets by tag.

#### Usage
```bash
go run QualysByTag/main.go -username <username> -password <password> -tagName <tag_name>
```

#### Flags
- `-url`: Qualys API URL (default: https://qualysapi.qg3.apps.qualys.com)
- `-username`: Qualys username (required)
- `-password`: Qualys password (required)
- `-tagName`: Tag name to search for (required)
- `-limit`: Number of results per page (default: 100, max: 1000)

#### Example
```bash
go run QualysByTag/main.go -username user -password pass -tagName "production-servers"
```

## Common Features

All programs in this directory:
- Support basic authentication
- Handle pagination automatically
- Write results to JSON files
- Include proper error handling
- Support custom API URLs
- Use the Qualys Asset Management API v2.0 