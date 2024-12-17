package accountment_service

import (
	"accountant_service/domain/accountment/accountment_provider"
	"accountant_service/framework/core"
)

type IServiceSales interface {
	core.IEventSubscribable
	accountment_provider.IPortApiSales
}
