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

type UnsafeIntQueue struct {
	Items []int
}

func NewUnsafeIntQueue() *UnsafeIntQueue {
	return &UnsafeIntQueue{Items: []int{}}
}

func (self *UnsafeIntQueue) Add(i int) {
	self.Items = append(self.Items, i)
}

func (self *UnsafeIntQueue) PopStart() int {
	i := self.Items[0]
	self.Items = self.Items[1:len(self.Items)]
	return i
}
func (self *UnsafeIntQueue) Front() int {
	i := self.Items[0]
	return i
}

func (self *UnsafeIntQueue) IsEmpty() bool {
	return len(self.Items) == 0
}

func (self *UnsafeIntQueue) Size() int {
	return len(self.Items)
}
