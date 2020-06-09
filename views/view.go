package views

import (
	"clipx/controllers"
	"clipx/models"
	"fmt"
	"log"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

type View struct {
	ctrl         controllers.Controller
	list         models.List
	listSize     uint
	listNotify   chan bool
	cursor       models.Cursor
	cursorNotify chan bool
	screen       tcell.Screen
	close        chan struct{}
}

func NewView(ctrl controllers.Controller, list models.List, cursor models.Cursor, close chan struct{}) *View {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err = screen.Init(); err != nil {
		panic(err)
	}
	return &View{ctrl, list, list.Size(), list.GetNotify(), cursor, cursor.GetNotify(), screen, close}
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
					return
				}
			}
		}
	}()

loop:
	for {
		select {
		case _, ok := <-this.listNotify:
			if !ok {
				break loop
			} else {
				this.render()
				break
			}
		case _, ok := <-this.cursorNotify:
			if !ok {
				break loop
			} else {
				this.render()
				break
			}
		case <-this.close:
			log.Printf("View.show: closed")
			this.screen.Fini()
			log.Printf("View.show: endloop")
			break loop
		}
	}
}

func (this *View) render() {
	this.screen.Clear()
	// title
	for i, c := range "-- [clipx] clipboard extention --" {
		this.screen.SetContent(i, 0, c, nil, tcell.StyleDefault)
	}

	indexFormat := "%x: "
	// 0: first contents
	// 1: second contents
	// ...
	// f: last contents
	for i := uint(0); i < this.listSize; i++ {
		row := int(i) + 1 // title
		style := tcell.StyleDefault
		if this.cursor.IsSelected(i) {
			style = tcell.StyleDefault.Reverse(true)
		}
		prefix := fmt.Sprintf(indexFormat, i)
		col := this.setLineContent(&prefix, row, 0, style)
		content := this.list.Get(i)
		for _, r := range *content {
			toReadable(&r)
			this.screen.SetContent(col, row, r, nil, style)
			col += runewidth.RuneWidth(r)
		}
	}
	this.screen.Show()
}

func (this *View) onKeyEvent(keyEvent *tcell.EventKey) (finished bool) {
	switch keyEvent.Key() {
	case tcell.KeyCtrlC:
		close(this.close)
		return true
	case tcell.KeyESC:
		this.ctrl.Disappear()
		return false
	case tcell.KeyUp:
		this.ctrl.Up()
		return false
	case tcell.KeyDown:
		this.ctrl.Down()
		return false
	case tcell.KeyEnter:
		err := this.ctrl.Paste()
		// TODO show status
		if err != nil {
			log.Printf("View.onKeyEvent: %v", err)
		}
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
	return col + i
}

func toReadable(r *rune) {
	// windows の clipboard は改行を \r\n でとってくる
	if *r == '\n' {
		*r = '⏎'
	} else if *r == '\r' {
		*r = '↵'
	}
}
