.PHONY: bin

BIN_DIR=bin
PROJECT_DIR=github.com/lerenn/telerdd-server
RELEASE_DIR=release
WEBCLIENT_DIR=webclient

all: bin

bin: fmt
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/telerdd-server .

clean:
	@rm -rf ./$(BIN_DIR)
	@rm -rf ./$(RELEASE_DIR)

docker:
	@sudo systemctl start docker
	@sudo docker build -t telerdd .

fmt:
	@go fmt $(PROJECT_DIR)
