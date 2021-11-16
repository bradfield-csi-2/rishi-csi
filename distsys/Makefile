server:
	go build -o dkvs-server cmd/server/main.go

client:
	go build -o dkvs-client cmd/client/main.go

proto:
	protoc --go_out=. command/command.proto
