#!/usr/bin/env bash

# // docker run --health-cmd='mysqladmin ping --silent' -d mysql
# // ping the database until it's ready
# https://stackoverflow.com/questions/30494050/how-do-i-pass-environment-variables-to-docker-containers

# setup mysql
if [[ -z $(docker ps -f name=chat_api_network -q) ]]
then
    echo "Creating chat_api_network network..."
    docker network create chat_api_network 
else 
    echo "Found chat_api_network network"
fi 

if [[ -z $(docker container ps -f name=^chat_api_mysqldb$ -q) ]]
then
	echo "Removing container..."
    docker container rm -f chat_api_mysqldb

    echo "Creating docker chat_api_mysqldb database..."
    docker run -d \
    -e MYSQL_ROOT_PASSWORD="root" \
    -e MYSQL_ROOT_USER="root" \
    -e MYSQL_ROOT_HOST="%" \
    --name=chat_api_mysqldb \
    -p=3310:3306  \
    --network=chat_api_network \
    mysql:8.0
fi

healthcheckResult=$(mysqladmin -uroot -proot ping -h localhost --port=3310 --protocol=tcp)
echo $healthcheckResult
while [[ $healthcheckResult != "mysqld is alive" ]]
do
    echo "Pinging mysqldb for health check..."
    sleep 1
    healthcheckResult=$(mysqladmin -uroot -proot ping -h localhost --port=3310 --protocol=tcp)
done
echo "mysql started"

mysql -uroot -proot --protocol=tcp -h localhost --port=3310 -e "
    drop database if exists chat_api_db;
    create database chat_api_db;
"

setup redis

