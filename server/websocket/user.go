package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"githum.com/mengdage/gochat/chatpb"
	"githum.com/mengdage/gochat/lib/cache"
	"githum.com/mengdage/gochat/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// SendUserMessage sends a message to a user
func SendUserMessage(fromUser *Client, toUserID string, message string) {
	toClient := clientManager.GetUserClient(toUserID)

	if toClient != nil {
		toClient.SendMessage([]byte(message))
		return
	}

	userOnline, err := cache.GetUserOnlineInfo(toUserID)
	if err != nil {
		log.Printf("Failed to get user online info: %v\n", err)
		res := model.NewResponse(string(SendType), http.StatusInternalServerError, fmt.Sprintf("User %s is not available", toUserID), nil)
		bs, _ := json.Marshal(res)
		fromUser.SendMessage(bs)
	}

	rpcAddr := fmt.Sprintf("%s:%s", userOnline.ServerIp, userOnline.RpcPort)

	log.Printf("Dial rpc addr: %s\n", rpcAddr)
	conn, err := grpc.Dial(rpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect: %v", err)
	}

	defer conn.Close()
	client := chatpb.NewChatServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := client.SendMessage(ctx, &chatpb.SendMessageRequest{
		SenderId:   fromUser.UserID,
		ReceiverId: toUserID,
		Content:    message,
	})

	if err != nil {
		errStatus, _ := status.FromError(err)
		log.Printf("Failed to send message through grpc: %v %v\n", errStatus.Code(), errStatus.Message())
		return
	}

	log.Printf("The message was sent to %s: %v %s", toUserID, resp.GetCode(), resp.GetMsg())
}

func SendUserMessageLocal(toUserID string, message string) error {
	toClient := clientManager.GetUserClient(toUserID)

	if toClient == nil {
		return fmt.Errorf("User %s does not exist", toUserID)
	}

	toClient.SendMessage([]byte(message))

	return nil
}
