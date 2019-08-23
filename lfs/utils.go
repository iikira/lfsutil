package lfs

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"github.com/iikira/BaiduPCS-Go/requester"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

func GetObjectByReader(r io.Reader) (*Object, error) {
	var (
		o = Object{
			Size: 0,
		}
		s      = sha256.New()
		buf    = make([]byte, 8192)
		handle = func(n int) {
			s.Write(buf[:n])
			o.Size += int64(n)
		}
	)
	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				handle(n)
				break
			}
			return nil, err
		}
		handle(n)
	}

	sum := s.Sum(nil)
	o.OID = hex.EncodeToString(sum)
	return &o, nil
}

func TestLink(link string) (ok bool, err error) {
	c := requester.NewHTTPClient()
	resp, err := c.Req("GET", link, nil, nil)
	if resp != nil {
		resp.Body.Close()
	}
	if err != nil {
		return
	}

	if resp.StatusCode/100 == 2 {
		ok = true
		return
	}
	return
}

// GetPtrContent support both local file name and remote url
func GetPtrContent(arg string) (rc io.ReadCloser, err error) {
	_, err = url.Parse(arg)
	if err == nil {
		c := requester.NewHTTPClient()
		resp, err := c.Req("GET", arg, nil, nil)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	}

	f, err := os.Open(arg)
	if err != nil {
		return
	}
	rc = f
	return
}

func ParsePtrContent(r io.Reader) (*Object, error) {
	var (
		oidKeyword  = []byte("oid sha256:")
		sizeKeyword = []byte("size ")
		b           = bufio.NewReader(r)
		object      = Object{}
	)
	for {
		line, _, err := b.ReadLine()
		if err != nil {
			break
		}
		if object.OID == "" && bytes.Index(line, oidKeyword) == 0 {
			object.OID = converter.ToString(line[len(oidKeyword):])
		}
		if object.Size == 0 && bytes.Index(line, sizeKeyword) == 0 {
			sizeStr := converter.ToString(line[len(sizeKeyword):])
			object.Size, err = strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				return nil, ErrParseObjectSize
			}
		}
	}
	if object.OID == "" {
		return nil, ErrParseObjectContentOIDNotFound
	}
	return &object, nil
}

func ParsePtr(arg string) (*Object, error) {
	rc, err := GetPtrContent(arg)
	if rc != nil {
		defer rc.Close()
	}
	if err != nil {
		return nil, err
	}
	return ParsePtrContent(rc)
}

func ParsePtrs(args ...string) ([]*Object, error) {
	objects := make([]*Object, 0, len(args))
	for _, arg := range args {
		o, err := ParsePtr(arg)
		if err != nil {
			return nil, err
		}

		objects = append(objects, o)
	}
	return objects, nil
}

func WritePtrSpecToWriter(w io.Writer, object *Object) error {
	t, err := template.New("v1").Parse(`version https://git-lfs.github.com/spec/v1
oid sha256:{{.OID}}
size {{.Size}}
`)
	if err != nil {
		panic(err)
	}
	return t.Execute(w, object)
}

func WritePtrSpecToFile(filename string, object *Object) (err error) {
	os.MkdirAll(filepath.Dir(filename), 0644)
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	return WritePtrSpecToWriter(f, object)
}
