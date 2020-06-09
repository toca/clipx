package controllers

import (
	"clipx/models"
)

// public interface
type Controller interface {
	Up()
	Down()
	Paste() error
	Appear()
	Disappear()
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
	this.window.Hide()
	data := this.list.Get(this.cursor.GetIndex())
	err := this.clipBoard.SetString(data)
	if err != nil {
		return err
	}
	err = this.window.SendPasteCommand()
	if err != nil {
		return err
	}
	// already clipboard added
	this.cursor.Down()
	this.list.Remove(this.cursor.GetIndex())
	return nil
}

func (this *controller) Appear() {
	this.cursor.Reset()
	this.window.Show()
}

func (this *controller) Disappear() {
	this.window.Hide()
}
