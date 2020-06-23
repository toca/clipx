package views

import (
	"clipx/controllers"
	"clipx/models"
	"fmt"
	"log"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

type View struct {
	ctrl        controllers.Controller
	list        models.List
	listSize    uint
	listEvent   chan struct{}
	cursor      models.Cursor
	cursorEvent chan struct{}
	screen      tcell.Screen
	closed      chan struct{}
	stop        chan struct{}
	status      string
}

const WINDOW_WIDTH = 64

func NewView(ctrl controllers.Controller, list models.List, cursor models.Cursor) *View {
	// prepare window
	err := ctrl.SetWindowSize(WINDOW_WIDTH, int16(list.Size()+1))
	if err != nil {
		panic(err)
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err = screen.Init(); err != nil {
		panic(err)
	}

	listEvent := make(chan struct{}, 1)
	list.AddListener(listEvent)
	cursorEvent := make(chan struct{}, 4)
	cursor.AddListener(cursorEvent)
	return &View{ctrl: ctrl,
		list:        list,
		listSize:    list.Size(),
		listEvent:   listEvent,
		cursor:      cursor,
		cursorEvent: cursorEvent,
		screen:      screen,
		closed:      make(chan struct{}),
		stop:        make(chan struct{}),
		status:      "--   Select:[↑↓]   Paste:[Enter]   Hide:[ECS]    Exit:[^C]   --",
	}
}

func (this *View) Show() {
	this.screen.EnableMouse()
	defer this.screen.DisableMouse()
	this.render()
	event := make(chan tcell.Event)
	go func() {
		for {
			e := this.screen.PollEvent()
			event <- e
		}
	}()

loop:
	for {
		select {
		case e := <-event:
			this.onEvent(e)
		case _, ok := <-this.listEvent:
			if !ok {
				break loop
			} else {
				this.render()
				break
			}
		case _, ok := <-this.cursorEvent:
			if !ok {
				break loop
			} else {
				this.render()
				break
			}
		case <-this.stop:
			this.screen.Fini()
			close(this.closed)
			break loop
		}
	}
}

func (this *View) GetClosedNotify() chan struct{} {
	return this.closed
}

func (this *View) render() {
	this.screen.Clear()

	// title
	newLine := this.setTitile(0)
	newLine = this.setListContent(newLine)
	newLine = this.setStatus(newLine)
	this.screen.Show()
}

func (this *View) onEvent(event tcell.Event) {
	switch e := event.(type) {
	case *tcell.EventKey:
		this.onKeyEvent(e)
	case *tcell.EventMouse:
		if e.Buttons() != tcell.Button1 {
			return
		}
		_, y := e.Position()
		if 1 <= y && uint(y) < this.listSize+1 {
			this.ctrl.SetCursor(uint(y - 1))
			err := this.ctrl.Paste()
			if err != nil {
				this.status = err.Error()
			}
		}
	case *tcell.EventResize:
		this.render()
	default:
		log.Printf("View.onEvent: %v", e)
	}
}

func (this *View) onKeyEvent(keyEvent *tcell.EventKey) {
	switch keyEvent.Key() {
	case tcell.KeyCtrlC:
		close(this.stop)
	case tcell.KeyESC:
		this.ctrl.Disappear()
	case tcell.KeyUp:
		this.ctrl.Up()
	case tcell.KeyDown:
		this.ctrl.Down()
	case tcell.KeyEnter:
		err := this.ctrl.Paste()
		if err != nil {
			this.status = err.Error()
		}
	default:
		num, err := strconv.ParseUint(string(keyEvent.Rune()), 16, 8)
		if err != nil {
			return
		}
		if 0 <= num && num < 16 {
			this.ctrl.SetCursor(uint(num))
			err := this.ctrl.Paste()
			if err != nil {
				this.status = err.Error()
			}
		}
		return
	}
}

// Set Title Line
func (this *View) setTitile(row int) int {
	for i, c := range "-----------------  <clipx> clipboard extention  -----------------" {
		this.screen.SetContent(i, row, c, nil, tcell.StyleDefault)
	}
	return row + 1
}

// Set Clipboard List Data
func (this *View) setListContent(row int) int {
	indexFormat := "%x: "
	// 0: first contents
	// 1: second contents
	// ...
	// e: 15th content
	// f: last contents
	for i := uint(0); i < this.listSize; i++ {
		style := tcell.StyleDefault
		if i%2 == 0 {
			style = tcell.StyleDefault.Background(tcell.ColorNavy)
		} else {
			// style = tcell.StyleDefault.Background(tcell.ColorBlue)
		}

		// style selected item
		if this.cursor.IsSelected(i) {
			style = tcell.StyleDefault.Reverse(true)
		}
		prefix := fmt.Sprintf(indexFormat, i)
		col := this.setLineContent(&prefix, row+int(i), 0, style)
		content := this.list.Get(i)
		for _, r := range *content {
			r := toReadable(&r)
			if r == nil {
				continue
			}
			this.screen.SetContent(col, row+int(i), *r, nil, style)
			col += runewidth.RuneWidth(*r)
		}
		for ; col <= WINDOW_WIDTH; col++ {
			this.screen.SetContent(col, row+int(i), ' ', nil, style)
		}
	}
	return row + int(this.listSize)
}

func (this *View) setStatus(row int) int {
	this.setLineContent(&this.status, row, 0, tcell.StyleDefault)
	// col := 0
	// for i := tcell.ColorBlack; i <= tcell.Color128; i++ {
	// 	style := tcell.StyleDefault.Background(i)
	// 	this.screen.SetContent(col, row, ' ', nil, style)
	// 	col++
	// }
	return row + 1
}

func (this *View) setLineContent(s *string, row int, col int, style tcell.Style) int {
	i := 0
	for _, r := range *s {
		this.screen.SetContent(col+i, row, r, nil, style)
		i += runewidth.RuneWidth(r)
	}
	return col + i
}

func toReadable(r *rune) *rune {
	mark := '⏎'
	// windows の clipboard は改行を \r\n でとってくる
	if *r == '\n' {
		return &mark
	} else if *r == '\r' {
		// *r = '⏎'
		return nil
	}
	return r
}
