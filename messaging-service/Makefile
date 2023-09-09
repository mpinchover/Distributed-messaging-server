test:
	cd src && ../scripts/test.sh

int-test:
	# docker compose run msgserver ../scripts/test.sh
	./scripts/test-integration.sh

migrate:
	mysql -uroot -proot --protocol=tcp -h localhost --port=3310 messaging < ./db/schema.sql

setup-dbs:
	cd dependencies && docker compose up -d

tear-down-dbs:
	cd dependencies && docker compose down

run-api:
	docker compose up

run-api-integration:
	docker compose -f docker-compose-integration.yml up

enter:
	docker exec -it messaging-service-msgserver-1 bash