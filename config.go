package main

import (
	"context"
	"flag"

	log "github.com/sirupsen/logrus"
)

// Number of tracks to fetch into a buffer.
const fetchingBuffer = 100

// Number of seconds to wait before rechecking the lazy fetcher.
const fetchingTimeout = 3

// Number of tracks to eagerly load.
const loadEager = 50

// Number of tracks to lazily load ahead of the cursor.
const loadLookahead = 50

// Number of seconds to wait before rechecking how many trackers are loaded
// ahead of the cursor.
const loadTimeout = 3

// Command line options.
var (
	verbose = flag.Bool("verbose", false, "Show debugging messages (same as '-log-level=trace')")
	quiet = flag.Bool("quiet", false, "Suppress messages (same as '-log-level=panic')")
	log_level = flag.String("log-level", "", "Set logging level to `[trace|debug|info|warning|error|fatal|panic]`")
	//TODO: log_file = flag.String("log-file", "", "Log to `file`")
	color = flag.Bool("color", true, "Display in color")
	port = flag.Int("port", 8080, "Spotify authenticator port")
	cache = flag.String("cache", "", "Cache `directory`")
	no_cache = flag.Bool("no-cache", false, "Do not use cached authentication, do not cache authentication")
	device = flag.String("device", "", "Spotify device `ID`")
	list_devices = flag.Bool("list-devices", false, "List available Spotify devices and exit")
)

// Credit to `https://stackoverflow.com/a/30716481`:
func Pointer[T any](v T) *T {
	return &v
}

func ConfiguredContext() context.Context {
	ctx := context.Background()

	flag.Parse()

	// If `-color` or `-color=true:
	if *color {
		log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: false})
	} else {
		log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: false})
	}

	// Prioritize explicit `-log-level`, then `-quiet`, then `-verbose`.
	if *log_level == "" {
		if *verbose {
			log_level = Pointer("trace")
		}
		if *quiet {
			log_level = Pointer("panic")
		}
	}
	switch *log_level {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}

	ctx = context.WithValue(ctx, "authport", *port)

	if *cache == "" {
		cache = Pointer(default_cache_dir())
	}

	// Signal `-no-cache` by forcing `-cache=''`.
	if *no_cache {
		cache = Pointer("")
	}
	ctx = context.WithValue(ctx, "cachedir", *cache)

	// Signal `-list-devices` by forcing `-device=''`.
	if *list_devices {
		device = Pointer("")
	}
	ctx = context.WithValue(ctx, "device", *device)

	return ctx
}

