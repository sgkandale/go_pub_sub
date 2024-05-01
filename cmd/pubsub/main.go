package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"pubsub/config"
	"pubsub/pkg/pubsub"
	"pubsub/server"
)

var (
	configFileLocation = flag.String("config", "config.yaml", "config file location")
)

func main() {
	flag.Parse()
	cfg := config.Read(*configFileLocation)
	globalCtx, cancelGlobalCtx := context.WithCancel(context.Background())
	globalWaitgroup := sync.WaitGroup{}
	activeServerCount := 0

	// pubsub
	pubsubHub := pubsub.New()

	if cfg.GrcpConfig.Use {
		grpcServer := server.NewGrpc(
			globalCtx,
			&server.NewServerRequest{
				Cfg:    &cfg.GrcpConfig,
				PubSub: pubsubHub,
			},
		)
		if grpcServer != nil {
			globalWaitgroup.Add(1)
			go func() {
				defer globalWaitgroup.Done()
				activeServerCount++
				err := grpcServer.Serve()
				activeServerCount--
				if err != nil {
					cancelGlobalCtx()
					log.Printf("[ERROR] starting grpc server : %s", err.Error())
				}
			}()
			go func() {
				<-globalCtx.Done()
				grpcServer.Stop()
			}()
		}
	}

	// monitor active server count
	go func() {
		iterateCount := 0
		for {
			time.Sleep(time.Millisecond * 500 * time.Duration(iterateCount))
			if activeServerCount == 0 {
				cancelGlobalCtx()
			}
		}
	}()

	go func() {
		// Wait for Control C to exit
		ch := make(chan os.Signal, 10)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		// Block until a signal is received
		signalCounter := 0
		for osSignal := range ch {
			signalCounter++
			log.Printf("[WARN] received os kill signal '%s'", osSignal.String())
			if signalCounter > 1 {
				log.Printf("[WARN] received os kill signal attempt : %d, killing program", signalCounter)
				os.Exit(2)
			}
			cancelGlobalCtx()
		}
	}()

	globalWaitgroup.Wait()
	log.Printf("[WARN] stopping pubsub")
	<-globalCtx.Done()
}
