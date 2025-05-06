## RUN APPLICATION
run:
	@echo -e "Running the application..."
	@go run main.go

## RUN TESTS
test:
	@echo -e "Running tests..."
	@go test -v ./tests/department_test.go