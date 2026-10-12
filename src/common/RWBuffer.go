package common

import (
	"sync"
)

type RWIntBuffer struct {
	buffer []int
	lock   sync.RWMutex
}

func NewRWIntBuffer(size int) *RWIntBuffer {
	return &RWIntBuffer{buffer: make([]int, size)}
}

func (self *RWIntBuffer) Append(i int) {
	self.lock.Lock()
	self.buffer = append(self.buffer, i)
	self.lock.Unlock()
}

func (self *RWIntBuffer) Insert(idx, i int) {
	self.lock.Lock()
	self.buffer[idx] = i
	self.lock.Unlock()
}

func (self *RWIntBuffer) Size() int {
	return len(self.buffer)
}

//////////////////////////////////////////
type RWFloat64Buffer struct {
	buffer []float64
	lock   sync.RWMutex
}

func NewRWFloat64Buffer(size int) *RWFloat64Buffer {
	return &RWFloat64Buffer{buffer: make([]float64, size)}
}

func (self *RWFloat64Buffer) Append(i float64) {
	self.lock.Lock()
	self.buffer = append(self.buffer, i)
	self.lock.Unlock()
}

func (self *RWFloat64Buffer) Insert(idx int, i float64) {
	self.lock.Lock()
	self.buffer[idx] = i
	self.lock.Unlock()
}

func (self *RWFloat64Buffer) Size() int {
	return len(self.buffer)
}
