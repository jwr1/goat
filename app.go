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

	rootElement := &Element{}

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
		rootConstraints := SizeInt(screenWidth, screenHeight).TightConstraints()

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
	width := canvas.size.Width.Int()

	for i, cell := range canvas.cells {
		style := tcell.StyleDefault.
			Foreground(tcell.NewRGBColor(int32(cell.Foreground.R), int32(cell.Foreground.G), int32(cell.Foreground.B))).
			Background(tcell.NewRGBColor(int32(cell.Background.R), int32(cell.Background.G), int32(cell.Background.B)))

		const uint8Midpoint = 0xFF / 2

		// If opacity is less than half, then use terminal default color
		if cell.Foreground.A < uint8Midpoint {
			style = style.Foreground(tcell.ColorDefault)
		}
		if cell.Background.A < uint8Midpoint {
			style = style.Background(tcell.ColorDefault)
		}

		if cell.TextStyle != nil {
			attr := tcell.AttrNone
			if cell.TextStyle != nil {
				if cell.TextStyle.Bold {
					attr |= tcell.AttrBold
				}
				if cell.TextStyle.Blink {
					attr |= tcell.AttrBlink
				}
				if cell.TextStyle.Dim {
					attr |= tcell.AttrDim
				}
				if cell.TextStyle.Italic {
					attr |= tcell.AttrItalic
				}
				if cell.TextStyle.Underline {
					attr |= tcell.AttrUnderline
				}
				if cell.TextStyle.StrikeThrough {
					attr |= tcell.AttrStrikeThrough
				}
			}

			style = style.
				Attributes(attr).
				Url(cell.TextStyle.Url).
				UrlId(cell.TextStyle.UrlId)
		}

		screen.SetContent(i%width, i/width, cell.Rune, nil, style)
	}

	screen.Show()
}
