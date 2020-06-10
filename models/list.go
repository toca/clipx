package models

import (
	"log"
)

type List interface {
	Size() uint
	Add(string)
	Get(uint) *string
	Remove(uint)
	AddListener(chan struct{})
	Dump()
}
type _list struct {
	list      []string
	size      uint
	cur       uint
	last      string
	listeners []chan struct{}
}

func NewList(size uint) List {
	list := make([]string, size, size)
	cur := -1
	return &_list{list, size, uint(cur), "", make([]chan struct{}, 0)}
}

func (this *_list) Size() uint {
	return uint(len(this.list))
}
func (this *_list) Add(str string) {
	if len(str) == 0 {
		return
	}
	this.last = this.list[(this.cur+1)%this.size]
	this.list[(this.cur+1)%this.size] = str
	this.cur++
	for i, _ := range this.listeners {
		this.listeners[i] <- struct{}{}
	}
}

func (this *_list) Get(i uint) *string {
	return &this.list[(this.cur-i)%this.size]
}

func (this *_list) Remove(index uint) {
	i := uint(index)
	for ; this.toAbsoluteIndex(0) != this.toAbsoluteIndex(i+1); i++ {
		// log.Printf("%v <- %v", this.list[this.toAbsoluteIndex(i)], this.list[this.toAbsoluteIndex(i+1)])
		this.list[this.toAbsoluteIndex(i)] = this.list[this.toAbsoluteIndex(i+1)]
	}
	// log.Printf("last %v <- %v", this.list[this.toAbsoluteIndex(i)], this.last)
	this.list[this.toAbsoluteIndex(i)] = this.last
}

func (this *_list) AddListener(c chan struct{}) {
	this.listeners = append(this.listeners, c)
}

func (this *_list) Dump() {
	for i := range this.list {
		log.Printf("%2d:[%s]\n", i, this.list[i])
	}
}
func (this *_list) toAbsoluteIndex(i uint) uint {
	return (this.cur - i) % this.size
}
