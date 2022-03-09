package network_define

import (
	"math/rand"
)

type IP int64

func NewIP() IP {
	return IP(rand.Int63())
}
