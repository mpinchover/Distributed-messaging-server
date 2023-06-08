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
