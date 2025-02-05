# USAGE

## Requirements:
- Docker
- Postman/curl(oprional for tests)

## Run

```shell
make run
```

It will up docker containers with postgres and server.

You can run
```shell
make swagger
```
and go to http://localhost:8090
for preview allowed methods and paths by swagger documentation.

## Request examples:

Deposit Money:
```shell
curl -X 'PATCH' 
  'https://localhost:8080/3fec06e9-29cc-4ff4-9ae7-fb0e7c757b61/deposit' 
  -H 'accept: application/json' 
  -H 'Content-Type: application/json' 
  -d '{
  "amount": 100000
}'
```

Transfer Money:
```shell
curl -X 'PATCH' 
  'https://localhost:8080/3fec06e9-29cc-4ff4-9ae7-fb0e7c757b61/transfer' 
  -H 'accept: application/json' 
  -H 'Content-Type: application/json' 
  -d '{
  "receiverID": "4178f61f-2ff9-4ab5-afa5-f30dc16e6ad9",
  "amount": 1
}'
```

Get 10 last user transactions:
```shell
curl -X 'GET' 
  'https://localhost:8080/3fec06e9-29cc-4ff4-9ae7-fb0e7c757b61/transactions' 
  -H 'accept: application/json'
```
