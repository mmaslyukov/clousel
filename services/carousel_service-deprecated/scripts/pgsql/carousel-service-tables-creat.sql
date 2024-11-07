-- Exported from QuickDBD: https://www.quickdatabasediagrams.com/
-- NOTE! If you have used non-SQL datatypes in your design, you will have to change these here.


CREATE TABLE "carousel-service-record" (
    "CarouselId" string   NOT NULL,
    "RoundTime" int   NULL,
    CONSTRAINT "pk_carousel-service-record" PRIMARY KEY (
        "CarouselId"
     )
);

CREATE TABLE "carousel-service-log" (
    "Time" DATETIME   NOT NULL,
    "CarouselId" string   NOT NULL,
    "RoundsChange" int   NOT NULL
);

CREATE TABLE "carousel-service-price" (
    "CarouselId" string   NOT NULL,
    "Rounds" int   NOT NULL,
    "RoundsPrice" int   NOT NULL
);

CREATE TABLE "carousel-service-status" (
    "CarouselId" string   NOT NULL,
    "Status" string   NOT NULL,
    "RoundsReady" int   NOT NULL
);

ALTER TABLE "carousel-service-log" ADD CONSTRAINT "fk_carousel-service-log_CarouselId" FOREIGN KEY("CarouselId")
REFERENCES "carousel-service-record" ("CarouselId");

ALTER TABLE "carousel-service-price" ADD CONSTRAINT "fk_carousel-service-price_CarouselId" FOREIGN KEY("CarouselId")
REFERENCES "carousel-service-record" ("CarouselId");

ALTER TABLE "carousel-service-status" ADD CONSTRAINT "fk_carousel-service-status_CarouselId" FOREIGN KEY("CarouselId")
REFERENCES "carousel-service-record" ("CarouselId");

