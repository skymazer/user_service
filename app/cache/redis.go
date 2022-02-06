package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/skymazer/user_service/loggerfx"
)

const (
	lifetime = 60
	host     = "cache"
	port     = "6379"
)

type Redis struct {
	Conn redis.Conn
	log  *loggerfx.Logger
}

var ErrNotFound = errors.New("no matching record")

func New(log *loggerfx.Logger) (*Redis, error) {
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, err
	}

	r := Redis{conn, log}
	return &r, nil
}

func (s *Redis) Set(key string, val []byte) error {
	s.log.Debug("SET", key, val, "EX", lifetime)
	_, err := s.Conn.Do("SET", key, val, "EX", lifetime)

	return err
}

func (s *Redis) Get(key string) ([]byte, error) {
	s.log.Debug("GET", key)
	r, err := redis.Bytes(s.Conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return []byte{}, ErrNotFound
		}

		return []byte{}, err
	}

	return r, nil
}

func (s *Redis) Remove(key string) error {
	s.log.Debug("DEL", key)
	_, err := s.Conn.Do("DEL", key)
	if err == redis.ErrNil {
		return ErrNotFound
	}

	return err
}
