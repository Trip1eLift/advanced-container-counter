package container_counter_system

import (
	"context"
	"fmt"
	"os"
	"time"

	ttlqueue "github.com/Trip1eLift/container-counter/cluster/container_counter_system/TTL_queue"
	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/broadcast"
	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/model"
	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/queue"
	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/utils"
)

const baseTime = 35 * time.Second

type Manager struct {
	self_id         string
	containers      ttlqueue.Client
	counter_channel string
	bc              *broadcast.Client
}

var manager Manager

func init() {
	manager.self_id = os.Getenv("CONTAINER_ID")
	manager.containers = *ttlqueue.New(baseTime)
	manager.counter_channel = utils.ReadClusterId() + "-container-counter"
	manager.bc = broadcast.New(manager.counter_channel, context.Background())
	manager.bc.Subscribe(context.Background())

	// subscribe to queue
	go func() {
		for {
			pack := queue.Pop_package() // only return when new pack arrives

			if pack.Container_id == manager.self_id {
				// I do not retrieve my own id from redis.
				continue
			}

			if manager.containers.UpdateContainer(pack.Container_id) {
				publish_enrollment()

				count := manager.containers.GetLength() + 1
				fmt.Printf("Update container count: %d\n", count)
			}

			go func() {
				if manager.containers.CleanupOnExpire() {
					count := manager.containers.GetLength() + 1
					fmt.Printf("Update container count: %d\n", count)
				}
			}()
		}
	}()
}

// Expected to be triggered by container health check every 30s
func OnHealth() {
	publish_enrollment()
}

func FirstPublish() {
	go func() {
		time.Sleep(1 * time.Second)
		publish_enrollment()
	}()
}

func publish_enrollment() {
	pack := model.Package{
		Container_id: manager.self_id,
	}
	manager.bc.Publish(context.Background(), pack)
}

func GetCount() int {
	count := manager.containers.GetLength() + 1
	fmt.Printf("container count: %d\n", count)
	return count
}
