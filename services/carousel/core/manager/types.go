package manager

type Carousel struct {
	CarId string `json:"CarouselId"`
	OwnId string `json:"OwnerId"`
}

type SnapshotData struct {
	CarId  string `json:"CarouselId"`
	Status string
	Rounds int
}
