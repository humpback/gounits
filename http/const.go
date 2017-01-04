package http

import "errors"

var (
	ErrFilePathInvalid     = errors.New("file path invalid.")
	ErrFileOpenException   = errors.New("file open error.")
	ErrFileStatException   = errors.New("file stat error.")
	ErrHttpNewRequest      = errors.New("http request error.")
	ErrHttpRequestFailed   = errors.New("http request failed.")
	ErrHttpRequestInvalid  = errors.New("http request object invalid.")
	ErrHttpResponseInvalid = errors.New("http response object invalid.")
	ErrHttpIOCopyFailed    = errors.New("http iocpoy response failed.")
)
