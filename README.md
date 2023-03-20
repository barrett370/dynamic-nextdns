# Dynamic-NextDNS


A simple cron service to poll a given nextdns dynamic dns linked IP endpoint.


Configurable via the following env vars:

- `NEXTDNS_TARGET`: target url to poll
- `NEXTDNS_INTERVAL_SECONDS`: (optional) interval in seconds at which to poll `NEXTDNS_TARGET`, defaults to 60
- `NEXTDNS_PROFILE`: (optional) string denoting nextdns profile name