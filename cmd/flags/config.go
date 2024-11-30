package flags

import (
	"strings"
)

var (
	RemotePath       string
	SessionPath      string
	CachePath        string
	Debug            bool
	Delete           bool
	DeleteAllSession bool
	RemoveStr        string
	RemoveReg        string
	UploadExtensions string
	IgnorePath       string
)

func GetRemoveStrs() []string {
	return strings.Split(RemoveStr, ",")
}

func GetExtensions() []string {
	return strings.Split(UploadExtensions, ",")
}

func GetIgnorePaths() []string {
	return strings.Split(IgnorePath, ",")
}
