package chat

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getConversation(db *gorm.DB, id uuid.UUID) (*Conversation, error) {
	conv := Conversation{}

	result := db.First(&conv, id)
	if err := result.Error; err != nil {
		return nil, err
	}

	return &conv, nil
}

func (l *HTTPListener) getConversations(user User) ([]*Conversation, error) {
	return user.Conversations, nil
}

func (l *HTTPListener) createConversation(name string, members []uuid.UUID) (conv Conversation, err error) {
	tx := l.db.Begin()

	defer func() {
		if err != nil {
			l.db.Rollback()
		} else {
			l.db.Commit()
		}
	}()

	conv.Name = name
	tx.Create(&conv)

	return
}
