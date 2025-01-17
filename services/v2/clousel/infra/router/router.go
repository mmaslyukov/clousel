package router

import (
	"bytes"
	"clousel/core/business"
	"clousel/core/client"
	"clousel/core/machine"
	"clousel/lib/fault"
	"encoding/json"
	"io"
	"net/http"
	"net/mail"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stripe/stripe-go/v72/webhook"
)

type ControllerApi struct {
	client      client.IClientRestController
	business    business.IBusinessRestController
	machine     machine.IMachineRestController
	tockenStore TockenStore
}
type MachineSelector struct {
	companyId *uuid.UUID
	machId    *uuid.UUID
	status    *machine.MachineStatus
	from      *time.Time
	till      *time.Time
}

func (m *MachineSelector) CompanyId() *uuid.UUID {
	return m.companyId
}
func (m *MachineSelector) MachId() *uuid.UUID {
	return m.machId
}
func (m *MachineSelector) Status() *machine.MachineStatus {
	return m.status
}
func (m *MachineSelector) TimeFrom() *time.Time {
	return m.from
}
func (m *MachineSelector) TimeTill() *time.Time {
	return m.till
}

func Listen(
	cfg IConfigRouter,
	uc client.IClientRestController,
	bc business.IBusinessRestController,
	mc machine.IMachineRestController,
	log *zerolog.Logger) {
	api := &ControllerApi{client: uc, business: bc, machine: mc, tockenStore: TockenStoreCreate()}
	http.HandleFunc("/business/register", func(w http.ResponseWriter, r *http.Request) { businessRegister(w, r, api, log) })
	http.HandleFunc("/business/login", func(w http.ResponseWriter, r *http.Request) { businessLogin(w, r, api, log) })
	http.HandleFunc("/business/skey", func(w http.ResponseWriter, r *http.Request) { businessAssignSecKey(w, r, api, log) })
	http.HandleFunc("/machine/add", func(w http.ResponseWriter, r *http.Request) { machineAdd(w, r, api, log) })
	http.HandleFunc("/machine/get", func(w http.ResponseWriter, r *http.Request) { machineGet(w, r, api, log) })
	http.HandleFunc("/machine/update", func(w http.ResponseWriter, r *http.Request) { machineUpdate(w, r, api, log) })
	http.HandleFunc("/machine/play", func(w http.ResponseWriter, r *http.Request) { machinePlay(w, r, api, log) })
	http.HandleFunc("/machine/poll", func(w http.ResponseWriter, r *http.Request) { machinePoll(w, r, api, log) })
	http.HandleFunc("/client/register", func(w http.ResponseWriter, r *http.Request) { clientRegister(w, r, api, log) })
	http.HandleFunc("/client/login", func(w http.ResponseWriter, r *http.Request) { clientLogin(w, r, api, log) })
	http.HandleFunc("/client/balance", func(w http.ResponseWriter, r *http.Request) { clientReadBalance(w, r, api, log) })
	http.HandleFunc("/client/buy", func(w http.ResponseWriter, r *http.Request) { clientCheckout(w, r, api, log) })
	http.HandleFunc("/client/price", func(w http.ResponseWriter, r *http.Request) { clientReadPrice(w, r, api, log) })
	http.HandleFunc("/webhook/{companyId}", func(w http.ResponseWriter, r *http.Request) { severWebhookWithCompanyId(w, r, api, log) })
	http.HandleFunc("/webhook/dev", func(w http.ResponseWriter, r *http.Request) { severWebhookDev(w, r, api, log, cfg) })
	log.Info().Str("URL", cfg.ServerAddress()).Msg("Listening...")
	http.ListenAndServe(cfg.ServerAddress(), nil)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		return
	}
}
func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Tocken, Authorization")
}

func businessRegister(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.businessRegister"
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		companyName := r.PostFormValue("CompanyName")
		email := r.PostFormValue("Email")
		password := r.PostFormValue("Password")
		if _, err := mail.ParseAddress(email); err != nil || len(password) < 4 {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			break
		}

		if err := api.business.Register(companyName, email, password); err != nil {
			log.Err(err).Msgf("%s: Fail to register new Business with, username:'%s', email:'%s'", fn,
				companyName, email)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		log.Info().Msgf("%s: Registered new business with, username:'%s', email:'%s'", fn,
			companyName, email)
		w.WriteHeader(http.StatusOK)
	}
}

func businessLogin(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.businessLogin"
	for ok := true; ok; ok = false {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		var err fault.IError
		var entry *business.BusinessEntry
		companyName := r.URL.Query().Get("CompanyName")
		password := r.URL.Query().Get("Password")

		if entry, err = api.business.Login(companyName, password); err != nil {
			log.Err(err).Msgf("%s: Fail to login for business with, companyName:'%s'", fn, companyName)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		t := TockenCreate(entry, TockenRoleBusiness)
		api.tockenStore.Add(t)
		log.Info().Msgf("%s: Logged in for business with, companyName:'%s'", fn, companyName)
		writeJSON(w, struct {
			Tocken string `json:"Tocken"`
		}{
			Tocken: t.Base64().Str(),
		})
	}
}

func businessAssignSecKey(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.businessAssignSecKey"
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		var tockenProvided ITockenBase64
		var err fault.IError
		log.Debug().Str("TockenQuery", r.PostFormValue("Tocken")).Send()
		if tockenProvided, err = TockenCreateFromBase64String(r.PostFormValue("Tocken")); err != nil {
			log.Error().Msgf("%s: Fail to parse tocken from query string, error:%s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		var tocken ITocken
		if tocken, err = api.tockenStore.Find(tockenProvided); err != nil {
			log.Err(err).Msgf("%s: Fail to find tocken with id:'%s'", fn, tockenProvided.Str())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		if tocken.Role() == TockenRoleClient {
			err = fault.New(ERouterNotAllowed).Msgf("No allowed to make changes")
			log.Err(err).Msgf("%s", fn)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		skey := r.PostFormValue("SKey")
		prodid := r.PostFormValue("ProdId")
		if len(skey) == 0 || len(prodid) == 0 {
			err = fault.New(ERouterNotFound).Msg("%s: Fail to parse Skey or ProdId")
			log.Err(err).Send()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		if err = api.business.AssignKeys(tocken.Auth().Id(), skey, prodid); err != nil {
			log.Err(err).Msgf("%s: Fail to assign secret key, Id:'%s", fn, tocken.Auth().Id().String())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		log.Info().Msgf("%s: Secret key for Id:'%s' has been assigned sucessfully", fn, tocken.Auth().Id().String())
		w.WriteHeader(http.StatusOK)
	}
}

func machineAdd(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.machineAdd"
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		var tockenProvided ITockenBase64
		var err fault.IError
		log.Debug().Str("TockenQuery", r.PostFormValue("Tocken")).Send()
		if tockenProvided, err = TockenCreateFromBase64String(r.PostFormValue("Tocken")); err != nil {
			log.Error().Msgf("%s: Fail to parse tocken from query string, error:%s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		var tocken ITocken
		if tocken, err = api.tockenStore.Find(tockenProvided); err != nil {
			log.Err(err).Msgf("%s: Fail to find tocken with id:'%s'", fn, tockenProvided.Str())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}

		var cost int
		if c, e := strconv.Atoi(r.PostFormValue("Cost")); e != nil {
			log.Err(e).Msgf("%s: Fail to parse game cost", fn)
			http.Error(w, e.Error(), http.StatusBadRequest)
			break
		} else {
			cost = c
		}

		var machid uuid.UUID
		if id, e := uuid.Parse(r.PostFormValue("MachId")); e != nil {
			log.Err(e).Msgf("%s: Fail to parse machine id", fn)
			http.Error(w, e.Error(), http.StatusBadRequest)
			break
		} else {
			machid = id
		}

		if err = api.machine.SaveNewMachineEntry(machid, tocken.Auth().Id(), cost); err != nil {
			log.Err(err).Msgf("%s: Fail to save new machine id:'%s", fn, machid.String())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		log.Info().Msgf("%s: Machine id:'%s' has been saved sucessfully", fn, machid.String())
		w.WriteHeader(http.StatusOK)
	}
}

func machineGet(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.machineGet"
	for ok := true; ok; ok = false {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}

		var err fault.IError
		var tockenProvided ITockenBase64
		if tockenProvided, err = TockenCreateFromBase64String(r.URL.Query().Get("Tocken")); err != nil {
			log.Error().Msgf("%s: Fail to parse tocken from query string, error:%s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		var tocken ITocken
		if tocken, err = api.tockenStore.Find(tockenProvided); err != nil {
			log.Err(err).Msgf("%s: Fail to find tocken with id:'%s'", fn, tockenProvided.Str())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		var machid *uuid.UUID = nil
		if queryMachId := r.URL.Query().Get("MachId"); len(queryMachId) != 0 {
			if id, e := uuid.Parse(r.URL.Query().Get("MachId")); e == nil {
				machid = &id
			} else {
				log.Err(e).Msgf("%s: Fail to parse machine id", fn)
			}
		}
		var status *machine.MachineStatus = nil
		if st := r.URL.Query().Get("Status"); len(st) != 0 {
			status = &st
		}
		// companyId, _ := uuid.Parse("a055963c-d4d8-4be3-8fac-f7a1d7cf1d59")
		companyId := tocken.Auth().Id()
		selector := &MachineSelector{
			companyId: &companyId,
			status:    status,
			machId:    machid,
		}
		var entries []*machine.MachineEntry
		if entries, err = api.machine.ReadMachineEntriesBySelector(selector); err != nil {
			log.Err(err).Msgf("%s: Fail to read machine entries", fn)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		log.Info().Msgf("%s: Successfully Prepared list of machines", fn)
		writeJSON(w, entries)
	}
}

func machineUpdate(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.machineUpdate"
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		var tockenProvided ITockenBase64
		var err fault.IError
		log.Debug().Str("TockenQuery", r.PostFormValue("Tocken")).Send()
		if tockenProvided, err = TockenCreateFromBase64String(r.PostFormValue("Tocken")); err != nil {
			log.Error().Msgf("%s: Fail to parse tocken from query string, error:%s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		var tocken ITocken
		if tocken, err = api.tockenStore.Find(tockenProvided); err != nil {
			log.Err(err).Msgf("%s: Fail to find tocken with id:'%s'", fn, tockenProvided.Str())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		if tocken.Role() == TockenRoleClient {
			err = fault.New(ERouterNotAllowed).Msgf("No allowed to make changes")
			log.Err(err).Msgf("%s", fn)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}

		var cost int
		if c, e := strconv.Atoi(r.PostFormValue("Cost")); e != nil {
			log.Err(e).Msgf("%s: Fail to parse game cost", fn)
			http.Error(w, e.Error(), http.StatusBadRequest)
			break
		} else {
			cost = c
		}

		var machid uuid.UUID
		if id, e := uuid.Parse(r.PostFormValue("MachId")); e != nil {
			log.Err(e).Msgf("%s: Fail to parse machine id", fn)
			http.Error(w, e.Error(), http.StatusBadRequest)
			break
		} else {
			machid = id
		}

		if err = api.machine.ChangeGameCost(machid, cost); err != nil {
			log.Err(err).Msgf("%s: Fail to change game cost for the machine id:'%s", fn, machid.String())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		log.Info().Msgf("%s: Machine id:'%s' has been updated sucessfully", fn, machid.String())
		w.WriteHeader(http.StatusOK)
	}
}

func machinePlay(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.machinePlay"
	for ok := true; ok; ok = false {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}

		var err fault.IError
		var tockenProvided ITockenBase64
		if tockenProvided, err = TockenCreateFromBase64String(r.URL.Query().Get("Tocken")); err != nil {
			log.Error().Msgf("%s: Fail to parse tocken from query string, error:%s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		var tocken ITocken
		if tocken, err = api.tockenStore.Find(tockenProvided); err != nil {
			log.Err(err).Msgf("%s: Fail to find tocken with id:'%s'", fn, tockenProvided.Str())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		var machid uuid.UUID
		if queryMachId := r.URL.Query().Get("MachId"); len(queryMachId) != 0 {
			var e error
			if machid, e = uuid.Parse(queryMachId); e != nil {
				log.Err(e).Msgf("%s: Fail to parse machine id", fn)
				http.Error(w, e.Error(), http.StatusBadRequest)
				break
			}
		}
		var eventId *uuid.UUID
		if eventId, err = api.machine.PlayRequest(machid, tocken.Auth().Id()); err != nil {
			log.Err(err).Msgf("%s: Fail to retrieve data", fn)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		if eventId == nil {
			err = fault.New(ERouterInvalidValue).Msgf("EventId is nil")
			log.Error().Msgf("%s: %s", fn, err.Full())
			http.Error(w, err.Full(), http.StatusInternalServerError)
			break
		}
		log.Info().Msgf("%s: Successfully performed play command", fn)
		writeJSON(w, struct {
			EventId string `json:"EventId"`
		}{
			EventId: eventId.String(),
		})
	}
}

func machinePoll(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.machinePoll"
	for ok := true; ok; ok = false {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}

		var err fault.IError
		var tockenProvided ITockenBase64
		if tockenProvided, err = TockenCreateFromBase64String(r.URL.Query().Get("Tocken")); err != nil {
			log.Error().Msgf("%s: Fail to parse tocken from query string, error:%s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		if _, err = api.tockenStore.Find(tockenProvided); err != nil {
			log.Err(err).Msgf("%s: Fail to find tocken with id:'%s'", fn, tockenProvided.Str())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		var eventid uuid.UUID
		if queryEventId := r.URL.Query().Get("EventId"); len(queryEventId) != 0 {
			var e error
			if eventid, e = uuid.Parse(queryEventId); e != nil {
				log.Err(e).Msgf("%s: Fail to parse event id", fn)
				http.Error(w, e.Error(), http.StatusBadRequest)
				break
			}
		}
		var status machine.GameStatus
		if status, err = api.machine.PollRequestStatus(eventid); err != nil {
			log.Err(err).Msgf("%s: Fail to read game status", fn)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		writeJSON(w, struct {
			EventId string `json:"EventId"`
			Status  string `json:"Status"`
		}{
			EventId: eventid.String(),
			Status:  status,
		})
	}
}

func clientRegister(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.clientRegister"
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		username := r.PostFormValue("Username")
		comapnyName := r.PostFormValue("Company")
		email := r.PostFormValue("Email")
		password := r.PostFormValue("Password")
		if _, err := mail.ParseAddress(email); err != nil || len(password) < 4 {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			break
		}

		if err := api.client.Register(username, email, password, comapnyName); err != nil {
			log.Err(err).Msgf("%s: Fail to register new client with, username:'%s', email:'%s'", fn,
				username, email)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		log.Info().Msgf("%s: Registered new client with, username:'%s', email:'%s'", fn,
			username, email)
		w.WriteHeader(http.StatusOK)
	}
}

func clientLogin(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.clientLogin"
	for ok := true; ok; ok = false {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		var err fault.IError
		var entry *client.ClientEntry
		username := r.URL.Query().Get("Username")
		password := r.URL.Query().Get("Password")

		if entry, err = api.client.Login(username, password); err != nil {
			log.Err(err).Msgf("%s: Fail to login for client with, username:'%s'", fn, username)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		t := TockenCreate(entry, TockenRoleClient)
		api.tockenStore.Add(t)
		log.Info().Msgf("%s: Logged in for client with, username:'%s'", fn, username)
		writeJSON(w, struct {
			Tocken string `json:"Tocken"`
		}{
			Tocken: t.Base64().Str(),
		})
	}
}

func clientReadBalance(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.clientReadBalance"
	for ok := true; ok; ok = false {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}

		var err fault.IError
		var tockenProvided ITockenBase64
		if tockenProvided, err = TockenCreateFromBase64String(r.URL.Query().Get("Tocken")); err != nil {
			log.Error().Msgf("%s: Fail to parse tocken from query string, error:%s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		var tocken ITocken
		if tocken, err = api.tockenStore.Find(tockenProvided); err != nil {
			log.Err(err).Msgf("%s: Fail to find tocken with id:'%s'", fn, tockenProvided.Str())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		if tocken.Role() != TockenRoleClient {
			err = fault.New(ERouterInvalidValue).Msgf("%s, Unexpected tocken role:'%s'", fn, tocken.Role())
			log.Error().Msgf("%s: %s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		var balance int
		if balance, err = api.client.ReadBalance(tocken.Auth().Id()); err != nil {
			log.Err(err).Msgf("%s: Fail to read balance", fn)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		writeJSON(w, struct {
			Balance int `json:"Balance"`
		}{
			Balance: balance,
		})
	}
}

func clientCheckout(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.clientCheckout"
	for ok := true; ok; ok = false {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}

		var err fault.IError
		var tockenProvided ITockenBase64
		if tockenProvided, err = TockenCreateFromBase64String(r.URL.Query().Get("Tocken")); err != nil {
			log.Error().Msgf("%s: Fail to parse tocken from query string, error:%s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		var tocken ITocken
		if tocken, err = api.tockenStore.Find(tockenProvided); err != nil {
			log.Err(err).Msgf("%s: Fail to find tocken with id:'%s'", fn, tockenProvided.Str())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		if tocken.Role() != TockenRoleClient {
			err = fault.New(ERouterInvalidValue).Msgf("%s: Unexpected tocken role:'%s'", fn, tocken.Role())
			log.Error().Msgf("%s: %s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		var cs client.ISession
		home := r.URL.Query().Get("Home")
		priceId := r.URL.Query().Get("PriceId")
		if cs, err = api.client.BuyTickets(tocken.Auth().Id(), priceId, home); err != nil {
			log.Error().Msgf("%s: Fail to execue buy tickets function %s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		http.Redirect(w, r, cs.Url(), http.StatusOK)
	}
}

func clientReadPrice(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	setupCORS(&w)
	const fn = "Router.clientReadPrice"
	for ok := true; ok; ok = false {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		var err fault.IError
		var tockenProvided ITockenBase64
		if tockenProvided, err = TockenCreateFromBase64String(r.URL.Query().Get("Tocken")); err != nil {
			log.Error().Msgf("%s: Fail to parse tocken from query string, error:%s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		var tocken ITocken
		if tocken, err = api.tockenStore.Find(tockenProvided); err != nil {
			log.Err(err).Msgf("%s: Fail to find tocken with id:'%s'", fn, tockenProvided.Str())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			break
		}
		if tocken.Role() != TockenRoleClient {
			err = fault.New(ERouterInvalidValue).Msgf("%s: Unexpected tocken role:'%s'", fn, tocken.Role())
			log.Error().Msgf("%s: %s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		var pts []client.PriceTag
		if pts, err = api.client.ReadPriceOptions(tocken.Auth().Id()); err != nil {
			log.Error().Msgf("%s: %s", fn, err.Full())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		writeJSON(w, pts)
	}
}

func severWebhookDev(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger, cfg IConfigRouter) {
	const fn = "Router.severWebhookDev"

	setupCORS(&w)
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Err(err).Msgf("%s:", fn)
			break
		}
		if res := severWebhookCommon(w, r, api, log, b, cfg.WebhookKey()); res == false {
			break
		}
		writeJSON(w, nil)
	}
}

func severWebhookWithCompanyId(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger) {
	const fn = "Router.severWebhookWithCompanyId"

	setupCORS(&w)
	for ok := true; ok; ok = false {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Err(err).Msgf("%s: Fail to read event", fn)
			break
		}

		var companyId uuid.UUID
		if id, e := uuid.Parse(r.PathValue("companyId")); e != nil {
			log.Err(e).Msgf("%s: Fail to parse company id", fn)
			http.Error(w, e.Error(), http.StatusBadRequest)
			break
		} else {
			companyId = id
		}

		whkey, err := api.business.ReadWhkey(companyId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Err(err).Msgf("%s: Fail to read webhook from the repo", fn)
			break
		}

		if res := severWebhookCommon(w, r, api, log, b, whkey); res == false {
			break
		}
		writeJSON(w, nil)
	}
}

func severWebhookCommon(w http.ResponseWriter, r *http.Request, api *ControllerApi, log *zerolog.Logger, array []byte, whkey string) bool {
	const fn = "Router.severWebhookCommon"
	for ok := true; ok; ok = false {

		event, err := webhook.ConstructEvent(array, r.Header.Get("Stripe-Signature"), whkey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Err(err).Msgf("%s: Fail to construct event", fn)
			return false
		}

		log.Debug().Str("Type", event.Type).Msgf("%s:", fn)
		if event.Type == "checkout.session.completed" {
			var sessionId string
			if id, b := event.Data.Object["id"].(string); b {
				sessionId = id
			}
			var paymentSatus string
			if ps, b := event.Data.Object["payment_status"].(string); b {
				paymentSatus = ps
			}
			var status = client.PaymentStatusPaid
			if paymentSatus != "paid" {
				status = paymentSatus
				log.Error().Str("payment_status", paymentSatus).Msgf("%s: Has unexpected value, but want 'paid'", fn)
			}
			// here supposed to be very small logic, that shouldn't have big impact of time of hook execution
			if ierr := api.client.ApplyPaymentResults(sessionId, status); ierr != nil {
				log.Err(ierr).Msgf("%s: Fail to apply payment status", fn)
			}
		}
	}
	return true
}
