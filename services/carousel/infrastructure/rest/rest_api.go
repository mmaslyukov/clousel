package rest

import (
	"carousel/core/manager"
	"carousel/core/operator"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/zerolog"
)

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func register(r *http.Request, manPort manager.IPortManagerControllerApi, log *zerolog.Logger) error {
	var err error
	var carousel manager.Carousel
	for ok := true; ok; ok = false {
		carousel.OwnId = r.PostFormValue("OwnerId")
		carousel.CarId = r.PostFormValue("CarouselId")
		if err = manPort.Register(carousel); err == nil {
			log.Info().Str("CarouselId", carousel.CarId).Msg("Rest.Register: Success")
		} else {
			log.Err(err).Str("CarouselId", carousel.CarId).Msg("Rest.Register: Fail")
		}
	}

	return err
}

func unregister(r *http.Request, manPort manager.IPortManagerControllerApi, log *zerolog.Logger) error {
	const qcNameCid = "CarouselId"
	const qcNameOid = "OwnerId"
	var carousel manager.Carousel //{Cid: "", Oid: ""}
	var err error
	if qvalue, ok := r.URL.Query()[qcNameCid]; ok {
		carousel.CarId = qvalue[0]
	}
	if qvalue, ok := r.URL.Query()[qcNameOid]; ok {
		carousel.OwnId = qvalue[0]
	}
	if err = manPort.Unregister(carousel); err == nil {
		log.Info().Str(qcNameOid, carousel.OwnId).Str(qcNameCid, carousel.CarId).Msg("Rest.Unregister: Success")
	} else {
		log.Err(err).Str(qcNameOid, carousel.OwnId).Str(qcNameCid, carousel.CarId).Msg("Rest.Unregister: Fail")
	}
	return err
}

func readOwned(r *http.Request, manPort manager.IPortManagerControllerApi, log *zerolog.Logger) ([]manager.Carousel, error) {
	const qcName = "OwnerId"
	var err error
	ownerId := "Unknown"
	var qvalue []string
	var ok bool
	var recordArray []manager.Carousel
	if qvalue, ok = r.URL.Query()[qcName]; ok {
		ownerId = qvalue[0]
		recordArray, err = manPort.Read(ownerId)
	} else {
		err = fmt.Errorf("Cannot find '%s' key in the query", qcName)
	}
	if err != nil {
		log.Err(err).Str("OwnerId", ownerId).Msg("Rest.ReadOwned: Fail")
	} else {
		log.Info().Str("OwnerId", ownerId).Msg("Rest.ReadOwned: Success")
	}
	return recordArray, err
}

func readSnapshot(r *http.Request, opPort operator.IPortOperatorControllerApi, log *zerolog.Logger) (*operator.SnapshotData, error) {
	const qcName = "CarouselId"
	var err error
	var qvalue []string
	var ok bool
	var carousel operator.Carousel
	var sd *operator.SnapshotData
	if qvalue, ok = r.URL.Query()[qcName]; ok {
		carousel.CarId = qvalue[0]
		sd, err = opPort.Read(carousel)
	} else {
		err = fmt.Errorf("Cannot find '%s' key in the query", qcName)
	}
	if err != nil {
		log.Err(err).Str("CarousleId", carousel.CarId).Msg("Rest.ReadSnapshot: Fail")
	} else {
		log.Info().Str("CarousleId", carousel.CarId).Msg("Rest.ReadSnapshot: Success")
	}
	return sd, err
}

func readPending(_ *http.Request, opPort operator.IPortOperatorControllerApi, log *zerolog.Logger) ([]operator.CompositeData, error) {
	var err error
	var cdArray []operator.CompositeData
	if cdArray, err = opPort.ReadPending(); err == nil {
		log.Info().Msg("Rest.readPending: Success")
	} else {
		log.Err(err).Msg("Rest.readPending: Fail")
	}
	return cdArray, err
}

func readByStatus(r *http.Request, opPort operator.IPortOperatorControllerApi, log *zerolog.Logger) ([]operator.SnapshotData, error) {
	const qcName = "Status"
	var status string
	var err error
	var qvalue []string
	var ok bool
	var sd []operator.SnapshotData
	if qvalue, ok = r.URL.Query()[qcName]; ok {
		status = qvalue[0]
		sd, err = opPort.ReadByStatus(status)
	} else {
		err = fmt.Errorf("Cannot find '%s' key in the query", qcName)
	}
	if err != nil {
		log.Err(err).Str(qcName, status).Msg("Rest.readByStatus: Fail")
	} else {
		log.Info().Str(qcName, status).Msg("Rest.readByStatus: Success")
	}
	return sd, err
}

func play(r *http.Request, opPort operator.IPortOperatorControllerApi, log *zerolog.Logger) error {
	const qcNameCid = "CarouselId"
	var c operator.Carousel
	var err error
	c.CarId = r.PostFormValue("CarouselId")
	if err = opPort.Play(c); err == nil {
		log.Info().Str(qcNameCid, c.CarId).Msg("Rest.Play: Success")
	} else {
		log.Err(err).Str(qcNameCid, c.CarId).Msg("Rest.Play: Fail")
	}
	return err
}

func refill(r *http.Request, opPort operator.IPortOperatorControllerApi, log *zerolog.Logger) error {
	var err error
	var carousel operator.Carousel
	for ok := true; ok; ok = false {
		var tickets int
		tickets, err = strconv.Atoi(r.PostFormValue("Tickets"))
		if err != nil {
			log.Err(err).Msg("Rest.Refill: Fail to parse 'Tickets' key in POST request")
			break
		}
		carousel.CarId = r.PostFormValue("CarouselId")
		if err = opPort.Refill(carousel, tickets); err == nil {
			log.Info().Str("CarouselId", carousel.CarId).Int("Tickets", tickets).Msg("Rest.Refill: Success")
		} else {
			log.Err(err).Str("CarouselId", carousel.CarId).Msg("Rest.Refill: Fail")
		}
	}
	return err
}

func New(manPort manager.IPortManagerControllerApi, opPort operator.IPortOperatorControllerApi, log *zerolog.Logger) *http.ServeMux {
	// var err error
	router := http.NewServeMux()
	router.HandleFunc("POST /carousel",
		func(w http.ResponseWriter, r *http.Request) {
			setupCORS(&w)
			if err := register(r, manPort, log); err == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		})
	router.HandleFunc("DELETE /carousel",
		func(w http.ResponseWriter, r *http.Request) {
			setupCORS(&w)
			if err := unregister(r, manPort, log); err == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		})
	router.HandleFunc("GET /carousel",
		func(w http.ResponseWriter, r *http.Request) {
			setupCORS(&w)
			if snapshot, err := readSnapshot(r, opPort, log); err == nil {
				json.NewEncoder(w).Encode(snapshot)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				// w.WriteHeader(http.StatusInternalServerError)
			}
			return
		})
	router.HandleFunc("GET /carousel/owned",
		func(w http.ResponseWriter, r *http.Request) {
			setupCORS(&w)
			if recordArray, err := readOwned(r, manPort, log); err == nil {
				json.NewEncoder(w).Encode(recordArray)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		})
	router.HandleFunc("GET /carousel/wstatus",
		func(w http.ResponseWriter, r *http.Request) {
			setupCORS(&w)
			if recordArray, err := readByStatus(r, opPort, log); err == nil {
				json.NewEncoder(w).Encode(recordArray)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		})
	router.HandleFunc("GET /carousel/pending",
		func(w http.ResponseWriter, r *http.Request) {
			setupCORS(&w)
			if recordArray, err := readPending(r, opPort, log); err == nil {
				json.NewEncoder(w).Encode(recordArray)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		})
	router.HandleFunc("POST /carousel/play",
		func(w http.ResponseWriter, r *http.Request) {
			setupCORS(&w)
			if err := play(r, opPort, log); err == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		})
	router.HandleFunc("POST /carousel/refill",
		func(w http.ResponseWriter, r *http.Request) {
			setupCORS(&w)
			if err := refill(r, opPort, log); err == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		})

	return router
}
