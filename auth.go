package main

// Prompt end users to authenticate with Spotify. Needed at startup.

import (
	"context"
	"net/http"
	"fmt"

	"github.com/zmb3/spotify/v2"
	auth "github.com/zmb3/spotify/v2/auth"
	log "github.com/sirupsen/logrus"
)

// These should be set by the linker.
// e.g. `go build -ldflags="-X 'main.CLIENTID=foobarbaz'"`
var (
	CLIENTID = ""
	CLIENTSECRET = ""
)

// Pull the cache directory information from the context.
func cache_info(ctx context.Context) string {
	dir, ok := ctx.Value("cachedir").(string)
	if !ok {
		dir = ""
	}

	return dir
}

// Pull the authentication URI information from the context.
func uri_info(ctx context.Context) (string, string) {
	port, ok := ctx.Value("authport").(int)
	if !ok {
		log.Warn("authenticating web server port is required; falling back to default :8080")
		port = 8080
	}

	full_uri := fmt.Sprintf("http://localhost:%d", port)
	short_uri := fmt.Sprintf(":%d", port)

	return full_uri, short_uri
}

// Serves the authenticator at `http://localhost:[authport]`. Prints that address
// and some instructions to STDOUT for the end user.
func ServeAuthenticator(ctx context.Context, ch chan<- *spotify.Client) *http.Server {
	full_uri, short_uri := uri_info(ctx)
	state := "nspotify"
	authenticator := auth.New(
		auth.WithClientID(CLIENTID),
		auth.WithClientSecret(CLIENTSECRET),
		auth.WithRedirectURL(full_uri),
		auth.WithScopes(auth.ScopeUserLibraryRead, auth.ScopeUserReadPlaybackState, auth.ScopeUserModifyPlaybackState))
	srv := &http.Server{Addr: short_uri}

	// Address and instructions for end user.
	fmt.Println("Log in to Spotify at:", authenticator.AuthURL(state))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tok, err := authenticator.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Failed to get token", http.StatusForbidden)
			log.WithError(err).Fatal("failed to get token")
		}

		if s := r.FormValue("state"); s != state {
			http.NotFound(w, r)
			log.Fatalf("invalid state: %s\n", s)
		}

		log.Debugf("login succeeded")

		client := spotify.New(authenticator.Client(ctx, tok))

		dir := cache_info(ctx)
		if dir != "" {
			WriteCache(dir, tok)
		}

		ch <- client
	})

	// Start server.
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.WithError(err).Fatal("authenticating web server died")
		}
	}()

	return srv
}

// Authenticate with a cached token from `[cachedir]/token.json`.
func CachedAuthentication(ctx context.Context) *spotify.Client {
	dir := cache_info(ctx)
	if dir == "" {
		return nil
	}

	tok, err := ReadCache(dir)
	if err != nil {
		return nil
	}

	authenticator := auth.New(auth.WithScopes(auth.ScopeUserLibraryRead))

	return spotify.New(authenticator.Client(ctx, tok))
	
}

// Authenticates the end user.
func Authenticate(ctx context.Context) *spotify.Client {
	// Try cache first.
	if cached := CachedAuthentication(ctx); cached != nil {
		return cached
	}

	// Start server.
	ch := make(chan *spotify.Client)
	srv := ServeAuthenticator(ctx, ch)

	cli := <-ch

	// Kill server.
	srv.Shutdown(ctx)

	return cli
}

