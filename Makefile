# Запуск для проверки
run:
	go run ./cmd/main.go

# Приминить миграцию
migrate-up:
	~/go/bin/goose --dir=./db/migrations sqlite3 ./storage/storage.db up

migrate-down:
	~/go/bin/goose --dir=./db/migrations sqlite3 ./storage/storage.db down

# Создание миграции
migrate-create:
	~/go/bin/goose --dir=./db/migrations sqlite3 ./storage/storage.db create $(table) sql