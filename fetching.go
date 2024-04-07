package main

// Manager for fetching tracks.

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
)

// Actually fetch tracks from Spotify.
func fetchingWorker(ctx context.Context, cli *spotify.Client, ch chan<- *spotify.FullTrack, done chan<- bool) {
	// Fetch first page.
	log.Trace("fetching first page...")
	page, err := cli.CurrentUsersTracks(ctx)
	if err != nil {
		// Cleanup.
		done <- true
		close(ch)

		// Terminate immediately.
		log.WithError(err).Fatal("failed to fetch first page")
	}
	log.Trace("fetched first page")

	// Fetch additional pages.
	for {
		select {

		// Context is cancelled.
		case <-ctx.Done():
			break

		// Continue fetching.
		default:
			for _, track := range page.Tracks {
				// If the channel buffer is full, this will
				// block. This is intentional. Manager will
				// ensure that the buffer is emptied if the
				// worker needs to terminate.
				ch <- &track.FullTrack
			}

			log.Trace("fetching a new page...")
			err = cli.NextPage(ctx, page)

			// Reached end of pages.
			if err == spotify.ErrNoMorePages {
				log.Debug("no more pages")
				break
			}

			// Other error?
			if err != nil {
				log.WithError(err).Error("failed to fetch a page")
				break
			}
		}

	}

	// Cleanup.
	done <- true
	close(ch)
}

// Manage fetching tracks from Spotify.
func FetchingManager(ctx context.Context, cli *spotify.Client, ch chan *spotify.FullTrack) {
	done := make(chan bool)
	go fetchingWorker(ctx, cli, ch, done)

	for {
		select {

		// Fetch worker is terminating faster than the manager.
		case <-done:
			log.Trace("fetch worker terminated, terminating fetch manager")
			break

		// Context is cancelled.
		case <-ctx.Done():
			log.Trace("context closed, cleaning up fetch worker...")

			// Drain channel to unblock worker. Worker should then
			// notice the context is cancelled, too.
			for _ = range ch {
			}

			log.Trace("cleanup complete, terminating fetch manager")
			break

		// Wait before re-checking the above situations.
		default:
			log.Trace("track manager & worker still running...")
			time.Sleep(fetchingTimeout * time.Second)
		}

	}

	// Block until worker is cleaned up.
	<- done
}

