.PHONY: dev build start clean

dev:
	go run cmd/server/main.go


build:
	go build -o bin/api cmd/server/main.go
	chmod +x bin/api

start:
	./bin/api

clean:
	rm -f bin/api bin/api.exe 

run: clean build start