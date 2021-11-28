package main

import (
	"context"
	"crypto/rand"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	// context to indicate about service shutdown
	exitctx context.Context
	exitfn  context.CancelFunc
	// wait group for all service goroutines
	exitwg sync.WaitGroup

	hashcount uint64
)

const hashpool = 64

// Command line parameters.
var (
	ncpu = runtime.NumCPU()
	nthr = flag.Int("n", ncpu, "number of threads to start")
	pdur = flag.Duration("d", 90*time.Minute, "duration of program working (in format '1d8h15m30s')")
	blen = flag.Int("b", 1024, "length of random bytes block to calculate for each hash")
	halg = flag.String("a", "sha256", "hash or signature algorithm, can be: md5, sha1, sha224, sha256, sha384, sha512, sha512/224, sha512/256, ecdsa, ed25519")
)

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

	msg = make([]byte, *blen)
	if _, err := rand.Read(msg); err != nil {
		panic(err)
	}

	for {
		for i := 0; i < hashpool; i++ {
			alg()
		}
		atomic.AddUint64(&hashcount, hashpool)

		// check on exit signal during running
		select {
		case <-exitctx.Done():
			return
		default:
		}
	}
}

// Init performs global data initialisation.
func Init() {
	log.Println("starts")

	flag.Parse()

	log.Printf("execution duration is %s", pdur.String())

	// create context and wait the break
	exitctx, exitfn = context.WithTimeout(context.Background(), *pdur)
	go func() {
		// Make exit signal on function exit.
		defer exitfn()

		var sigint = make(chan os.Signal, 1)
		var sigterm = make(chan os.Signal, 1)
		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM (Ctrl+/)
		// SIGKILL, SIGQUIT will not be caught.
		signal.Notify(sigint, syscall.SIGINT)
		signal.Notify(sigterm, syscall.SIGTERM)
		// Block until we receive our signal.
		select {
		case <-exitctx.Done():
			if errors.Is(exitctx.Err(), context.DeadlineExceeded) {
				log.Println("shutting down by timeout")
			} else if errors.Is(exitctx.Err(), context.Canceled) {
				log.Println("shutting down by cancel")
			} else {
				log.Printf("shutting down by %s", exitctx.Err().Error())
			}
		case <-sigint:
			log.Println("shutting down by break")
		case <-sigterm:
			log.Println("shutting down by process termination")
		}
		signal.Stop(sigint)
		signal.Stop(sigterm)
	}()

	DetectAlg()
}

// Run launches working threads.
func Run() {
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
}

// Done performs graceful network shutdown,
// waits until all server threads will be stopped.
func Done() {
	var start = time.Now()
	// wait for exit signal
	<-exitctx.Done()
	// wait until all server threads will be stopped.
	exitwg.Wait()
	var rundur = time.Since(start)
	log.Printf("calculated %d entities for message of %d bytes\n", hashcount, *blen)
	log.Printf("average speed %4.f entities per second", float64(hashcount)/float64(rundur)*float64(time.Second))
	log.Println("shutting down complete.")
}

func main() {
	Init()
	Run()
	Done()
}
