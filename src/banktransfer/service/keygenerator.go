package service

import (
	"fmt"
	"os"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

// KeyGenerator generates keys in the pattern {hostname}-{counter}
type KeyGenerator struct {
	counter  int32
	hostName string
}

func NewKeyGenerator() *KeyGenerator {
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return &KeyGenerator{counter: 0, hostName: hostName}
}

func (g *KeyGenerator) getUniqueId() string {
	return fmt.Sprintf("%s-%d", g.hostName, atomic.AddInt32(&g.counter, 1))
}
