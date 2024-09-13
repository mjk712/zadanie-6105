package models

import (
	"github.com/google/uuid"
	"time"
)

type OrganizationType string

const (
	IE  OrganizationType = "IE"
	LLC OrganizationType = "LLC"
	JSC OrganizationType = "JSC"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Firstname string    `json:"first_name"`
	Lastname  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Organization struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OrganizationType
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrganizationResponsible struct {
	Id             uuid.UUID `json:"id"`
	OrganizationId uint64    `json:"organization_id"`
	UserId         uint64    `json:"user_id"`
}

type Tender struct {
	Id              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Description     string    `json:"description" db:"description"`
	ServiceType     string    `json:"serviceType" db:"service_type"`
	Status          string    `json:"status" db:"status"`
	OrganizationId  uuid.UUID `json:"organizationId" db:"organization_id"`
	CreatorUsername string    `json:"creatorUsername" db:"creator_username"`
	Version         uint64    `json:"version" db:"version"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

type Bid struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	TenderId    uuid.UUID `json:"tenderId" db:"tender_id"`
	AuthorType  string    `json:"authorType" db:"author_type"`
	AuthorId    uuid.UUID `json:"authorId" db:"author_id"`
	Version     uint64    `json:"version" db:"version"`
	Decision    string    `json:"decision" db:"decision"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type BidFeedback struct {
	BidID        uuid.UUID `db:"id"`
	TenderID     uuid.UUID `db:"tender_id"`
	AuthorID     uuid.UUID `db:"author_id"`
	FeedbackText string    `db:"feedback_text"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
