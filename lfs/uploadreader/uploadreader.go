package uploadreader

import (
	"errors"
	"github.com/iikira/iikira-go-utils/requester/rio"
	"github.com/iikira/iikira-go-utils/requester/rio/speeds"
	"io"
	"os"
	"sync/atomic"
	"time"
)

type (
	UploadReader struct {
		readed    int64
		speedStat speeds.Speeds
		r         io.Reader
		length    int64
		callback  UploadCallback
		interval  time.Duration
		stopChan  chan struct{}
		isStart   uint32
		isStop    uint32
		nowTime   *time.Time
	}

	UploadCallback func(state TransferState)
)

func New() *UploadReader {
	return &UploadReader{}
}

func NewWithReader(r io.Reader, length int64) *UploadReader {
	return &UploadReader{
		r:      r,
		length: length,
	}
}

func (ur *UploadReader) lazyInit() {
	if ur.interval <= 0 {
		ur.interval = time.Second
	}
	if ur.stopChan == nil {
		ur.stopChan = make(chan struct{})
	}
}

func (ur *UploadReader) SetReader(r io.Reader, length int64) {
	ur.r = r
	ur.length = length
}

func (ur *UploadReader) SetFile(f *os.File) error {
	if f == nil {
		return errors.New("f is nil")
	}

	ur.r = f
	info, err := f.Stat()
	if err != nil {
		return err
	}
	ur.length = info.Size()
	return nil
}

func (ur *UploadReader) SetReaderLen64(rl rio.ReaderLen64) {
	ur.r = rl
	ur.length = rl.Len()
}

func (ur *UploadReader) SetCallback(callback UploadCallback) {
	ur.callback = callback
}

func (ur *UploadReader) SetInterval(interval time.Duration) {
	if interval > 0 {
		ur.interval = interval
	}
}

func (ur *UploadReader) Read(p []byte) (n int, err error) {
	n, err = ur.r.Read(p)
	if err != nil {
		ur.Stop()
		return
	}
	ni := int64(n)
	ur.speedStat.Add(ni)
	atomic.AddInt64(&ur.readed, ni)
	return
}

func (ur *UploadReader) Len() int64 {
	return ur.length
}

func (ur *UploadReader) Start() (stopChan <-chan struct{}) {
	c := make(chan struct{})
	stopChan = c
	if atomic.AddUint32(&ur.isStart, 1) != 1 { // 只允许执行一次
		close(c)
		return
	}

	ur.lazyInit()
	if ur.callback == nil {
		close(c)
		return
	}
	ticker := time.NewTicker(ur.interval)
	nowTime := time.Now()
	ur.nowTime = &nowTime
	go func() {
		for {
			select {
			case <-ticker.C:
				ur.callback(ur)
			case <-ur.stopChan:
				ur.callback(ur)
				close(c)
				return
			}
		}
	}()
	return
}

func (ur *UploadReader) Stop() {
	if ur.stopChan == nil {
		return
	}
	if atomic.AddUint32(&ur.isStop, 1) != 1 { // 只允许执行一次
		return
	}
	close(ur.stopChan)
}
