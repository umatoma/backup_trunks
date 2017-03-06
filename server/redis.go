package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/redsync.v1"
)

type RedisClient struct {
	host     string
	password string
	db       int
	pool     *redis.Pool
	// If set, path to a socket file overrides hostname
	socketPath string
	redsync    *redsync.Redsync
}

func NewRedisClient(host, password, socketPath string, db int) *RedisClient {
	return &RedisClient{
		host: host,
		db: db,
		password: password,
		socketPath: socketPath,
	}
}

func (cli *RedisClient) GetAllTasks() ([]interface{}, error) {
	conn := cli.open()
	defer conn.Close()

	keysReply, err := redis.Values(conn.Do("KEYS", "task_*"))
	if err != nil {
		return nil, err
	}

	strsReply, err := redis.Values(conn.Do("MGET", keysReply...))
	if err != nil {
		return nil, err
	}

	return strsReply, nil
}

// Returns / creates instance of Redis connection
func (cli *RedisClient) open() redis.Conn {
	if cli.pool == nil {
		cli.pool = cli.newPool()
	}

	if cli.redsync == nil {
		var pools = []redsync.Pool{cli.pool}
		cli.redsync = redsync.New(pools)
	}

	return cli.pool.Get()
}

// Returns a new pool of Redis connections
func (cli *RedisClient) newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			var (
				c redis.Conn
				err error
				opts = make([]redis.DialOption, 0)
			)

			if cli.password != "" {
				opts = append(opts, redis.DialPassword(cli.password))
			}

			if cli.socketPath != "" {
				c, err = redis.Dial("unix", cli.socketPath, opts...)
			} else {
				c, err = redis.Dial("tcp", cli.host, opts...)
			}

			if cli.db != 0 {
				_, err = c.Do("SELECT", cli.db)
			}

			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
