package chat

import (
	"time"

	"github.com/google/uuid"
)

type ModelBase struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ModelBase
	Username string `gorm:"uniqueIndex" json:"username"`
	Name     string `gorm:"not null" json:"name"`

	Conversations []*Conversation `gorm:"many2many:conversation_memberships;foreignKey:ID;joinForeignKey:MemberID;" json:"conversations"`
}

type Conversation struct {
	ModelBase
	Name       string    `gorm:"not null" json:"name"`
	IsPublic   bool      `gorm:"not null;default:false" json:"is_public"`
	InviteLink string    `gorm:"uniqueIndex" json:"invite_link"`
	CreatorID  uuid.UUID `gorm:"index" json:"creator_id"`
	OwnerID    uuid.UUID `gorm:"index" json:"owner_id"`

	Creator User `json:"creator"`
	Owner   User `json:"owner"`

	Members []*User `gorm:"many2many:conversation_memberships;references:ID;joinReferences:ConversationID;"`
}

type ConversationMembership struct {
	ConversationID uuid.UUID `gorm:"primaryKey"`
	MemberID       uuid.UUID `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Conversation Conversation
	Member       User
}

type Message struct {
	ModelBase
	ConversationID    uuid.UUID
	Conversation      Conversation
	SenderID          uuid.UUID
	Sender            User
	ContentType       string    `gorm:"default:text"`
	ContentFormatting string    `gorm:"default:markdown"`
	Content           string    `gorm:"not null"`
	SentAt            time.Time `gorm:"autoCreateTime"`
}
