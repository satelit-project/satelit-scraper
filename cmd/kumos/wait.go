package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func WaitPort(port int, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	addr := fmt.Sprintf(":%d", port)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if _, err := net.Dial("tcp", addr); err == nil {
				return nil
			}

			time.Sleep(1 * time.Second)
		}

	}
}
