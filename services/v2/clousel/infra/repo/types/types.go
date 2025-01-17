package types

/*
 - Changed time.RFC3339 to time.DateTime, due to incompatible time fromat in sqlite and in app
*/

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Named[T any] struct {
	name  string
	value T
}

func (n *Named[T]) Name() string {
	return n.name
}
func (n *Named[T]) Value() T {
	return n.value
}
func (n *Named[T]) ValuePtr() *T {
	return &n.value
}
func (n *Named[T]) SetValue(value T) {
	n.value = value
}

func NamedCreateDefault[T any](name string) Named[T] {
	return Named[T]{name: name}
}

type NamedOpt[T any] struct {
	name  string
	value *T
}

func (n *NamedOpt[T]) Name() string {
	return n.name
}
func (n *NamedOpt[T]) Value() *T {
	return n.value
}
func (n *NamedOpt[T]) ValuePtr() **T {
	return &n.value
}

func (n *NamedOpt[T]) SetValue(value *T) {
	n.value = value
}
func NamedOptCreateDefault[T any](name string) NamedOpt[T] {
	return NamedOpt[T]{name: name}
}

type UUIDString struct {
	str string
}

func (u *UUIDString) Uuid() uuid.UUID {
	uid := uuid.Max
	if pid, err := uuid.Parse(u.str); err == nil {
		uid = pid
	}
	return uid
}

func (u *UUIDString) SetUuid(uid uuid.UUID) {
	u.str = uid.String()
}

func (u *UUIDString) Str() string {
	return u.str
}
func (u *UUIDString) Ptr() *string {
	return &u.str
}
func (u *UUIDString) SetStr(str string) {
	u.str = str
}

type TimeString struct {
	str string
}

func (u *TimeString) Str() string {
	return u.str
}
func (u *TimeString) Ptr() *string {
	return &u.str
}
func (u *TimeString) SetStr(str string) {
	u.str = str
}
func (u *TimeString) Time() time.Time {
	t, _ := time.Parse(time.DateTime, u.str)
	// t, _ := time.Parse(time.RFC3339, u.str)
	return t
}

func (u *TimeString) SetTime(tm time.Time) {
	u.str = fmt.Sprint(tm.Format(time.DateTime))
	// u.str = fmt.Sprint(tm.Format(time.RFC3339))
}
