build-all:
	go build ./client/... ./server/... ./broker/...

test-all:
	go test ./common/... ./client/... ./server/... ./broker/...

tidy-all:
	cd common && go mod tidy
	cd client && go mod tidy
	cd server && go mod tidy
	cd broker && go mod tidy
