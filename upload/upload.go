package upload

import (
	"PanUpload/cmd/flags"
	"fmt"
	"github.com/caiguanhao/opencc/configs/t2s"
	"github.com/hefeiyu2025/cloudreve-client"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const cloudreveUrl = "https://pan.huang1111.cn"

var cloudreveClient *cloudreve.CloudreveClient

func initCloudreveSession() string {
	session, err := os.Open(flags.SessionPath)
	if err != nil {
		fmt.Println("open session file error:", flags.SessionPath, err)
		os.Exit(1)
	}
	defer session.Close()

	data, err := io.ReadAll(session)
	if err != nil {
		fmt.Println("Error reading session file:", flags.SessionPath, err)
		os.Exit(1)
	}

	if len(data) == 0 {
		fmt.Println("Error reading session file:", flags.SessionPath, ",file is empty")
		os.Exit(1)
	}
	fmt.Println("read session file:", flags.SessionPath, " success, session:", string(data))
	return string(data)
}

func refreshCloudreveSession(cloudreveSession string) {
	session, err := os.OpenFile(flags.SessionPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("open session file error:", flags.SessionPath, err)
		os.Exit(1)
	}
	defer session.Close()

	// 写入数据到文件
	_, err = session.WriteString(cloudreveSession)
	if err != nil {
		fmt.Println("Error writing to session file:", flags.SessionPath, err)
		os.Exit(1)
	}
	fmt.Println("Success refresh session file:", flags.SessionPath)
}

func initClient() {
	cloudreveClient = cloudreve.NewClientWithRefresh(cloudreveUrl, initCloudreveSession(), func(session string) {
		refreshCloudreveSession(session)
	})
	//cloudreveClient = cloudreve.NewClient(cloudreveUrl, initCloudreveSession())
}

func exitByError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func StartUpload(file string) {
	initClient()
	if flags.DeleteAllSession {
		resp, err := DeleteAllSession()
		if err != nil {
			fmt.Println(resp, err)
		} else {
			fmt.Println("delete all session success")
			err := os.RemoveAll(os.TempDir() + "/cloudreve-cache")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("delete cloudreve-cache success", os.TempDir()+"/cloudreve-cache")
			}

		}
	}
	directoryResp, err := cloudreveClient.ListDirectory(flags.RemotePath)
	exitByError(err)
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
			err = cloudreveClient.UploadFile(cloudreve.OneStepUploadFileReq{
				LocalFile:  file,
				RemotePath: flags.RemotePath,
				PolicyId:   directoryResp.Data.Policy.ID,
				Resumable:  true,
				SuccessDel: flags.Delete,
				RemoteTransfer: func(remotePath, remoteName string) (string, string) {
					return t2s.Dicts.Convert(remotePath), t2s.Dicts.Convert(remoteName)
				},
			})
			exitByError(err)
			os.Exit(2)
		}
	}
	_, sessionName := filepath.Split(flags.SessionPath)
	err = cloudreveClient.UploadPath(cloudreve.OneStepUploadPathReq{
		LocalPath:   root,
		RemotePath:  flags.RemotePath,
		PolicyId:    directoryResp.Data.Policy.ID,
		Resumable:   true,
		SkipFileErr: true,
		SuccessDel:  flags.Delete,
		IgnorePaths: flags.GetIgnorePaths(),
		IgnoreFiles: []string{sessionName},
		Extensions:  flags.GetExtensions(),
		RemoteTransfer: func(remotePath, remoteName string) (string, string) {
			newFileName := remoteName
			newRemotePath := remotePath
			for _, removeStr := range flags.GetRemoveStrs() {
				newFileName = strings.ReplaceAll(newFileName, removeStr, "")
				newRemotePath = strings.ReplaceAll(newRemotePath, removeStr, "")
			}
			// 使用正则表达式替换字符串
			re := regexp.MustCompile(flags.RemoveReg)
			newRemotePath = re.ReplaceAllString(newRemotePath, "")

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

func DeleteAllSession() (*cloudreve.Resp, error) {
	initClient()
	return cloudreveClient.FileUploadDeleteAllUploadSession()
}
