# Myaktion-Go Microservices

## Prepare development environment

1.  Install `protoc` compiler

        # brew install protoc

2.  Install go plugins for protoc

    Make sure to call the following command in this directory (to make sure the protoc generators are installed
    globally and not added to a `go.mod` file of an individual microservice):

        # go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
        # go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

After that, follow the individual development instructions for each microservice.

## Build

Build and run the docker containers:

      # docker compose up --build

Shut down the docker containers:

      # docker compose down

## Use the REST API

### Add a campaign

      # curl -H "Content-Type: application/json" -d '{"name":"Covid","organizerName":"Martin","donationMinimum":2,"targetAmount":100,"account":{"name":"Martin","bankName":"DKB","number":"123456"}}' localhost:8000/campaign

To check if the campaign was persisted, call:

      # docker exec -it myaktion-go_mariadb_1 mysql -uroot -proot -e 'USE myaktion; SELECT * FROM campaigns'

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

      # curl -X PUT -H "Content-Type: application/json"  -d '{"name":"Corona","organizerName":"Marcus","donationMinimum":2,"targetAmount":100,"Account":{"Name":"Marcus","BankName":"DKB","Number":"123456"}}' localhost:8000/campaigns/1

### Add a donation to a campaign

This command adds a donation to the campaign with ID 1:

      # curl -H "Content-Type: application/json" -d '{"Amount":20,"donorName":"Martin","receiptRequested":true,"status":"IN_PROCESS","account":{"name":"Martin","bankName":"DKB","number":"123456"}}' localhost:8000/campaigns/1/donation
