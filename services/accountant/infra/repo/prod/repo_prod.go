package prod

import (
	"accountant/infra/repo/driver"
	"accountant/infra/repo/types"
	"database/sql"
	"fmt"

	"accountant/core/owner"
	erro "accountant/core/owner/error"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	TableProduct = "product"
)

type ProductColumn struct {
	OwnerId types.Named[string]
	CarId   types.Named[string]
	ProdId  types.NamedOpt[string]
}

func ProductColumnDefault() ProductColumn {
	return ProductColumn{
		OwnerId: types.NamedCreateDefault[string]("OwnerId"),
		CarId:   types.NamedCreateDefault[string]("CarouselId"),
		ProdId:  types.NamedOptCreateDefault[string]("ProductId"),
	}
}

type RepositoryProduct struct {
	drv driver.IDBDriver
	log *zerolog.Logger
}

func RepositoryProductCreate(drv driver.IDBDriver, log *zerolog.Logger) *RepositoryProduct {
	return &RepositoryProduct{drv: drv, log: log}
}

func (r *RepositoryProduct) OwnerAddCarousel(ownerId owner.Owner, carId owner.Carousel) erro.IError {
	var prompt string
	var ierr erro.IError
	c := ProductColumnDefault()
	c.OwnerId.Value = ownerId.String()
	c.CarId.Value = carId.String()
	prompt = fmt.Sprintf("insert into '%s' (%s, %s) values ('%s', '%s')", TableProduct,
		c.OwnerId.Name(), c.CarId.Name(),
		c.OwnerId.Value, c.CarId.Value)
	if e := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Product.OwnerAddCarousel: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Msg("Repository.Product.OwnerAddCarousel: Fail to Register")
		}
		return err
	}); e != nil {
		ierr = erro.New(erro.ECRepoExecPrompt).Msg(e.Error())
	}
	return ierr
}

func (r *RepositoryProduct) OwnerDeleteCarousel(carId owner.Carousel) erro.IError {
	var ierr erro.IError
	c := ProductColumnDefault()
	c.CarId.Value = carId.String()
	prompt := fmt.Sprintf("delete from '%s' where %s='%s'", TableProduct, c.CarId.Name(), c.CarId.Value)
	if err := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Product.OwnerDeleteCarousel: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str(c.CarId.Name(), c.CarId.Value).Msg("Repository.Product.OwnerDeleteCarousel: Fail to Remove Carousel")
		}
		return err
	}); err != nil {
		ierr = erro.New(erro.ECCarouselDeleteFailure).Msgf("Repository.Product.OwnerDeleteCarousel: Failed with error:%v", err)
	}

	return ierr
}

func (r *RepositoryProduct) OwnerAssignStripeProductId(carId owner.Carousel, prodId owner.Product) erro.IError {
	var ierr erro.IError
	var e error
	c := ProductColumnDefault()
	c.CarId.Value = carId.String()
	c.ProdId.Value = &prodId
	prompt := fmt.Sprintf("update '%s' set %s='%s' where %s='%s'", TableProduct,
		c.ProdId.Name(), *c.ProdId.Value, c.CarId.Name(), c.CarId.Value)
	if e = r.drv.Session(func(db *sql.DB) error {
		if _, e = db.Exec(prompt); e == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Product.OwnerAssignStripeProductId: Success")
		} else {
			r.log.Err(e).Str("SQL", prompt).
				Str(c.CarId.Name(), c.CarId.Value).
				Str(c.ProdId.Name(), *c.ProdId.Value).
				Msg("Repository.Product.OwnerAssignStripeProductId: Failure")
		}
		return e
	}); e != nil {
		ierr = erro.New(erro.ECRepoExecPrompt).Msgf("Repository.Product.OwnerAssignStripeProductId: Fail to assign ProductId %s", e.Error())
	}

	return ierr
}

func (r *RepositoryProduct) OwnerReadProdEntry(carousel owner.Carousel) (owner.ProductEntry, erro.IError) {
	var ierr erro.IError
	var e error
	c := ProductColumnDefault()
	c.CarId.Value = carousel.String()
	prompt := fmt.Sprintf("select * from '%s' where %s='%s'", TableProduct, c.CarId.Name(), c.CarId.Value)
	if e = r.drv.Session(func(db *sql.DB) error {
		if e = db.QueryRow(prompt).Scan(
			&c.OwnerId.Value,
			&c.CarId.Value,
			&c.ProdId.Value); e == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Product.OwnerReadProdEntry: Success")
		}
		return e
	}); e != nil {
		ierr = erro.New(erro.ECRepoExecPrompt).Msg(e.Error())
		r.log.Err(e).Str("SQL", prompt).Msgf("Repository.Product.OwnerReadProdEntry: Fail to Read from '%s' table", TableProduct)
	}
	oid, e := uuid.Parse(c.OwnerId.Value)
	if e != nil {
		ierr = erro.New(erro.ECGeneralFailure).Msgf("Fail to parse OwnerId as uuid, error: %s", e.Error())
	}
	cid, e := uuid.Parse(c.CarId.Value)
	if e != nil {
		ierr = erro.New(erro.ECGeneralFailure).Msgf("Fail to parse CarouselId as uuid, error:%s", e.Error())
	}
	return owner.ProductEntry{
		OwnerId: oid,
		CarId:   cid,
		ProdId:  c.ProdId.Value,
	}, ierr
}

func (r *RepositoryProduct) OwnerReadProdEntries(ownerId owner.Owner) ([]owner.ProductEntry, erro.IError) {
	var pe []owner.ProductEntry
	return pe, nil
}
