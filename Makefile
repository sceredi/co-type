build:
	go build ./client/... ./server/... ./broker/...

test:
	go test -v -race ./common/... ./client/... ./server/... ./broker/...

tidy:
	cd common && go mod tidy
	cd client && go mod tidy
	cd server && go mod tidy
	cd broker && go mod tidy
