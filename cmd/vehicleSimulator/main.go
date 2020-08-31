package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/vehicleSimulator/vehicle"
)

func main() {
	runtime.GOMAXPROCS(6)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call: %+v \n", oscall)
		cancel()
	}()

	err := startService(ctx)
	if err != nil {
		cancel()
		panic(err)
	}
}

func startService(ctx context.Context) error {
	vehicle.Start(ctx)

	return nil

}
