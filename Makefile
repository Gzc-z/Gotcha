
build:
	@go build -o bin/client src/client/client.go
	@go build -o bin/server src/server/server.go

s: build
	@bin/server

c: build
	@bin/client

clean:
	@rm -r bin
