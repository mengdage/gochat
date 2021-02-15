package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"githum.com/mengdage/gochat/pkg/api"
)

type DBStorage interface {
	CreateUser(ctx context.Context, req api.NewUserRequest) (*api.User, error)
	GetUserByID(ctx context.Context, userID int) (*api.User, error)
	GetUserByName(ctx context.Context, userName string) (*api.User, error)
	GetAllUsers(ctx context.Context) ([]*api.User, error)
	SaveChatMessage(ctx context.Context, chatMessage *api.ChatMessage) error
	GetConversation(ctx context.Context, cID string, user string) (*api.ChatConversaion, error)
	GetConversationMessages(ctx context.Context, cID string) ([]*api.ChatMessage, error)
	CreateConversation(ctx context.Context, cID string, users ...string) (*api.ChatConversaion, error)
}

type dbStorage struct {
	db *pgxpool.Pool
}

func NewDBStorage(pool *pgxpool.Pool) DBStorage {
	return &dbStorage{
		db: pool,
	}
}

func (d *dbStorage) CreateUser(ctx context.Context, req api.NewUserRequest) (*api.User, error) {
	user := api.User{}
	err := d.db.QueryRow(ctx, `
	INSERT INTO chat_user(name)
		VALUES ($1)
	RETURNING id, name;
	`, req.Name).Scan(&user.ID, &user.Name)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, fmt.Errorf("User name %s already exists", req.Name)
			}
		}

		return nil, err
	}

	return &user, nil
}

func (d *dbStorage) GetUserByID(ctx context.Context, userID int) (*api.User, error) {
	user := api.User{}
	err := d.db.QueryRow(ctx, `
	SELECT id, name
	FROM chat_user
	WHERE id = $1;
	`, userID).Scan(&user.ID, &user.Name)

	if err != nil {
		log.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("The user %d does not exist", userID)

			return nil, fmt.Errorf("The user does not exist")
		}
		return nil, err
	}

	return &user, nil
}

func (d *dbStorage) GetUserByName(ctx context.Context, userName string) (*api.User, error) {
	user := api.User{}
	err := d.db.QueryRow(ctx, `
	SELECT id, name
	FROM chat_user
	WHERE name = $1;
	`, userName).Scan(&user.ID, &user.Name)

	if err != nil {
		log.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("The user %s does not exist", userName)

			return nil, fmt.Errorf("The user does not exist")
		}
		return nil, err
	}

	return &user, nil
}

func (d *dbStorage) GetAllUsers(ctx context.Context) ([]*api.User, error) {
	rows, err := d.db.Query(ctx, `
	SELECT id, name
	FROM chat_user;
	`)

	if err != nil {
		log.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No user exists")

			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	users := make([]*api.User, 0)
	for rows.Next() {
		user := api.User{}
		rows.Scan(&user.ID, &user.Name)
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error while reading users from DB: %v", err)
		return nil, err
	}

	return users, nil
}

func (d *dbStorage) SaveChatMessage(ctx context.Context, chatMessage *api.ChatMessage) error {
	c, err := d.db.Exec(ctx, `
	INSERT INTO chat_message(conversation_id, from_user, to_user, content, created_at) VALUES
	($1, $2, $3, $4, $5)
	`, chatMessage.ConversationID, chatMessage.FromUser, chatMessage.ToUser, chatMessage.Content, chatMessage.CreatedAt)
	if err != nil {
		log.Printf("Error saving chat message to db: %v", err)
		return err
	}

	log.Printf("Successfully save chat message to db: %s", c.String())
	return nil
}

func (d *dbStorage) GetConversation(ctx context.Context, cID string, user string) (*api.ChatConversaion, error) {
	conv := api.ChatConversaion{}
	err := d.db.QueryRow(ctx, `
	SELECT c.id, c.created_at
	FROM user_conversation uc JOIN conversation c ON uc.conversation_id = c.id
	WHERE uc.user_name = $1 AND uc.conversation_id = $2
	`, user, cID).Scan(&conv.ID, &conv.CreateAt)

	if err != nil {
		log.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("The conversation %s does not exist", cID)

			return nil, nil
		}
		return nil, err
	}

	return &conv, nil
}

func (d *dbStorage) GetConversationMessages(ctx context.Context, cID string) ([]*api.ChatMessage, error) {
	rows, err := d.db.Query(ctx, `
	SELECT conversation_id, from_user, to_user, content, created_at
	FROM chat_message
	WHERE conversation_id = $1;
	`, cID)
	if err != nil {
		log.Printf("Error while querying the db: %v\n", err)
		return nil, err
	}
	// .Scan(&conv.ID, &conv.CreateAt)
	chatMsgs := []*api.ChatMessage{}
	for rows.Next() {
		cm := &api.ChatMessage{}

		rows.Scan(&cm.ConversationID, &cm.FromUser, &cm.ToUser, &cm.Content, &cm.CreatedAt)
		chatMsgs = append(chatMsgs, cm)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error while reading data: %v", err)
		return nil, err
	}

	return chatMsgs, nil
}

func (d *dbStorage) CreateConversation(ctx context.Context, cID string, users ...string) (*api.ChatConversaion, error) {
	createdAt := time.Now()
	tx, err := d.db.Begin(ctx)
	if err != nil {
		log.Printf("Error while starting a transaction: %v", err)
		return nil, err
	}

	defer tx.Rollback(ctx)

	c, err := tx.Exec(ctx, `
	INSERT INTO conversation(id, created_at) VALUES
	($1, $2);
	`, cID, createdAt)
	if err != nil {
		log.Printf("Error saving conversation %s to db: %v", cID, err)
		return nil, err
	}
	log.Printf("Create conversation: %s", c)

	insertUserConvQueryParts := []string{"INSERT INTO user_conversation(conversation_id, user_name) VALUES"}
	vals := []interface{}{cID}

	for i, user := range users {
		var t string

		if i == len(users)-1 {
			t = fmt.Sprintf("($1, $%d)", i+2)
		} else {
			t = fmt.Sprintf("($1, $%d),", i+2)

		}
		insertUserConvQueryParts = append(insertUserConvQueryParts, t)
		vals = append(vals, user)
	}
	insertUserConvQueryParts = append(insertUserConvQueryParts, ";")
	insertUserConvQuery := strings.Join(insertUserConvQueryParts, "\n")

	log.Printf("Insert user conversations query:\n %s", insertUserConvQuery)

	c, err = tx.Exec(ctx, insertUserConvQuery, vals...)
	if err != nil {
		log.Printf("Error saving user conversations to db: %v", err)
		return nil, err
	}
	log.Printf("Create user conversations: %s", c)

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Error while committing changes: %v", err)
		return nil, err

	}

	log.Printf("Successfully create conversation in db")

	cc := &api.ChatConversaion{
		ID:       cID,
		CreateAt: createdAt,
	}

	return cc, nil
}
