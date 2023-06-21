# To run the service

## Create the docket network

```
docker network create external-example
```

## Get the databases up and running first

```
cd dependencies
docker compose up --build
```

## Get the service up and running

```
cd messaging-service
docker compose up --build
```

## Run integration tests

```
cd messaging-service/integration-tests
go test
```

# Misc.

```
 docker build --tag messaging-server . --no-cache
```

```
 docker run -d -p 5002:9090 messaging-server
```

Dockerfile uses host.docker.internal instead of localhost to connect to the localhost of the host machine.

Build all the images first

```
docker compose up -d --build
```

connect to locally running dockerized mysql

```
mysql -h 127.0.0.1 -P 3308 -u root -p
```

Restart just one container in docker compose

```
docker-compose restart msgserver
```

Start a single service in docker compose

```
docker-compose up -d --no-deps --build <service_name>

```

## TODO's

- Add validation to endpoints
- controltower should be parent of
  - realtime controller (sockets)
  - message controller (sync flows for messages)
- FromUUID should be extracted from header
- Add enum status to roles
- default to "PARTICIPANT" role
- Add a "seen" attribute to the message
  - If a message wasn't seen by the client, push it to the top
- Leave room should be an event sent to other clients
- Allow event and messages to be sent from the room itself
- Dependency injection or singleton needed for redis and mysql db.
- Separate out socket and redis events
- Create an error struct response for API's
- Run processMessage in go routine and inform the client if a message fails. This will let the message be routed directly to the client
