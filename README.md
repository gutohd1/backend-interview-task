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