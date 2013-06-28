package main

import (
	"log"
	"strings"

	"github.com/lxn/walk"
)

func main() {
	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}

	if err := mw.SetTitle("GoCalSync"); err != nil {
		log.Fatal(err)
	}

	// load icon from file
	icon, err := walk.NewIconFromFile("assets/caldb.ico")
	if err != nil {
		log.Fatal(err)
	}

	mw.SetIcon(icon)

	mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		*canceled = true
		mw.Hide()
	})

	mw.SizeChanged().Attach(func() {
		log.Println("size changed", mw.Size())
	})

	size := walk.Size{600, 400}
	mw.SetMinMaxSize(size, size)
	mw.SetSize(size)

	if err := mw.SetLayout(walk.NewVBoxLayout()); err != nil {
		log.Fatal(err)
	}

	// horizontal splitter
	splitter, err := walk.NewHSplitter(mw)
	if err != nil {
		log.Fatal(err)
	}
	mw.Children().Add(splitter)

	// let side
	inTE, err := walk.NewTextEdit(splitter)
	if err != nil {
		log.Fatal(err)
	}
	splitter.Children().Add(inTE)

	// right side
	outTE, err := walk.NewTextEdit(splitter)
	if err != nil {
		log.Fatal(err)
	}
	outTE.SetReadOnly(true)
	splitter.Children().Add(outTE)

	// add button
	push, err := walk.NewPushButton(mw)
	if err != nil {
		log.Fatal(err)
	}
	push.SetText("SCREAM")
	push.Clicked().Attach(func() {
		outTE.SetText(strings.ToUpper(inTE.Text()))
	})
	mw.Children().Add(push)

	// create systray
	si := CreateSystrayIcon(mw, icon)
	defer si.Dispose()

	mw.Run()
}

func CreateSystrayIcon(mw *walk.MainWindow, icon *walk.Icon) *walk.NotifyIcon {

	// create notify icon
	ni, err := walk.NewNotifyIcon()
	if err != nil {
		log.Fatal(err)
	}

	// set icon and tooltip
	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}

	if err := ni.SetToolTip("Click for settings or use context menu to exit."); err != nil {
		log.Fatal(err)
	}

	// show settings on left mouse button
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		mw.Show()
	})

	// add settings action
	settingsAction := walk.NewAction()
	if err := settingsAction.SetText("S&ettings"); err != nil {
		log.Fatal(err)
	}

	settingsAction.Triggered().Attach(func() {
		mw.Show()
	})

	if err := ni.ContextMenu().Actions().Add(settingsAction); err != nil {
		log.Fatal(err)
	}

	// separator
	if err := ni.ContextMenu().Actions().Add(walk.NewSeparatorAction()); err != nil {
		log.Fatal(err)
	}

	// add exit action
	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		log.Fatal(err)
	}

	exitAction.Triggered().Attach(func() {
		walk.App().Exit(0)
	})

	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}

	// show notify icon
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	if err := ni.ShowInfo("GoCalSync is enabled", "Click the icon to display settings"); err != nil {
		log.Fatal(err)
	}

	return ni
}
