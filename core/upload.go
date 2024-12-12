package core

import (
	"PanUpload/internal"
	"github.com/caiguanhao/opencc/configs/t2s"
	client "github.com/hefeiyu2025/pan-client"
	"github.com/hefeiyu2025/pan-client/pan"
	"regexp"
	"strings"
)

func uploadInit() pan.Driver {
	internal.InitConfig()
	c, err := client.GetClient(pan.DriverType(internal.Config.Upload.UploadClient))
	if err != nil {
		panic(err)
	}
	return c
}

func StartUpload() {
	driver := uploadInit()
	uploadConfig := internal.Config.Upload
	remoteTransfer := func(remote string) string {
		newRemote := remote
		for _, removeStr := range uploadConfig.RemoveStr {
			newRemote = strings.ReplaceAll(newRemote, removeStr, "")
		}
		// 使用正则表达式替换字符串
		if uploadConfig.RemoveReg != "" {
			re := regexp.MustCompile(uploadConfig.RemoveReg)
			newRemote = re.ReplaceAllString(newRemote, "")
		}

		newRemote = strings.TrimSpace(newRemote)
		newRemote = t2s.Dicts.Convert(newRemote)

		return newRemote
	}
	paths := uploadConfig.LocalPath
	for _, path := range paths {
		err := driver.UploadPath(pan.UploadPathReq{
			LocalPath:          path,
			RemotePath:         uploadConfig.RemotePath,
			Resumable:          true,
			SkipFileErr:        true,
			OnlyFast:           uploadConfig.OnlyFast,
			SuccessDel:         uploadConfig.SuccessDelete,
			IgnorePaths:        uploadConfig.IgnorePath,
			Extensions:         uploadConfig.UploadExtension,
			RemotePathTransfer: remoteTransfer,
			RemoteNameTransfer: remoteTransfer,
		})

		if err != nil {
			panic(err)
		}
	}

}
