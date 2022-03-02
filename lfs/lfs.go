package lfs

import (
	"github.com/iikira/iikira-go-utils/requester"
	"net/url"
	"time"
)

type (
	LFS struct {
		lfsURL *url.URL
		auth   string
		c      *requester.HTTPClient
	}
)

func NewLFS() *LFS {
	return &LFS{}
}

func (l *LFS) SetInfo(lfsURL, auth string) (err error) {
	l.lfsURL, err = url.Parse(lfsURL)
	if err != nil {
		return
	}
	l.auth = auth
	return
}

func (l *LFS) lazyInit() {
	if l.c == nil {
		l.c = requester.NewHTTPClient()
		l.c.SetTimeout(0)
		l.c.SetResponseHeaderTimeout(1 * time.Hour)
	}
}

func (l *LFS) genBatchURL() *url.URL {
	u := *l.lfsURL
	u.Path += "/objects/batch"
	return &u
}

func (l *LFS) makeHeader(isJSON bool) map[string]string {
	var accept string
	if isJSON {
		accept = "application/vnd.git-lfs+json"
	} else {
		accept = "application/vnd.git-lfs"
	}
	return map[string]string{
		"Accept":        accept,
		"Authorization": l.auth,
	}
}
