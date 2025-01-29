SRC=$(shell find . -type f -name '*.go')

clean:
	rm -f go.mod go.sum nspotify

go.mod:
	go mod init git.sr.ht/~dricottone/nspotify
	go get github.com/zmb3/spotify/v2
	go get github.com/sirupsen/logrus
	go get golang.org/x/oauth2
	go get github.com/rivo/tview

GO_LDFLAGS:=
GO_LDFLAGS+=-X main.CLIENTID=$(file < clientid.txt)
GO_LDFLAGS+=-X main.CLIENTSECRET=$(file < clientsecret.txt)
# TODO: GO_LDFLAGS+=-X main.VERSION=0.0.1

nspotify: go.mod $(SRC)
	go build -ldflags "$(GO_LDFLAGS)" .

build: nspotify

.PHONY: clean build
