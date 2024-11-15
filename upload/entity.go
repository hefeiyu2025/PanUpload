package upload

import (
	"io"
	"time"
)

type Resp[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *T     `json:"data"`
}

type Object struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Path          string    `json:"path"`
	Pic           string    `json:"pic"`
	Size          int       `json:"size"`
	Type          string    `json:"type"`
	Date          time.Time `json:"date"`
	CreateDate    time.Time `json:"create_date"`
	SourceEnabled bool      `json:"source_enabled"`
}
type Policy struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	MaxSize  int      `json:"max_size"`
	FileType []string `json:"file_type"`
}

type DirectoryResp struct {
	Parent  string   `json:"parent"`
	Objects []Object `json:"objects"`
	Policy  Policy   `json:"policy"`
}
type USessionInfo struct {
	SessionID  string   `json:"sessionID"`
	ChunkSize  int      `json:"chunkSize"`
	Expires    int      `json:"expires"`
	UploadURLs []string `json:"UploadURLs"`
}

type USessionReq struct {
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	Name         string `json:"name"`
	PolicyId     string `json:"policy_id"`
	LastModified int64  `json:"last_modified"`
}

type ChunkData struct {
	StartSize   int
	EndSize     int
	ChunkSize   int
	TotalSize   int
	ChunkNum    int
	chunkReader io.Reader
}

// ObjectExistError 自定义错误类型
type ObjectExistError struct {
	Message string
}

// Implement the Error method to satisfy the error interface.
func (e ObjectExistError) Error() string {
	return e.Message
}
