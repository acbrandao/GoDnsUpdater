# GoDnsUpdater
Go cli application to Update Dynamic DNS Records  for NameCheap API

To use this script:

Save the Go code as ddns-updater.go

Enable ** Dynamic DNS **in your Namecheap account:
-  Go to Domain List > Manage > Advanced DNS
-   Enable Dynamic DNS and note the password
-   Create A+Dynamic DNS records for your hosts

## Build then run  Go Program 
```go build -o ddns-updater
./ddns-updater
```

modify this sample **config.json** to match your domains 
Create a config.json file with your domain details:
domain: Your domain name
host: The hostname (e.g., "@" for root domain, "www" for subdomain)
password: The Dynamic DNS password from Namecheap's Advanced DNS settings
check_interval_minutes: How often to check IP (in minutes)
log_file: Where to store the log
```
{
    "domains": [
        {
            "domain": "example.com",
            "host": "@",
            "password": "your-ddns-password-1"
        },
        {
            "domain": "example.com",
            "host": "www",
            "password": "your-ddns-password-1"
        },
        {
            "domain": "otherdomain.com",
            "host": "@",
            "password": "your-ddns-password-2"
        }
    ],
    "check_interval_minutes": 5,
    "log_file": "ddns.log"
}
```


