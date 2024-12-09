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
	c, err := client.GetClient(pan.Cloudreve)
	if err != nil {
		panic(err)
	}
	return c
}

func StartUpload() {
	cloudreveDriver := uploadInit()
	uploadConfig := internal.Config.Upload
	err := cloudreveDriver.UploadPath(pan.UploadPathReq{
		LocalPath:   uploadConfig.LocalPath,
		RemotePath:  uploadConfig.RemotePath,
		Resumable:   true,
		SkipFileErr: true,
		SuccessDel:  uploadConfig.SuccessDelete,
		IgnorePaths: uploadConfig.IgnorePath,
		Extensions:  uploadConfig.UploadExtension,
		RemoteTransfer: func(remotePath, remoteName string) (string, string) {
			newFileName := remoteName
			newRemotePath := remotePath
			for _, removeStr := range uploadConfig.RemoveStr {
				newFileName = strings.ReplaceAll(newFileName, removeStr, "")
				newRemotePath = strings.ReplaceAll(newRemotePath, removeStr, "")
			}
			// 使用正则表达式替换字符串
			if uploadConfig.RemoveReg != "" {
				re := regexp.MustCompile(uploadConfig.RemoveReg)
				newRemotePath = re.ReplaceAllString(newRemotePath, "")
			}

			newFileName = strings.TrimSpace(newFileName)
			newFileName = t2s.Dicts.Convert(newFileName)
			newRemotePath = strings.TrimSpace(newRemotePath)
			newRemotePath = t2s.Dicts.Convert(newRemotePath)

			return newRemotePath, newFileName
		},
	})

	if err != nil {
		panic(err)
	}

}
