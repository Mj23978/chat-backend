package main

import (
	_ "net/http/pprof"

	nprotoo "github.com/mj23978/chat-backend-x/broker/nats"
	discovery "github.com/mj23978/chat-backend-x/discovery/etcd"
	log "github.com/mj23978/chat-backend-x/logger/zerolog"
	conf "github.com/mj23978/chat-backend/pkg/conf"
)

func main() {
	log.Init(conf.Log.Level)
	log.Infof("--- Starting ChatServer Node ---")

	serviceNode := discovery.NewServiceNode(conf.Registry.Addrs, conf.Global.Dc)
	serviceNode.RegisterNode("chat-server", "node-chat-server")
	//game_server.Init(conf.Global.Dc, serviceNode, conf.Broker.URL)

	serviceWatcher := discovery.NewServiceWatcher(conf.Registry.Addrs, conf.Global.Dc)
	go serviceWatcher.WatchServiceNode("chat", NotImp)
	protoo := nprotoo.NewNatsProtoo(conf.Broker.URL)
	protoo.OnRequest(serviceNode.GetRPCChannel(), func(request nprotoo.Request, accept nprotoo.RespondFunc, reject nprotoo.RejectFunc) {
		log.Infof("method => %s, data => %v", request.Method, request.Data)
		if request.Method == "offer" {
			accept("We Accept Your Offer")
		}
		reject(404, "Not found")
	})
	node := gameNode.NewGameNode(conf.Global.Dc, conf.Broker.URL, conf.Registry.Addrs)
	node.Init()
	// serviceNode := discovery.NewServiceNode(conf.Registry.Addrs, conf.Global.Dc)
	// serviceNode.RegisterNode("game", "node-game")

	// serviceWatcher := discovery.NewServiceWatcher(conf.Registry.Addrs, conf.Global.Dc)
	// go serviceWatcher.WatchServiceNode("game-server", NotImp)
	// protoo := nprotoo.NewNatsProtoo(conf.Broker.URL)

	listener, err := net.Listen("tcp", conf.Global.Port)
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	pb.RegisterAppServiceServer(srv, node.Server)
	reflection.Register(srv)
	log.Infof("Server Listening at Port %v", conf.Global.Port)
	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}

func NotImp(service string, state discovery.NodeStateType, nodes discovery.Node) {
	if state == 0 {
		log.Debugf("%v is Up", service)
	} else if state == 1 {
		log.Debugf("%v is Down", service)
	}
	log.Debugf("our Nodes : ", nodes)
}
