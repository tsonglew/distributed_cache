package clients

// import "github.com/go-redis/redis"

// type redisClient struct {
// 	*redis.Client
// }

// func (r *redisClient) get(key string) (string, error) {
// 	resp, err := r.Get(key).Result()
// 	if err == redis.Nil {
// 		return "", nil
// 	}
// 	return resp, err
// }

// func (r *redisClient) set(key, value string) error {
// 	return r.Set(key, value, 0).Err()
// }

// func (r *redisClient) del(key string) error {
// 	return r.Del(key).Err()
// }

// func (r *redisClient) Run(c *Cmd) {
// 	switch c.Name {
// 	case "get":
// 		c.Value, c.Error = r.get(c.Key)
// 	case "set":
// 		c.Error = r.set(c.Key, c.Value)
// 	case "del":
// 		c.Error = r.del(c.Key)
// 	default:
// 		panic("unkown cmd name " + c.Name)
// 	}
// }

// func (r *redisClient) PipelineRun(cmds []*Cmd) {
// 	if len(cmds) == 0 {
// 		return
// 	}
// 	pipe := r.Pipeline()
// 	cmders := make([]redis.Cmder, len(cmds))
// 	for i, c := range cmds {
// 		switch c.Name {
// 		case "get":
// 			cmders[i] = pipe.Get(c.Key)
// 		case "set":
// 			cmders[i] = pipe.Set(c.Key, c.Value, 0)
// 		case "del":
// 			cmders[i] = pipe.Del(c.Key)
// 		default:
// 			panic("unkown cmd name " + c.Name)
// 		}
// 	}
// 	_, err := pipe.Exec()
// 	if err != nil && err != redis.Nil {
// 		panic(err)
// 	}
// 	for i, c := range cmds {
// 		if c.Name == "get" {
// 			v, err := cmders[i].(*redis.StringCmd).Result()
// 			if err == redis.Nil {
// 				v, err = "", nil
// 			}
// 			c.Value, c.Error = v, err
// 		} else {
// 			c.Error = cmders[i].Err()
// 		}
// 	}
// }

// func newRedisClient(server string) *redisClient {
// 	return &redisClient{redis.NewClient(&redis.Options{Addr: server + ":6379", ReadTimeout: -1})}
// }
