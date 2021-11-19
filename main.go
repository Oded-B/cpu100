package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var (
	// context to indicate about service shutdown
	exitctx context.Context
	exitfn  context.CancelFunc
	// wait group for all service goroutines
	exitwg sync.WaitGroup
)

var (
	ncpu = runtime.NumCPU()
	nthr = flag.Int("n", ncpu, "number of threads to start")
	pdur = flag.String("d", "1h30m", "duration of program working (in format '1d8h15m30s')")
)

// WaitInterrupt returns shutdown signal was recivied and cancels some context.
func WaitInterrupt(cancel context.CancelFunc) {
	// Make exit signal on function exit.
	defer cancel()

	var sigint = make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM (Ctrl+/)
	// SIGKILL, SIGQUIT will not be caught.
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Block until we receive our signal.
	<-sigint
	log.Println("shutting down by break")
}

// Loader makes loading on any one CPU core.
func Loader(ithr int) {
	defer exitwg.Done()
	var err error
	defer func() {
		if err != nil {
			log.Printf("%d thread failed with error %s\n", ithr, err.Error())
		}
	}()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var buf = make([]byte, 1024)
	var h = sha256.New()
	for {
		for i := 0; i < 16; i++ {
			if _, err = rand.Read(buf); err != nil {
				return
			}
			if _, err = h.Write(buf); err != nil {
				return
			}
		}

		// check on exit signal during running
		select {
		case <-exitctx.Done():
			return
		default:
		}
	}
}

// Run launches working threads.
func Run() {
	flag.Parse()

	var err error
	var dur time.Duration
	if dur, err = time.ParseDuration(*pdur); err != nil {
		log.Fatalln(err)
		return
	}
	log.Printf("execution duration is %s", dur.String())

	// create context and wait the break
	exitctx, exitfn = context.WithTimeout(context.Background(), dur)
	go WaitInterrupt(exitfn)

	if *nthr > ncpu {
		log.Printf("recieved %d threads to start, it can be used maximum %d by number of CPU cores", *nthr, ncpu)
		*nthr = ncpu
	}
	log.Printf("runs %d threads", *nthr)
	for i := 0; i < *nthr; i++ {
		var i = i // localize
		exitwg.Add(1)
		go Loader(i)
	}

	exitwg.Wait()
}

func main() {
	log.Println("starts")
	Run()
	log.Println("done")
}
