package scheduler

import (
	"fmt"
	"sync"
)

type Lookup struct {
	table map[int]struct{}
	m     *sync.RWMutex
}

func NewLookup() *Lookup {
	return &Lookup{
		table: make(map[int]struct{}),
		m:     &sync.RWMutex{},
	}
}

func (l *Lookup) Add(id int) error {
	l.m.Lock()
	defer l.m.Unlock()

	if _, ok := l.table[id]; ok {
		return fmt.Errorf("key %d already exists", id)
	}

	l.table[id] = struct{}{}
	return nil
}
