package actions

import (
	"context"
	"time"
)

// RunDeletion is a job that calls DeleteUnprocessed every 10 minutes
func (a actions) RunDeletion() {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)

			select {
			case <-ctx.Done():
				a.storage.DeleteUnprocessed()
				cancel()

			case <-a.StopDeletion:
				cancel()
				return
			}
		}
	}()
}
