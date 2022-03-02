package uploadreader

import (
	"sync/atomic"
	"time"
)

type (
	TransferState interface {
		Transfered() int64
		Total() int64
		SpeedPerSecond() int64
		Elapsed() time.Duration
	}
)

func (ur *UploadReader) Transfered() int64 {
	return atomic.LoadInt64(&ur.readed)
}
func (ur *UploadReader) Total() int64 {
	return ur.length
}
func (ur *UploadReader) SpeedPerSecond() int64 {
	return ur.speedStat.GetSpeeds()
}
func (ur *UploadReader) Elapsed() time.Duration {
	if ur.nowTime == nil {
		return 0
	}
	return time.Since(*ur.nowTime)
}
