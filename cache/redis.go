package cache

import (
	"log"

	"github.com/keimoon/gore"
)

type Redis struct {
	pool *gore.Pool
}

func NewCache(connectionString, password string) (redis *Redis, err error) {
	pool := &gore.Pool{
		InitialConn: 5,
		MaximumConn: 10,
		Password:    password,
	}
	err = pool.Dial(connectionString)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	redis = &Redis{pool: pool}
	return
}

func (r Redis) ClosePool() {
	r.pool.Close()
}

func (r Redis) Get(key string) (string, error) {
	conn, _ := r.pool.Acquire()
	defer conn.Close()

	reply, err := gore.NewCommand("GET", key).Run(conn)
	parsedRep, _ := reply.String()

	return parsedRep, err
}

func (r Redis) Set(key, value string) (bool, error) {
	conn, _ := r.pool.Acquire()
	defer conn.Close()

	reply, err := gore.NewCommand("SET", key, value).Run(conn)
	parsedRep, _ := reply.Bool()

	return parsedRep, err
}
