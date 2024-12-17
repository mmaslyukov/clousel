package event

type IRepositoryCarousel interface {
	ReadCarouselsIds() ([]string, error)
}
