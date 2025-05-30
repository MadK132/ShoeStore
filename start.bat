@echo off
echo Starting Redis...
start cmd /k "C:\redis\redis-server.exe"
timeout /t 5

echo Starting NATS...
start cmd /k "C:\nats-server\nats-server.exe"
timeout /t 5

echo Setting up Email Service environment...
call email-service/cmd/set-env.bat
timeout /t 2

echo Starting microservices...
cd email-service/cmd && start cmd /k "go run main.go" && cd ../..
timeout /t 3
cd user-service/cmd && start cmd /k "go run main.go" && cd ../..
timeout /t 3
cd product-service/cmd && start cmd /k "go run main.go" && cd ../..
timeout /t 3
cd order-service/cmd && start cmd /k "go run main.go" && cd ../..
timeout /t 3
cd api-gateway && start cmd /k "go run main.go" && cd ..

echo Installing frontend dependencies...
cd frontend
call install.bat
echo Starting frontend...
start cmd /k "npm start"
cd ..

echo All services started! 
echo Frontend: http://localhost:3000
echo API Gateway: http://localhost:8080 