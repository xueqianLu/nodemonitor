package async

import (
	"github.com/ethereum/go-ethereum/common"
	"sync"
)

type SafeMap struct {
	lock   sync.RWMutex
	data   map[common.Address]int
	blocks map[common.Address][]int
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data:   make(map[common.Address]int, 1000),
		blocks: make(map[common.Address][]int, 1000),
	}
}

func (s *SafeMap) AddLose(addr common.Address, block int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if lose, exist := s.data[addr]; exist {
		s.data[addr] = lose + 1
		blocks := s.blocks[addr]
		blocks = append(blocks, block)
		s.blocks[addr] = blocks
	} else {
		s.data[addr] = 1
		blocks := make([]int, 0)
		blocks = append(blocks, block)
		s.blocks[addr] = blocks
	}
}

func (s *SafeMap) GetLose(addr common.Address) int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	lose, _ := s.data[addr]
	return lose
}

func (s *SafeMap) Range(rfunc func(addr common.Address, lose int, blocks []int)) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k, v := range s.data {
		b := s.blocks[k]
		if rfunc != nil {
			rfunc(k, v, b)
		}
	}
}
