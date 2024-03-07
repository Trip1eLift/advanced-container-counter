package broadcast

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/model"
	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/queue"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack"
)

type Broadcast interface {
	Subscribe(ctx context.Context)
	Publish(ctx context.Context, v interface{})
	Close() error
}

type Client struct {
	Broadcast
	Redis   *redis.Client
	channel string
}

func New(channel_name string, ctx context.Context) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       atoi(os.Getenv("REDIS_DB")),
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		client.Close()
		panic(fmt.Errorf("Redis connection err: %v", err))
	} else {
		return &Client{Redis: client, channel: channel_name}
	}
}

func (c *Client) Subscribe(ctx context.Context) {
	go func() {
		pubsub := c.Redis.Subscribe(ctx, c.channel)
		defer pubsub.Close()

		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				panic(fmt.Errorf("pubsub receive error: %v\n", err))
			}

			// Assume everything from broadcast channel is model.Package
			pack := model.Package{}
			err = msgpack.Unmarshal([]byte(msg.Payload), &pack)
			if err != nil {
				panic(fmt.Errorf("msg unpacking error: %v\n", err))
			}
			queue.Push_package(pack)
		}
	}()
}

func (c *Client) Publish(ctx context.Context, pack model.Package) {
	payload, err := msgpack.Marshal(pack)
	if err != nil {
		panic(fmt.Errorf("msg packing error: %v\n", err))
	}

	// redis can handle concurrent publish as long as the events aren't too many
	// utils.RandSleepS(1) // rand sleep reduce chance of redis race condition in local

	err = c.Redis.Publish(ctx, c.channel, payload).Err()
	if err != nil {
		panic(fmt.Errorf("pubsub publish error: %v\n", err))
	}

	return
}

func (c *Client) Close() error {
	return c.Redis.Close()
}

func atoi(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}
