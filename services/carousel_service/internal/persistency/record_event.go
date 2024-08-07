package persistency

import (
	_ "carousel_service/internal/logger"
	// . "carousel_service/internal/ports/port_persistency"
	. "carousel_service/internal/utils"
	"fmt"
)

const (
	table_event = "carousel-service-evt-queue"
)

type EventRecordAccessor struct {
	session *Session
}

func (c *EventRecordAccessor) Create(any) error {
	var err error
	// if record.RoundsReady.Valid() && record.Status.Valid() {
	// 	prompt := fmt.Sprintf("insert into '%s' (CarouselId, RoundsReady, Status) values ('%s', %d, '%s')", TABLE_STATUS, record.CarouselId, record.RoundsReady.Get(), record.Status.Get())
	// 	logger.Debug.Printf("SQL prompt: %s", prompt)
	// 	_, err = c.session.db.Exec(prompt)
	// } else {
	err = fmt.Errorf("Invalid argument(s)")
	// }
	return err
}

func (c *EventRecordAccessor) Update(any) error {
	return fmt.Errorf("Update function isn't implemented")
}

func (c *EventRecordAccessor) Read(carouselId string) (Optional[any], error) {
	// var r StatusRecord
	// r.Status.Set("")
	// r.RoundsReady.Set(0)
	// prompt := fmt.Sprintf("select * from '%s' where CarouselId='%s'", TABLE_STATUS, carouselId)
	// logger.Debug.Printf("SQL prompt: %s", prompt)
	// err := c.session.db.QueryRow(prompt).Scan(
	// 	&r.CarouselId,
	// 	&r.Time,
	// 	r.Status.Ptr(),
	// 	r.RoundsReady.Ptr(),
	// )
	// if err != nil {
	return NewOptionalNil[any](), nil
	// }
	// return NewOptionalValue[StatusRecord](r), nil
}

func (c *EventRecordAccessor) Delete(carouselId string) error {
	return fmt.Errorf("Delete function isn't implemented")
}
