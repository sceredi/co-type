module github.com/sceredi/co-type/broker

go 1.26.0

replace github.com/sceredi/co-type/common => ../common

require (
	github.com/sceredi/co-type/common v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.85.0
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260817212433-ac3dfec99bb1 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)
