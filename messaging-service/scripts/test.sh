#!/bin/bash

# // check to see if chat_api_mysqldb_unit_test exists or it's stopped
# // https://stackoverflow.com/questions/38576337/how-to-execute-a-bash-command-only-if-a-docker-container-with-a-given-name-does

# mysqladmin -uroot -proot ping -h localhost --port=3310 --protocol=tcp
mysql -uroot -proot --protocol=tcp -h localhost --port=3310 -e "
    drop database if exists chat_test_db;
    create database chat_test_db;
"
mysql -uroot -proot --protocol=tcp -h localhost --port=3310 chat_test_db < ../db/schema.sql

cd ../src
go test ./... -v -count=1

