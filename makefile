create-network:
	@docker network create tuduit-network

dev:
	@docker-compose down
	@docker-compose up -d

down:
	@docker-compose down

up:
	@docker-compose up -d

rebuild:
	@docker-compose up --build -d

log-go:
	@docker logs --follow tuduit_golang --tail 1000

log-http:
	@docker logs --follow tuduit_golang --tail 1000 | awk '!/^\[.*\]/ || /^\[http\]/' | awk '{sub(/\[http\]\t/,""); print}'

log-mq:
	@docker logs --follow tuduit_golang --tail 1000 | awk '!/^\[.*\]/ || /^\[mq\]/' | awk '{sub(/\[mq\]\t/,""); print}'
	
log-cron:
	@docker logs --follow tuduit_golang --tail 1000 | awk '!/^\[.*\]/ || /^\[cron\]/' | awk '{sub(/\[cron\]\t/,""); print}'

log-grpc:
	@docker logs --follow tuduit_golang --tail 1000 | awk '!/^\[.*\]/ || /^\[grpc\]/' | awk '{sub(/\[grpc\]\t/,""); print}'

log-redis:
	@docker logs --follow tuduit_redis --tail 50

log-pg:
	@docker logs --follow tuduit_pg --tail 50

log-nginx:
	@docker logs --follow tuduit_nginx --tail 50

debug-on:
	@docker exec -t tuduit_golang bash -c "echo 'export USE_DEBUG=true' >> ~/.bashrc"

debug-off:
	@docker exec -t tuduit_golang bash -c "echo 'export USE_DEBUG=false' >> ~/.bashrc"

doc:
	@docker exec -t tuduit_golang autodoc
	@git add ./autodoc/openapi.yaml

new-migration:
	@read -p "enter migration name: " mig; \
	echo "migration name: $${mig}"; \
	docker exec -t tuduit_golang migrate create -ext sql -dir db/migrations -seq $${mig}
migrate-up:
	docker exec -t tuduit_golang migrate -source file://./db/migrations -database "sqlite3:///var/lib/sqlite3/tuduit.db" up

migrate-down:
	docker exec -t tuduit_golang migrate -source file://./db/migrations -database "sqlite3:///var/lib/sqlite3/tuduit.db" down

migrate-drop:
	docker exec -t tuduit_golang migrate -path ./db/migrations -database "sqlite3:///var/lib/sqlite3/tuduit.db" drop -f
