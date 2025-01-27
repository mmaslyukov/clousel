rd /s /q bin
mkdir bin
set GOARCH=arm64
set GOOS=linux
go build -o bin/clousel
sqlite3.exe bin/clousel_new.db ".read .\scripts\sql\create_clousel_tables.sql"
