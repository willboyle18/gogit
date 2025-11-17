BINARY_NAME=gogit
OPERATING_SYSTEM=linux # change this variable depending on your OS (linux, darwin, windows)
TEST_DIRECTORY=test # change this depending on what directory you want to create the repository in

build:
	@GOARCH=amd64 GOOS=${OPERATING_SYSTEM} go build -o ./target/${BINARY_NAME}-${OPERATING_SYSTEM} cmd/gogit/main.go

init: build
	@cd ${TEST_DIRECTORY} && ../target/${BINARY_NAME}-${OPERATING_SYSTEM} init

clean:
	@go clean
	@rm target/${BINARY_NAME}-${OPERATING_SYSTEM}