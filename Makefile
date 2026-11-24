BINARY_NAME=rehand
DIST_DIR=dist

.PHONY: build start stop clean restart migrate-up migrate-down

build:
	@echo "Building..."
	@go build -o $(DIST_DIR)/$(BINARY_NAME) ./cmd/api
	@echo "Built!"
## start: builds and runs the application in the background
start: build
	@echo "Starting..."
	@./$(DIST_DIR)/$(BINARY_NAME) &
	@echo "Running!"
## stop: stops the running application
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "$(DIST_DIR)/$(BINARY_NAME)"
	@echo "Stopped!"
## restart: stops and restarts the application
restart: stop start
## clean: removes compiled binaries
clean:
	@echo "Cleaning..."
	@-rm -rf $(DIST_DIR)
	@go clean
	@echo "Cleaned!"

migrate-up:
	@echo "Migrating..."
	@go run ./cmd/migrate/main.go up
	@echo "Migrated!"

migrate-down:
	@echo "Migrating..."
	@go run ./cmd/migrate/main.go down
	@echo "Migrated!"