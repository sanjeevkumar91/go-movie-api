rollbackCount ?= 1

install_deps:
	chmod +x scripts/install.sh
	scripts/install.sh
	go mod tidy
	liquibase update
	
start:
	go run cmd/main.go

run_tests:
	go test ./...

run_format:
	go fmt ./...

run_lint:
	golangci-lint run

run_prepush:
	$(MAKE) run_format
	$(MAKE) run_tests
	$(MAKE) run_lint

run_liquibase_in_go:
	go run cmd/main.go

run_migrations:
	liquibase update

rollback_migrations:
	liquibase rollbackCount $(rollbackCount)

run_tests_with_coverage:
	go test -cover -coverprofile=coverage.out ./...

coverage_report:	
	go tool cover -html=coverage.out
