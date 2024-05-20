package main

import (
	"log"
	"time"
)

type Retry struct {
	timeout         time.Duration
	max_timeout     time.Duration
	max_retry_count int
	retry_count     int

	cause_func func(error)

	timeout_strat func(count int, timeout time.Duration) time.Duration

	retry_strat func(count int, max_count int) bool
}

func NewRetry(max_retry_count int, max_timeout time.Duration) Retry {
	return Retry{
		retry_count:     0,
		max_retry_count: max_retry_count,

		timeout:     time.Duration(0),
		max_timeout: max_timeout,

		cause_func: func(err error) {
			log.Println(err)
		},

		timeout_strat: func(count int, timeout time.Duration) time.Duration {
			if max_retry_count == 0 {
				return max_timeout
			}

			multiplier := min(count, max_retry_count)

			next_timeout := 2 * time.Duration(multiplier) * time.Second

			return min(next_timeout, max_timeout)
		},

		retry_strat: func(count, max_count int) bool {
			if max_count == 0 {
				return true
			}

			if count > max_count {
				return false
			}
			return true
		},
	}
}

type RetryFunc func() (cause error, retry bool)

func (r *Retry) Do(f RetryFunc) {
	for {
		cause, retry := f()
		if cause != nil {
			r.Cause(cause)
		}

		if retry && r.CanRetry() {
			r.retry_count++
		} else {
			break
		}

		r.timeout = r.NextTimeout()

		<-time.After(r.timeout)
	}
}

func (r *Retry) NextTimeout() time.Duration {
	return r.timeout_strat(r.retry_count, r.timeout)
}

func (r *Retry) Cause(cause error) {
	r.cause_func(cause)
}

func (r *Retry) CanRetry() bool {
	return r.retry_strat(r.retry_count, r.max_retry_count)
}
