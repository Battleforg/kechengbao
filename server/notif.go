package main

import (
	"encoding/json"
	"gopkg.in/redis.v3"
)

type Notif struct {
	Date string
	Text string
}

const (
	DEFAULT_NUM = 10
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	//client.FlushDb()
}

func AppendMessage(cid, msg string) error {
	return client.LPush(cid, msg).Err()
}

func RetrieveMessage(cid string) ([]byte, error) {
	res, err := client.LRange(cid, 0, DEFAULT_NUM-1).Result()
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return b, nil
}
