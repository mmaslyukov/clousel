package business

import (
	"clousel/lib/fault"
	"clousel/lib/pswd"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Business struct {
	cfg    IBusinessConfigAdapter
	repo   IBusinessRepoAdapter
	stripe IBusinessStripeAdapter
	log    *zerolog.Logger
}

func BusinessCreate(
	cfg IBusinessConfigAdapter,
	repo IBusinessRepoAdapter,
	stripe IBusinessStripeAdapter,
	log *zerolog.Logger,
) *Business {
	return &Business{
		cfg:    cfg,
		repo:   repo,
		stripe: stripe,
		log:    log,
	}
}

func (b *Business) Register(companyName string, email string, password string) fault.IError {
	const fn = "Core.Business.Register"
	companyId := uuid.New()
	psswd := pswd.PasswordPlainCreate(password)
	err := b.repo.SaveNewBusinessEntry(companyId, companyName, email, psswd.Hash().Encode().Str())
	if err == nil {
		b.log.Info().Msgf("%s: Success, company %s has been registred", fn, companyName)
	} else {
		b.log.Err(err).Msgf("%s: Fail to register company %s", fn, companyName)
	}
	return err
}

func (b *Business) Login(companyName string, password string) (entry *BusinessEntry, err fault.IError) {
	const fn = "Core.Business.Login"
	for ok := true; ok; ok = false {

		if len(companyName) == 0 || len(password) == 0 {
			err = fault.New(EBusinessInvalidArgument).Msg("Either companyName or password is empty")
			break
		}

		if entry, err = b.repo.ReadBusinessEntryByName(companyName); err != nil {
			break
		}
		pknown := pswd.PasswordHasedBase64Create(entry.Password)
		if !pswd.PasswordPlainCreate(password).Hash().Encode().Equal(pknown) {
			err = fault.New(EBusinessPasswordMismatch).Msg("Wrong password")
			entry = nil
			b.log.Err(err).Msgf("%s: Fail %s to login", fn, companyName)
			break
		}
		b.log.Info().Msgf("%s: Success %s logged in", fn, companyName)
	}

	return entry, err
}

func (b *Business) AssignKeys(companyId uuid.UUID, skey string, prodId string) fault.IError {
	const fn = "Core.Business.AssignKeys"
	var err fault.IError
	var whid, whkey string
	var entry *BusinessEntry
	url := b.cfg.WebhookUrl(companyId.String())

	if entry, err = b.repo.ReadBusinessEntryById(companyId); err != nil {
		return err
	}
	if entry.Whid == nil {
		whid, whkey, err = b.stripe.WebhookRegister(url, skey)
	} else {
		err = b.stripe.WebhookUpdateUrl(url, skey, *entry.Whid)
		whid = *entry.Whid
		whkey = *entry.Whkey
	}
	if len(skey) > 0 &&
		len(whkey) > 0 &&
		len(whid) > 0 &&
		len(prodId) > 0 {
		err = b.repo.AssignKeys(companyId, skey, prodId, whid, whkey)
	} else {
		err = fault.New(EBusinessInvalidArgument).Msgf("%s: Empty key value %d %d %d %d", fn, len(skey), len(whkey), len(whid), len(prodId))
	}

	if err == nil {
		b.log.Info().Msgf("%s: Success, Keys have been assigned for %s", fn, companyId.String())
	} else {
		b.log.Err(err).Msgf("%s: Failure, Keys have not been assigned for %s", fn, companyId.String())
	}
	return err
}

func (b *Business) ReadWhkey(companyId uuid.UUID) (whkey string, err fault.IError) {
	const fn = "Core.Business.ReadWhkey"
	var entry *BusinessEntry
	if entry, err = b.repo.ReadBusinessEntryById(companyId); err != nil {
		if entry.Whkey != nil {
			whkey = *entry.Whkey
		} else {
			err = fault.New(EBusinessNilPointer).Msgf("%s: Whkey pointer is nil", fn)
		}
	}
	return whkey, err
}

/*
IClientBusinessAdapter
*/
func (b *Business) ClientReadKeys(companyName string) (skey string, prodId string, err fault.IError) {
	var entry *BusinessEntry
	if entry, err = b.repo.ReadBusinessEntryByName(companyName); err == nil {
		skey = *entry.Skey
		prodId = *entry.ProdId
	}
	return skey, prodId, err
}

func (b *Business) IsCompanyExists(companyName string) (exists bool, err fault.IError) {
	var entry *BusinessEntry = nil
	entry, err = b.repo.ReadBusinessEntryByName(companyName)
	exists = (err == nil) && (entry != nil)
	return exists, err
}

// func (b *Business) ClientReadId(companyName string) (id uuid.UUID, err fault.IError) {
// 	var entry *BusinessEntry
// 	if entry, err = b.repo.ReadBusinessEntryByName(companyName); err == nil {
// 		id = entry.CompanyId
// 	}
// 	return id, err
// }
