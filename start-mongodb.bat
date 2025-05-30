@echo off
echo Starting MongoDB...
start cmd /k "C:\Programm Files\MongoDB\Server\6.0\bin\mongod.exe"
timeout /t 5

echo Seeding database...
cd product-service/scripts && go run seed.go && cd ../..

echo MongoDB is ready! 