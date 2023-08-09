# To run the service

## Create the docker network

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

## Read files

```

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
  t, err := template.ParseFiles(wd + "/assets/templates/template.html")

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

### P1

- make a makefile with a testfile that checks for the test db and restarts it
- use mysql long/lat to detect distance?
- get all the tracked questions this user has liked. Search for all tracked questions with likes on those q's.
  then get those users. Then just do a mapping and see if the largest numbers have over the threshold.
- need mappers for room
- test out matching + recv events
- change the delete to deleteByRoomId
- change uuid to the test name
- custom error mapping, last step should just translate the error befor sending it back customer err code -> http code
- clean up naming
- don't use middleware.New for http route
- don't return all messages in getroomsbyuserUUID
- ensure secret for JWT is not hardcoded

### P2

- ensure passwords match in signup
- ensure password, email are correct formats
- update tests to run in go routines to mimic high, concurrent volumes
- separate chat notifications from messages, or just fix messages to relay messagse from server
- separate out channel and client events
- allow someone to be invited to the chat
- Add a "typing" attribute/event to the message
- Run processMessage in go routine and inform the client if a message fails. This will let the message be routed directly to the client
- run UpdateMessageToSeen in go routine
- on the server side, mark the messages as seen so the client doesn't have to

### P3

- implement delete, update authprofile
- ResetPassword, GeneratePasswordLink test later
- controltower should be parent of
  - realtime controller (sockets)
  - message controller (sync flows for messages)
- LeaveRoom should also save messages that someone has left the chat
- app should show if someone screenshotted
- add in a user permissions table to link to member table

### ranking

- show how many people also liked that card
  - get this before you send it back
- for each person that you swipe on, event back to server to check and see if this is a match. If it is, send an event letting the app know there's been a match.
- run through every question and show the matches based off the top 3 or so from each category.
- store the question in an excel sheet and and in s3 and version it so you can roll it back. Store the UUID in there too and save to DB.

# Testing

Run a single test
`go test -v -run TestRoomAndMessagesPaginationByApiKey`

Run all tests
`go test ./... -count=1`

# Mocks

To generate mocks for all interfaces, go to `messaging-service` and run
`mockery --all --keeptree`
