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

### Authenticate

First, you have to [get an Auth token](https://auth0.com/docs/flows/call-your-api-using-resource-owner-password-flow#request-tokens) by calling:

         payload="{\"client_id\":\"9l4Fp3tK9mtRGSdS8WdSVHiE8QFToMOD\",\"client_secret\":\"$CLIENT_SECRET\",\"audience\":\"http://localhost:8000\",\"grant_type\":\"password\",\"username\":\"$API_USERNAME\",\"password\":\"$API_PASSWORD\"}"
         TOKEN=$(curl --request POST --url https://dev-elymxg0u.us.auth0.com/oauth/token --header 'content-type: application/json' --data "$payload" | jq -r .access_token)

The token that will be returned, has to be used in the private API calls below.

To get the credentials, you have to configure an application in your [Auth0 account](https://manage.auth0.com/).

### Add a campaign

      # curl -H "Content-Type: application/json" -d '{"name":"Covid","organizerName":"Martin","donationMinimum":2,"targetAmount":100,"account":{"name":"Martin","bankName":"DKB","number":"123456"}}' localhost:8000/campaign --header "authorization: Bearer $TOKEN"

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

      # curl -X DELETE localhost:8000/campaigns/1  --header "authorization: Bearer $TOKEN"

### Update a campaign

The following command updates the campaign with ID 1:

      # curl -X PUT -H "Content-Type: application/json"  -d '{"name":"Corona","organizerName":"Marcus","donationMinimum":2,"targetAmount":100,"Account":{"Name":"Marcus","BankName":"DKB","Number":"123456"}}' localhost:8000/campaigns/1  --header "authorization: Bearer $TOKEN"

### Add a donation to a campaign

This command adds a donation to the campaign with ID 1:

      # curl -H "Content-Type: application/json" -d '{"Amount":20,"donorName":"Martin","receiptRequested":true,"status":"IN_PROCESS","account":{"name":"Martin","bankName":"DKB","number":"123456"}}' localhost:8000/campaigns/1/donation

## Non-Functional Tests

### Load test adding Donations

To test the retry mechanism, we have to send a lot of transactions. In one terminal run:

      # docker compose down && docker compose up --build > compose.log

In another just show the error logs of docker compose:

      # tail -f compose.log | grep "level=error"

And in a third terminal, first create a new campaign:

      # curl -H "Content-Type: application/json" -d '{"name":"Covid","organizerName":"Martin","donationMinimum":2,"targetAmount":100,"account":{"name":"Martin","bankName":"DKB","number":"123456"}}' localhost:8000/campaign

And then continuously send new donations to it:

      # while sleep 0.1; do curl -H "Content-Type: application/json" -d '{"Amount":20,"donorName":"Martin","receiptRequested":true,"status":"IN_PROCESS","account":{"name":"Martin","bankName":"DKB","number":"123456"}}' localhost:8000/campaigns/1/donation; done

Look for the strings "Requeuing transaction" and "Requeued transaction"

### Scaling banktransfer

Now, we're scaling the banktransfer service:

     # docker compose down && docker compose up --build --scale banktransfer=2

After adding a campaign:

     # curl -H "Content-Type: application/json" -d '{"name":"Covid","organizerName":"Martin","donationMinimum":2,"targetAmount":100,"account":{"name":"Martin","bankName":"DKB","number":"123456"}}' localhost:8000/campaign

We do a couple of donations by calling:

      # curl -H "Content-Type: application/json" -d '{"Amount":20,"donorName":"Martin","receiptRequested":true,"status":"IN_PROCESS","account":{"name":"Martin","bankName":"DKB","number":"123456"}}' localhost:8000/campaigns/1/donation

If we're using the Kafka version, all donations are processed as the messages are shared between banktransfer instances.

If we're using the version based on Go channels, the result is that only about half of all donations are processed.

This happens if the transaction monitor (`myaktion/monitor.go`) is connected to another instance of the banktransfer service
as the `AddDonation` method (`maktion/service/donation.go`).
Then the goroutine in the `processTransaction` method (`banktransfer/service/banktransfer.go`) is waiting endlessly for
the goroutine in `ProcessTransactions` to receive the transaction. That won't happen as the goroutine is running
in another service instance.

## Run on K8S

1.  Locally build the docker images

         # docker compose build

2.  Install application on K8S

         # kubectl apply -f ./kubernetes-manifests

3.  Forward port 8000 of myaktion to local port 8000

         # kubectl -n myaktion port-forward service/myaktion 8000

4.  Watch logs of banktransfer service (dedicated terminal)

         # kubectl logs -f -n myaktion -l run=banktransfer

5.  Watch logs of myaktion service (dedicated terminal)

         # kubectl logs -f -n myaktion -l run=myaktion

Then you can use curl to test the application as described above. Once you're finished, you can delete the application by calling:

         # kubectl delete ns myaktion
