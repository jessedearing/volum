package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	var metrics = map[int]int{}
	var sizemetrics = map[int]int{}
	mut := sync.Mutex{}
	var closer = make(chan os.Signal)
	signal.Notify(closer, os.Interrupt)

	go func() {
		var bytevol int
	infloop:
		for {
			var bs = make([]byte, 512, 512)
			numbytes, err := os.Stdin.Read(bs)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
				close(closer)
				return
			}

			for i := 0; i < numbytes; i++ {
				if bs[i] == '\n' {
					bytevol += i
					mut.Lock()
					secondbucket := time.Now().Second()
					metrics[secondbucket]++
					sizemetrics[secondbucket] += bytevol
					bytevol = 0
					mut.Unlock()
					continue infloop
				}
			}
			bytevol += numbytes
		}
	}()

	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()

		for linet := range t.C {
			prevbucket := time.Now().Second() - 1
			mut.Lock()
			linepersec := metrics[prevbucket]
			bytespersec := sizemetrics[prevbucket]
			mut.Unlock()
			fmt.Printf("%s - %d lines/sec %d kilobytes/sec\n", linet.String(), linepersec, bytespersec/1000)
		}
	}()

	<-closer
}
