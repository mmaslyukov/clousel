-- .read scripts/sqlite/carousel-tables-creat.sql

CREATE TABLE "carousel-record" (
    "CarouselId" string PRIMARY KEY,
    "OwnerId" string NOT NULL,
    "Active" int
);

CREATE TABLE "carousel-event" (
    "CarouselId" string NOT NULL,
    "EventId" string UNIQUE NOT NULL,
    "Time" datetime DEFAULT CURRENT_TIMESTAMP,
    -- "Reason" string NOT NULL,
    "Status" string,
    "Tickets" int NOT NULL,
    "Pending" int,
    "Error" string,
    "Extra" string,
    FOREIGN KEY (
        "CarouselId"
    )
    REFERENCES "carousel-record" ( CarouselId )
);

CREATE TABLE "carousel-snapshot" (
    "CarouselId" string UNIQUE NOT NULL,
    "Status" int NOT NULL,
    "Tickets" int NOT NULL,
    "Extra" string,
    FOREIGN KEY (
        "CarouselId"
    )
    REFERENCES "carousel-record" ( CarouselId )
);

