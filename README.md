# To run the service

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

TODO's
 – Implement delete room and leave room
 – add validation to all endpoints
 – controltower should be parent of
– realtime controller (sockets)
– message controller (getting saving updating messages etc)
– FromUUID should be extracted from headers
