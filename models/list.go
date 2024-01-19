package models

import (
	"log"
	linked_list "container/list"
)

type List interface {
	Size() uint
	Push(string)
	Get(uint) *string
	GetData() []string
	Pop(uint) *string
	AddListener(chan struct{})
	Dump()
}
type list struct {
	data      *linked_list.List
	size      uint
	listeners []chan struct{}
}

func NewList(size uint) List {
	linkedList := linked_list.New()
	for i := uint(0); i < size; i++ {
		linkedList.PushBack("")
	}
	return &list{linkedList, size, make([]chan struct{}, 0)}
}
func NewListWithData(size uint, data []string) List {
	linkedList := linked_list.New()
	log.Printf("NewListWithData %v", len(data))
	for _, each := range data {
		linkedList.PushBack(each)
	}
	for uint(linkedList.Len()) < size {
		linkedList.PushBack("")
	}
	return &list{linkedList, size, make([]chan struct{}, 0)}
}

func (this *list) Size() uint {
	return this.size
}

func (this *list) Push(str string) {
	if len(str) == 0 {
		return
	}
	if str == *this.Get(0) {
		return
	}
	this.data.PushBack(str)
	this.data.Remove(this.data.Front())
	this.notify()
}

func (this *list) Get(index uint) *string {
	s := ""
	cur := this.data.Back()
	for i := uint(0); i < this.size; i++ {
		if cur == nil {
			return &s
		}
		if index == i {
			s = cur.Value.(string)
			return &s
		}
		cur = cur.Prev()
	}
	log.Printf("List.Get() ERROR: out of index")
	panic("List.Get Error: data not found")
}

func (this *list) GetData() []string {
	result := make([]string, this.size, this.size)
	cur := this.data.Front()
	for i := uint(0); i < this.size; i++ {
		if cur == nil {
			return result
		}
		result[i] = cur.Value.(string)
		cur = cur.Next()
	}
	log.Printf("List: %v", result)
	return result
}

func (this *list) Pop(index uint) *string {
	s := ""
	cur := this.data.Back()
	for i := uint(0); i < this.size; i++ {
		if index == i {
			s = cur.Value.(string)
			this.data.MoveToBack(cur)
			this.notify()
			return &s
		}
		cur = cur.Prev()
	}
	log.Printf("List.Pop error")
	panic("List.Pop Data not found")
}

func (this *list) AddListener(c chan struct{}) {
	this.listeners = append(this.listeners, c)
}

func (this *list) Dump() {
	log.Printf("List.Dump: size %v", this.size)
	log.Printf("data.len %v", this.data.Len())
	cur := this.data.Back()
	for i := uint(0); i < this.size; i++ {
		if cur == nil {
			break
		}
		log.Printf("List: %02d:[%s]", i, cur.Value.(string))
		cur = cur.Prev()
	}
}

func (this *list) notify() {
	for i, _ := range this.listeners {
		this.listeners[i] <- struct{}{}
	}
}