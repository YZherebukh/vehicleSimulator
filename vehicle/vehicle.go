package vehicle

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/vehicleSimulator/vehicle/data"
)

const (
	min = 1.000001
	max = 1.000005
)

// Vehicle is a vehicle struct
type Vehicle struct {
	ID  string
	Lat float64
	Lon float64
}

// new creates new vehicle
func new(vid string) Vehicle {
	return Vehicle{
		ID:  vid,
		Lat: rand.Float64(),
		Lon: rand.Float64(),
	}
}

// Start creates and moves vehicles
// and sends updates about their position to fleetState
func Start(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(data.VidList))

	for i := range data.VidList {
		v := new(data.VidList[i])

		go func(ctx context.Context, v Vehicle, wg *sync.WaitGroup) {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				case <-time.Tick(time.Second):
					v.move()
					v.sendPosition()
				}
			}
		}(ctx, v, &wg)
	}
	wg.Wait()
}

// Move change vehicle geolocation
func (v *Vehicle) move() {
	v.Lat = min + rand.Float64()*(max-min)
	v.Lon = min + rand.Float64()*(max-min)
}

// sendPosition sends position update to FleetState
func (v *Vehicle) sendPosition() error {
	var req struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	req.Lat = v.Lat
	req.Lon = v.Lon

	b, err := json.Marshal(req)
	if err != nil {
		log.Printf("failed to marshal request %s\n", err)
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, "http://fleet_state:8080/v1/vehicle/"+v.ID, bytes.NewBuffer(b))
	if err != nil {
		log.Printf("failed to create request %s\n", err)
		return err
	}

	c := http.Client{Timeout: time.Second * 10}
	_, err = c.Do(httpReq)
	return err
}
