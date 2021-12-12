package chat

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	m "github.com/mj23978/chat-backend/pkg/models"
	pb "github.com/mj23978/chat-backend/pkg/proto"
)

type chatHandler struct {
	pb.UnimplementedChatServiceServer
	chatService ChatService
}

func NewChatHandler(chatService ChatService) *chatHandler {
	return &chatHandler{
		chatService: chatService,
	}
}

func (h *chatHandler) GetChat(ctx context.Context, req *pb.GetChatRequest) (*pb.Chat, error) {
	chat, err := h.chatService.Get(req.GetId())
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, status.Error(codes.NotFound, "Cannot Find Any Chat with this ID")
	}
	res, err2 := turnChatToChatGrpc(chat)
	if err2 != nil {
		log.Fatal().Msgf("err : %v", err2)
	}
	return res, nil
}

func turnChatToChatGrpc(chat *m.Chat) (*pb.Chat, error) {
	chatRes := &pb.Chat{}
	if err := copier.Copy(chatRes, chat); err != nil {
		return nil, err
	}
	//chatRes.Id = chat.ID
	return chatRes, nil
}

func turnChatGrpcToChat(chat *pb.Chat) (*m.Chat, error) {
	chatRes := &m.Chat{}
	if err := copier.Copy(chatRes, chat); err != nil {
		return nil, err
	}
	//chatRes.ID = chat.Id
	return chatRes, nil
}
