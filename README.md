#backend Go test

## Running via Docker
go to app and run 
```bash 
docker compose up -d
```
this command will download the images and run the server listening on the port 9000 and the database.

## Running on command line:
rename ```.env.example``` file to ```.env```

then run: 
```bash 
go run main.go
```

## Tests:
go to app/handlers and run 
```bash
go test handlers_test.go -v
```

## Client test
in order to check the functionality of the application I have implemented one Client
go to ```app/client``
these are the options to run this client
```bash
  -actor string
    	-actor=1 | id to call specific actor user (default "1")
  -function string
    	-function=ListLikedYou, ListNewLikedYou, CountLikedYou and PutDecision
  -like
    	-like=false | Can only be used on PutDecision (default true)
  -page string
    	-page=1 | page number to paginate likes (default "1")
  -recipient string
    	-recipient=1 | id to call specific recipient user (default "1")
```

Example:
```bash
go run client.go -function ListLikedYou -recipient=1 -page=1
```

## Database
### Migrations
go to ```app/database/migrations```
and run 
```bash
go run main.go -migrate=up
```
or down
```bash
go run main.go -migrate=down
```
### Seeds
Though on create a seeder to create random data but in order to keep it simple decided against and just create one SQL file available in
```app/database```
