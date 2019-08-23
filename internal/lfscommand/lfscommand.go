package lfscommand

import (
	"github.com/iikira/lfsutil/lfs"
	"log"
	"os"
)

const (
	EnvRepoURL = "LFS_REPO_URL"
	EnvAuth    = "LFS_AUTH"
)

var (
	LFS *lfs.LFS
)

func lazyInitCMD() {
	if LFS != nil {
		return
	}
	initCMD()
}

func initCMD() {
	lfsURL, ok := os.LookupEnv(EnvRepoURL)
	if !ok {
		log.Fatalf("env %s not set\n", EnvRepoURL)
	}
	auth, ok := os.LookupEnv(EnvAuth)
	if !ok {
		log.Fatalf("env %s not set\n", EnvAuth)
	}
	LFS = lfs.NewLFS()
	err := LFS.SetInfo(lfsURL, auth)
	if err != nil {
		log.Fatalln(err)
	}
}
