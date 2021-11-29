//go:build integration
// +build integration

package wnom

import (
	"net"
	"time"
)

const defaultRetryTimeout = time.Second

// awaitForPort awaits for the port within the timeout.
func awaitForPort(host, port string, timeout time.Duration) error {
	return retryWithTimeout(func() error {
		_, err := net.Dial("tcp", net.JoinHostPort(host, port))
		return err
	}, timeout)
}

// retryWithTimeout retries the operation within the timout.
func retryWithTimeout(operation func() error, timeout time.Duration) error {
	startTime := time.Now().UnixNano()
	for {
		err := operation()
		if err == nil {
			return nil
		}
		if time.Now().UnixNano()-startTime > timeout.Nanoseconds() {
			return err
		}
		time.Sleep(defaultRetryTimeout)
	}
}
