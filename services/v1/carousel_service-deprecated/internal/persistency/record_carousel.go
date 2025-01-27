package persistency

import (
	"carousel_service/internal/logger"
	. "carousel_service/internal/ports/port_persistency"
	. "carousel_service/internal/utils"
	"fmt"
)

const (
	table_carousel = "carousel-service-record"
)

type CarouselRecordAccessor struct {
	session *Session
}

func (c *CarouselRecordAccessor) Create(record CarouselRecord) error {
	var prompt string
	if record.RoundTime.Valid() {
		prompt = fmt.Sprintf("insert into '%s' (CarouselId, OwnerId, RoundTime) values ('%s', '%s', %d)", table_carousel, record.CarouselId, record.OwnerId, record.RoundTime.Get())
	} else {
		prompt = fmt.Sprintf("insert into '%s' (CarouselId) values ('%s')", table_carousel, record.CarouselId)
	}
	logger.Debug.Printf("SQL prompt: %s", prompt)
	_, err := c.session.db.Exec(prompt)
	return err
}

func (c *CarouselRecordAccessor) Update(record CarouselRecord) error {
	var err error
	if record.RoundTime.Valid() {
		prompt := fmt.Sprintf("update '%s' set RoundTime=%d where CarouselId='%s'", table_carousel, record.RoundTime.Get(), record.CarouselId)
		logger.Debug.Printf("SQL prompt: %s", prompt)
		_, err = c.session.db.Exec(prompt)
	} else {
		err = fmt.Errorf("Can't update '%s' table, becaouse provided RoundTime isn't valid", table_carousel)
	}
	return err
}

func (c *CarouselRecordAccessor) Read(carouselId string) (Optional[CarouselRecord], error) {
	var r CarouselRecord
	var roundTime int
	prompt := fmt.Sprintf("select * from '%s' where CarouselId='%s'", table_carousel, carouselId)
	logger.Debug.Printf("SQL prompt: %s", prompt)
	err := c.session.db.QueryRow(prompt).Scan(
		&r.CarouselId,
		&r.OwnerId,
		&roundTime,
	)
	if err != nil {
		return NewOptionalNil[CarouselRecord](), err
	}
	if roundTime != 0 {
		r.RoundTime = NewOptionalValue[int](roundTime)
	} else {
		r.RoundTime = NewOptionalNil[int]()
	}
	return NewOptionalValue[CarouselRecord](r), err
}

func (c *CarouselRecordAccessor) ReadOneBy(where func() string) (Optional[CarouselRecord], error) {
	var r CarouselRecord
	var roundTime int
	prompt := fmt.Sprintf("select * from '%s' where %s", table_carousel, where())
	logger.Debug.Printf("SQL prompt: %s", prompt)
	err := c.session.db.QueryRow(prompt).Scan(
		&r.CarouselId,
		&r.OwnerId,
		&roundTime,
	)
	if err != nil {
		return NewOptionalNil[CarouselRecord](), err
	}
	if roundTime != 0 {
		r.RoundTime = NewOptionalValue[int](roundTime)
	} else {
		r.RoundTime = NewOptionalNil[int]()
	}
	return NewOptionalValue[CarouselRecord](r), err
}

func (c *CarouselRecordAccessor) ReadManyBy(where func() string) (Optional[[]CarouselRecord], error) {
	var roundTime int
	prompt := fmt.Sprintf("select * from '%s' where %s", table_carousel, where())
	logger.Debug.Printf("SQL prompt: %s", prompt)
	rows, err := c.session.db.Query(prompt)
	if err != nil {
		return NewOptionalNil[[]CarouselRecord](), err
	}
	defer rows.Close()

	var crArray []CarouselRecord
	for rows.Next() {
		r := NewDefaultCarouselRecord()
		err := rows.Scan(
			&r.CarouselId,
			&r.OwnerId,
			&roundTime,
		)
		if err != nil {
			return NewOptionalNil[[]CarouselRecord](), err
		}
		if roundTime != 0 {
			r.RoundTime = NewOptionalValue[int](roundTime)
		} else {
			r.RoundTime = NewOptionalNil[int]()
		}
		crArray = append(crArray, r)
	}
	return NewOptionalValue[[]CarouselRecord](crArray), err
}

func (c *CarouselRecordAccessor) ReadAll() (Optional[[]CarouselRecord], error) {
	var roundTime int
	prompt := fmt.Sprintf("select * from '%s'", table_carousel)
	logger.Debug.Printf("SQL prompt: %s", prompt)
	rows, err := c.session.db.Query(prompt)
	if err != nil {
		return NewOptionalNil[[]CarouselRecord](), err
	}
	defer rows.Close()

	var crArray []CarouselRecord
	for rows.Next() {
		r := NewDefaultCarouselRecord()
		err := rows.Scan(
			&r.CarouselId,
			&r.OwnerId,
			&roundTime,
		)
		if err != nil {
			return NewOptionalNil[[]CarouselRecord](), err
		}
		if roundTime != 0 {
			r.RoundTime = NewOptionalValue[int](roundTime)
		} else {
			r.RoundTime = NewOptionalNil[int]()
		}
		crArray = append(crArray, r)
	}
	return NewOptionalValue[[]CarouselRecord](crArray), err
}

func (c *CarouselRecordAccessor) Delete(carouselId string) error {
	prompt := fmt.Sprintf("delete from '%s' where CarouselId='%s'", table_carousel, carouselId)
	logger.Debug.Printf("SQL prompt: %s", prompt)
	_, err := c.session.db.Exec(prompt)
	return err
}
