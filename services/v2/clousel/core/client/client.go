package client

import (
	"clousel/lib/fault"
	"clousel/lib/pswd"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Client struct {
	// cfg    IBusinessConfigAdapter
	repoBalance  IClientRepoBalanceChangeAdapter
	repoCheckout IClientRepoCheckoutSessionAdapter
	repoGen      IClientRepoGeneralAdapter
	business     IClientBusinessAdapter
	stripe       IClientStripeAdapter
	log          *zerolog.Logger
}

func ClientCreate(
	repoBalance IClientRepoBalanceChangeAdapter,
	repoCheckout IClientRepoCheckoutSessionAdapter,
	repoGen IClientRepoGeneralAdapter,
	business IClientBusinessAdapter,
	stripe IClientStripeAdapter,
	log *zerolog.Logger,
) *Client {
	return &Client{
		repoBalance:  repoBalance,
		repoCheckout: repoCheckout,
		repoGen:      repoGen,
		business:     business,
		stripe:       stripe,
		log:          log,
	}
}

/*
IMachineRestController
*/
func (c *Client) Register(username string, email string, password string, companyName string) (err fault.IError) {
	const fn = "Core.Client.Register"
	userId := uuid.New()
	psswd := pswd.PasswordPlainCreate(password)

	for ok := true; ok; ok = false {
		if exists, e := c.business.IsCompanyExists(companyName); err != nil || !exists {
			c.log.Err(e).Msgf("%s: Comapny name '%s' is not exists", fn, companyName)
			err = e
			break
		}
		if err = c.repoGen.SaveNewClientEntry(userId, username, email, psswd.Hash().Encode().Str(), companyName); err != nil {
			c.log.Err(err).Msgf("%s: Fail to register company %s", fn, username)
			break
		}
		c.log.Info().Msgf("%s: Success, company %s has been registred", fn, username)
	}
	return err
}

func (c *Client) Login(username string, password string) (*ClientEntry, fault.IError) {
	const fn = "Core.Client.Login"
	var err fault.IError
	var entry *ClientEntry
	entry, err = c.repoGen.ReadClientEntryByName(username)
	pknown := pswd.PasswordHasedBase64Create(entry.Password)
	if !pswd.PasswordPlainCreate(password).Hash().Encode().Equal(pknown) {
		err = fault.New(EClientPasswordMismatch).Msg("Wrong password")
		entry = nil
		c.log.Err(err).Msgf("%s: Fail %s to login", fn, username)
	} else {
		c.log.Info().Msgf("%s: Success %s logged in", fn, username)
	}
	return entry, err
}

func (c *Client) BuyTickets(userId uuid.UUID, priceId string, afterSellVisitUrl string) (session ISession, err fault.IError) {
	const fn = "Core.Client.BuyTickets"
	for ok := true; ok; ok = false {
		if len(priceId) == 0 {
			err = fault.New(EClientInvalidValue).Msgf("%s: priceId is empty", fn)
			break
		}
		if len(afterSellVisitUrl) == 0 {
			err = fault.New(EClientInvalidValue).Msgf("%s: afterSellVisitUrl is empty", fn)
			break
		}
		var entry *ClientEntry
		if entry, err = c.repoGen.ReadClientEntryById(userId); err != nil {
			c.log.Err(err).Str("UserId", userId.String()).Msgf("%s: Failed to read client entry", fn)
			break
		}
		var skey string
		if skey, _, err = c.business.ClientReadKeys(entry.CompanyName); err != nil {
			c.log.Err(err).Str("UserId", userId.String()).Msgf("%s: Failed to read keys", fn)
			break
		}
		var pt PriceTag
		if pt, err = c.stripe.ReadPriceDetails(skey, priceId); err != nil {
			c.log.Err(err).Str("UserId", userId.String()).Msgf("%s: Failed to read price details", fn)
			break
		}
		query_success := url.Values{}
		query_success.Add("type", "popup_success")
		query_success.Add("msg", "Payment has been confirmed")
		query_error := url.Values{}
		query_error.Add("type", "popup_error")
		query_error.Add("msg", "Something went wrong")
		purls := PaymentResultUrls{
			Success: fmt.Sprintf("%s?%s", afterSellVisitUrl, query_success.Encode()),
			Cancel:  fmt.Sprintf("%s?%s", afterSellVisitUrl, query_error.Encode()),
		}
		if session, err = c.stripe.GenCheckoutSessionUrl(entry.Email, skey, priceId, purls); err != nil {
			c.log.Err(err).Str("UserId", userId.String()).Msgf("%s: Failed to generate checkout session", fn)
			break
		}
		eventId := uuid.New()
		if err = c.repoCheckout.SaveNewCheckoutEntry(eventId, userId, session.Id(), pt.Amount, pt.Tickets); err != nil {
			c.log.Err(err).Str("UserId", userId.String()).Msgf("%s: Failed to save checkout entry", fn)
			break
		}
		c.log.Info().Str("UserId", userId.String()).Msgf("%s: Successfully generated checkout session", fn)
	}
	return session, err
}

func (c *Client) ApplyPaymentResults(sessionId string, status PaymentStatus) (err fault.IError) {
	const fn = "Core.Client.ApplyPaymentResults"

	for ok := true; ok; ok = false {
		var ce *CheckoutEntry
		if ce, err = c.repoCheckout.ReadCheckoutEntriesBySessionId(sessionId); err != nil {
			c.log.Err(err).Msgf("%s: Fail to read checkout entry by session '%s'", fn, sessionId)
			break
		}
		if ce.SessionId == PaymentStatusPaid {
			c.log.Warn().Msgf("%s: Session '%s' already has statsu '%s'", fn, sessionId, PaymentStatusPaid)
			break
		}
		if err = c.repoCheckout.UpdateCheckoutStatus(sessionId, status); err != nil {
			c.log.Err(err).Msgf("%s: Fail to update checkout status of session '%s", fn, sessionId)
			break
		}
		if status != PaymentStatusPaid {
			break
		}
		if err = c.repoBalance.SaveNewBalanceChangeEntry(ce.EventId, ce.UserId, ce.Tickets); err != nil {
			break
		}
		c.log.Info().Str("SessionId", sessionId).Str("Status", status).Msgf("%s: Result has been updated sucessffuly", fn)
	}

	return err
}

func (c *Client) ReadPriceOptions(userId uuid.UUID) (tags []PriceTag, err fault.IError) {
	const fn = "Core.Client.ReadPriceOptions"
	for ok := true; ok; ok = false {
		var entry *ClientEntry
		if entry, err = c.repoGen.ReadClientEntryById(userId); err != nil {
			c.log.Err(err).Str("UserId", userId.String()).Msgf("%s: Failed to read client entry", fn)
			break
		}
		var skey, prodId string
		if skey, prodId, err = c.business.ClientReadKeys(entry.CompanyName); err != nil {
			c.log.Err(err).Str("CompanyName", entry.CompanyName).Msgf("%s: Failed to read keys", fn)
			break
		}
		if tags, err = c.stripe.ReadPriceListByProdId(skey, prodId, 10); err != nil {
			c.log.Err(err).Str("CompanyName", entry.CompanyName).Msgf("%s: Failed to price list", fn)
			break
		}
		c.log.Debug().Msgf("%s: %s Successfully read prices %d entries ", fn, userId.String(), len(tags))
	}
	return tags, err
}

func (c *Client) ReadBalance(userId uuid.UUID) (balance int, err fault.IError) {
	const fn = "Core.Client.ReadBalance"
	var entries []*TicketsBalanceEntry
	if entries, err = c.repoBalance.ReadBalanceEntriesByUserId(userId); err == nil {
		for _, e := range entries {
			balance += e.Change
		}
	}
	return balance, err
}

/*
IMachineClientAdapter
*/
func (c *Client) ApplyGameCost(eventId uuid.UUID, userId uuid.UUID, tickets int) (err fault.IError) {
	const fn = "Core.Client.ApplyGameCost"
	if err = c.repoBalance.SaveNewBalanceChangeEntry(eventId, userId, -tickets); err != nil {
		return err
	}
	var balance int
	var entries []*TicketsBalanceEntry
	if entries, err = c.repoBalance.ReadBalanceEntriesByUserId(userId); err == nil {
		for _, e := range entries {
			balance += e.Change
		}
	}
	if balance == 0 {
		for _, e := range entries {
			if err = c.repoBalance.RemoveBalanceByEventId(e.EventId); err != nil {
				c.log.Err(err).Str("EventId", e.EventId.String()).Msgf("%s: Fail to remove balance entry %s", fn, e.EventId.String())
			}
		}
	}
	return err
}

func (c *Client) IsCanPay(userId uuid.UUID, cost int) (can bool, err fault.IError) {
	const fn = "Core.Client.IsCanPay"
	var balance int
	can = false
	if balance, err = c.ReadBalance(userId); err == nil {
		if balance >= cost {
			can = true
		}
	}
	return can, err
}
