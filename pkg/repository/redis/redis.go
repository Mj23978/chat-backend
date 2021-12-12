package redis

import (
	"context"
	"fmt"
	"time"

	log "github.com/mj23978/chat-backend-x/logger/zerolog"

	"github.com/go-redis/redis/v8"
)

var (
	lockExpire = 3 * time.Second
)

type Config struct {
	Addrs []string `mapstructure:"addrs"`
	Pwd   string   `mapstructure:"password"`
	DB    int      `mapstructure:"db"`
}

type Redis struct {
	cluster     *redis.ClusterClient
	single      *redis.Client
	clusterMode bool
	ctx         context.Context
}

func NewRedis(c Config) *Redis {
	if len(c.Addrs) == 0 {
		log.Errorf("invalid addrs: %v", c.Addrs)
		return nil
	}

	r := &Redis{}
	r.ctx = context.Background()
	if len(c.Addrs) == 1 {
		r.single = redis.NewClient(
			&redis.Options{
				Addr:         c.Addrs[0], // use default Addr
				Password:     c.Pwd,      // no password set
				DB:           c.DB,       // use default DB
				DialTimeout:  3 * time.Second,
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 5 * time.Second,
			})
		if err := r.single.Ping(r.ctx).Err(); err != nil {
			log.Errorf(err.Error())
			return nil
		}
		r.single.Do(r.ctx, "CONFIG", "SET", "notify-keyspace-events", "AKE")
		r.clusterMode = false
		log.Infof("redis new client single mode: %v", r)
		return r
	}

	r.cluster = redis.NewClusterClient(
		&redis.ClusterOptions{
			Addrs:        c.Addrs,
			Password:     c.Pwd,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		})
	if err := r.cluster.Ping(r.ctx).Err(); err != nil {
		log.Errorf(err.Error())
	}
	r.cluster.Do(r.ctx, "CONFIG", "SET", "notify-keyspace-events", "AKE")
	r.clusterMode = true
	log.Infof("redis new client cluster mode: %v", r)
	return r
}

func (r *Redis) Close() {
	if r.single != nil {
		r.single.Close()
	}
	if r.cluster != nil {
		r.cluster.Close()
	}
}

func (r *Redis) Set(key string, value interface{}, t time.Duration) error {
	r.Acquire(key)
	defer r.Release(key)
	if r.clusterMode {
		return r.cluster.Set(r.ctx, key, value, t).Err()
	}
	return r.single.Set(r.ctx, key, value, t).Err()
}

func (r *Redis) Get(k string) string {
	r.Acquire(k)
	defer r.Release(k)
	if r.clusterMode {
		return r.cluster.Get(r.ctx, k).Val()
	}
	return r.single.Get(r.ctx, k).Val()
}

func (r *Redis) HSet(key, field string, value interface{}) error {
	r.Acquire(key)
	defer r.Release(key)
	if r.clusterMode {
		return r.cluster.HSet(r.ctx, key, field, value).Err()
	}
	return r.single.HSet(r.ctx, key, field, value).Err()
}

func (r *Redis) HGet(key, field string) *redis.StringCmd {
	r.Acquire(key)
	defer r.Release(key)
	if r.clusterMode {
		return r.cluster.HGet(r.ctx, key, field)
	}
	return r.single.HGet(r.ctx, key, field)
}

func (r *Redis) HMSet(key string, values ...interface{}) *redis.BoolCmd {
	r.Acquire(key)
	defer r.Release(key)
	if r.clusterMode {
		return r.cluster.HMSet(r.ctx, key, values)
	}
	return r.single.HMSet(r.ctx, key, values)
}

func (r *Redis) HMGet(key string, fields ...string) *redis.SliceCmd {
	r.Acquire(key)
	defer r.Release(key)
	if r.clusterMode {
		return r.cluster.HMGet(r.ctx, key, fields...)
	}
	return r.single.HMGet(r.ctx, key, fields...)
}

func (r *Redis) HGetAll(key string) map[string]string {
	r.Acquire(key)
	defer r.Release(key)
	if r.clusterMode {
		return r.cluster.HGetAll(r.ctx, key).Val()
	}
	return r.single.HGetAll(r.ctx, key).Val()
}

func (r *Redis) HDel(key string, fields ...string) *redis.IntCmd {
	r.Acquire(key)
	defer r.Release(key)
	if r.clusterMode {
		return r.cluster.HDel(r.ctx, key, fields...)
	}
	return r.single.HDel(r.ctx, key, fields...)
}

func (r *Redis) Incr(key string) *redis.IntCmd {
	r.Acquire(key)
	defer r.Release(key)
	if r.clusterMode {
		return r.cluster.Incr(r.ctx, key)
	}
	return r.single.Incr(r.ctx, key)
}

func (r *Redis) Expire(key string, t time.Duration) error {
	r.Acquire(key)
	defer r.Release(key)
	if r.clusterMode {
		return r.cluster.Expire(r.ctx, key, t).Err()
	}

	return r.single.Expire(r.ctx, key, t).Err()
}

func (r *Redis) HSetTTL(t time.Duration, key, field string, value interface{}) error {
	r.Acquire(key)
	defer r.Release(key)
	if r.clusterMode {
		if err := r.cluster.HSet(r.ctx, key, field, value).Err(); err != nil {
			return err
		}
		return r.cluster.Expire(r.ctx, key, t).Err()
	}
	if err := r.single.HSet(r.ctx, key, field, value).Err(); err != nil {
		return err
	}
	return r.single.Expire(r.ctx, key, t).Err()
}

func (r *Redis) HMSetTTL(t time.Duration, k string, values ...interface{}) error {
	r.Acquire(k)
	defer r.Release(k)
	if r.clusterMode {
		if err := r.cluster.HMSet(r.ctx, k, values...).Err(); err != nil {
			return err
		}
		return r.cluster.Expire(r.ctx, k, t).Err()
	}

	if err := r.single.HMSet(r.ctx, k, values...).Err(); err != nil {
		return err
	}
	return r.single.Expire(r.ctx, k, t).Err()
}

func (r *Redis) Keys(key string) []string {
	r.Acquire(key)
	defer r.Release(key)
	if r.clusterMode {
		return r.cluster.Keys(r.ctx, key).Val()
	}
	return r.single.Keys(r.ctx, key).Val()
}

func (r *Redis) Del(k string) *redis.IntCmd {
	r.Acquire(k)
	defer r.Release(k)
	if r.clusterMode {
		return r.cluster.Del(r.ctx, k)
	}
	return r.single.Del(r.ctx, k)
}

func (r *Redis) TxPipeline() redis.Pipeliner {
	if r.clusterMode {
		return r.cluster.TxPipeline()
	}
	return r.single.TxPipeline()
}

// Watch http://redisdoc.com/topic/notification.html
func (r *Redis) Watch(ctx context.Context, key string) <-chan interface{} {
	r.Acquire(key)
	defer r.Release(key)
	var pubsub *redis.PubSub
	if r.clusterMode {
		pubsub = r.cluster.PSubscribe(r.ctx, fmt.Sprintf("__key*__:%s", key))
	} else {
		pubsub = r.single.PSubscribe(r.ctx, fmt.Sprintf("__key*__:%s", key))
	}

	res := make(chan interface{})
	go func() {
		for {
			select {
			case msg := <-pubsub.Channel():
				op := msg.Payload
				log.Infof("key => %s, op => %s", key, op)
				res <- op
			case <-ctx.Done():
				pubsub.Close()
				close(res)
				return
			}
		}
	}()

	return res
}

func (r *Redis) lock(key string) bool {
	// Tips: use ("lock-"+key) as lock key is better than (key+"-lock")
	// this avoid "keys /xxxxx/*" to get this lock
	if r.clusterMode {
		ok, _ := r.cluster.SetNX(r.ctx, "lock-"+key, 1, lockExpire).Result()
		return ok
	}
	ok, _ := r.single.SetNX(r.ctx, "lock-"+key, 1, lockExpire).Result()
	return ok
}

func (r *Redis) unlock(key string) {
	if r.clusterMode {
		_, _ = r.cluster.Del(r.ctx, "lock-"+key).Result()
	}
	_, _ = r.single.Del(r.ctx, "lock-"+key).Result()
}

// Acquire a destributed lock
func (r *Redis) Acquire(key string) {
	// retry if lock failed
	for !r.lock(key) {
		time.Sleep(time.Millisecond * 100)
	}
}

// Release a destributed lock
func (r *Redis) Release(key string) {
	r.unlock(key)
}
