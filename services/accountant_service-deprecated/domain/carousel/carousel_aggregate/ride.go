package carousel_aggregate

import (
	"accountant_service/domain/carousel"
	"accountant_service/framework/utils"
	"time"

	"github.com/google/uuid"
)

type Ride struct {
	carouselId    uuid.UUID
	rides         uint
	delivered     bool
	timeCreated   time.Time
	timeDelivered utils.Optional[time.Time]
}

func RideCreateFull(carouselId uuid.UUID, rides uint, delivered bool, timeCreated time.Time, timeDelivered time.Time) Ride {
	return Ride{
		carouselId:    carouselId,
		rides:         rides,
		delivered:     delivered,
		timeCreated:   timeCreated,
		timeDelivered: utils.OptionalValueCreate[time.Time](timeDelivered),
	}
}

func RideCreate(carouselId uuid.UUID, rides uint) Ride {
	return Ride{
		carouselId:    carouselId,
		rides:         rides,
		delivered:     false,
		timeCreated:   time.Now(),
		timeDelivered: utils.OptionalNilCreate[time.Time](),
	}
}

func (r *Ride) ToMinimal() carousel.RideMinimal {
	return carousel.RideMinimal{
		CarouselId: r.carouselId,
		Rides:      r.rides,
	}
}

func (r *Ride) CarouselId() uuid.UUID {
	return r.carouselId
}

func (r *Ride) Rides() uint {
	return r.rides
}

func (r *Ride) SetDelivered() {
	r.delivered = true
	r.timeDelivered.Replace(time.Now())
}

func (r *Ride) IsDelivered() bool {
	return r.delivered
}

func (r *Ride) TimeCreated() time.Time {
	return r.timeCreated
}

func (r *Ride) TimeDelivered() utils.Optional[time.Time] {
	return r.timeDelivered
}
