package main

// Manager for events.

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
)

// Types of player events to expect.
type EventType int

const (
	// Player event requesting that playback begin with a specific URI.
	PlayURI EventType = 0

	// Player event requesting that a specific URI be enqueued for future
	// playback.
	QueueURI EventType = 1

	// Player event requesting that playback play (resume).
	Play EventType = 2

	// Player event requesting that playback pause.
	Pause EventType = 3

	// Player event requesting that playback toggle (play/pause).
	Toggle EventType = 4

	// Player event requesting that the previous song play.
	PlayPrevious EventType = 5

	// Player event requesting that the next song play.
	PlayNext EventType = 6
)

// Convert a player event to a printable (debug-able) string.
func debugEvent(ev EventType) string {
	switch ev {
	case PlayURI:
		return "PlayURI"

	case QueueURI:
		return "QueueURI"

	case Play:
		return "Play"

	case Pause:
		return "Pause"

	case Toggle:
		return "Toggle"

	case PlayPrevious:
		return "PlayPrevious"

	case PlayNext:
		return "PlayNext"
	}
	
	return fmt.Sprintf("%d", ev)
}

// A player event.
type Event struct {
	Type EventType
	URI  spotify.URI
}

// Creates an `Event` of type `PlayURI`.
func RequestPlayURI(uri spotify.URI) *Event {
	return &Event{
		Type: PlayURI,
		URI: uri,
	}
}

// Creates an `Event` of type `QueueURI`.
func RequestQueueURI(uri spotify.URI) *Event {
	return &Event{
		Type: QueueURI,
		URI: uri,
	}
}

// Creates an `Event` of type `Play`.
func RequestPlay() *Event {
	return &Event{
		Type: Play,
		URI: spotify.URI(""),
	}
}

// Creates an `Event` of type `Pause`.
func RequestPause() *Event {
	return &Event{
		Type: Pause,
		URI: spotify.URI(""),
	}
}

// Creates an `Event` of type `Toggle`.
func RequestToggle() *Event {
	return &Event{
		Type: Toggle,
		URI: spotify.URI(""),
	}
}

// Creates an `Event` of type `PlayPrevious`.
func RequestPlayPrevious() *Event {
	return &Event{
		Type: PlayPrevious,
		URI: spotify.URI(""),
	}
}

// Creates an `Event` of type `PlayNext`.
func RequestPlayNext() *Event {
	return &Event{
		Type: PlayNext,
		URI: spotify.URI(""),
	}
}

// Reusable handler for an `Event` of type `Play`.
func handlePlayEvent(ctx context.Context, cli *spotify.Client) {
	err := cli.Play(ctx)
	if err != nil {
		log.WithError(err).Error("request to play failed")
	}
}

// Reusable handler for an `Event` of type `Pause`.
func handlePauseEvent(ctx context.Context, cli *spotify.Client) {
	err := cli.Pause(ctx)
	if err != nil {
		log.WithError(err).Error("request to pause failed")
	}
}

// Manager for player events.
func EventsManager(ctx context.Context, cli *spotify.Client, ch <-chan *Event) {
	sdev, ok := ctx.Value("device").(string)
	dev := spotify.ID(sdev)
	if !ok || (dev == "") {
		log.Errorf("invalid device: %s", sdev)
	}

	for ev := range ch {
		switch ev.Type {
		case PlayURI:
			opts := &spotify.PlayOptions{
				DeviceID: &dev,
				URIs: []spotify.URI{ev.URI},
			}
			err := cli.PlayOpt(ctx, opts)
			if err != nil {
				log.WithError(err).Error("request to play URI failed")
			}
	
		case QueueURI:
			suri := string(ev.URI)
			sid := strings.Replace(suri, "spotify:track:", "", 1)
			id := spotify.ID(sid)
	
			err := cli.QueueSong(ctx, id)
			if err != nil {
				log.WithError(err).Error("request to queue URI failed")
			}
	
		case Play:
			handlePlayEvent(ctx, cli)
	
		case Pause:
			handlePauseEvent(ctx, cli)
	
		case Toggle:
			status, err := cli.PlayerCurrentlyPlaying(ctx)
			if err != nil {
				log.WithError(err).Error("request to determine playback status, so assuming that status is playing")
				handlePauseEvent(ctx, cli)
			} else {
				if status.Playing {
					handlePauseEvent(ctx, cli)
				} else {
					handlePlayEvent(ctx, cli)
				}
			}

		case PlayNext:
			err := cli.Next(ctx)
			if err != nil {
				log.WithError(err).Error("request to play next failed")
			}
	
		case PlayPrevious:
			err := cli.Previous(ctx)
			if err != nil {
				log.WithError(err).Error("request to play previous failed")
			}
	
		default:
			log.Errorf("unhandled event: %s", debugEvent(ev.Type))
		}
	}

	log.Trace("no more events, terminating events manager")
}

