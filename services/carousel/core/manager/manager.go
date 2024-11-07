package manager

import (
	"github.com/rs/zerolog"
)

type Manager struct {
	repoMan IPortManagerAdapterRepository
	log     *zerolog.Logger
}

func New(repoMan IPortManagerAdapterRepository, log *zerolog.Logger) *Manager {
	return &Manager{repoMan: repoMan, log: log}
}

func (m *Manager) Register(c Carousel) error {
	var err error
	for ok := true; ok; ok = false {
		if err = m.repoMan.AddCarousel(c); err != nil {
			break
		}
		if exists, e := m.repoMan.IsCarouselExistsInEvents(c); e != nil || exists {
			err = e
			break
		}
		if err = m.repoMan.AddEventWithStatusNew(c); err != nil {
			break
		}
	}
	return err
}

func (m *Manager) Unregister(c Carousel) error {
	return m.repoMan.Remove(c)
}
func (m *Manager) Read(ownerId string) ([]Carousel, error) {
	return m.repoMan.ReadOwned(ownerId)
}
