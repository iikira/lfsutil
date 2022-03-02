package lfs

import (
	"bytes"
	"errors"
	"github.com/iikira/iikira-go-utils/requester/rio"
	"github.com/iikira/iikira-go-utils/utils/jsonhelper"
	"github.com/iikira/lfsutil/lfs/uploadreader"
	"io/ioutil"
	"os"
)

var (
	ErrInfoNil         = errors.New("info is nil")
	ErrObjectNil       = errors.New("object is nil")
	ErrObjectVerifyNil = errors.New("object verify is nil")
)

func (l *LFS) UploadWithLinkByFile(filename string, info *TransferInfo, callback uploadreader.UploadCallback) (msg []byte, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	rl := rio.NewFileReaderLen64(f)
	return l.UploadWithLinkByReader(rl, info, callback)
}

func (l *LFS) UploadWithLinkByReader(rl rio.ReaderLen64, info *TransferInfo, callback uploadreader.UploadCallback) (msg []byte, err error) {
	if info == nil {
		err = ErrInfoNil
		return
	}

	l.lazyInit()
	ur := uploadreader.New()
	ur.SetReaderLen64(rl)
	ur.SetCallback(callback)
	ur.Start()
	resp, err := l.c.Req("PUT", info.Href, ur, info.Header)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}

	msg, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode/100 != 2 {
		err = errors.New(resp.Status)
		return
	}
	return
}

func (l *LFS) Verify(o *Object, v *ObjectVerify) (msg []byte, err error) {
	if o == nil {
		err = ErrObjectNil
		return
	}
	if v == nil {
		err = ErrObjectVerifyNil
		return
	}
	l.lazyInit()
	newO := Object{
		OID:  o.OID,
		Size: o.Size,
	}
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	err = jsonhelper.MarshalData(buf, &newO)
	if err != nil {
		panic(err)
	}

	resp, err := l.c.Req("POST", v.Href, buf, v.Header)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}

	msg, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode/100 != 2 {
		err = errors.New(resp.Status)
		return
	}
	return
}
