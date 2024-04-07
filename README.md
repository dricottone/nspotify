# nspotify

Minimal TUI spotify client.
Read-only interface to a user's saved songs.
Playback sold separately,
see [spotifyd](https://github.com/Spotifyd/spotifyd) or
[go-librespot](https://github.com/devgianlu/go-librespot).


## Building

Depends on a small set of libraries:

 + [zmb3's spotify API wrapper](github.com/zmb3/spotify)
 + [logrus](github.com/sirupsen/logrus) for logging
 + [tview](github.com/rivo/tview) and [tcell](https://github.com/gdamore/tcell)
   for the TUI
 + external `oauth2` package

A Spotify API client ID and secret are required.
See [here](https://github.com/zmb3/spotify?tab=readme-ov-file#authentication)
for more details.
Then try:

```
go build -ldflags "-X main.CLIENTID=yourclientid -X main.CLIENTSECRET=yourclientsecret" .
```


## Licensing

I share the contents of this repository under the BSD 3 clause license.


