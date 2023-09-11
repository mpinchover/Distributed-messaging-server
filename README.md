# What is this service?

This is a distributed messaging server built in Golang (server), Redis (caching and channels) and MySQL (database).
Two more devices can enter the same or disparate servers for real time chat communication through websockets.

# Architecture

- user sends message to server over websocket
- server
  - saves message to database
  - broadcasts the event to the chat room channel
- all servers subscribed to the chat room channel receive the message
- if the other users in the chat room are connected to the receiving server, the message is broadcast to their web sockets

# To run the service

```
make run-api
```

# Enter the container

```
make enter
```

# Testing

### Run integration tests

```
make run-api-integration
```

And in another window run

```
make int-test
```

### Run unit tests

```
make test
```

### Mocks

To generate mocks for all interfaces, go to `messaging-service` and run
`mockery --all --keeptree`

# Assets

### Automated emails

assets should be queried from root directory as the starting point.

```

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
  t, err := template.ParseFiles(wd + "/assets/templates/template.html")

```

# Database

Setup mysql and redis

```
make setup-dbs
```

connect to locally running dockerized mysql

```
mysql -h 127.0.0.1 -P 3310 -u root -p
```
