package main

import (
	"clousel/core/business"
	"clousel/core/client"
	"clousel/core/machine"
	"clousel/infra/cfg"
	"clousel/infra/ipc"
	"clousel/infra/log"
	"clousel/infra/repo"
	"clousel/infra/router"
	"clousel/infra/stripe"

	"github.com/joho/godotenv"
)

func main() {
	log := log.New()
	err := godotenv.Load()
	if err != nil {
		log.Err(err).Msg("Error loading .env file")
	}
	config := cfg.New()
	drv := repo.DriverSQLite.New(cfg.New().DatabseUrl())
	repoCompany := repo.Company.New(drv, log)
	repoMachine := repo.Machine.New(drv, log)
	repoUser := repo.User.New(drv, log)
	stripe := stripe.StripeGatewayCreate(log)
	ipc := ipc.IpcCreate(config, log)
	// machBaron := machine.BarounessCreate(log)
	domainBusiness := business.BusinessCreate(config, repoCompany, stripe, log)
	domainClient := client.ClientCreate(repoUser, repoUser, repoUser, domainBusiness, stripe, log)
	domainMachine := machine.MachineCreate(repoMachine, repoMachine, domainClient, config, ipc, log)
	router.Listen(config, domainClient, domainBusiness, domainMachine, log)
}
