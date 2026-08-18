
build:
	@go build -o client.o bin/client.go
	@go build -o server.o bin/server.go

s: # server
	@bin/server.o

c: # client
	@bin/client.o

clean:
	rm -f bin client.o
