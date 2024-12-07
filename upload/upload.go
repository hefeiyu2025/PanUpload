package upload

import (
	"PanUpload/internal"
	"fmt"
	"github.com/caiguanhao/opencc/configs/t2s"
	client "github.com/hefeiyu2025/pan-client"
	"github.com/hefeiyu2025/pan-client/pan"
	"os"
	"regexp"
	"strings"
)

var cloudreveDriver pan.Driver

func initClient() {
	internal.InitConfig()
	c, err := client.GetClient(pan.Cloudreve)
	if err != nil {
		panic(err)
	}
	cloudreveDriver = c
}

func exitByError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func StartUpload(file string) {
	initClient()

	// 默认为当前目录
	root := "./"
	//判断参数是否为空，为空则上传当前目录，否则上传指定文件
	if file != "" {
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Println("file read error,", file, err)
			return
		}
		if fileInfo.IsDir() {
			root = file
		} else {
			err = cloudreveDriver.UploadFile(pan.OneStepUploadFileReq{
				LocalFile:  file,
				RemotePath: internal.Config.RemotePath,
				Resumable:  true,
				SuccessDel: internal.Config.SuccessDelete,
				RemoteTransfer: func(remotePath, remoteName string) (string, string) {
					return t2s.Dicts.Convert(remotePath), t2s.Dicts.Convert(remoteName)
				},
			})
			exitByError(err)
			return
		}
	}
	err := cloudreveDriver.UploadPath(pan.OneStepUploadPathReq{
		LocalPath:   root,
		RemotePath:  internal.Config.RemotePath,
		Resumable:   true,
		SkipFileErr: true,
		SuccessDel:  internal.Config.SuccessDelete,
		IgnorePaths: internal.Config.IgnorePath,
		Extensions:  internal.Config.UploadExtension,
		RemoteTransfer: func(remotePath, remoteName string) (string, string) {
			newFileName := remoteName
			newRemotePath := remotePath
			for _, removeStr := range internal.Config.RemoveStr {
				newFileName = strings.ReplaceAll(newFileName, removeStr, "")
				newRemotePath = strings.ReplaceAll(newRemotePath, removeStr, "")
			}
			// 使用正则表达式替换字符串
			if internal.Config.RemoveReg != "" {
				re := regexp.MustCompile(internal.Config.RemoveReg)
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
