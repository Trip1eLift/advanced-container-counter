package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Trip1eLift/container-counter/cluster/container_counter_system"
)

const ip = "0.0.0.0"
const port = "8000"

func main() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			container_counter_system.GetCount()
		}
	}()

	http.HandleFunc("/health", func(write http.ResponseWriter, request *http.Request) {
		container_counter_system.OnHealth()
		fmt.Fprintf(write, "Healthy golang server.\n")
	})

	log.Printf("Listening on %s:%s\n", ip, port)
	container_counter_system.FirstPublish()
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil))
}
