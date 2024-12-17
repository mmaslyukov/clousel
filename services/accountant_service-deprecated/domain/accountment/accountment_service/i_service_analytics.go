package accountment_service

import (
	"accountant_service/domain/accountment/accountment_provider"
)

type IServiceAnalytics interface {
	accountment_provider.IPortApiAnalytics
}
