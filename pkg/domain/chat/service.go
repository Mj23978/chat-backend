package chat

import (
	"github.com/go-redis/redis"
	"github.com/mj23978/chat-backend-x/utils"
	m "github.com/mj23978/chat-backend/pkg/models"
	errs "github.com/pkg/errors"
	"gopkg.in/dealancer/validate.v2"
)

type ChatService interface {
	FindById(id int) (*m.Chat, error)
	Get(id string) (*m.Chat, error)
	Create(chat *m.Chat) error
	Update(id string, fields map[string]interface{}) error
	PubSub(channels ...string) (*redis.PubSub, error)
	Request() error
	Accept() error
	Block() error
	Hide() error
	Unhide() error
	Lock() error
	Unlock() error
	Listen() error
	Filter() error
}

type chatService struct {
	cacheRepo ChatCacheRepo
	dbRepo    ChatDBRepo
}

func NewChatService(cacheRepo ChatCacheRepo, dbRepo ChatDBRepo) ChatService {
	return &chatService{
		cacheRepo,
		dbRepo,
	}
}

func (svc *chatService) Request() error {
	return nil
}

func (svc *chatService) Accept() error {
	return nil
}

func (svc *chatService) Block() error {
	return nil
}

func (svc *chatService) Hide() error {
	return nil
}

func (svc *chatService) Unhide() error {
	return nil
}

func (svc *chatService) Lock() error {
	return nil
}

func (svc *chatService) Unlock() error {
	return nil
}

func (svc *chatService) Listen() error {
	return nil
}

func (svc *chatService) Filter() error {
	return nil
}

func (svc *chatService) Create(chat *m.Chat) error {
	if e := validate.Validate(chat); e != nil {
		return errs.Wrap(e, "service.Chat.Store")
	}
	// chat.toml.Password = repo.EncryptPassword(data.Password)
	return svc.dbRepo.Create(chat)
}

func (svc *chatService) FindById(id int) (*m.Chat, error) {
	filter := map[string]interface{}{"id": id}
	res, e := svc.dbRepo.FindBy(filter)
	if e != nil {
		return res, errs.Wrap(e, "service.News.GetById")
	}
	return res, nil
}

func (svc *chatService) PubSub(channels ...string) (*redis.PubSub, error) {
	return svc.cacheRepo.PubSub(channels...)
}

func (svc *chatService) Get(id string) (*m.Chat, error) {
	res, err := svc.cacheRepo.GetAll(id)
	if err := utils.ErrorWrap(err, "Chat.service.Get"); err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *chatService) Update(id string, fields map[string]interface{}) error {
	if err := svc.cacheRepo.Update(id, fields); err != nil {
		return utils.ErrorWrap(err, "Chat.service.Get")
	}
	return nil
}

// func (svc *chatService) Start(username, lobbyName string, challengeTime, playerTime int, roles []m.MafiaRole) (*m.Game, error) {
// 	lobbyRes := &m.Lobby{}
// 	svc.pubsub.SendToChannel("lobby", "get", lobbyName)
// 	lobbyRaw, err := svc.pubsub.ListenResChan(fmt.Sprintf("lobby-get-%v", lobbyName))
// 	if err := utils.MapToStruct(lobbyRes, "lobby", lobbyRaw); err != nil {
// 		return nil, errs.Wrap(err, "GameService.Start.service")
// 	}
// 	if err != nil {
// 		if err.Error() == "We Couldn't Find Any Result" {
// 			return nil, err
// 		}
// 		return nil, errs.Wrap(err, "GameService.Start.cacheRepo")
// 	}
// 	if username != lobbyRes.Admin {
// 		return nil, errs.Errorf("only admin (%v) can start the chat.toml you (%v) are not the admin", lobbyRes.Admin, username)
// 	}
// 	if len(lobbyRes.Members) != len(roles) {
// 		return nil, errs.Errorf("your lobby has %v members but you provide %v roles", len(lobbyRes.Members), len(roles))
// 	}
// 	rolesRes, t, err2 := player.TurnMafiaRoleToRole(roles)
// 	if err2 != nil {
// 		return nil, err2
// 	}
// 	game := &m.Game{}
// 	game.ID = lobbyName
// 	game.PlayerTime = playerTime
// 	game.Stage = "Day"
// 	game.ChallengeTime = challengeTime
// 	rand.Seed(time.Now().UnixNano())
// 	rand.Shuffle(len(lobbyRes.Members), func(i, j int) { lobbyRes.Members[i], lobbyRes.Members[j] = lobbyRes.Members[j], lobbyRes.Members[i] })
// 	rand.Shuffle(len(rolesRes), func(i, j int) { rolesRes[i], rolesRes[j] = rolesRes[j], rolesRes[i] })
// 	for index, v := range lobbyRes.Members {
// 		playerRes := &m.Player{}
// 		playerRes.ID = fmt.Sprintf("%v-%v", lobbyName, v.Username)
// 		playerRes.MafiaRole = rolesRes[index].Role
// 		playerRes.Number = uint8(index)
// 		playerRes.IsTurn = false
// 		playerRes.Injured = false
// 		playerRes.IsDead = false
// 		game.Players = append(game.Players, playerRes.ID)
// 		svc.pubsub.SendToChannel("player", "create", playerRes)
// 		_, err := svc.pubsub.ListenResChan(fmt.Sprintf("player-create-%v", playerRes.ID))
// 		if err != nil {
// 			return nil, errs.Wrap(err, "GameService.Start.service")
// 		}
// 	}
// 	if err := svc.cacheRepo.Create(game); err != nil {
// 		return nil, err
// 	}
// 	if err := svc.dbRepo.Create(game); err != nil {
// 		return nil, err
// 	}
// 	svc.pubsub.SendToChannel("lobby", "delete", lobbyRes.Name)
// 	_, err3 := svc.pubsub.ListenResChan(fmt.Sprintf("lobby-delete-%v", lobbyRes.Name))
// 	if err3 != nil {
// 		return nil, errs.Wrap(err3, "GameService.Start.service")
// 	}
// 	fmt.Println(rolesRes, t)
// 	fmt.Println(lobbyRes)
// 	return game, nil
// }
