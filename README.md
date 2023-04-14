## Start

Start the server:

      # cd src/myaktion && go run main.go

## Use the REST API

### Add a campaign

      # curl -H "Content-Type: application/json" -d '{"name":"Covid","organizerName":"Martin","donationMinimum":2,"targetAmount":100,"account":{"name":"Martin","bankName":"DKB","number":"123456"}}' localhost:8000/campaign

### Get campaign data

To retrieve all the persisted campaign objects, send this GET request:

      # curl localhost:8000/campaigns

For retrieving just a specific campaign, append the campaign's ID. E.g. for retrieving the
campaign with ID 1, call:

      # curl localhost:8000/campaigns/1

### Delete a campaign

The following command deletes the campaign with ID 1:

      # curl -X DELETE localhost:8000/campaigns/1

### Update a campaign

The following command updates the campaign with ID 1:

      # curl -X PUT -H "Content-Type: application/json"  -d '{"name":"Corona","organizerName":"Marcus","donationMinimum":2,"targetAmount":100}' localhost:8000/campaigns/1

### Add a donation to a campaign

This command adds a donation to the campaign with ID 1:

      # curl -H "Content-Type: application/json" -d '{"Amount":20,"donorName":"Martin","receiptRequested":true,"status":"IN_PROCESS","account":{"name":"Martin","bankName":"DKB","number":"123456"}}' localhost:8000/campaigns/1/donation
