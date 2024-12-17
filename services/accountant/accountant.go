package main

import (
	"accountant/core/owner"
	"accountant/core/store"
	"accountant/infra/config"
	external "accountant/infra/external/carousel"
	"accountant/infra/external/stripe"
	"accountant/infra/logger"
	"accountant/infra/repo"
	"accountant/infra/rest"

	"github.com/joho/godotenv"
)

func main() {
	log := logger.New()
	err := godotenv.Load()
	if err != nil {
		log.Err(err).Msg("Error loading .env file")
	}

	cfg := config.New()
	drv := repo.DriverSQLite.New(cfg.DatabseUrl())
	repoProfile := repo.Profile.New(drv, &log)
	repoProduct := repo.Product.New(drv, &log)
	repoBook := repo.Book.New(drv, &log)
	cg := external.CarouselGatewayCreate(cfg)
	sg := stripe.StripeGatewayCreate()

	od := owner.OwnerDomainCreate(cg, repoProduct, repoProfile, sg, cfg, &log)
	sd := store.StoreDomainCreate(od, cg, repoBook, sg, &log)

	rest.Listen(cfg, od, sd, &log)
}
