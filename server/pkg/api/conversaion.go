package api

import (
	"context"
	"log"
	"time"

	"githum.com/mengdage/gochat/helper"
)

type ChatConversaion struct {
	ID       string
	CreateAt time.Time
}

type ChatMessage struct {
	ConversationID string    `json:"conversationId"`
	FromUser       string    `json:"fromUser"`
	ToUser         string    `json:"toUser"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
}

type ConversationService interface {
	GetConversationHistory(ctx context.Context, currentUser string, userIDs ...string) ([]*ChatMessage, error)
}

type ConversationStorage interface {
	GetConversation(ctx context.Context, string, currentUser string) (*ChatConversaion, error)
	GetConversationMessages(ctx context.Context, cID string) ([]*ChatMessage, error)
}

type conversationService struct {
	storage ConversationStorage
}

func NewConversationService(storage ConversationStorage) ConversationService {
	return &conversationService{
		storage: storage,
	}
}

func (c *conversationService) GetConversationHistory(ctx context.Context, currentUser string, userIDs ...string) ([]*ChatMessage, error) {
	cID := helper.CreateConversationID(userIDs...)

	chatConv, err := c.storage.GetConversation(ctx, cID, currentUser)
	if err != nil {
		log.Printf("Error while getting converstion: %v", err)
		return nil, err
	}

	if chatConv == nil {
		log.Println("The conversation does not exists")
		return []*ChatMessage{}, nil
	}

	msgs, err := c.storage.GetConversationMessages(ctx, cID)
	if err != nil {
		log.Printf("Error while getting converstion messages: %v", err)
		return nil, err
	}

	return msgs, nil
}
