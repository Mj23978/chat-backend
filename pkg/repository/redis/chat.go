package redis

import (
	"encoding/json"
	"fmt"
	log "github.com/mj23978/chat-backend-x/logger/zerolog"
	"github.com/mj23978/chat-backend/pkg/domain/chat"
	m "github.com/mj23978/chat-backend/pkg/models"
	"github.com/pkg/errors"
	"time"
)

type chatCacheRepo struct {
	client *Redis
}

func ChatRedisRepo(redis *Redis) (chat.ChatCacheRepo, error) {
	repo := &chatCacheRepo{}
	repo.client = redis
	return repo, nil
}

func (g *chatCacheRepo) GetAll(id string) (*m.Chat, error) {
	res, err := g.client.HMGet(g.generateKey(id)).Result()
	if err != nil {
		return nil, errors.Wrap(err, "cacheRepo.Chat.Get")
	}
	chat := &m.Chat{}
	log.Infof("Redis Update %v", res)
	return chat, nil
}

func (g *chatCacheRepo) Incr(id string, field uint8) (int64, error) {
	res, err := g.client.Incr(id).Result()
	if err != nil {
		return 0, errors.Wrap(err, "cacheRepo.Chat.Create")
	}
	return res, nil
}

func (g *chatCacheRepo) Publish(id string, message interface{}) error {
	return nil
}

func (g *chatCacheRepo) generateKey(code string) string {
	return fmt.Sprintf("chat:%s", code)
}

func (g *chatCacheRepo) Create(chat *m.Chat) error {
	key := g.generateKey(chat.Name)
	data := map[string]interface{}{
		"image_url": chat.ImageUrl,
	}
	pipe := g.client.TxPipeline()
	res, err := pipe.HMSet(g.client.ctx, key, data).Result()
	if err != nil {
		return errors.Wrap(err, "cacheRepo.User.Create")
	}
	if _, err := pipe.Expire(g.client.ctx, key, 4*time.Hour).Result(); err != nil {
		return errors.Wrap(err, "repository.Game.Start")
	}
	if _, err := pipe.Exec(g.client.ctx); err != nil {
		return errors.Wrap(err, "cacheRepo")
	}
	log.Infof("Redis Create %v", res)
	return nil
}

func (g *chatCacheRepo) Delete(id string, fields ...string) error {
	res, err := g.client.HDel(g.generateKey(id), fields...).Result()
	if err != nil {
		return errors.Wrap(err, "cacheRepo.Chat.Delete")
	}
	log.Infof("Redis Delete %v", res)
	return nil
}

func (g *chatCacheRepo) Update(id string, fields map[string]interface{}) error {
	res, err := g.client.HMSet(g.generateKey(id), fields).Result()
	if err != nil {
		return errors.Wrap(err, "cacheRepo.Chat.Update")
	}
	log.Infof("Redis Update %v", res)
	return nil
}

func (g *chatCacheRepo) Get(id string, fields ...string) (*m.Chat, error) {
	res, err := g.client.HMGet(g.generateKey(id), fields...).Result()
	if err != nil {
		return nil, errors.Wrap(err, "cacheRepo.Chat.Get")
	}
	input := make(map[string]interface{})
	for i := 0; i < len(fields); i++ {
		input[fields[i]] = res[i]
	}
	chat := &m.Chat{}
	rawMsg, err2 := json.Marshal(input)
	if err2 != nil {
		return nil, errors.Wrap(err, "cacheRepo.Chat.Get")
	}
	if e := json.Unmarshal(rawMsg, &chat); e != nil {
		return nil, errors.Wrap(e, "cacheRepo.Chat.Get")
	}
	log.Infof("Redis Update %v", res)
	return chat, nil
}
