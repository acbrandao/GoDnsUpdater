# GoDnsUpdater
Go cli application to Update Dynamic DNS Records  for NameCheap API

modify this sample config.json to match your dowmain
`

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
`
