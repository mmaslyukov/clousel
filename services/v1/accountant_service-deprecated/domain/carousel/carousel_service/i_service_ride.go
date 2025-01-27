package carousel_service

import (
	"accountant_service/domain/carousel/carousel_provider"
	"accountant_service/framework/core"
)

type IServiceRide interface {
	core.IEventListener
	carousel_provider.IPortApiRide
}
