package grpc_server

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/spf13/viper"
	"githum.com/mengdage/gochat/chatpb"
	"githum.com/mengdage/gochat/server/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	chatpb.UnimplementedChatServiceServer
}

func (s *server) SendMessage(c context.Context, req *chatpb.SendMessageRequest) (*chatpb.SendMessageResponse, error) {
	userID := req.GetReceiverId()
	content := req.GetContent()
	if err := websocket.SendUserMessageLocal(userID, content); err != nil {
		pberr := status.Errorf(codes.InvalidArgument, "Failed to send message: "+err.Error())
		return nil, pberr
	}

	resp := &chatpb.SendMessageResponse{
		Code: http.StatusOK,
		Msg:  "Success",
	}

	return resp, nil
}

// Start starts the grpc server
func Start() {
	rpcPort := viper.GetViper().GetString("app.rpc_port")

	lis, err := net.Listen("tcp", ":"+rpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	chatpb.RegisterChatServiceServer(s, &server{})
	log.Printf("grpc server listening on :%s", rpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
