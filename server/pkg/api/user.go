package api

import (
	"context"
	"errors"
	"log"
)

type UserService interface {
	New(context.Context, NewUserRequest) (*User, error)
	Login(context.Context, LoginRequest) (*User, error)
	GetUserByID(ctx context.Context, userID int) (*User, error)
	GetUserByName(context.Context, string) (*User, error)
	GetAllUsers(ctx context.Context, currentUserName string) ([]*User, error)
}

type UserStorage interface {
	CreateUser(context.Context, NewUserRequest) (*User, error)
	GetUserByID(context.Context, int) (*User, error)
	GetUserByName(context.Context, string) (*User, error)
	GetAllUsers(context.Context) ([]*User, error)
	ClearAllUsersCache(ctx context.Context) error
}

type UserCache interface {
	SaveUserInfo(ctx context.Context, user *User) error
}

type userService struct {
	storage UserStorage
}

func NewUserService(storage UserStorage) UserService {
	return &userService{
		storage: storage,
	}
}

func (u *userService) New(ctx context.Context, req NewUserRequest) (*User, error) {
	if req.Name == "" {
		return nil, errors.New("Name cannot be empty")
	}

	user, err := u.storage.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := u.storage.ClearAllUsersCache(ctx); err != nil {
		log.Println("A new user is created but the alluser cache is not cleared")
	}

	return user, nil
}

func (u *userService) Login(ctx context.Context, req LoginRequest) (*User, error) {
	user, err := u.storage.GetUserByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) GetUserByID(ctx context.Context, userID int) (*User, error) {
	user, err := u.storage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) GetUserByName(ctx context.Context, userName string) (*User, error) {
	user, err := u.storage.GetUserByName(ctx, userName)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) GetAllUsers(ctx context.Context, currentUserName string) ([]*User, error) {
	log.Printf("[GetAllUsers] currentUserName %s\n", currentUserName)
	users, err := u.storage.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	excludeCurrentUser := make([]*User, 0, len(users)-1)
	for _, user := range users {
		if user.Name == currentUserName {
			continue
		}

		excludeCurrentUser = append(excludeCurrentUser, user)
	}

	return excludeCurrentUser, nil
}
