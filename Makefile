go.mod:
	go mod init git.dominic-ricottone.com/~dricottone/nspotify
	go get github.com/zmb3/spotify/v2
	go get github.com/sirupsen/logrus
	go get golang.org/x/oauth2
	go get github.com/rivo/tview

GO_SRC!=find * -type f -name '*.go'

GO_LDFLAGS:=
GO_LDFLAGS+=-X main.CLIENTID=$(file < clientid.txt)
GO_LDFLAGS+=-X main.CLIENTSECRET=$(file < clientsecret.txt)
# TODO: GO_LDFLAGS+=-X main.VERSION=0.0.1

nspotify: go.mod $(GO_SRC)
	go build -ldflags "$(GO_LDFLAGS)" .

build: nspotify

clean:
	rm -f go.mod nspotify

.PHONY: run build clean
