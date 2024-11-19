package manager

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Manager struct {
	crRepo IPortManagerAdapterCarouselRepository
	snRepo IPortManagerAdapterSnapshotRepository
	log    *zerolog.Logger
}

func New(crRepo IPortManagerAdapterCarouselRepository, snRepo IPortManagerAdapterSnapshotRepository, log *zerolog.Logger) *Manager {
	return &Manager{crRepo: crRepo, snRepo: snRepo, log: log}
}

func (m *Manager) Register(c Carousel) error {
	var err error
	for ok := true; ok; ok = false {
		if err = m.crRepo.ManagerAddCarousel(c); err != nil {
			break
		}
		if err = m.snRepo.ManagerStoreNewSnapshot(c.CarId); err != nil {
			break
		}
	}
	return err
}

func (m *Manager) Unregister(c Carousel) error {
	if len(c.CarId) > 0 {
		return m.crRepo.ManagerRemoveCarousel(c.CarId)
	} else if len(c.OwnId) > 0 {
		return m.crRepo.ManagerRemoveOwner(c.OwnId)
	} else {
		return fmt.Errorf("Invalid argument")
	}
}

func (m *Manager) Read(ownerId string) ([]Carousel, error) {
	return m.crRepo.ManagerReadOwnedCarousel(ownerId)
}
