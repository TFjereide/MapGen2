package common

import (
	"sync"
)

type IntQueue struct {
	items []int
	lock  sync.RWMutex
}

func NewIntQueue() *IntQueue {
	return &IntQueue{items: []int{}}
}

func (self *IntQueue) Add(i int) {
	self.lock.Lock()
	self.items = append(self.items, i)
	self.lock.Unlock()
}

func (self *IntQueue) PopStart() int {
	self.lock.Lock()
	i := self.items[0]
	self.items = self.items[1:len(self.items)]
	self.lock.Unlock()
	return i
}
func (self *IntQueue) Front() int {
	self.lock.RLock()
	i := self.items[0]
	self.lock.RUnlock()
	return i
}

func (self *IntQueue) IsEmpty() bool {
	return len(self.items) == 0
}

func (self *IntQueue) Size() int {
	return len(self.items)
}
