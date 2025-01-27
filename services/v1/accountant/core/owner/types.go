package owner

import (
	"time"

	"github.com/google/uuid"
)

type Session = string
type Product = string
type Owner = uuid.UUID
type Carousel = uuid.UUID

type Token = uuid.UUID
type UserRole int

const (
	UserRoleRegular   UserRole = 0
	UserRoleAdmin     UserRole = 100
	UserRoleArchitect UserRole = 999
)
const (
	tokenExpireTime = 30 * time.Minute
)

type TokenDetails struct {
	ownerId Owner
	token   Token
	time    time.Time
}

func (t *TokenDetails) IsExpired() bool {
	return time.Since(t.time) > tokenExpireTime
}
func (t *TokenDetails) Refresh() {
	t.time = time.Now()
}

type OwnerEntry struct {
	OwnerId    Owner
	Email      string
	Password   string
	Role       UserRole
	SecretKey  *string
	PublishKey *string
	WebhookId  *string
	WebhookKey *string
}

type ProductEntry struct {
	OwnerId Owner
	CarId   Carousel
	ProdId  *Product
}
