package prod_test

import (
	"accountant/infra/logger"
	"accountant/infra/repo"
	"accountant/infra/repo/driver"
	"accountant/infra/repo/prod"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

func CreateTableProduct(drv driver.IDBDriver) error {
	prompt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS '%s' ('OwnerId' string NOT NULL,'CarouselId' string UNIQUE NOT NULL,'ProductId' string)", prod.TableProduct)
	return drv.Session(func(db *sql.DB) error {
		var err error
		_, err = db.Exec(prompt)
		return err
	})
}

func TestAddCarouselAndAssignProductId(t *testing.T) {
	const dbPath = "test.db"
	const email = "test@mail.com"
	const password = "123qweasd"
	const prodId = "pr_12h3123kl"
	var ownerId = uuid.New()
	var carId = uuid.New()
	log := logger.New()
	drv := repo.DriverSQLite.New(dbPath)
	repoProd := repo.Product.New(drv, &log)

	for ok := true; ok; ok = false {
		if err := CreateTableProduct(drv); err != nil {
			t.Errorf("Fail to create table '%s', err:%s", prod.TableProduct, err)
			break
		}
		if err := repoProd.OwnerAddCarousel(ownerId, carId); err != nil {
			t.Errorf("Fail to Add carousel, err: %s", err.Full())
			break
		}
		entryProduct, err := repoProd.OwnerReadProdEntry(carId)
		if err != nil {
			t.Errorf("Fail to read Product entry: %s", err.Full())
			break
		}
		if entryProduct.CarId != carId {
			t.Errorf("Fail to verify carId, %s != %s", entryProduct.CarId, carId)
			break
		}
		if err = repoProd.OwnerAssignStripeProductId(carId, prodId); err != nil {
			t.Errorf("Fail to Assign Product id %s", err.Full())
			break
		}
		entryProduct, err = repoProd.OwnerReadProdEntry(carId)
		if err != nil {
			t.Errorf("Fail to read Product entry: %s", err.Full())
			break
		}
		if entryProduct.CarId != carId {
			t.Errorf("Fail to verify carId, %s != %s", entryProduct.CarId, carId)
			break
		}
		if entryProduct.ProdId == nil {
			t.Errorf("Fail to get valid prodId, is nil, expects %s", prodId)
			break
		}
		if *entryProduct.ProdId != prodId {
			t.Errorf("Fail to verify prodId, %s != %s", *entryProduct.ProdId, prodId)
			break
		}
	}
	os.Remove(dbPath)

}
