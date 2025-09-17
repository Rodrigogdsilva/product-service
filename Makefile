include .env
export

.PHONY: start stop logs migrate-up migrate-down create-migration

start:
	@echo "Iniciando os containers Docker em segundo plano..."
	@docker-compose up -d --build

stop:
	@echo "Parando os containers Docker..."
	@docker-compose down -v --remove-orphans

# Mostra os logs da aplicação em tempo real
logs:
	@echo "Acompanhando os logs da aplicação... (Pressione Ctrl+C para sair)"
	@docker-compose logs -f app

migrate-up:
	@echo "Aplicando migrations..."
	@docker-compose run --rm migrator -database "$(DATABASE_URL)" -path "/migrations" up

migrate-down:
	@echo "Revertendo a última migration..."
	@docker-compose run --rm migrator -database "$(DATABASE_URL)" -path "/migrations" down 1

create-migration:
	@read -p "Digite o nome da migration: " name; \
	migrate create -ext sql -dir ./database -seq $$name

vulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

lint:
	golangci-lint run

gitleaks:
	docker run --rm -v $(CURDIR):/path zricethezav/gitleaks:latest detect --source=/path --verbose

test:
	@echo "Rodando os testes..."
	@go test -v ./...