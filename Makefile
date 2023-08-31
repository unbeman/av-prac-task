swagger:
	swag init --parseDependency  --parseInternal -g cmd/main.go

mock:
	mockgen -destination=internal/database/mock/database_mock.go "github.com/unbeman/av-prac-task/internal/database" IDatabase

test:
	go test --cover ./...

dc-up:
	docker-compose up -d