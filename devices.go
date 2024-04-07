package main

// Spotify device interactions.

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
)

// Fetch and report the devices available to Spotify.
func ListDevices(ctx context.Context, cli *spotify.Client) {
	dev, err := cli.PlayerDevices(ctx)
	if err != nil {
		log.WithError(err).Fatal("failed to list devices")
	}

	fmt.Printf("%-30s %s\n", "Device (*=active)", "ID (pass this with `-device=ID`)")

	for _, device := range dev {
		id := device.ID.String()
		if id == "" {
			id = "?"
		}
		if device.Restricted {
			id = "X"
		}

		name := device.Name
		if device.Active {
			name += " (*)"
		}

		fmt.Printf("%-30s %s\n", name, id)
	}
}

