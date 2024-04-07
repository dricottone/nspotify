package main

// Manager for the track listing.

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
	"github.com/rivo/tview"
)

// Actually append tracks to the listing.
//
// NOTE: `ch` is for receiving tracks from the `FetchingManager`.
//       `quit` is for receiving a termination signal from the `ListingManager`.
//       `done` is for the reverse, sending a termination signal to the `ListingManager`.
func listingWorker(listing *tview.Table, ch <-chan *spotify.FullTrack, quit <-chan bool, done chan<- bool) {
	for {
		select {

		// Terminate loader.
		case <-quit:
			break

		default:
			cursor, _ := listing.GetSelection()
			length := listing.GetRowCount()

			if (length - cursor) < loadLookahead {
				track, ok := <- ch

				// Channel is closed; terminate now.
				if !ok {
					log.Trace("no more tracks to load")
					break
				}

				log.Tracef("loading %s...", track.Name)
				for col, cell := range IntoCells(track) {
					listing.SetCell(length, col, cell)
				}

				continue
			}

			// Wait before retrying
			time.Sleep(loadTimeout * time.Second)
		}
	}

	done <- true
}

// Manager for appending tracks to the listing.
func ListingManager(ctx context.Context, listing *tview.Table, ch <-chan *spotify.FullTrack) {
	// Load first N tracks eagerly.
	log.Tracef("loading %d tracks...", loadEager)
	for i := 0; i < loadEager; i++ {
		track := <- ch
		for col, cell := range IntoCells(track) {
			listing.SetCell(i, col, cell)
		}
	}
	log.Tracef("loaded %d tracks", loadEager)

	// Load more tracks lazily.
	quit := make(chan bool)
	done := make(chan bool)
	go listingWorker(listing, ch, quit, done)

	for {
		select {

		// Somehow track loader is terminating faster than this loop.
		case <-done:
			log.Trace("track renderer stopped running")
			break

		// Context is cancelled; terminate track loader.
		case <-ctx.Done():
			break

		}
	}

	quit <- true
}

