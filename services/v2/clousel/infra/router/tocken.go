package router

import (
	"clousel/lib/fault"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	TockenRoleClient   = "client"
	TockenRoleBusiness = "business"
	// TokenRoleAdmin    = "admin"
)

type IAuth interface {
	Id() uuid.UUID
}

type ITocken interface {
	Id() uuid.UUID
	Refresh()
	IsExpired() bool
	Auth() IAuth
	Base64() ITockenBase64
	Role() Role
}

type ITockenBase64 interface {
	Id() uuid.UUID
	Str() string
}

// type ITockenData interface {
// 	Id() uuid.UUID
// }

// type ITockenInfo interface {
// 	Uid() uuid.UUID
// 	Base64String() string
// }

// type ITocken interface {
// 	ITockenInfo
// 	Data() ITockenData
// 	Prolong()
// 	IsExpired() bool
// }

// type TockenInfo struct {
// 	uid uuid.UUID
// }

// func TockenInfoCreateFromBase64String(b64 string) (ITockenInfo, fault.IError) {
// 	var err error = nil
// 	var ferr fault.IError = nil
// 	var tocken ITockenInfo = nil
// 	for ok := true; ok; ok = false {
// 		var raw []byte
// 		if raw, err = base64.StdEncoding.DecodeString(b64); err != nil {
// 			ferr = fault.New(ERouterTockenDecode).Msg(err.Error())
// 			break
// 		}
// 		var id uuid.UUID
// 		if id, err = uuid.ParseBytes(raw); err != nil {
// 			ferr = fault.New(ERouterTockenParse).Msg(err.Error())
// 			break
// 		}
// 		tocken = &TockenInfo{uid: id}
// 	}
// 	return tocken, ferr
// }

const (
	tokenExpireTime = 30 * time.Minute
)

type TokenBase64 struct {
	id  uuid.UUID
	buf []byte
	b64 string
}

func TockenCreateFromBase64String(b64 string) (ITockenBase64, fault.IError) {
	if t, err := decode([]byte(b64)); err != nil {
		return nil, fault.New(ERouterTockenDecode)
	} else {
		t.b64 = b64
		return t, nil
	}
}
func (t *TokenBase64) Id() uuid.UUID {
	return t.id
}
func (t *TokenBase64) Str() string {
	return string(t.b64)
}

type Tocken struct {
	data IAuth
	role Role
	id   uuid.UUID
	ts   time.Time
}

func TockenCreate(auth IAuth, role Role) ITocken {
	return &Tocken{id: uuid.New(), data: auth, role: role, ts: time.Now()}
}
func (t *Tocken) Id() uuid.UUID {
	return t.id
}

func (t *Tocken) Refresh() {
	t.ts = time.Now()
}
func (t *Tocken) IsExpired() bool {
	return time.Since(t.ts) > tokenExpireTime
}
func (t *Tocken) Auth() IAuth {
	return t.data
}
func (t *Tocken) Base64() ITockenBase64 {
	return encode(t.id)
}
func (t *Tocken) Role() Role {
	return t.role
}

// type Auth struct {
// 	id uuid.UUID
// }

// func (a *Auth) Id() uuid.UUID {
// 	return a.id
// }

func encode(id uuid.UUID) *TokenBase64 {
	bid, _ := id.MarshalBinary()
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(bid)))
	base64.StdEncoding.Encode(buf, bid)
	return &TokenBase64{id: id, buf: buf, b64: string(buf)}
}

func decode(b64 []byte) (*TokenBase64, error) {
	id := uuid.New()
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(b64)))
	if l, err := base64.StdEncoding.Decode(buf, b64); err == nil {
		buf = buf[:l]
		err = id.UnmarshalBinary(buf)
		return &TokenBase64{id: id, buf: buf}, err
	} else {
		return nil, err
	}
}

// type Tocken struct {
// 	TockenInfo
// 	data ITockenData
// 	ts   time.Time
// }

// func TockenDataCreate(data ITockenData) ITocken {
// 	return &Tocken{
// 		TockenInfo: TockenInfo{uid: uuid.New()},
// 		data:       data,
// 		ts:         time.Now(),
// 	}
// }

// func (t *TockenInfo) Base64String() string {
// 	return base64.StdEncoding.EncodeToString(t.uid[:])
// }

// func (t *TockenInfo) Uid() uuid.UUID {
// 	return t.uid
// }
// func (t *Tocken) Data() ITockenData {
// 	return t.data
// }

// func (t *Tocken) Prolong() {
// 	t.ts = time.Now()
// }

// func (t *Tocken) IsExpired() bool {
// 	return time.Since(t.ts) > tokenExpireTime
// }

type ITockenStore interface {
	Find(t ITockenBase64) (ITocken, fault.IError)
	Add(t ITocken) ITockenStore
	Cleanup() ITockenStore
}

type TockenStore struct {
	store map[uuid.UUID]ITocken
}

func TockenStoreCreate() TockenStore {
	return TockenStore{
		store: make(map[uuid.UUID]ITocken),
	}
}

func (ts *TockenStore) Find(t ITockenBase64) (ITocken, fault.IError) {
	if ts.store[t.Id()] == nil {
		return nil, fault.New(ERouterNotFound).Msgf("%s Not Found", t.Id().String())
	} else if ts.store[t.Id()].IsExpired() {
		delete(ts.store, t.Id())
		return nil, fault.New(ERouterNotFound).Msgf("%s Expired and removed", t.Id().String())
	} else {
		ts.store[t.Id()].Refresh()
		return ts.store[t.Id()], nil
	}
}

func (ts *TockenStore) Add(t ITocken) ITockenStore {
	ts.store[t.Id()] = t
	return ts
}

func (ts *TockenStore) Cleanup() ITockenStore {
	for k, v := range ts.store {
		if v.IsExpired() {
			delete(ts.store, k)
		}
	}
	return ts
}
