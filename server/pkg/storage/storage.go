package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/cache/v8"
	"githum.com/mengdage/gochat/pkg/api"
)

type Storage interface {
	CreateUser(ctx context.Context, req api.NewUserRequest) (*api.User, error)
	GetUserByID(ctx context.Context, userID int) (*api.User, error)
	GetUserByName(ctx context.Context, userName string) (*api.User, error)
	GetAllUsers(ctx context.Context) ([]*api.User, error)
	ClearAllUsersCache(ctx context.Context) error
	GetUserServer(ctx context.Context, userName string) (string, error)
	SaveUserServer(ctx context.Context, userName string, serverAddr string) error
	SaveChatMessage(ctx context.Context, chatMessage *api.ChatMessage) error
	GetConversation(ctx context.Context, cID string, user string) (*api.ChatConversaion, error)
	GetConversationMessages(ctx context.Context, cID string) ([]*api.ChatMessage, error)
	CreateConversation(ctx context.Context, cID string, users ...string) (*api.ChatConversaion, error)
}

type storage struct {
	db    DBStorage
	cache *cache.Cache
}

func NewStorage(db DBStorage, cache *cache.Cache) Storage {
	return &storage{
		db:    db,
		cache: cache,
	}
}

func (s *storage) CreateUser(ctx context.Context, req api.NewUserRequest) (*api.User, error) {
	return s.db.CreateUser(ctx, req)
}
func (s *storage) GetUserByName(ctx context.Context, userName string) (*api.User, error) {
	log.Printf("[GetUserByName] userName %s", userName)

	key := fmt.Sprintf("user:%s", userName)
	user := new(api.User)
	err := s.cache.Once(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: user,
		TTL:   15 * time.Minute,
		Do: func(i *cache.Item) (interface{}, error) {
			log.Printf("%v", i)
			u, err := s.db.GetUserByName(ctx, userName)
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	})
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (s *storage) GetAllUsers(ctx context.Context) ([]*api.User, error) {
	log.Println("[GetAllUsers]")

	key := "user:all_users"
	users := make([]*api.User, 0)
	err := s.cache.Once(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: &users,
		TTL:   15 * time.Minute,
		Do: func(i *cache.Item) (interface{}, error) {
			u, err := s.db.GetAllUsers(ctx)
			if err != nil {
				log.Printf("Error while getting all users from db: %v", err)
				return nil, err
			}
			return u, nil
		},
	})
	if err != nil {
		fmt.Printf("Error while getting all users from storage: %v", err)
		return nil, err
	}

	return users, nil
}

func (s *storage) GetUserByID(ctx context.Context, userID int) (*api.User, error) {
	log.Printf("[GetUserByID] userID %d", userID)

	key := fmt.Sprintf("user:%d", userID)
	user := new(api.User)
	err := s.cache.Once(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: user,
		TTL:   15 * time.Minute,
		Do: func(i *cache.Item) (interface{}, error) {
			u, err := s.db.GetUserByID(ctx, userID)
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *storage) ClearAllUsersCache(ctx context.Context) error {
	key := "user:all_users"
	err := s.cache.Delete(ctx, key)
	if err != nil {
		log.Printf("Error while clearing cache %s: %v\n", key, err)
		return err
	}

	return nil
}

func (s *storage) SaveUserServer(ctx context.Context, userName string, serverAddr string) error {
	key := fmt.Sprintf("user:%s:server", userName)

	err := s.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: serverAddr,
		TTL:   2 * time.Minute,
	})

	if err != nil {
		log.Printf("Error while saving user server info")
		return err
	}

	return nil
}

func (s *storage) GetUserServer(ctx context.Context, userName string) (string, error) {
	key := fmt.Sprintf("user:%s:server", userName)
	var serverAddr string

	err := s.cache.Get(ctx, key, &serverAddr)

	if err != nil {
		log.Printf("Error while saving user server info")
		return "", err
	}

	return serverAddr, nil
}
func (s *storage) SaveChatMessage(ctx context.Context, chatMessage *api.ChatMessage) error {
	err := s.db.SaveChatMessage(ctx, chatMessage)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("conversation:%s:messages", chatMessage.ConversationID)
	err = s.cache.Delete(ctx, cacheKey)
	if err != nil {
		log.Printf("Error while invalidating conversation messages: %v", err)
		return nil
	}

	return nil
}

func (s *storage) GetConversation(ctx context.Context, cID string, user string) (*api.ChatConversaion, error) {
	conv := api.ChatConversaion{}
	key := fmt.Sprintf("conversation:%s:user:%s", cID, user)

	log.Printf("[GetConversation] cache key: %s\n", key)
	err := s.cache.Once(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: &conv,
		TTL:   1 * time.Hour,
		Do: func(i *cache.Item) (interface{}, error) {
			log.Println("[GetConversation] Get convesation from db")
			u, err := s.db.GetConversation(ctx, cID, user)

			if err != nil {
				return nil, err
			}
			return u, nil
		},
	})
	if err != nil {
		return nil, err
	}

	if conv.ID == "" {
		return nil, nil
	}

	return &conv, nil
}

func (s *storage) GetConversationMessages(ctx context.Context, cID string) ([]*api.ChatMessage, error) {
	msgs := []*api.ChatMessage{}
	key := fmt.Sprintf("conversation:%s:messages", cID)
	err := s.cache.Once(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: &msgs,
		TTL:   1 * time.Hour,
		Do: func(i *cache.Item) (interface{}, error) {
			u, err := s.db.GetConversationMessages(ctx, cID)
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	})
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (s *storage) CreateConversation(ctx context.Context, cID string, users ...string) (*api.ChatConversaion, error) {
	cc, err := s.db.CreateConversation(ctx, cID, users...)

	return cc, err
}
