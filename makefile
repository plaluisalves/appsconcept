# Mark phony targets
.PHONY: test cover update install run air

test: # test will run all tests
	@echo ==== Running tests ====
	go test ./... -v -count=1
	@echo ==========================================================================================
cover: # run the coverage tool
	@echo ==== Running coverage ====
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -html=coverage.out -o coverage.html
	@echo ==========================================================================================
update: # update will update go.mod dependencies
	@echo ==== Updating go.mod dependencies ====
	go get -u ./...
	go get tool
	go mod tidy
	@echo ==========================================================================================
install:
	@echo ==== Installing dependencies by go tool ====
	go install tool
	@echo ==========================================================================================
run: # run will run the main application
	@echo ==== Running the Main application ====
	-go run cmd/web/main.go
	@echo ==========================================================================================
air: # will run air tool (need a .air.toml file)
	@echo ==== Running air ====
	-go tool air
	@echo ==========================================================================================
