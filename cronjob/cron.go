package cronjob

import (
	"context"
	"github.com/shengwenjin/auxiliary/cronjob/grpool"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	wg  sync.WaitGroup
	ctx context.Context
)

type Cronjob struct {
	RunFunc func()
	Timeout time.Duration
}

func RunCronJobs(cronjobs ...Cronjob) {
	defer func() {
		_ = grpool.Release()
	}()
	for _, cronjob := range cronjobs {
		wg.Add(1)
		runJob := cronjob.RunFunc
		err := grpool.Submit(func() {
			defer wg.Done()
			runWithTimeout(runJob, cronjob.Timeout)
		})
		if err != nil {
			log.Errorf("submit job error: %v", err)
		}
	}
	wg.Wait()
}

func runWithTimeout(runFunc func(), duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	// create a done channel to tell the request it's done
	done := make(chan struct{})
	go func() {
		defer func() {
			if c := recover(); c != nil {
				log.Errorf("response request panic: %v", c)
			}
			close(done)
		}()
		runFunc()
	}()
	// non-blocking select on two channels see if the request
	// times out or finishes
	select {
	// if the context is done it timed out or was canceled
	// so don't return anything
	case <-ctx.Done():
		log.Errorf("handler timeout")
		return
	// if the request finished then finish the request by
	// writing the response
	case <-done:
		return
	}
}
