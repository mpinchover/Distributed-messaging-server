FROM golang:1.20.4
WORKDIR /app
ENV REDIS_URL="host.docker.internal:6379" 
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /messaging-server
EXPOSE 9090
CMD ["/messaging-server"]