proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./common/proto/**/*.proto

build: proto
	go build ./client/... ./server/... ./broker/...

test:
	go test -v -race ./common/... ./client/... ./server/... ./broker/...

tidy:
	cd common && go mod tidy
	cd client && go mod tidy
	cd server && go mod tidy
	cd broker && go mod tidy
