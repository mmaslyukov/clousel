package persistency

import (
	"fmt"

	"carousel_service/internal/logger"
	. "carousel_service/internal/ports/port_persistency"
	. "carousel_service/internal/utils"
)

const (
	table_status = "carousel-service-status"
)

type StatusRecordAccessor struct {
	session *Session
}

func (c *StatusRecordAccessor) Create(record StatusRecord) error {
	var err error
	if record.RoundsReady.Valid() && record.Status.Valid() {
		prompt := fmt.Sprintf("insert into '%s' (CarouselId, RoundsReady, Status) values ('%s', %d, '%s')", table_status, record.CarouselId, record.RoundsReady.Get(), record.Status.Get())
		logger.Debug.Printf("SQL prompt: %s", prompt)
		_, err = c.session.db.Exec(prompt)
	} else {
		err = fmt.Errorf("Invalid argument(s)")
	}
	return err
}

func (c *StatusRecordAccessor) Update(record StatusRecord) error {
	var err error
	var prompt string
	if record.RoundsReady.Valid() && record.Status.Valid() {
		prompt = fmt.Sprintf("update '%s' set RoundsReady=%d, Status='%s', Time=CURRENT_TIMESTAMP where CarouselId='%s'",
			table_status,
			record.RoundsReady.Get(),
			record.Status.Get(),
			record.CarouselId)
	} else if record.RoundsReady.Valid() {
		prompt = fmt.Sprintf("update '%s' set RoundsReady=%d, Time=CURRENT_TIMESTAMP  where CarouselId='%s'",
			table_status,
			record.RoundsReady.Get(),
			record.CarouselId)
	} else if record.Status.Valid() {
		prompt = fmt.Sprintf("update '%s' set Status='%s', Time=CURRENT_TIMESTAMP  where CarouselId='%s'",
			table_status,
			record.Status.Get(),
			record.CarouselId)
	}
	_, err = c.session.db.Exec(prompt)
	logger.Debug.Printf("SQL prompt: %s", prompt)
	return err
}

func (c *StatusRecordAccessor) Read(carouselId string) (Optional[StatusRecord], error) {
	r := NewDefaultStatusRecord()
	prompt := fmt.Sprintf("select * from '%s' where CarouselId='%s'", table_status, carouselId)
	logger.Debug.Printf("SQL prompt: %s", prompt)
	err := c.session.db.QueryRow(prompt).Scan(
		&r.CarouselId,
		&r.Time,
		r.Status.Ptr(),
		r.RoundsReady.Ptr(),
	)
	if err != nil {
		return NewOptionalNil[StatusRecord](), err
	}
	return NewOptionalValue[StatusRecord](r), nil
}

func (c *StatusRecordAccessor) ReadOneBy(where func() string) (Optional[StatusRecord], error) {
	r := NewDefaultStatusRecord()
	prompt := fmt.Sprintf("select * from '%s' where %s", table_status, where())
	logger.Debug.Printf("SQL prompt: %s", prompt)
	err := c.session.db.QueryRow(prompt).Scan(
		&r.CarouselId,
		&r.Time,
		r.Status.Ptr(),
		r.RoundsReady.Ptr(),
	)
	if err != nil {
		return NewOptionalNil[StatusRecord](), err
	}
	return NewOptionalValue[StatusRecord](r), nil
}

func (c *StatusRecordAccessor) ReadManyBy(where func() string) (Optional[[]StatusRecord], error) {
	prompt := fmt.Sprintf("select * from '%s' where %s", table_status, where())
	logger.Debug.Printf("SQL prompt: %s", prompt)
	rows, err := c.session.db.Query(prompt)
	if err != nil {
		return NewOptionalNil[[]StatusRecord](), err
	}
	defer rows.Close()

	var srArray []StatusRecord
	for rows.Next() {
		r := NewDefaultStatusRecord()
		err := rows.Scan(
			&r.CarouselId,
			&r.Time,
			r.Status.Ptr(),
			r.RoundsReady.Ptr(),
		)
		if err != nil {
			return NewOptionalNil[[]StatusRecord](), err
		}
		srArray = append(srArray, r)
	}
	return NewOptionalValue[[]StatusRecord](srArray), nil
}

func (c *StatusRecordAccessor) ReadAll() (Optional[[]StatusRecord], error) {
	prompt := fmt.Sprintf("select * from '%s'", table_status)
	logger.Debug.Printf("SQL prompt: %s", prompt)
	rows, err := c.session.db.Query(prompt)
	if err != nil {
		return NewOptionalNil[[]StatusRecord](), err
	}
	defer rows.Close()

	var srArray []StatusRecord
	for rows.Next() {
		r := NewDefaultStatusRecord()
		err := rows.Scan(
			&r.CarouselId,
			&r.Time,
			r.Status.Ptr(),
			r.RoundsReady.Ptr(),
		)
		if err != nil {
			return NewOptionalNil[[]StatusRecord](), err
		}
		srArray = append(srArray, r)
	}
	return NewOptionalValue[[]StatusRecord](srArray), nil
}

func (c *StatusRecordAccessor) Delete(carouselId string) error {
	prompt := fmt.Sprintf("delete from '%s' where CarouselId='%s'", table_status, carouselId)
	logger.Debug.Printf("SQL prompt: %s", prompt)
	_, err := c.session.db.Exec(prompt)
	return err
}
