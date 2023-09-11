#!/bin/bash
mysql -uroot -proot --protocol=tcp -h localhost --port=3310 -e "
    drop database if exists chat_test_db;
    create database chat_test_db;
"
# redis-cli flushall
redis-cli -h localhost -p 6380 flushall

# running this from root dir
mysql -uroot -proot --protocol=tcp -h localhost --port=3310 chat_test_db < ./db/schema.sql
docker exec -it messaging-service-msgserver-1 bash -c "cd integration-tests && go test ./... -v -count=1 -failfast"
# docker compose -f docker-compose-integration.yml run msgserver bash -c "cd integration-tests && go test ./... -v -count=1 -failfast"
