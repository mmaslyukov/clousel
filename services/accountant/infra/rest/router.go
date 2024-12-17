package rest

import (
	"accountant/core/owner"
	erro "accountant/core/owner/error"
	"accountant/core/store"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	//"github.com/sanity-io/litter"
	"github.com/stripe/stripe-go/v72/webhook"
)

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func ownerRegister(
	w http.ResponseWriter,
	r *http.Request,
	oapi owner.IPortOwnerControllerOwnerApi,
	sapi store.IPortBookControllerApi,
	log *zerolog.Logger) {
	setupCORS(&w)
	for ok := true; ok; ok = false {
		log.Debug().Msgf("ownerRegister %s", r.Method)
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		email := r.PostFormValue("Email")
		password := r.PostFormValue("Password")
		if _, err := mail.ParseAddress(email); err != nil || len(password) < 4 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusNotAcceptable)
			break
		}

		if ierr := oapi.Register(email, password); ierr != nil {
			http.Error(w, fmt.Sprintf("Fail to register %v", ierr.Error()), http.StatusInternalServerError)
			break
		}
		w.WriteHeader(http.StatusOK)
	}
}

func ownerLogin(
	w http.ResponseWriter,
	r *http.Request,
	oapi owner.IPortOwnerControllerOwnerApi,
	sapi store.IPortBookControllerApi,
	log *zerolog.Logger) {
	var ierr owner.IError
	var token owner.Token
	setupCORS(&w)
	for ok := true; ok; ok = false {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		email := r.URL.Query().Get("Email")
		password := r.URL.Query().Get("Password")
		if token, ierr = oapi.Login(email, password); ierr != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			break
		}
		writeJSON(w, struct {
			Token owner.Token `json:"Token"`
		}{
			Token: token,
		})
	}
}

func ownerPkey(
	w http.ResponseWriter,
	r *http.Request,
	oapi owner.IPortOwnerControllerOwnerApi,
	sapi store.IPortBookControllerApi,
	log *zerolog.Logger) {
	setupCORS(&w)
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			log.Warn().Msg("StatusMethodNotAllowed")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		token, err := uuid.Parse(r.PostFormValue("Token"))
		if err != nil {
			log.Err(err).Msg("Fail to parse Token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			break
		}
		pkey := r.PostFormValue("PublisKey")

		if ierr := oapi.AssignPkeys(token, pkey); ierr != nil {
			log.Err(ierr).Msg("Fail AssignPSkeys")
			switch ierr.Code() {
			case erro.ECUserTokenNotFoundOrExpired:
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				break
			default:
				http.Error(w, fmt.Sprintf("Fail to stripe keys: %v", ierr.Error()), http.StatusInternalServerError)
				break
			}
			break
		}
	}
}

func ownerSkey(
	w http.ResponseWriter,
	r *http.Request,
	oapi owner.IPortOwnerControllerOwnerApi,
	sapi store.IPortBookControllerApi,
	log *zerolog.Logger) {
	setupCORS(&w)
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			log.Warn().Msg("StatusMethodNotAllowed")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		token, err := uuid.Parse(r.PostFormValue("Token"))
		if err != nil {
			log.Err(err).Msg("Fail to parse Token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			break
		}
		skey := r.PostFormValue("SecretKey")

		if ierr := oapi.AssignSkeys(token, skey); ierr != nil {
			log.Err(ierr).Msg("Fail AssignPSkeys")
			switch ierr.Code() {
			case erro.ECUserTokenNotFoundOrExpired:
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				break
			default:
				http.Error(w, fmt.Sprintf("Fail to stripe keys: %v", ierr.Error()), http.StatusInternalServerError)
				break
			}
			break
		}
	}
}

func ownerWebhookRefresh(
	w http.ResponseWriter,
	r *http.Request,
	oapi owner.IPortOwnerControllerOwnerApi,
	sapi store.IPortBookControllerApi,
	log *zerolog.Logger) {

	setupCORS(&w)
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		token, err := uuid.Parse(r.PostFormValue("Token"))
		if err != nil {
			log.Err(err).Msg("Fail to parse Token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			break
		}
		if ierr := oapi.RefreshWebhook(token); ierr != nil {
			log.Err(err)
			http.Error(w, ierr.Error(), http.StatusInternalServerError)
			break
		}

		writeJSON(w, nil)
	}
}

func carouselAdd(
	w http.ResponseWriter,
	r *http.Request,
	oapi owner.IPortOwnerControllerOwnerApi,
	sapi store.IPortBookControllerApi,
	log *zerolog.Logger) {
	setupCORS(&w)
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			log.Warn().Msg("StatusMethodNotAllowed")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		token, err := uuid.Parse(r.PostFormValue("Token"))
		if err != nil {
			log.Err(err).Msg("Fail to parse Token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			break
		}
		carId, err := uuid.Parse(r.PostFormValue("CarouselId"))
		if err != nil {
			log.Err(err).Msg("Fail to parse CarouselId")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			break
		}
		prodId := r.PostFormValue("ProdId")
		var ierr owner.IError
		if len(prodId) > 0 {
			ierr = oapi.AddCarousel(token, carId, &prodId)
		} else {
			ierr = oapi.AddCarousel(token, carId, nil)
		}
		if ierr != nil {
			switch ierr.Code() {
			case erro.ECUserTokenNotFoundOrExpired:
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				break
			default:
				http.Error(w, fmt.Sprintf("Fail to Add carousel: %v", ierr.Error()), http.StatusInternalServerError)
				break
			}
			break
		}
		w.WriteHeader(http.StatusOK)
	}
}
func carouselProdId(
	w http.ResponseWriter,
	r *http.Request,
	oapi owner.IPortOwnerControllerOwnerApi,
	sapi store.IPortBookControllerApi,
	log *zerolog.Logger) {
	setupCORS(&w)
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			log.Warn().Msg("StatusMethodNotAllowed")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		token, err := uuid.Parse(r.PostFormValue("Token"))
		if err != nil {
			log.Err(err).Msg("Fail to parse Token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			break
		}
		carId, err := uuid.Parse(r.PostFormValue("CarouselId"))
		if err != nil {
			log.Err(err).Msg("Fail to parse CarouselId")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			break
		}
		prodId := r.PostFormValue("ProdId")

		if ierr := oapi.AssignProdId(token, carId, prodId); ierr != nil {
			switch ierr.Code() {
			case erro.ECUserTokenNotFoundOrExpired:
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				break
			default:
				http.Error(w, fmt.Sprintf("Fail to Assign prodId: %v", ierr.Error()), http.StatusInternalServerError)
				break
			}
			break
		}
	}
}

// ownerRegister
// ownerLogin
// carouselAdd
// carouselProdId
// carouselSPkey
// clientCheckout
// clientPrices

// CARD NUM 4242424242424242
func clientCheckout(
	w http.ResponseWriter,
	r *http.Request,
	oapi owner.IPortOwnerControllerOwnerApi,
	sapi store.IPortBookControllerApi,
	log *zerolog.Logger) {
	setupCORS(&w)
	for ok := true; ok; ok = false {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		carId, err := uuid.Parse(r.URL.Query().Get("CarouselId"))
		if err != nil {
			log.Err(err).Msg("Fail to parse CarouselId")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			break
		}
		priceId := r.URL.Query().Get("PriceId")
		homeUrl, err := url.QueryUnescape(r.URL.Query().Get("HomeUrl"))
		log.Debug().Str("PriceId", priceId).Str("HomeUrl", homeUrl).Str("CarId", carId.String()).Send()
		if err != nil {
			http.Error(w, fmt.Sprintf("Fail to parse HomeUrl %v", err.Error()), http.StatusInternalServerError)
			break
		}
		cs, ierr := sapi.Checkout(carId, priceId, homeUrl)
		if ierr != nil {
			http.Error(w, fmt.Sprintf("Fail to create checkout session %v", ierr.Error()), http.StatusInternalServerError)
			break
		}
		// writeJSON(w, struct {
		// 	Url string `json:"Url"`
		// }{
		// 	Url: cs.Url(),
		// })
		http.Redirect(w, r, cs.Url(), http.StatusOK)
	}
}
func clientPrices(
	w http.ResponseWriter,
	r *http.Request,
	oapi owner.IPortOwnerControllerOwnerApi,
	sapi store.IPortBookControllerApi,
	log *zerolog.Logger) {
	setupCORS(&w)
	for ok := true; ok; ok = false {
		log.Debug().Msg("call clientPrices")
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		carId, err := uuid.Parse(r.URL.Query().Get("CarouselId"))
		if err != nil {
			log.Err(err).Msg("Fail to parse CarouselId")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusMethodNotAllowed)
			break
		}
		priceTags, ierr := sapi.ReadPriceOptions(carId)
		if ierr != nil {
			http.Error(w, fmt.Sprintf("Fail to Read prices: %v", ierr.Error()), http.StatusInternalServerError)
			break
		}
		type PriceResp struct {
			Amount  int    `json:"Amount"`
			Tickets int    `json:"Tickets"`
			PriceId string `json:"PriceId"`
		}
		var priceResp []PriceResp
		for _, p := range priceTags {
			priceResp = append(priceResp, PriceResp{Amount: p.Amount, Tickets: p.Tickets, PriceId: p.PriceId})
		}
		writeJSON(w, priceResp)
	}
}

func severWebhookCommon(
	w http.ResponseWriter,
	r *http.Request,
	oapi owner.IPortOwnerControllerOwnerApi,
	sapi store.IPortBookControllerApi,
	log *zerolog.Logger,
	array []byte,
	whkey string) bool {

	event, err := webhook.ConstructEvent(array, r.Header.Get("Stripe-Signature"), whkey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Err(err).Msg("Router.Webhook: webhook.ConstructEvent")
		return false
	}

	log.Debug().Str("Type", event.Type).Msg("Router.Webhook:")
	if event.Type == "checkout.session.completed" {
		var sessionId string
		if id, b := event.Data.Object["id"].(string); b {
			sessionId = id
		}
		var paymentSatus string
		if ps, b := event.Data.Object["payment_status"].(string); b {
			paymentSatus = ps
		}
		var status = store.BookOrderStatusPaid
		if paymentSatus != "paid" {
			status = paymentSatus
			log.Error().Str("payment_status", paymentSatus).Msg("Router.Webhook: Has unexpected value, but want 'paid'")
		}
		// here supposed to be very small logic, that shouldn't have big impact of time of hook execution
		// however ApplyPaymentResults fuction contains remote call of carousel service
		// Better would be to put market into the book repo and leter run refill in separate thread
		// There retries could be handled as well
		if ierr := sapi.ApplyPaymentResults(sessionId, status); ierr != nil {
			log.Err(ierr).Msg("Router.Webhook: Fail to apply payment status")
		}
	}
	return true
}

func severWebhook(
	w http.ResponseWriter,
	r *http.Request,
	oapi owner.IPortOwnerControllerOwnerApi,
	sapi store.IPortBookControllerApi,
	log *zerolog.Logger,
	cfg IConfigRouter) {

	setupCORS(&w)
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Err(err).Msg("Router.Webhook:")
			break
		}
		if res := severWebhookCommon(w, r, oapi, sapi, log, b, cfg.WebhookKey()); res == false {
			break
		}
		writeJSON(w, nil)
	}
}

func severWebhookWithOwnerId(
	w http.ResponseWriter,
	r *http.Request,
	oapi owner.IPortOwnerControllerOwnerApi,
	sapi store.IPortBookControllerApi,
	log *zerolog.Logger,
	cfg IConfigRouter) {

	setupCORS(&w)
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		ownerid := r.PathValue("ownerid")
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Err(err).Msg("Router.Webhook: Fail to read event")
			break
		}

		whkey, err := oapi.ReadWhkey(ownerid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Err(err).Msg("Router.Webhook: Fail to read webhook from the repo")
			break
		}

		if res := severWebhookCommon(w, r, oapi, sapi, log, b, whkey); res == false {
			break
		}
		writeJSON(w, nil)
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		// log.Printf("io.Copy: %v", err)
		return
	}
}

func Listen(cfg IConfigRouter, oapi owner.IPortOwnerControllerOwnerApi, sapi store.IPortBookControllerApi, log *zerolog.Logger) {
	http.HandleFunc("/owner/register", func(w http.ResponseWriter, r *http.Request) { ownerRegister(w, r, oapi, sapi, log) })
	http.HandleFunc("/owner/login", func(w http.ResponseWriter, r *http.Request) { ownerLogin(w, r, oapi, sapi, log) })
	http.HandleFunc("/owner/pkey", func(w http.ResponseWriter, r *http.Request) { ownerPkey(w, r, oapi, sapi, log) })
	http.HandleFunc("/owner/skey", func(w http.ResponseWriter, r *http.Request) { ownerSkey(w, r, oapi, sapi, log) })
	http.HandleFunc("/owner/whook", func(w http.ResponseWriter, r *http.Request) { ownerWebhookRefresh(w, r, oapi, sapi, log) })
	http.HandleFunc("/carousel/add", func(w http.ResponseWriter, r *http.Request) { carouselAdd(w, r, oapi, sapi, log) })
	http.HandleFunc("/carousel/prodid", func(w http.ResponseWriter, r *http.Request) { carouselProdId(w, r, oapi, sapi, log) })
	http.HandleFunc("/client/checkout", func(w http.ResponseWriter, r *http.Request) { clientCheckout(w, r, oapi, sapi, log) })
	http.HandleFunc("/client/prices", func(w http.ResponseWriter, r *http.Request) { clientPrices(w, r, oapi, sapi, log) })
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) { severWebhook(w, r, oapi, sapi, log, cfg) })
	http.HandleFunc("/webhook/{ownerid}", func(w http.ResponseWriter, r *http.Request) { severWebhookWithOwnerId(w, r, oapi, sapi, log, cfg) })
	log.Info().Str("URL", cfg.ServerAddress()).Msg("Listening...")
	// http.ListenAndServeTLS(cfg.ServerAddress(), cfg.ServerCertPath(), cfg.ServerKeyPath(), nil)
	http.ListenAndServe(cfg.ServerAddress(), nil)
}
