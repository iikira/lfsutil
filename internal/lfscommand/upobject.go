package lfscommand

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"github.com/iikira/lfsutil/lfs"
	"github.com/iikira/lfsutil/lfs/uploadreader"
	"log"
	"os"
)

type (
	UpObjectOption struct {
		Args      []string
		PtrDir    string
		PtrSuffix string
		NoVerify  bool
	}
)

func UpObject(opt *UpObjectOption) {
	lazyInitCMD()
	_, statErr := os.Stat(opt.PtrDir)
	if statErr != nil {
		fmt.Printf("%s\n", statErr)
	}
	for j, arg := range opt.Args {
		r, err := LFS.GetUploadObjectsResponseByFilename(arg)
		if err != nil {
			log.Fatalln(err)
		}

		for k, o := range r.Objects {
			if o.Actions == nil {
				fmt.Printf("[%d] - [%d] [%s] action not found, mayby file already uploaded, oid: %s, size: %d\n", j, k, o.FilePath, o.OID, o.Size)
				afterDownloadLink(j, k, o, statErr, arg, opt)
				continue
			}

			if o.Actions.Upload == nil {
				fmt.Printf("[%d] - [%d] [%s] upload url not found, mayby file already uploaded, oid: %s, size: %d", j, k, o.FilePath, o.OID, o.Size)
				afterDownloadLink(j, k, o, statErr, arg, opt)
				continue
			}

			fmt.Printf("[%d] - [%d] [%s] upload start\n", j, k, o.FilePath)
			msg, err := LFS.UploadWithLinkByFile(o.FilePath, o.Actions.Upload, getUploadCallback(k))
			if err != nil {
				fmt.Printf("[%d] - [%d] [%s] failed, msg: %s, err: %s\n", j, k, o.FilePath, msg, err)
				continue
			}
			fmt.Printf("[%d] - [%d] [%s] success, oid: %s, size: %d", j, k, o.FilePath, o.OID, o.Size)
			afterDownloadLink(j, k, o, statErr, arg, opt)
			// Verify
			if !opt.NoVerify {
				if o.Actions.ObjectVerify == nil {
					fmt.Printf("[%d] - [%d] [%s] server unsupport verify\n", j, k, o.FilePath)
					continue
				}
				msg, err := LFS.Verify(o, o.Actions.ObjectVerify)
				if err != nil {
					fmt.Printf("[%d] - [%d] [%s] verify failed, msg: %s, err: %s\n", j, k, o.FilePath, msg, err)
					continue
				}
			}
		}
	}
}

func afterDownloadLink(j, k int, o *lfs.Object, statErr error, arg string, opt *UpObjectOption) {
	var err error
	if o.Actions == nil || o.Actions.Upload == nil {
		// 无上传对象, 写入ptr
		if statErr == nil {
			err = lfs.WritePtrSpecToFile(ParseOutName(arg, o.FilePath, opt.PtrDir)+opt.PtrSuffix, o)
		}
	}

	if o.Actions != nil && o.Actions.Download != nil {
		fmt.Printf(", download link: %s", o.Actions.Download.Href)
	}
	fmt.Println()
	if err != nil {
		fmt.Printf("[%d] - [%d] [%s] write error: %s\n", j, k, o.FilePath, err)
	}
}

func getUploadCallback(k int) uploadreader.UploadCallback {
	return func(state uploadreader.TransferState) {
		fmt.Printf("[%d] ↑ %s/%s %s/s in %s\n", k,
			converter.ConvertFileSize(state.Transfered(), 2),
			converter.ConvertFileSize(state.Total(), 2),
			converter.ConvertFileSize(state.SpeedPerSecond(), 2),
			state.Elapsed()/1e7*1e7,
		)
	}
}
