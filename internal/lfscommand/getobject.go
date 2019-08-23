package lfscommand

import (
	"fmt"
	"github.com/iikira/lfsutil/lfs"
	"log"
)

type (
	GetObjectOption struct {
		Args      []string
		By        string
		NoTest    bool
		PtrDir    string // TODO
		PtrSuffix string // TODO
	}
)

func GetObject(opt *GetObjectOption) {
	lazyInitCMD()
	var (
		r   *lfs.BatchResponse
		err error
	)
	switch opt.By {
	case "oid":
		r, err = LFS.GetDownloadObjectsResponse(ParseInputObject(opt.Args...)...)
	case "file":
		r, err = LFS.GetDownloadObjectsResponseByFilename(opt.Args...)
	case "ptr":
		objects, err := lfs.ParsePtrs(opt.Args...)
		if err != nil {
			log.Fatalln(err)
		}
		r, err = LFS.GetDownloadObjectsResponse(objects...)
	default:
		log.Fatalf("unknown by: %s\n", opt.By)
	}
	if err != nil {
		log.Fatalln(err)
	}
	for k, o := range r.Objects {
		var filenamePrefix string
		if o.FilePath != "" {
			filenamePrefix = "[" + o.FilePath + "] "
		}

		if o.Actions == nil {
			fmt.Printf("[%d] %saction not found\n", k, filenamePrefix)
			continue
		}

		if o.Actions.Download == nil {
			fmt.Printf("[%d] %sdownload action not found\n", k, filenamePrefix)
			continue
		}

		var ok bool
		if !opt.NoTest {
			ok, err = lfs.TestLink(o.Actions.Download.Href)
			if err != nil {
				fmt.Printf("[%d] %shref test error: %s\n", k, filenamePrefix, err)
				continue
			}
		} else {
			ok = true
		}

		if ok {
			fmt.Printf("[%d] %soid: %s, size: %d, link: %s\n", k, filenamePrefix, o.OID, o.Size, o.Actions.Download.Href)
			continue
		}
		fmt.Printf("[%d] %shref test failed\n", k, filenamePrefix)
	}
	return
}
