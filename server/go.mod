module github.com/sceredi/co-type/server

go 1.25.7

replace github.com/sceredi/co-type/common => ../common

require (
	github.com/sceredi/co-type/common v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.82.1
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
