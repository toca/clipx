package models

type Cursor interface {
	Up()
	Down()
	Reset()
	GetIndex() uint
	IsSelected(uint) bool
	GetNotify() chan bool
}
type cursor struct {
	index  uint
	length uint
	notify chan bool
}

// max: uint 0 origin
func NewCursor(max uint) Cursor {
	return &cursor{0, max, make(chan bool)}
}

func (this *cursor) Up() {
	this.index--
	this.notify <- true
}

func (this *cursor) Down() {
	this.index++
	this.notify <- true
}

func (this *cursor) Reset() {
	this.index = 0
	this.notify <- true
}
func (this *cursor) GetIndex() uint {
	return (this.index % this.length)
}

func (this *cursor) IsSelected(i uint) bool {
	return (this.index % this.length) == i
}

func (this *cursor) GetNotify() chan bool {
	return this.notify
}
