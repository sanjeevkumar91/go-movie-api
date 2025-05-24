install_deps:
	chmod +x scripts/install.sh
	scripts/install.sh
	
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
