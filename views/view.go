package views

import (
	"clipx/models"
	"fmt"
	"sync"

	"github.com/gdamore/tcell"
)

type View struct {
	list     *models.List
	listSize uint
	screen   tcell.Screen
	once     sync.Once
	notify   chan bool
	close    chan struct{}
	cursor   uint
	rerender chan struct{}
}

func NewView(list *models.List, close chan struct{}) *View {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err = screen.Init(); err != nil {
		panic(err)
	}
	return &View{list, list.Size(), screen, sync.Once{}, list.GetNotify(), close, 0, make(chan struct{})}
}

// TODO
// index access
// window resize
// date
// help?

func (this *View) Show() {
	this.render()
	go func() {
		for {
			e := this.screen.PollEvent()
			switch e := e.(type) {
			case *tcell.EventKey:
				terminate := this.onKeyEvent(e)
				if terminate {
					break
				}
			}
		}
	}()

loop:
	for {
		select {
		case _, ok := <-this.notify:
			if !ok {
				break loop
			} else {
				this.render()
				break
			}
		case <-this.rerender:
			this.render()
		case <-this.close:
			this.screen.Fini()
			break loop
		}
	}
}

func (this *View) render() {
	this.screen.Clear()
	// title
	for i, c := range "-- [clipx] clipboard extention for windows --" {
		this.screen.SetContent(i, 0, c, nil, tcell.StyleDefault)
	}

	indexFormat := "%x: "
	// 1: first contents
	// 2: second contents
	// ...
	// f: last contents
	for row := uint(1); row < this.listSize; row++ {
		style := tcell.StyleDefault
		if this.cursor%(this.listSize-1) == row-1 {
			style = tcell.StyleDefault.Reverse(true)
		}
		index := fmt.Sprintf(indexFormat, row)
		colIndex := this.setLineContent(&index, int(row), 0, style)
		content := this.list.Get(row)
		for col, r := range *content {
			toReadable(&r)
			this.screen.SetContent(col+colIndex, int(row), r, nil, style)
		}
	}
	this.screen.Show()
}

func (this *View) onKeyEvent(keyEvent *tcell.EventKey) (finished bool) {
	switch keyEvent.Key() {
	case tcell.KeyESC:
		close(this.close)
		return true
	case tcell.KeyCtrlC:
		close(this.close)
		return true
	case tcell.KeyUp:
		this.cursor--
		this.rerender <- struct{}{}
		return false
	case tcell.KeyDown:
		this.cursor++
		this.rerender <- struct{}{}
		return false
	default:
		return false
	}
}

func (this *View) setLineContent(s *string, row int, col int, style tcell.Style) int {
	i := 0
	rStr := []rune(*s)
	for ; i < len(*s); i++ {
		this.screen.SetContent(col+i, row, rStr[i], nil, style)
	}
	return col + i + 1
}

func toReadable(r *rune) {
	// windows の clipboard は改行を \r\n でとってくる
	if *r == '\n' {
		*r = '⏎'
	} else if *r == '\r' {
		*r = '↵'
	}
}
