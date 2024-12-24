rd /s /q bin
mkdir bin
set GOARCH=arm64
set GOOS=linux
go build -o bin/carousel
sqlite3.exe bin/carousel.db ".read .\scripts\sqlite\carousel-tables-creat.sql"
