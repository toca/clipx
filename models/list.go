package models

import (
	"fmt"
)

type List struct {
	list   []string
	size   uint
	cur    uint
	notify chan bool
}

func NewList(size uint) *List {
	list := make([]string, size, size)
	cur := -1
	notify := make(chan bool, 1)
	return &List{list, size, uint(cur), notify}
}

func (this *List) Add(str string) {
	if len(str) == 0 {
		return
	}
	this.list[(this.cur+1)%this.size] = str
	this.cur++
	// this.notify <- true
}

func (this *List) Get(i int) *string {
	return &this.list[i]
}

func (this *List) GetNotify() chan bool {
	return this.notify
}

func (this *List) Dump() {
	for i := range this.list {
		fmt.Printf("%2d:[%s]\n", i, this.list[i])
	}
}
