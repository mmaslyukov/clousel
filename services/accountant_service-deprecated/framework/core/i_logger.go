package core

import (
	"log"
)

type ILogger interface {
	Err() *log.Logger
	Wrn() *log.Logger
	Inf() *log.Logger
	Dbg() *log.Logger
}