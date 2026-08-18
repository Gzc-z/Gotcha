build-client:
	@go build -o client.o client/client.go
build-server:
	@go build -o server.o server/server.go

s: build-server
	@./server.o

c: build-client
	@./client.o

clean:
	rm -f server.o client.o
