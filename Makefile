BINARY_NAME	:=	rv32i-as

MODULE_PATH	:=	github.com/ayase-mstk/go32as
#SRC         := $(shell find . -name '*.go' ! -path './src/*')
SRC					:= src/main.go
TEST_DIR		:=	test
TEST_NAME		:=	ztail_test.go
RM					:=	rm -rf

all: init build

init:
	@if [ ! -f go.mod ]; then \
		go mod init ${MODULE_PATH}; \
	fi

build:
	go build -o ${BINARY_NAME} ${SRC}

run: all
	./${BINARY_NAME}

fmt:
	gofmt -s -w .

test: init
	go test -v $(TEST_DIR)/$(TEST_NAME) --short

clean:
	go clean
	${RM} go.mod $(TEST_DIR)/test.test

fclean: clean
	${RM} ${BINARY_NAME} $(TEST_DIR)/${BINARY_NAME}

re: fclean all

.PHONY: init, build, run, fmt, clean, fclean, re
