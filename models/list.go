package models

import (
	"log"
)

type List interface {
	Size() uint
	Push(string)
	Get(uint) *string
	GetData() []string
	Remove(uint)
	AddListener(chan struct{})
	Dump()
}
type list struct {
	data      []string
	size      uint
	cur       uint
	last      string
	listeners []chan struct{}
}

func NewList(size uint) List {
	data := make([]string, size, size)
	cur := -1
	return &list{data, size, uint(cur), "", make([]chan struct{}, 0)}
}
func NewListWithData(size uint, data []string) List {
	cur := len(data) - 1
	fixedData := make([]string, size)
	copy(fixedData, data)
	return &list{fixedData, size, uint(cur), "", make([]chan struct{}, 0)}
}

func (this *list) Size() uint {
	return uint(len(this.data))
}

func (this *list) Push(str string) {
	if len(str) == 0 {
		return
	}
	if str == *this.Get(0) {
		return
	}
	this.last = this.data[(this.cur+1)%this.size]
	this.data[(this.cur+1)%this.size] = str
	this.cur++
	for i, _ := range this.listeners {
		this.listeners[i] <- struct{}{}
	}
}

func (this *list) Get(i uint) *string {
	return &this.data[(this.cur-i)%this.size]
}

func (this *list) GetData() []string {
	result := make([]string, 0, this.size)
	for i, _ := range this.data {
		index := this.toAbsoluteIndex(this.size - 1 - uint(i)) // from old
		if len(this.data[index]) == 0 {
			continue
		}
		result = append(result, this.data[index])
	}
	log.Println(result)
	return result
}

func (this *list) Remove(index uint) {
	i := uint(index)
	for ; this.toAbsoluteIndex(0) != this.toAbsoluteIndex(i+1); i++ {
		// log.Printf("%v <- %v", this.list[this.toAbsoluteIndex(i)], this.list[this.toAbsoluteIndex(i+1)])
		this.data[this.toAbsoluteIndex(i)] = this.data[this.toAbsoluteIndex(i+1)]
	}
	// log.Printf("last %v <- %v", this.list[this.toAbsoluteIndex(i)], this.last)
	this.data[this.toAbsoluteIndex(i)] = this.last
}

func (this *list) AddListener(c chan struct{}) {
	this.listeners = append(this.listeners, c)
}

func (this *list) Dump() {
	for i := range this.data {
		log.Printf("%2d:[%s]\n", i, this.data[i])
	}
}
func (this *list) toAbsoluteIndex(i uint) uint {
	return (this.cur - i) % this.size
}
