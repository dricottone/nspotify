package main

// Spotify tracks interactions.

import (
	"fmt"
	"strings"

	"github.com/zmb3/spotify/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Format an array of Spotify artists into a `&`-delimited list of artist
// names.
func FormatArtists(artists []spotify.SimpleArtist) string {
	buff := []string{}
	for _, artist := range artists {
		buff = append(buff, artist.Name)
		//artist_id := artist.ID
		//artist_uri := artist.URI
	}

	return strings.Join(buff, " & ")
}

// Format a duration in milliseconds into `[HH]:[MM]:[SS]` (or `[MM]:[SS]` if
// duration is less than 1 hour).
func FormatDuration(ms int) string {
	ms /= 1000
	s := ms % 60
	ms /= 60
	m := ms % 60
	h := ms / 60
	if h == 0 {
		return fmt.Sprintf("%d:%02d", m, s)
	}
	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

// Create a row of table cells from a Spotify track.
func IntoCells(track *spotify.FullTrack) []*tview.TableCell {
	name := tview.NewTableCell(track.Name).SetTextColor(tcell.ColorWhite).SetReference(track.URI)
	artist := tview.NewTableCell(FormatArtists(track.Artists)).SetTextColor(tcell.ColorWhite)
	album := tview.NewTableCell(track.Album.Name).SetTextColor(tcell.ColorWhite)
	duration := tview.NewTableCell(FormatDuration(track.Duration)).SetAlign(tview.AlignRight).SetTextColor(tcell.ColorWhite)

	//id := track.ID
	//number := track.TrackNumber
	//album_id := track.Album.ID
	//album_uri := track.Album.URI
	//album_date := track.Album.ReleaseDate

	return []*tview.TableCell{name, artist, album, duration}
}

