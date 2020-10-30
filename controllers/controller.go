package controllers

import (
	"clipx/models"
	"fmt"
)

// public interface
type Controller interface {
	Up()
	Down()
	Paste() error
	Appear()
	Disappear()
	SetWindowSize(int16, int16) error
	SetCursor(uint)
}

// internal implement
type controller struct {
	cursor    models.Cursor
	clipBoard models.Clipboard
	list      models.List
	window    models.Window
}

func NewController(cursor models.Cursor, cb models.Clipboard, list models.List) Controller {
	return &controller{cursor, cb, list, models.NewWindow()}
}

func (this *controller) Up() {
	this.cursor.Up()
}

func (this *controller) Down() {
	this.cursor.Down()
}

func (this *controller) Paste() error {
	data := this.list.Pop(this.cursor.GetIndex())
	if len(*data) == 0 {
		return fmt.Errorf("data is empty")
	}

	this.window.Hide()

	err := this.clipBoard.SetString(data)
	if err != nil {
		return err
	}
	err = this.window.SendPasteCommand()
	if err != nil {
		return err
	}
	
	// this.list.Remove(this.cursor.GetIndex() + 1) // + 1 => already clipboard added
	return nil
}

func (this *controller) Appear() {
	this.cursor.Reset()
	this.window.Show()
}

func (this *controller) Disappear() {
	this.window.Hide()
}

func (this *controller) SetWindowSize(w int16, h int16) error {
	return this.window.ResizeWindow(w, h)
}

func (this *controller) SetCursor(i uint) {
	this.cursor.Set(i)
}
