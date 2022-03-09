package network_define

import (
	"errors"
	"math/rand"
)

var (
	ErrNodeIdGetDistanceFailed = errors.New("ErrNodeIdGetDistanceFailed")
)

type NodeId8 int32

func NewNodeId8() NodeId8 {
	return NodeId8(rand.Int31())
}

func (self NodeId8) GetDistance(id NodeId8) (int64, error) {
	return int64(int32(self) ^ int32(id)), nil
}
