package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"githum.com/mengdage/gochat/chatpb"
	"githum.com/mengdage/gochat/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageOperatorStorage interface {
	GetConversation(ctx context.Context, cID string, currentUser string) (*ChatConversaion, error)
	GetUserServer(ctx context.Context, userName string) (string, error)
	CreateConversation(ctx context.Context, cID string, users ...string) (*ChatConversaion, error)
	SaveChatMessage(ctx context.Context, chatMessage *ChatMessage) error
}

type MessageOperator struct {
	clientManager *ClientManager
	remoteConns   map[string]*grpc.ClientConn
	storage       MessageOperatorStorage
}

func NewMessageOperator(clientManager *ClientManager, storage MessageOperatorStorage) *MessageOperator {
	return &MessageOperator{
		clientManager: clientManager,
		remoteConns:   make(map[string]*grpc.ClientConn),
		storage:       storage,
	}
}

func (o *MessageOperator) SendMessageLocal(toUser string, msg *WSMessage) error {
	msgBs, _ := json.Marshal(msg)

	client := o.clientManager.GetClient(toUser)
	if client == nil {
		return fmt.Errorf("The user %s does not exist", toUser)
	}

	client.Send <- msgBs

	return nil
}

func (o *MessageOperator) SendMessageRemote(toUser string, msg *ChatMessage) error {
	addr, err := o.storage.GetUserServer(context.Background(), toUser)
	if err != nil {
		return err
	}

	conn, err := o.dialServer(addr)
	if err != nil {
		return err
	}

	client := chatpb.NewChatServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	resp, err := client.SendMessage(ctx, &chatpb.SendMessageRequest{
		ConversationId: msg.ConversationID,
		SenderId:       msg.FromUser,
		ReceiverId:     msg.ToUser,
		Content:        msg.Content,
		CreatedAt:      timestamppb.New(msg.CreatedAt),
	})

	if err != nil {
		errStatus, _ := status.FromError(err)
		log.Printf("Failed to send message through grpc: %v %v\n", errStatus.Code(), errStatus.Message())
		return err
	}

	log.Printf("The message was sent to %s: %v %s", msg.ToUser, resp.GetCode(), resp.GetMsg())
	return nil

}

func (o *MessageOperator) dialServer(addr string) (*grpc.ClientConn, error) {
	if conn, ok := o.remoteConns[addr]; ok {
		if conn.GetState() == connectivity.Ready {
			return conn, nil
		}
		conn.Close()
	}

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return nil, err
	}

	return conn, nil
}

func (o *MessageOperator) ProcessMessage(msg []byte, fromClient *Client) error {
	wsMsg := WSMessage{}
	if err := json.Unmarshal(msg, &wsMsg); err != nil {
		return err
	}

	if wsMsg.Cmd == "send" {
		body := WSSendBody{}
		if err := json.Unmarshal(wsMsg.Body, &body); err != nil {
			return err
		}

		toUser := body.UserName
		cID := helper.CreateConversationID(fromClient.UserName, body.UserName)
		conversation, err := o.storage.GetConversation(context.Background(), cID, fromClient.UserName)
		if err != nil {
			return err
		}

		if conversation == nil {
			conversation, err = o.storage.CreateConversation(context.Background(), cID, fromClient.UserName, body.UserName)
			if err != nil {
				return err
			}
		}

		chatMessage := &ChatMessage{
			ConversationID: conversation.ID,
			FromUser:       fromClient.UserName,
			ToUser:         body.UserName,
			Content:        body.Content,
			CreatedAt:      time.Now(),
		}

		if err := o.storage.SaveChatMessage(context.Background(), chatMessage); err != nil {
			log.Printf("Failed to store chat message in the storage: %v", err)
			return err
		}
		chatMessageBs, _ := json.Marshal(chatMessage)

		recvMsg := &WSMessage{
			Cmd:  "recv",
			Body: chatMessageBs,
		}

		o.SendMessageLocal(fromClient.UserName, recvMsg)

		if o.clientManager.Exists(toUser) {

			o.SendMessageLocal(toUser, recvMsg)
		} else {
			o.SendMessageRemote(toUser, chatMessage)
		}

	}

	return nil
}
