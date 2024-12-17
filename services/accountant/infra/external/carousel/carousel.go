package external

import (
	"accountant/core/owner"
	erro "accountant/core/owner/error"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"accountant/core/store"
	errs "accountant/core/store/error"

	"github.com/rs/zerolog"
)

// Code for interation with the carousel service
// - Register carousel
// - Refill

type CarouselGateway struct {
	cfg IConfigCarousel
	log *zerolog.Logger
}

func CarouselGatewayCreate(cfg IConfigCarousel) *CarouselGateway {
	return &CarouselGateway{cfg: cfg}
}

func (cg *CarouselGateway) Refill(carId store.Carousel, tickets int) errs.IError {
	var ierr errs.IError
	form := url.Values{}
	form.Add("CarouselId", carId.String())
	form.Add("Tickets", fmt.Sprint(tickets))
	req, err := http.NewRequest("POST", cg.cfg.ExternalServiceCarouselRefillUrl(), strings.NewReader(form.Encode()))
	if err != nil {
		cg.log.Err(err).Msg("Faul to prepare register")
		ierr = errs.New(errs.ECRemoteServiceCarouselRefill).Msgf("Fail to rpc:refill, error:%v", err)
		return ierr
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	hc := http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		ierr = errs.New(errs.ECRemoteServiceCarouselRefill).Msgf("Fail to rpc:refill, error:%v", err)
		return ierr
	}
	if resp.StatusCode != http.StatusOK {
		ierr = errs.New(errs.ECRemoteServiceCarouselRefill).Msgf("Fail to rpc:refill, remote status code:%d", resp.StatusCode)
		return ierr
	}

	return nil
}
func (cg *CarouselGateway) Register(ownerId owner.Owner, carId owner.Carousel) erro.IError {
	var ierr erro.IError
	form := url.Values{}
	form.Add("CarouselId", carId.String())
	form.Add("OwnerId", ownerId.String())
	req, err := http.NewRequest("POST", cg.cfg.ExternalServiceCarouselRegisterUrl(), strings.NewReader(form.Encode()))
	if err != nil {
		cg.log.Err(err).Msg("Faul to prepare register")
		ierr = erro.New(erro.ECRemoteServiceCarouselRegister).Msgf("Fail to rpc:register, error:%v", err)
		return ierr
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	hc := http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		ierr = erro.New(erro.ECRemoteServiceCarouselRegister).Msgf("Fail to rpc:register, error:%v", err)
		return ierr
	}
	if resp.StatusCode != http.StatusOK {
		ierr = erro.New(erro.ECRemoteServiceCarouselRegister).Msgf("Fail to rpc:register, remote status code:%d", resp.StatusCode)
		return ierr
	}

	return nil
}
