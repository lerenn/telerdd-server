.PHONY: bin

BIN_DIR=bin
PACKAGE=github.com/lerenn/telerdd-server
RELEASE_DIR=release
SRC_DIR=src
WEBCLIENT_DIR=webclient

all: bin

bin: fmt
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/telerdd-server ./$(SRC_DIR)

clean:
	@rm -rf ./$(BIN_DIR)
	@rm -rf ./$(RELEASE_DIR)

docker:
	@sudo systemctl start docker
	@sudo docker build -t telerdd .

fmt:
	@go fmt $(PACKAGE)/${SRC_DIR}/...
