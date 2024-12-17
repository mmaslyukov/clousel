-- .read scripts/sqlite/accountant-tables-creat.sql
-- sqlite3.exe .\accountant.db ".read .\scripts\accountant-tables-creat.sql"

CREATE TABLE IF NOT EXISTS "profile" (
    "OwnerId" string PRIMARY KEY,
    "Email" string UNIQUE NOT NULL,
    "Password" string NOT NULL,
    "SecretKey" string UNIQUE,
    "PublishKey" string UNIQUE,
    "WebhookId" string,
    "WebhookKey" string,
    "Role" int ,
    "Time" datetime DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS "product" (
    "OwnerId" string NOT NULL,
    "CarouselId" string PRIMARY KEY,
    "ProductId" string,
    FOREIGN KEY ("OwnerId") REFERENCES "profile" ( OwnerId )
);

CREATE TABLE IF NOT EXISTS 'book' (
    'Time' datetime DEFAULT CURRENT_TIMESTAMP, 
    'SessionId' string PRIMARY KEY,
    'CarouselId' string NOT NULL,
    'Amount' int, 
    'Tickets' int, 
    'Status' string NOT NULL, 
    'Error' string,
    FOREIGN KEY ("CarouselId") REFERENCES "product" ( CarouselId )
)
