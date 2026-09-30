DIR ?= wctool

build:
	@cd ./${DIR} && go build -o ./dist/ .

run: build
	@cd ./${DIR} && ./dist/${DIR} ${ARGS}