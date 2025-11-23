# school-21-workshop

## 23.11.2025 Пирамида тестирования, тесты в Go

Пример урезанного проекта по автоматизации инвестирования в Т-Инвестиции. Содержит только ручку проверки и создания портфолио.

Ссылки:
- https://golang.testcontainers.org/modules/postgres/

Команды:
- `go get github.com/testcontainers/testcontainers-go/modules/postgres`
- `go get -tool go.uber.org/mock/mockgen`
- `//go:generate go tool mockgen -source=$GOFILE -destination contract_mock.go -package $GOPACKAGE`
- `//go:build integration`
- `go test -tags=integration ./...`
- `go test ./... -coverprofile=coverage.txt`
- `go tool cover -html coverage.txt -o index.html`