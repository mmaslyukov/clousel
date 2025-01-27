package owner

import (
	erro "accountant/core/owner/error"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type OwnerDomain struct {
	carousel IPortStoreAdapterCarouselService
	prodRepo IPortOwnerAdapterProductRepo
	profRepo IPortOwnerAdapterProfileRepo
	stripe   IPortOwnerAdapterStripeService
	cfg      IPortOwnerAdapterProfileConfig
	log      *zerolog.Logger
	tokens   map[Token]*TokenDetails
}

func OwnerDomainCreate(
	carousel IPortStoreAdapterCarouselService,
	prodRepo IPortOwnerAdapterProductRepo,
	profRepo IPortOwnerAdapterProfileRepo,
	stripe IPortOwnerAdapterStripeService,
	cfg IPortOwnerAdapterProfileConfig,
	log *zerolog.Logger) *OwnerDomain {
	return &OwnerDomain{
		carousel: carousel,
		prodRepo: prodRepo,
		profRepo: profRepo,
		stripe:   stripe,
		cfg:      cfg,
		log:      log,
		tokens:   make(map[Token]*TokenDetails)}
}

func (od *OwnerDomain) Register(email string, password string) IError {
	var err IError
	pp := PasswordPlainCreate(password)
	// pp := PasswordPlain{}
	// pp.data = password
	// pp.Hash()
	// pp.Hash()
	err = od.profRepo.OwnerRegister(email, pp.Hash().Base64().Str(), UserRoleRegular)
	if err == nil {
		od.log.Info().Msgf("Owner.Register: User %s has been registered", email)
	}
	return err
}

func (od *OwnerDomain) Login(email string, password string) (Token, IError) {
	var err IError
	var token Token
	var entry OwnerEntry
	for ok := true; ok; ok = false {
		if entry, err = od.profRepo.OwnerReadEntry(email); err != nil {
			err = erro.New(erro.ECUserNotFound).Msgf("User '%s' is not found", email)
			break
		}
		passwordProvided := PasswordPlainCreate(password)
		passwordRegistered := PasswordHashedCreate(entry.Password)
		if !passwordProvided.Hash().Base64().Eq(passwordRegistered) {
			err = erro.New(erro.ECUserPasswordMissmatch).Msgf("Password is invalid")
			break
		}
		token = uuid.New()
		od.tokens[token] = &TokenDetails{ownerId: entry.OwnerId, token: token, time: time.Now()}
	}
	if err == nil {
		od.log.Info().Msgf("Owner.Login: User %s has been logged in, token is %s", email, token)
	}
	return token, err

}

func (od *OwnerDomain) AddCarousel(token Token, carId Carousel, prodId *string) IError {
	var ierr IError
	var td *TokenDetails
	for ok := true; ok; ok = false {
		if td, ierr = od.readToken(token); ierr != nil || td == nil || td.IsExpired() {
			ierr = erro.New(erro.ECProductAssignFailure).
				Msgf("Owner.AddCarousel: Fail to verify token").
				Err(ierr)
			break
		}
		if ierr = od.prodRepo.OwnerAddCarousel(td.ownerId, carId); ierr != nil {
			ierr = erro.New(erro.ECCarouselRegisterFailure).
				Msgf("Owner.AddCarousel: Fail to register in the local databse, ownerId:%s, CarId:%s", td.ownerId, carId).
				Err(ierr)
			break
		}
		if prodId != nil {
			if ierr = od.prodRepo.OwnerAssignStripeProductId(td.ownerId, *prodId); ierr != nil {
				ierr = erro.New(erro.ECProductAssignFailure).
					Msgf("Owner.AddCarousel: Fail to assign product id in the local databse, ownerId:%s, CarId:%s", td.ownerId, carId).
					Err(ierr)
				break
			}
		}
		if ierr = od.carousel.Register(td.ownerId, carId); ierr != nil {
			ierr = erro.New(erro.ECCarouselRegisterFailure).
				Msgf("Owner.AddCarousel: Fail to register in the remote service, ownerId:%s, CarId:%s", td.ownerId, carId).
				Err(ierr)
				// delete carousel from the local
			if err := od.prodRepo.OwnerDeleteCarousel(carId); err != nil {
				ierr = erro.New(erro.ECCarouselRegisterFailure).
					Err(ierr).
					Msgf("Owner.AddCarousel: Fail to cleanup semi-registered carousel, CarId:%s", carId.String())
			}
			break
		}
		od.log.Info().Msgf("Owner.AddCarousel: Owner %s has added Carousel %s", td.ownerId.String(), carId)
	}
	return ierr
}

func (od *OwnerDomain) AssignProdId(token Token, carId Carousel, prodId Product) IError {
	var err IError
	var td *TokenDetails
	for ok := true; ok; ok = false {
		if td, err = od.readToken(token); err != nil || td == nil || td.IsExpired() {
			err = erro.New(erro.ECProductAssignFailure).
				Msgf("Owner.AssignProdId: Fail to verify token").
				Err(err)
			break
		}

		if err = od.prodRepo.OwnerAssignStripeProductId(carId, prodId); err != nil {
			err = erro.New(erro.ECProductAssignFailure).
				Msgf("Owner.AssignProdId: Fail to assign product id into the local databse").
				Err(err)
			break
		}
		od.log.Info().Msgf("Owner.AssignProdId: Owner %s has assigned Product %s for Carousel %s", td.ownerId.String(), prodId, carId)
	}
	return err
}

func (od *OwnerDomain) AssignSkeys(token Token, skey string) IError {
	var err IError
	for ok := true; ok; ok = false {
		var td *TokenDetails
		if td, err = od.readToken(token); err != nil || td == nil || td.IsExpired() {
			err = erro.New(erro.ECUserTokenNotFoundOrExpired).
				Msgf("Owner.AssignSkeys: Fail to verify token").
				Err(err)
			break
		}
		if err = od.profRepo.OwnerAssignStripeKeys(td.ownerId, nil, &skey); err != nil {
			err = erro.New(erro.ECProductAssignFailure).
				Msgf("Owner.AssignSkeys: Fail to assign product id into the local databse").
				Err(err)
			break
		}
		od.log.Info().Msgf("Owner.AssignSkeys: Owner %s has assigned a Secret key", td.ownerId.String())

		if err = od.registerWebhook(td.ownerId, skey); err != nil {
			err = erro.New(erro.ECProductAssignFailure).
				Msgf("Owner.AssignSkeys: Fail to register webhook key for owner %s", td.ownerId).
				Err(err)
			break
		}
		od.log.Info().Msgf("Owner.AssignSkeys: Owner %s has registered a Webhook key", td.ownerId.String())
	}
	return err
}

func (od *OwnerDomain) AssignPkeys(token Token, pkey string) IError {
	var err IError
	for ok := true; ok; ok = false {
		var td *TokenDetails
		if td, err = od.readToken(token); err != nil || td == nil || td.IsExpired() {
			err = erro.New(erro.ECUserTokenNotFoundOrExpired).
				Msgf("Owner.AssignPkeys: Fail to verify token").
				Err(err)
			break
		}
		if err = od.profRepo.OwnerAssignStripeKeys(td.ownerId, &pkey, nil); err != nil {
			err = erro.New(erro.ECProductAssignFailure).
				Msgf("Owner.AssignPkeys: Fail to assign product id into the local databse").
				Err(err)
			break
		}
		od.log.Info().Msgf("Owner.AssignPkeys: Owner %s has assigned a Publish key", td.ownerId.String())
	}
	return err
}

func (od *OwnerDomain) registerWebhook(ownerId Owner, skey string) IError {
	var ierr IError
	for ok := true; ok; ok = false {
		url := od.cfg.WebhookUrl(ownerId.String())
		od.log.Info().Msgf("Owner.registerWebhook: Assign webhook url to '%s'", url)
		var whid, whkey string
		if whid, whkey, ierr = od.stripe.WebhookRegister(url, skey); ierr != nil {
			ierr = erro.New(erro.ECStripeRegWebhook).
				Msg("Owner.registerWebhook: Fail to register webhook").
				Err(ierr)
			break
		}
		if ierr = od.profRepo.OwnerAssignWebhook(ownerId, whid, whkey); ierr != nil {
			ierr = erro.New(erro.ECStripeRegWebhook).
				Msg("Owner.registerWebhook: Fail to assign webhook").
				Err(ierr)
			break
		}
	}
	return ierr
}

func (od *OwnerDomain) RefreshWebhook(token Token) IError {
	var ierr IError
	for ok := true; ok; ok = false {
		var td *TokenDetails
		if td, ierr = od.readToken(token); ierr != nil || td == nil || td.IsExpired() {
			ierr = erro.New(erro.ECUserTokenNotFoundOrExpired).
				Msg("Owner.RefreshWebhook: Fail to verify token").
				Err(ierr)
			break
		}
		var entry OwnerEntry
		if entry, ierr = od.profRepo.OwnerReadEntryByOwner(td.ownerId.String()); ierr != nil {
			ierr = erro.New(erro.ECGeneralFailure).
				Msgf("Owner.RefreshWebhook: Fail to read owner entry of owner %s", td.ownerId.String()).
				Err(ierr)
			break
		}
		if entry.SecretKey == nil {
			ierr = erro.New(erro.ECUserKeyNil).
				Msg("Owner.RefreshWebhook: Secret key is nil").
				Err(ierr)
			break
		}

		// WebhookRegister(url string, skey string) (string, string, IError)
		// WebhookUpdateUrl(url string, skey string, whkeyId string) IError

		url := od.cfg.WebhookUrl(td.ownerId.String())
		if entry.WebhookId == nil {
			var whid, whkey string
			if whid, whkey, ierr = od.stripe.WebhookRegister(url, *entry.SecretKey); ierr != nil {
				ierr = erro.New(erro.ECUserKeyNil).
					Msgf("Owner.RefreshWebhook: Fail to register new webhook for owner %s", td.ownerId.String()).
					Err(ierr)
				break
			}
			if ierr = od.profRepo.OwnerAssignWebhook(td.ownerId, whid, whkey); ierr != nil {
				ierr = erro.New(erro.ECStripeRegWebhook).
					Msgf("Owner.RefreshWebhook: Fail to assign webhook for owner %s", td.ownerId.String()).
					Err(ierr)
				break
			}
		} else {
			if ierr = od.stripe.WebhookUpdateUrl(url, *entry.SecretKey, *entry.WebhookId); ierr != nil {
				ierr = erro.New(erro.ECUserKeyNil).
					Msgf("Owner.RefreshWebhook: Fail to refresh webhook key for owner %s", td.ownerId.String()).
					Err(ierr)
				break
			}
		}

		od.log.Info().Msgf("Owner.RefreshWebhook: Owner %s has refreshed Webhook url to '%s'", td.ownerId.String(), url)
	}
	return ierr
}

func (od *OwnerDomain) ReadWhkey(ownerId string) (string, IError) {
	var whkey string
	var entry OwnerEntry
	var ierr IError
	if entry, ierr = od.profRepo.OwnerReadEntryByOwner(ownerId); ierr == nil {
		if entry.WebhookKey != nil {
			whkey = *entry.WebhookKey
		} else {
			ierr = erro.New(erro.ECUserKeyNil).Err(ierr)
		}
	} else {
		ierr = erro.New(erro.ECGeneralFailure).Err(ierr)
	}
	return whkey, ierr
}

func (od *OwnerDomain) ReadKeys(carId uuid.UUID) (pkey string, skey string, prodId string, err error) {
	var ierr IError
	var prodEntry ProductEntry
	var ownerEntry OwnerEntry
	for ok := true; ok; ok = false {
		if prodEntry, err = od.prodRepo.OwnerReadProdEntry(carId); err != nil {
			err = erro.New(erro.ECCarouselNotFound).Msgf("Carousel '%s' is not found", carId)
			break
		}
		if ownerEntry, err = od.profRepo.OwnerReadEntryByOwner(prodEntry.OwnerId.String()); err != nil {
			err = erro.New(erro.ECCarouselNotFound).Msgf("Carousel '%s' is not found", carId)
			break
		}
		// Publish key is optional
		// if ownerEntry.PublishKey == nil {
		// 	err = erro.New(erro.ECUserKeyNil).Msgf("Publish key of the owner '%s' is nil", ownerEntry.OwnerId.String())
		// 	break
		// }
		// pkey = *ownerEntry.PublishKey
		if ownerEntry.SecretKey == nil {
			err = erro.New(erro.ECUserKeyNil).Msgf("Secret key of the owner '%s' is nil", ownerEntry.OwnerId.String())
			break
		}
		skey = *ownerEntry.SecretKey
		if prodEntry.ProdId == nil {
			err = erro.New(erro.ECCarouselProdIdNil).Msgf("Product id of the carousel '%s' is nil", prodEntry.CarId.String())
			break
		}
		prodId = *prodEntry.ProdId
	}
	// od.log.Debug().Str("pkey", pkey).Str("skey", skey).Str("prodid", prodId).Send()
	return pkey, skey, prodId, ierr
}

func (od *OwnerDomain) readToken(token Token) (*TokenDetails, IError) {
	var err IError
	var tokenDetails *TokenDetails
	for ok := true; ok; ok = false {
		if od.tokens[token] == nil {
			err = erro.New(erro.ECUserTokenNotFoundOrExpired).Msgf("Token '%s' doesn't exists", token)
			break
		}
		if od.tokens[token].IsExpired() {
			err = erro.New(erro.ECUserTokenNotFoundOrExpired).Msgf("Token '%s' is expired", token)
			break
		}
		od.tokens[token].Refresh()
		tokenDetails = od.tokens[token]
	}
	return tokenDetails, err
}
