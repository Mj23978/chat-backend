package main

import (
	"github.com/mj23978/chat-backend/pkg/domain/auth"
	"github.com/mj23978/chat-backend/pkg/domain/chat"
	"github.com/mj23978/chat-backend/pkg/domain/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"

	log "github.com/mj23978/chat-backend-x/logger/zerolog"
	"github.com/mj23978/chat-backend/pkg/conf"
	pb "github.com/mj23978/chat-backend/pkg/proto"
)

func main() {
	log.Init(conf.Log.Level)
	log.Infof("--- Starting ChatServer Node ---")
	listener, err := net.Listen("tcp", conf.Global.Port)
	if err != nil {
		panic(err)
	}
	jwtManager := auth.NewJWTManager("Secret_Key", 15*time.Minute, "Refresh_Secret_Key", 14*24*time.Hour)
	authInter := auth.NewAuthInterceptor(jwtManager)
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(authInter.Unary()),
		grpc.StreamInterceptor(authInter.Stream()),
	)
	userServer, authServer := UserServer(jwtManager, conf.Cache.URL, conf.Database.URL, conf.Broker.URL)
	chatServer := ChatServer(conf.Cache.URL, conf.Database.URL, conf.Broker.URL)
	pb.RegisterChatServiceServer(srv, chatServer)
	pb.RegisterUserServiceServer(srv, userServer)
	pb.RegisterAuthServiceServer(srv, authServer)
	reflection.Register(srv)
	log.Infof("Server Listening at Port %v", conf.Global.Port)
	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}

func UserServer(jwtManager *auth.JWTManager, redisUrl, cassandraUrl, natsUrl string) (pb.UserServiceServer, pb.AuthServiceServer) {
	//redisRepo, err := redis.GameRedisRepo(redisUrl)
	//if err := utils.ErrorWrap(err, "game.init"); err != nil {
	//	utils.ErrorFatal(err)
	//}
	//cassandraRepo, err := cassandra.GameMongoRepository(cassandraUrl, "games", 60)
	//if err := utils.ErrorWrap(err2, "game.init"); err != nil {
	//	utils.ErrorFatal(err)
	//}
	userService := user.NewUserService(cassandraRepo, redisRepo)
	userServer := user.NewUserHandler(userService)
	authServer := auth.NewAuthHandler(userService, jwtManager)
	return userServer, authServer
}

func ChatServer(redisUrl, cassandraUrl, natsUrl string) pb.ChatServiceServer {
	//redisRepo, err := redis.GameRedisRepo(redisUrl)
	//if err := utils.ErrorWrap(err, "game.init"); err != nil {
	//	utils.ErrorFatal(err)
	//}
	//cassandraRepo, err := cassandra.GameMongoRepository(cassandraUrl, "games", 60)
	//if err := utils.ErrorWrap(err2, "game.init"); err != nil {
	//	utils.ErrorFatal(err)
	//}
	chatService := chat.NewChatService(redisRepo, cassandraRepo, pubsub)
	chatServer := chat.NewChatHandler(chatService)
	return chatServer
}
