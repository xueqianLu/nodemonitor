package async

import (
	"github.com/ethereum/go-ethereum/common"
	"sync"
)

type SafeMap struct {
	lock sync.RWMutex
	data map[common.Address]int
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[common.Address]int, 1000),
	}
}

func (s *SafeMap) AddLose(addr common.Address) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if lose, exist := s.data[addr]; exist {
		s.data[addr] = lose + 1
	} else {
		s.data[addr] = 1
	}
}

func (s *SafeMap) GetLose(addr common.Address) int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	lose, _ := s.data[addr]
	return lose
}

func (s *SafeMap) Range(rfunc func(addr common.Address, lose int)) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k, v := range s.data {
		if rfunc != nil {
			rfunc(k, v)
		}
	}
}
