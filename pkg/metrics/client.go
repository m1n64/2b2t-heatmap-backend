package metrics

import (
	"github.com/alexcesaro/statsd"
)

func NewStatsDClient(addr string) (*statsd.Client, error) {
	return statsd.New(
		statsd.Address(addr),
		statsd.Prefix("tile"),
	)
}
