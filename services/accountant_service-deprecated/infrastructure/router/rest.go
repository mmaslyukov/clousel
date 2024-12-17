package infrastructure

import (
	"accountant_service/domain/accountment"
	"accountant_service/domain/accountment/accountment_provider"
	"accountant_service/domain/carousel"
	"accountant_service/domain/carousel/carousel_provider"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Router struct {
	sales     accountment_provider.IPortApiSales
	analytics accountment_provider.IPortApiAnalytics
	ride      carousel_provider.IPortApiRide
}

func RouterCreate(
	sales accountment_provider.IPortApiSales,
	analytics accountment_provider.IPortApiAnalytics,
	ride carousel_provider.IPortApiRide,
) Router {
	return Router{
		sales:     sales,
		analytics: analytics,
		ride:      ride,
	}
}

func (rt *Router) ServeMux() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("POST /payment/initiate",
		func(w http.ResponseWriter, r *http.Request) {
		})
	router.HandleFunc("POST /payment/complete",
		func(w http.ResponseWriter, r *http.Request) {
		})
	router.HandleFunc("POST /payment/failure",
		func(w http.ResponseWriter, r *http.Request) {
		})

	router.HandleFunc("POST /carousel/pricetags",
		func(w http.ResponseWriter, r *http.Request) {
			var err error
			pt := accountment.PriceTagsDetails{}
			if err = rt.sales.WritePriceTags(pt); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		})

	router.HandleFunc("GET /carousel/pricetags",
		func(w http.ResponseWriter, r *http.Request) {
			var err error
			var ptd accountment.PriceTagsDetails
			carouselId := uuid.New()
			if ptd, err = rt.sales.ReadPriceTags(carouselId); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err = json.NewEncoder(w).Encode(ptd); err != nil {
				http.Error(w, err.Error(), http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
	router.HandleFunc("GET /rides/undelivered",
		func(w http.ResponseWriter, r *http.Request) {
			var err error
			var rides []carousel.RideMinimal
			carouselId := uuid.New()
			if rides, err = rt.ride.ReadUndeliveredRides(carouselId); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err = json.NewEncoder(w).Encode(rides); err != nil {
				http.Error(w, err.Error(), http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusOK)

		})
	return router
}

func delmeRouter() *http.ServeMux {
	router := http.NewServeMux()

	// decoder := json.NewDecoder(r.Body)
	// var t Register
	// err := decoder.Decode(&t)
	// if err != nil {
	// 	logger.Error.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// regData := pch.RegisterData{
	// 	CarouselId: pch.CarouselId{Id: t.CarouselId},
	// 	OwnerId:    t.OwnerId,
	// 	RoundTime:  NewOptionalValue[int](*t.RoundTime),
	// }
	// if t.RoundTime != nil {
	// 	regData.RoundTime.Set(*t.RoundTime)
	// }
	// err = carousel.Register(regData)
	// if err != nil {
	// 	logger.Error.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// w.WriteHeader(http.StatusOK)
	// return
	router.HandleFunc("POST /carousel/refill",
		func(w http.ResponseWriter, r *http.Request) {
			// decoder := json.NewDecoder(r.Body)
			// var t Refill
			// _ = decoder.Decode(&t)

			// rd := pch.RefillData{
			// 	CarouselId:  pch.CarouselId{Id: t.CarouselId},
			// 	RoundsReady: t.Rounds,
			// }

			// var err error
			// if err = carousel.Refill(rd); err != nil {
			// 	logger.Error.Println(err)
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	w.Write([]byte(err.Error()))
			// 	return
			// }
			// w.WriteHeader(http.StatusOK)

		})

	return router
}
