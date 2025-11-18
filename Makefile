BINARY_NAME=gogit
OPERATING_SYSTEM=linux # change this variable depending on your OS (linux, darwin, windows)
TEST_DIRECTORY=test # change this depending on what directory you want to create the repository in
ARCHITECTURE=amd64 # change this depending on what computer architecture you have (amd64 for intel/amd cpus, arm64 for apple silicon)

build:
	@GOARCH=${ARCHITECTURE} GOOS=${OPERATING_SYSTEM} go build -o ./target/${BINARY_NAME}-${OPERATING_SYSTEM} cmd/gogit/main.go

init: build
	@cd ${TEST_DIRECTORY} && ../target/${BINARY_NAME}-${OPERATING_SYSTEM} init


add: build
	@cd ${TEST_DIRECTORY} && ../target/${BINARY_NAME}-${OPERATING_SYSTEM} ${ARGS}


ARGS := $(filter-out $@,$(MAKECMDGOALS))

%:
	@:

clean:
	@go clean
	@rm target/${BINARY_NAME}-${OPERATING_SYSTEM}

reset:
	@cd ${TEST_DIRECTORY} && rm -rf .gogit

setup:
	@mkdir ${TEST_DIRECTORY}
