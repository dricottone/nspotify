package main

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

func main() {
	ctx := ConfiguredContext()
	ctx, cancel := context.WithCancel(ctx)

	// TODO: Display version mode.
	// version, ok := ctx.Value("version").(string)
	// if ok && (version != "") {
	// 	fmt.Println("nspotify %s", version)
	// 	cancel()
	// 	return
	// }

	// Authenticate with Spotify.
	cli := Authenticate(ctx)
	// TODO: incorporate rate limiting? "set the AutoRetry field on the Client struct to true"

	// List devices mode.
	sdev, ok := ctx.Value("device").(string)
	dev := spotify.ID(sdev)
	if !ok || (dev == "") {
		ListDevices(ctx, cli)
		cancel()
		return
	}

	// Fetch user tracks with Spotify client. Will continue to run in
	// background.
	fetchCh := make(chan *spotify.FullTrack, fetchingBuffer)
	go FetchingManager(ctx, cli, fetchCh)

	evCh := make(chan *Event)
	go EventsManager(ctx, cli, evCh)

	// Run terminal application. Will block until application terminates.
	Start(ctx, fetchCh, evCh)

	cancel()
}

