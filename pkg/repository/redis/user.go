package redis

import (
	"encoding/json"
	"fmt"
	log "github.com/mj23978/chat-backend-x/logger/zerolog"
	"github.com/mj23978/chat-backend/pkg/domain/user"
	m "github.com/mj23978/chat-backend/pkg/models"
	"github.com/pkg/errors"
	"time"
)

type userCacheRepo struct {
	client *Redis
}

func UserRedisRepo(redis *Redis) (user.UserCacheRepo, error) {
	repo := &userCacheRepo{}
	//client, err := userRedisClient(redisURL)
	//if err != nil {
	//	return nil, errors.Wrap(err, "cacheRepo.NewRedisRepository")
	//}
	repo.client = redis
	return repo, nil
}

func (g *userCacheRepo) generateKey(code string) string {
	return fmt.Sprintf("user:%s", code)
}

func (g userCacheRepo) Create(user *m.User) error {
	key := g.generateKey(user.Username)
	data := map[string]interface{}{
		"display_name": user.DisplayName,
		"last_active":  user.LastActive,
		"status":       int32(user.Status),
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

func (g *userCacheRepo) Delete(id string, fields ...string) error {
	res, err := g.client.HDel(g.generateKey(id), fields...).Result()
	if err != nil {
		return errors.Wrap(err, "cacheRepo.User.Delete")
	}
	log.Infof("Redis Delete %v", res)
	return nil
}

func (g *userCacheRepo) Update(id string, fields map[string]interface{}) error {
	res, err := g.client.HMSet(g.generateKey(id), fields).Result()
	if err != nil {
		return errors.Wrap(err, "cacheRepo.User.Update")
	}
	log.Infof("Redis Update %v", res)
	return nil
}

func (g *userCacheRepo) Get(id string, fields ...string) (*m.User, error) {
	res, err := g.client.HMGet(g.generateKey(id), fields...).Result()
	if err != nil {
		return nil, errors.Wrap(err, "cacheRepo.User.Get")
	}
	input := make(map[string]interface{})
	for i := 0; i < len(fields); i++ {
		input[fields[i]] = res[i]
	}
	user := &m.User{}
	rawMsg, err2 := json.Marshal(input)
	if err2 != nil {
		return nil, errors.Wrap(err, "cacheRepo.User.Get")
	}
	if e := json.Unmarshal(rawMsg, &user); e != nil {
		return nil, errors.Wrap(e, "cacheRepo.User.Get")
	}
	log.Infof("Redis Update %v", res)
	return user, nil
}
