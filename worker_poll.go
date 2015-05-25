package main

import (
	"sync"
)

type WorkerCallback func(s string) (*Robots)

func worker(linkChan chan string, wg *sync.WaitGroup, cb WorkerCallback) {
   defer wg.Done()

   for url := range linkChan {
     cb(url)
   }
}

func master(collection []string, workerFn WorkerCallback, workers int) {
	lCh := make(chan string)
    wg := new(sync.WaitGroup)

    // Adding routines to workgroup and running then
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go worker(lCh, wg, workerFn)
    }

    // Processing all links by spreading them to `free` goroutines
    for _, link := range collection {
        lCh <- link
    }

    // Closing channel (waiting in goroutines won't continue any more)
    close(lCh)

    // Waiting for all goroutines to finish (otherwise they die as main routine dies)
    wg.Wait()
}