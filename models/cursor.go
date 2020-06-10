package models

type Cursor interface {
	Up()
	Down()
	Set(uint)
	Reset()
	GetIndex() uint
	IsSelected(uint) bool
	AddListener(chan struct{})
}
type cursor struct {
	index     uint
	length    uint
	listeners []chan struct{}
}

// max: uint 0 origin
func NewCursor(max uint) Cursor {
	return &cursor{0, max, make([]chan struct{}, 0)}
}

func (this *cursor) Up() {
	this.index--
	this.notifyAll()
}

func (this *cursor) Down() {
	this.index++
	this.notifyAll()
}

func (this *cursor) Set(i uint) {
	this.index = i
	this.notifyAll()
}

func (this *cursor) Reset() {
	this.index = 0
	this.notifyAll()
}

func (this *cursor) GetIndex() uint {
	return (this.index % this.length)
}

func (this *cursor) IsSelected(i uint) bool {
	return (this.index % this.length) == i
}

func (this *cursor) AddListener(c chan struct{}) {
	this.listeners = append(this.listeners, c)
}

func (this *cursor) notifyAll() {
	for i, _ := range this.listeners {
		this.listeners[i] <- struct{}{}
	}
}
