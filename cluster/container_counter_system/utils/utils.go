package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"
)

func RandTimeS(timeRangeSecond int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(1000 * timeRangeSecond)
	return time.Duration(num) * time.Millisecond
}

func RandSleepS(timeRangeSecond int) {
	time.Sleep(RandTimeS(timeRangeSecond))
}

func ReadClusterId() string {
	path := os.Getenv("ECS_CONTAINER_METADATA_FILE")
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("unable to read file: %v", err))
	}
	var body map[string]interface{}
	if err := json.Unmarshal(data, &body); err != nil {
		panic(fmt.Errorf("unable to decode json file: %v", err))
	}
	cluster_name := body["Cluster"].(string) // cluster-prod-e_x
	x_stack, _ := regexp.MatchString("-e_x|-w_x", cluster_name)
	y_stack, _ := regexp.MatchString("-e_y|-w_y", cluster_name)
	var cluster_id string
	if x_stack {
		cluster_id = cluster_name[:len(cluster_name)-4] + "-x"
	} else if y_stack {
		cluster_id = cluster_name[:len(cluster_name)-4] + "-y"
	} else {
		panic(fmt.Errorf("cluster does not match x or y stack pattern: %v", err))
	}
	return cluster_id
}
