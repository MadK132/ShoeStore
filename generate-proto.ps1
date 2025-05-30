Set-Location -Path "proto"
protoc --go_out=. --go_opt=paths=source_relative user.proto
protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto
Set-Location -Path ".." 