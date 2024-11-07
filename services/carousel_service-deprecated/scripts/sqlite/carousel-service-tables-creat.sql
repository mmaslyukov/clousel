-- .read scripts/sqlite/carousel-service-tables-creat.sql

CREATE TABLE "carousel-service-record" (
    "CarouselId" string PRIMARY KEY,
    "OwnerId" string NOT NULL,
    "RoundTime" int  NULL
);

CREATE TABLE "carousel-service-log" (
    "CarouselId" string UNIQUE NOT NULL,
    "EventId" string NOT NULL,
    "Time" datetime DEFAULT CURRENT_TIMESTAMP,
    "StatusChange" int,
    "RoundsChange" int,
    "Rounds" int NOT NULL,
    "Error" string,
    FOREIGN KEY (
        "CarouselId"
    )
    REFERENCES "carousel-service-record" ( CarouselId )
);

CREATE TABLE "carousel-service-status" (
    "CarouselId" string UNIQUE NOT NULL ,
    "Time" datetime DEFAULT CURRENT_TIMESTAMP,
    "Status" string NOT NULL,
    "RoundsReady" int NOT NULL,
    FOREIGN KEY (
        "CarouselId"
    )
    REFERENCES "carousel-service-record" ( CarouselId )
);

CREATE TABLE "carousel-service-evt-queue" (
    "EventId" string UNIQUE NOT NULL ,
    "Time" datetime DEFAULT CURRENT_TIMESTAMP,
    "Type" string NOT NULL,
    "Data" string NOT NULL
);

-- CREATE TABLE test (
--     id int,
--     one int,
--     two int NULL
-- )
