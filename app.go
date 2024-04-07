package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Start the core application and return once it terminates.
func Start(ctx context.Context, rx <-chan *spotify.FullTrack, tx chan<- *Event) {
	ctx, cancel := context.WithCancel(ctx)
	app := tview.NewApplication()
	pages := tview.NewPages()

	//TODO: help := tview.NewTextView().SetScrollable(false).ScrollToEnd()
	//fmt.Fprint(help, `help
	//text
	//here`)
	//pages.AddPage("help", help, true, false)

	logs := tview.NewTextView().SetDynamicColors(true).SetScrollable(false).ScrollToEnd()
	log.SetOutput(tview.ANSIWriter(logs))
	pages.AddPage("logs", logs, true, true)

	// End user pressed one of `Escape`, `Enter`, `Tab`, or `Backtab` on the logs
	// page.
	logs.SetDoneFunc(func(key tcell.Key) {
		pages.SwitchToPage("listing")
	})

	// Redraw if logs have been written.
	logs.SetChangedFunc(func() {
		app.Draw()
	})

	listing := tview.NewTable().SetSelectable(true, false).Select(0, 0)
	go ListingManager(ctx, listing, rx)
	pages.AddPage("listing", listing, true, false)

	// End user pressed `Enter` on the listing page.
	listing.SetSelectedFunc(func(row, _ int) {
		uri, ok := listing.GetCell(row, 0).GetReference().(spotify.URI)
		if !ok {
			log.Error("invalid URI")
		} else {
			tx <- RequestPlayURI(uri)
		}
	})

	// End user pressed one of `Escape`, `Tab`, or `Backtab` on the listing page.
	// Also triggers if user pressed `Enter` when nothing is selected.
	listing.SetDoneFunc(func(key tcell.Key) {
	})

	listing.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if ev.Key() == tcell.KeyRune {
			switch ev.Rune() {

			// End user pressed `p` on the listing.
			case 'p':
				tx <- RequestToggle()

			// End user pressed `f` on the listing.
			case 'f':
				tx <- RequestPlayNext()

			// End user pressed `b` on the listing.
			case 'b':
				tx <- RequestPlayPrevious()

			// End user pressed `Space` on the listing.
			case ' ':
				row, _ := listing.GetSelection()
				uri, ok := listing.GetCell(row, 0).GetReference().(spotify.URI)
				if !ok {
					log.Error("invalid URI")
				} else {
					tx <- RequestQueueURI(uri)
				}

			}
		}

		// Pass event back to primitive for other handlers.
		return ev
	})

	app.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {

		// End user pressed `F1` anywhere.
		case tcell.KeyF1:
			pages.SwitchToPage("listing")
			return nil

		// End user pressed `F2` anywhere.
		case tcell.KeyF2:
			pages.SwitchToPage("logs")
			return nil

		// End user pressed `F3` anywhere.
		case tcell.KeyF3:
			// TODO: show help

		case tcell.KeyRune:

			switch ev.Rune() {

			// End user pressed `?` anywhere.
			case '?':
				// TODO: show help

			// End user pressed `q` anywhere.
			case 'q':
				app.Stop()
				return nil
			}
		}

		// Pass event back to application, will be routed to focused
		// primitive.
		return ev
	})

	// This will block until the application dies.

	err := app.SetRoot(pages, true).Run()
	cancel()

	log.SetOutput(os.Stdout)
	if err != nil {
		log.WithError(err).Fatal("application died")
	}

	// TODO: for long-running programs, maybe need to periodically run cli.RefreshToken?
}

