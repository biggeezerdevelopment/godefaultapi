# Qualys Add User Tool

This tool automates the process of adding users to Qualys through their API. It reads user information from a CSV file and creates users in Qualys, with options to customize the user creation process.

## CSV File Format

The input CSV file must contain the following columns:
- `user`: Full name of the user (must contain both first and last name)
- `email`: Email address of the user
- `support_group`: Support group identifier for the user

Example CSV content:
```csv
user,email,support_group
John Doe,john.doe@example.com,group1
Jane Smith,jane.smith@example.com,group2
```

## Command Line Arguments

| Argument | Description | Default Value |
|----------|-------------|---------------|
| `-url` | Qualys API URL | https://qualysapi.qg3.apps.qualys.com |
| `-username` | Qualys username | (required) |
| `-password` | Qualys password | (required) |
| `-input` | Input CSV file path | (required) |
| `-doemail` | Send welcome email to users | false |
| `-address1` | User's address line 1 | JCI |
| `-city` | User's city | Milwaukee |
| `-zipcode` | User's zip code | 53202 |
| `-state` | User's state | Wisconsin |

## Example Usage

Basic usage with required fields:
```
QualysAddUser.exe -username your_username -password your_password -input users.csv
```
***SPECIAL NOTE***
If you want the users to recieve an email for registration, set the doemail option to true

Full usage with all options:
```
QualysAddUser.exe \
  -username your_username \
  -password your_password \
  -input users.csv \
  -doemail true \
  -address1 "123 Main St" \
  -city "Chicago" \
  -zipcode "60601" \
  -state "Illinois"
```

## Output

The program creates an output CSV file named `user_add_results_YYYYMMDD_HHMMSS.csv` containing:
- All original columns from the input CSV
- An additional `username` column with the result of the user creation:
  - The created username on success
  - "Error: [error message]" on failure
  - "Skipped: Email already exists" if the email is already registered

## Notes

- The program checks for existing users by email before attempting to create new ones
- A 1-second delay is enforced between API requests to prevent rate limiting
- All users are created with the "reader" role by default
- The program includes a progress bar to show the status of user creation
- Failed user creations are logged in the output file but don't stop the program 