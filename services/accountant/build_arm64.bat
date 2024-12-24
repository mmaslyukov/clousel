rd /s /q bin
mkdir bin
set GOARCH=arm64
set GOOS=linux
go build -o bin/accountant
sqlite3.exe bin/accountant.db ".read .\scripts\sqlite\accountant-tables-creat.sql"
copy .env bin
