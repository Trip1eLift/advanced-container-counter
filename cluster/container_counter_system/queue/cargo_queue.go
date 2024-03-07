package queue

import (
	"sync"

	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/model"
)

type package_queue struct {
	boxes []model.Package
	mu    sync.RWMutex

	on_wait    bool
	on_wait_mu sync.RWMutex
	wait       sync.WaitGroup

	mu_level2 sync.RWMutex
}

var Packages package_queue

func init() {
	Packages.boxes = make([]model.Package, 0)
}

func Push_package(pack model.Package) {
	Packages.mu.Lock()
	Packages.boxes = append(Packages.boxes, pack)
	Packages.mu.Unlock()

	// TODO: Consider optimize on_wait using channel
	Packages.on_wait_mu.Lock()
	if Packages.on_wait == true {
		Packages.on_wait = false
		Packages.wait.Done()
	}
	Packages.on_wait_mu.Unlock()
}

// Assume Pop_package is only called in one place
func Pop_package() model.Package {
	// prevent Pop_package being called multiple times at once
	// Without level2 lock, when 2 Pop_package are waiting and a package arrive,
	// they can rush to stage 2 together and break the code.
	//Packages.mu_level2.Lock() // cause sever delay issue...

	// Stage 1: Waiting room
	for {
		Packages.mu.RLock()
		if len(Packages.boxes) > 0 {
			Packages.mu.RUnlock()
			break
		} else {
			Packages.mu.RUnlock()

			// TODO: Consider optimize on_wait using channel
			Packages.on_wait_mu.Lock()
			Packages.wait.Add(1)
			Packages.on_wait = true
			Packages.on_wait_mu.Unlock()

			Packages.wait.Wait()
		}
	}

	// Stage 2: Pop queue's head
	Packages.mu.Lock()
	pack := Packages.boxes[0]
	Packages.boxes = Packages.boxes[1:]
	Packages.mu.Unlock()
	//Packages.mu_level2.Unlock()

	return pack
}
