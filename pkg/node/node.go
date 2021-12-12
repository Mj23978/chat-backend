package gameNode

import (
	"fmt"

	nprotoo "github.com/mj23978/chat-backend-x/broker/nats"
	discovery "github.com/mj23978/chat-backend-x/discovery/etcd"
	log "github.com/mj23978/chat-backend-x/logger/zerolog"
	pb "github.com/mj23978/chat-backend/pkg/proto"
	ep "github.com/mj23978/chat-backend/tests"
)

type gameNode struct {
	Server        pb.AppServiceServer
	ServiceNode   *discovery.ServiceNode
	dc            string
	ServiceName   string
	natsURL       string
	etcdAddrs     []string
	emitters      map[string]*ep.Emitter
	groups        map[string]*ep.Group
	protoo        *nprotoo.NatsProtoo
	Broadcaster   *nprotoo.Broadcaster
	services      map[string]discovery.Node
	Requestor     *nprotoo.Requestor
	WatchServices []string
}

func NewGameNode(dc, natsURL string, etcdAddr []string) *gameNode {
	res := &gameNode{
		dc:        dc,
		emitters:  make(map[string]*ep.Emitter),
		groups:    make(map[string]*ep.Group),
		natsURL:   natsURL,
		etcdAddrs: etcdAddr,
	}

	//res.Server = testServer(res.newEmitter("api"))

	return res
}

// Init func
func (gn *gameNode) Init() {
	gn.ServiceName = "game"
	gn.registerNode()
	gn.protoo = nprotoo.NewNatsProtoo(gn.natsURL)
	gn.Broadcaster = gn.protoo.NewBroadcaster(gn.ServiceNode.GetEventChannel())
	gn.Requestor = gn.protoo.NewRequestor(gn.ServiceNode.GetRPCChannel())
	gn.WatchServices = []string{"game-chat"}
	gn.services = make(map[string]discovery.Node)
	gn.handleRequest()
	gn.handleEmitters()
}

func (gn *gameNode) registerNode() {
	gn.ServiceNode = discovery.NewServiceNode(gn.etcdAddrs, gn.dc)
	gn.ServiceNode.RegisterNode(gn.ServiceName, fmt.Sprintf("%s-node-0", gn.ServiceName))

	for _, service := range gn.WatchServices {
		serviceWatcher := discovery.NewServiceWatcher(gn.etcdAddrs, gn.dc)
		go serviceWatcher.WatchServiceNode(service, func(serviceName string, state discovery.NodeStateType, node discovery.Node) {
			if state == discovery.UP {
				log.Infof("Service UP [%s] => %v", serviceName, node)
				gn.services[serviceName] = node
			} else if state == discovery.DOWN {
				log.Infof("Service DOWN [%s] => %v", serviceName, node)
				delete(gn.services, serviceName)
			}
		})
	}
}

//func testServer(em *ep.Emitter) pb.AppServiceServer {
//	gameServer := gameCore.NewTestHandler(em)
//	return gameServer
//}
