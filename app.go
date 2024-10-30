package goat

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gdamore/tcell/v2"
)

func RunApp(w Widget) error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	err = screen.Init()
	if err != nil {
		return err
	}

	screen.EnableMouse()
	screen.EnablePaste()
	screen.Clear()

	quit := func() {
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}

	defer quit()

	eventChan := make(chan tcell.Event)
	quitEventChan := make(chan struct{})
	quitRenderChan := make(chan struct{})

	go screen.ChannelEvents(eventChan, quitEventChan)

	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range osChan {
			quitRenderChan <- struct{}{}
			return
		}
	}()

	// Tree layout and event handling cannot happen at the same time
	// due to widget Build() methods altering event listeners
	treeLock := sync.Mutex{}

	handleEvent := func(event tcell.Event) {
		switch event := event.(type) {
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				quitRenderChan <- struct{}{}
				return
			}
		}

		treeLock.Lock()
		defer treeLock.Unlock()

		for elem, listeners := range globalHookEventListeners {
			for _, listener := range listeners {
				go listener(EventContext{
					Event:      event,
					RenderPos:  elem.renderAbsPos,
					RenderSize: elem.size,
				})
			}
		}
	}

	go func() {
		defer quit()

		for event := range eventChan {
			handleEvent(event)
		}
	}()

	rootElement := &element{}

	// debugFile, _ := os.Create("debug-tree.txt")
	// debugFileSize := 0

	// debugTree := func() {
	// 	var builder strings.Builder
	// 	stringifyTree(rootElement, &builder, 0)
	// 	debugFile.Truncate(int64(debugFileSize))
	// 	n, _ := debugFile.WriteAt([]byte(builder.String()), 0)
	// 	debugFileSize = n
	// }

	for {
		select {
		case <-quitRenderChan:
			destroyTree(rootElement)
			return nil
		default:
		}

		screenWidth, screenHeight := screen.Size()
		rootConstraints := Size{screenWidth, screenHeight}.TightConstraints()

		treeLock.Lock()

		err = rebuildTree(w, rootElement, rootConstraints)
		if err != nil {
			return fmt.Errorf("build error: %w", err)
		}

		canvas, err := renderTree(rootElement)
		if err != nil {
			return fmt.Errorf("render error: %w", err)
		}

		treeLock.Unlock()

		drawCanvasToScreen(canvas, screen)

		// debugTree()
	}
}

func drawCanvasToScreen(canvas Canvas, screen tcell.Screen) {
	width := canvas.size.Width

	for i, cell := range canvas.cells {
		screen.SetContent(i%width, i/width, cell.rune, nil, cell.style)
	}

	screen.Show()
}
