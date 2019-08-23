package lfs

import (
	"errors"
)

var (
	ErrObjectLengthNotEqual          = errors.New("object length not equal")
	ErrParseObjectContentOIDNotFound = errors.New("parse object: oid not found")
	ErrParseObjectSize               = errors.New("parse object: convert size error")
)
