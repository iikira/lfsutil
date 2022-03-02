package lfs

import (
	"bytes"
	"errors"
	"github.com/iikira/iikira-go-utils/utils/jsonhelper"
	"io"
	"os"
	"path/filepath"
)

type (
	Object struct {
		FilePath string        `json:"-"` // custom
		OID      string        `json:"oid"`
		Size     int64         `json:"size"`
		Actions  *ObjectAction `json:"actions,omitempty"`
	}
	Ref struct {
		Name string `json:"name,omitempty"`
	}
	BatchPost struct {
		Operation string    `json:"operation,omitempty"`
		Objects   []*Object `json:"objects,omitempty"`
		Ref       *Ref      `json:"ref,omitempty"`
	}
	BatchResponse struct {
		Objects []*Object `json:"objects,omitempty"`
	}

	ObjectAction struct {
		Download     *TransferInfo `json:"download,omitempty"`
		Upload       *TransferInfo `json:"upload,omitempty"`
		ObjectVerify *ObjectVerify `json:"verify,omitempty"`
	}
	TransferInfo struct {
		Href      string            `json:"href,omitempty"`
		Header    map[string]string `json:"header,omitempty"`
		ExpiresAt string            `json:"expires_at,omitempty"`
	}
	ObjectVerify struct {
		Href   string            `json:"href,omitempty"`
		Header map[string]string `json:"header,omitempty"`
	}
)

func (l *LFS) getResponse(req *BatchPost) (r *BatchResponse, err error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	err = jsonhelper.MarshalData(buffer, &req)
	if err != nil {
		panic(err)
	}
	resp, err := l.c.Req("POST", l.genBatchURL().String(), buffer, l.makeHeader(true))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}

	if resp.StatusCode/100 != 2 {
		err = errors.New(resp.Status)
		return
	}

	batchResp := BatchResponse{}
	err = jsonhelper.UnmarshalData(resp.Body, &batchResp)
	if err != nil {
		return
	}
	return &batchResp, nil
}

func (l *LFS) GetDownloadObjectsResponse(objects ...*Object) (r *BatchResponse, err error) {
	l.lazyInit()
	req := BatchPost{
		Operation: "download",
		Objects:   objects,
		Ref: &Ref{
			Name: "refs/heads/master",
		},
	}

	return l.getResponse(&req)
}

func (l *LFS) GetUploadObjectsResponse(objects ...*Object) (r *BatchResponse, err error) {
	l.lazyInit()
	req := BatchPost{
		Operation: "upload",
		Objects:   objects,
		Ref: &Ref{
			Name: "refs/heads/master",
		},
	}

	return l.getResponse(&req)
}

func (l *LFS) getObjectsByReader(rs ...io.Reader) (objects []*Object, err error) {
	objects = make([]*Object, 0, len(rs))
	for _, r := range rs {
		object, err := GetObjectByReader(r)
		if err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}
	return objects, nil
}

func (l *LFS) getObjectsByFilename(filenames ...string) (objects []*Object, err error) {
	objects = make([]*Object, 0, len(filenames))
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		// 检查目录
		// 递归获取目录下的文件
		finfo, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if finfo.IsDir() {
			names, err := f.Readdirnames(-1)
			f.Close()
			if err != nil {
				return nil, err
			}
			for k := range names {
				names[k] = filepath.Join(filename, names[k])
			}
			subObjects, err := l.getObjectsByFilename(names...)
			if err != nil {
				return nil, err
			}
			objects = append(objects, subObjects...)
			continue
		}

		object, err := GetObjectByReader(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		f.Close()
		object.FilePath = filename // custom
		objects = append(objects, object)
	}
	return objects, nil
}

func (l *LFS) GetDownloadObjectsResponseByReader(rs ...io.Reader) (r *BatchResponse, err error) {
	objects, err := l.getObjectsByReader(rs...)
	if err != nil {
		return nil, err
	}
	return l.GetDownloadObjectsResponse(objects...)
}

func (l *LFS) customObjects(remoteObjects []*Object, localObjects []*Object) error {
	if len(remoteObjects) != len(localObjects) {
		return ErrObjectLengthNotEqual
	}
	for k := range localObjects {
		remoteObjects[k].FilePath = localObjects[k].FilePath
	}
	return nil
}

func (l *LFS) GetDownloadObjectsResponseByFilename(filenames ...string) (r *BatchResponse, err error) {
	objects, err := l.getObjectsByFilename(filenames...)
	if err != nil {
		return nil, err
	}
	r, err = l.GetDownloadObjectsResponse(objects...)
	if err != nil {
		return nil, err
	}
	err = l.customObjects(r.Objects, objects)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (l *LFS) GetUploadObjectsResponseByReader(rs ...io.Reader) (r *BatchResponse, err error) {
	objects, err := l.getObjectsByReader(rs...)
	if err != nil {
		return nil, err
	}
	return l.GetUploadObjectsResponse(objects...)
}

func (l *LFS) GetUploadObjectsResponseByFilename(filenames ...string) (r *BatchResponse, err error) {
	objects, err := l.getObjectsByFilename(filenames...)
	if err != nil {
		return nil, err
	}
	r, err = l.GetUploadObjectsResponse(objects...)
	if err != nil {
		return nil, err
	}
	err = l.customObjects(r.Objects, objects)
	if err != nil {
		return nil, err
	}
	return r, nil
}
