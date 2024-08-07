## URL

`https://app.quickdatabasediagrams.com/#/`


# SCHEMA
````
carousel-service-record as csr
---
CarouselId  PK string 
RoundTime int NULL

carousel-service-log
---
Time DATETIME
CarouselId string FK >- csr.CarouselId
RoundsChange int

carousel-service-price
---
CarouselId string FK >- csr.CarouselId
Rounds int
RoundsPrice int

carousel-service-status
---
CarouselId string FK >- csr.CarouselId
Status string
RoundsReady int
```
