# fireamp_api_consumer
Consumes event data from Cisco FireAMP API to push to Splunk
* FireAMP --> log file --> Splunk forwarder --> Splunk indexer

# Usage
```
go build && ./fireamp_api_consumer -c=[client id] -a=[api key]
```

# Dependencies
Access to Cisco FireAMP API