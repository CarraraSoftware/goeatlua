all: gen
	go run .

build: gen
	go build .

gen:
	go generate _interp/lua.go > interp/lua.c.go
