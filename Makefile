BIN_DIR=bin
DOWNLOADS_DIR=downloads
PROJECT_DIR=github.com/nightwall/nightwall-server
TOOLS_DIR=tools
WEBCLIENT_DIR=webclient

all: api tools

api: fmt
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/nightwall-server .

clean:
	@rm -rf ./$(BIN_DIR)
		@rm -rf ./$(DOWNLOADS_DIR)

docker: api
	@sudo systemctl start docker
	@sudo docker build -t nightwall .

fmt:
	@go fmt $(PROJECT_DIR)

tools: fmt
	@mkdir -p $(BIN_DIR)/$(TOOLS_DIR)
	@go build -o $(BIN_DIR)/$(TOOLS_DIR)/nightwall-image-download ./$(TOOLS_DIR)/imgdwnld
