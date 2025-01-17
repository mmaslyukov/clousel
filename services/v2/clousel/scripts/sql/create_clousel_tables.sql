-- sqlite3.exe .\clousel.db ".read .\scripts\sql\create_clousel_tables.sql"

CREATE TABLE IF NOT EXISTS "company" (
    "Id" string PRIMARY KEY,
    "Name" string UNIQUE NOT NULL,
    "Email" string UNIQUE NOT NULL,
    "Password" string NOT NULL,
    "ProductId" string UNIQUE,
    "SecKey" string UNIQUE,
    "WhId" string,
    "WhKey" string,
    "Enabled" int
);

CREATE TABLE IF NOT EXISTS "user" (
    "Id" string PRIMARY KEY,
    "Companyname" string NOT NULL,
    "Username" string UNIQUE NOT NULL,
    "Email" string UNIQUE NOT NULL,
    "Password" string NOT NULL,
    -- 'Balance' int, 
    'Time' datetime 
    -- 'Time' datetime DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "checkout" (
    "EventId" string PRIMARY KEY,
    "SessionId" string UNIQUE NOT NULL,
    "UserId" string NOT NULL,
    'Price' int, 
    'Tickets' int, 
    'Status' string, 
    'Time' datetime , 
    -- 'Time' datetime DEFAULT CURRENT_TIMESTAMP, 
    FOREIGN KEY ("UserId") REFERENCES "user" ( Id )
);

CREATE TABLE IF NOT EXISTS "balance" (
    "EventId" string PRIMARY KEY,
    "UserId" string NOT NULL,
    'Change' int, 
    FOREIGN KEY ("UserId") REFERENCES "user" ( Id )
);

CREATE TABLE IF NOT EXISTS "machine" (
    "Id" string PRIMARY KEY,
    "CompanyId" string NOT NULL,
    'Cost' int, 
    'Status' string, 
    'Fee' int,
    'Updated' datetime, 
    -- 'Updated' datetime DEFAULT CURRENT_TIMESTAMP, 
    FOREIGN KEY ("CompanyId") REFERENCES "company" ( Id )
);

CREATE TABLE IF NOT EXISTS "game" (
    "Id" string PRIMARY KEY,
    "MachId" string NOT NULL,
    "UserId" string NOT NULL,
    'Cost' int, 
    'Status' string, 
    'Started' datetime, 
    -- 'Started' datetime DEFAULT CURRENT_TIMESTAMP, 
    FOREIGN KEY ("MachId") REFERENCES "machine" ( Id ),
    FOREIGN KEY ("UserId") REFERENCES "user" ( Id )
);
