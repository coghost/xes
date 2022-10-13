package xes

import (
	"context"
	"errors"
	"net/http"
	"syscall"
	"time"

	"github.com/olivere/elastic/v7"
)

// MyRetrier this is copied from https://github.com/olivere/elastic/wiki/Retrier-and-Backoff
type MyRetrier struct {
	backoff elastic.Backoff
	retries int
}

func newMyRetrier(retries int) *MyRetrier {
	return &MyRetrier{
		backoff: elastic.NewExponentialBackoff(10*time.Millisecond, 8*time.Second),
		retries: retries,
	}
}

func (r *MyRetrier) Retry(ctx context.Context, retry int, req *http.Request, resp *http.Response, err error) (time.Duration, bool, error) {
	// Fail hard on a specific error
	if err == syscall.ECONNREFUSED {
		return 0, false, errors.New("elasticsearch or network down")
	}

	// Stop after 5 retries
	if retry >= r.retries {
		return 0, false, nil
	}

	// Let the backoff strategy decide how long to wait and whether to stop
	wait, stop := r.backoff.Next(retry)
	return wait, stop, nil
}
