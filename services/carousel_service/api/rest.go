package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"carousel_service/internal/logger"
	// . "carousel_service/internal/ports"
	. "carousel_service/internal/utils"

	pch "carousel_service/internal/ports/port_carousel"
)

// {
// 	"CarouselId": "550e8400-e29b-41d4-a716-446655440000",
// 	"RoundTime": 30,
// 	"RoundsReady": 1
// }

type Register struct {
	CarouselId  string `json:"CarouselId"`
	OwnerId     string `json:"OwnerId"`
	RoundTime   *int   `json:"RoundTime,omitempty"`
	RoundsReady *int   `json:"RoundsReady,omitempty"`
}
type Refill struct {
	CarouselId string `json:"CarouselId"`
	Rounds     int    `json:"Rounds"`
}
type Play struct {
	CarouselId string `json:"CarouselId"`
}
type Carousel struct {
	CarouselId  string `json:"CarouselId"`
	RoundTime   *int   `json:"RoundTime,omitempty"`
	RoundsReady *int   `json:"RoundsReady,omitempty"`
	Status      string `json:"RoundsReady,omitempty"`
	Time        string `json:"RoundsReady,omitempty"`
}

func Router(carousel pch.CarouselInterface) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /carousel",
		func(w http.ResponseWriter, r *http.Request) {
			//TODO: go func(w http.ResponseWriter, r *http.Request) {...}
			var qvalue []string
			var ok bool
			if qvalue, ok = r.URL.Query()["CarouselId"]; !ok {
				err := fmt.Errorf("Query is invalid")
				logger.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return

			}
			data, err := carousel.Read(pch.CarouselId{Id: qvalue[0]})
			// data, err := carousel.NewCarouselHandler().Read(pcl.CarouselId{Id: qvalue[0]})
			if err != nil {
				logger.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			c := Carousel{
				CarouselId:  data.Ptr().Id,
				RoundTime:   data.Ptr().RoundTime.Ptr(),
				RoundsReady: data.Ptr().RoundsReady.Ptr(),
				Status:      data.Ptr().Status,
				Time:        data.Ptr().Time,
			}
			json.NewEncoder(w).Encode(c)
			w.WriteHeader(http.StatusOK)
		})
	router.HandleFunc("GET /carousel/owned",
		func(w http.ResponseWriter, r *http.Request) {
			//TODO: go func(w http.ResponseWriter, r *http.Request) {...}
			var qvalue []string
			var ok bool
			if qvalue, ok = r.URL.Query()["OwnerId"]; !ok {
				err := fmt.Errorf("Query is invalid")
				logger.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return

			}
			data, err := carousel.ReadByOwner(qvalue[0])
			// data, err := carousel.NewCarouselHandler().Read(pcl.CarouselId{Id: qvalue[0]})
			if err != nil {
				logger.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			var ownedCarouselsArray []Carousel
			for _, owned := range data.Get() {
				c := Carousel{
					CarouselId:  owned.Id,
					RoundTime:   owned.RoundTime.Ptr(),
					RoundsReady: owned.RoundsReady.Ptr(),
					Status:      owned.Status,
					Time:        owned.Time,
				}
				ownedCarouselsArray = append(ownedCarouselsArray, c)

			}
			json.NewEncoder(w).Encode(ownedCarouselsArray)
			w.WriteHeader(http.StatusOK)
		})
	router.HandleFunc("GET /carousel/play",
		func(w http.ResponseWriter, r *http.Request) {
			var qvalue []string
			var ok bool
			// var err error
			if qvalue, ok = r.URL.Query()["CarouselId"]; !ok {
				err := fmt.Errorf("Query is invalid")
				logger.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			// p := Play{
			// 	CarouselId: qvalue[0],
			// }
			// port.Receiver() <- p //temporary
			if err := carousel.Play(pch.CarouselId{Id: qvalue[0]}); err != nil {
				logger.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			w.WriteHeader(http.StatusOK)
		})
	router.HandleFunc("POST /carousel",
		func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			var t Register
			err := decoder.Decode(&t)
			if err != nil {
				logger.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			regData := pch.RegisterData{
				CarouselId: pch.CarouselId{Id: t.CarouselId},
				OwnerId:    t.OwnerId,
				RoundTime:  NewOptionalValue[int](*t.RoundTime),
			}
			if t.RoundTime != nil {
				regData.RoundTime.Set(*t.RoundTime)
			}
			err = carousel.Register(regData)
			if err != nil {
				logger.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		})
	router.HandleFunc("DELETE /carousel",
		func(w http.ResponseWriter, r *http.Request) {
			var qvalue []string
			var ok bool
			if qvalue, ok = r.URL.Query()["CarouselId"]; !ok {
				err := fmt.Errorf("Query is invalid")
				logger.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			err := carousel.Delete(pch.CarouselId{Id: qvalue[0]})
			if err != nil {
				logger.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		})
	router.HandleFunc("POST /carousel/refill",
		func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			var t Refill
			_ = decoder.Decode(&t)

			rd := pch.RefillData{
				CarouselId:  pch.CarouselId{Id: t.CarouselId},
				RoundsReady: t.Rounds,
			}

			var err error
			if err = carousel.Refill(rd); err != nil {
				logger.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusOK)

		})

	return router
}
