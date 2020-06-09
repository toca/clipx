package models

import (
	"log"
)

type List interface {
	Size() uint
	Add(string)
	Get(uint) *string
	Remove(uint)
	GetNotify() chan bool
	Dump()
}
type _list struct {
	list   []string
	size   uint
	cur    uint
	notify chan bool
	last   string
}

func NewList(size uint) List {
	list := make([]string, size, size)
	cur := -1
	notify := make(chan bool, 1)
	return &_list{list, size, uint(cur), notify, ""}
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
	this.notify <- true
}

func (this *_list) Get(i uint) *string {
	return &this.list[(this.cur-i)%this.size]
}

func (this *_list) Remove(index uint) {
	i := uint(index)
	for ; this.toAbsoluteIndex(0) < this.toAbsoluteIndex(i+1); i++ {
		this.list[this.toAbsoluteIndex(i)] = this.list[this.toAbsoluteIndex(i+1)]
	}
	this.list[this.toAbsoluteIndex(i)] = this.last
}

func (this *_list) GetNotify() chan bool {
	return this.notify
}

func (this *_list) Dump() {
	for i := range this.list {
		log.Printf("%2d:[%s]\n", i, this.list[i])
	}
}
func (this *_list) toAbsoluteIndex(i uint) uint {
	return (this.cur - i) % this.size
}
