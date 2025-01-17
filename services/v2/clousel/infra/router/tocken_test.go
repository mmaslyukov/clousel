package router_test

import (
	"clousel/infra/router"
	"testing"

	"github.com/google/uuid"
)

type Auth struct {
	id uuid.UUID
}

func (a *Auth) Id() uuid.UUID {
	return a.id
}

func TestTocken(t *testing.T) {
	id := uuid.New()
	tn := router.TockenCreate(&Auth{id: id}, router.TockenRoleClient)
	t.Log(tn.Id())
	t.Log(tn.Base64().Id())
	t.Log(tn.Base64().Str())
	td, err := router.TockenCreateFromBase64String(tn.Base64().Str())
	if err != nil {
		t.Error(err)
	}
	t.Log(td.Id())
	t.Log(td.Str())
}
