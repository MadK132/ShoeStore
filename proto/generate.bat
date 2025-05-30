@echo off
REM Генерация Go-кода из proto-файлов
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto
 
echo Proto files generated successfully! 