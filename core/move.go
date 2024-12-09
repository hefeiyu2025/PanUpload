package core

import (
	"PanUpload/internal"
	"fmt"
	"github.com/caiguanhao/opencc/configs/t2s"
	client "github.com/hefeiyu2025/pan-client"
	"github.com/hefeiyu2025/pan-client/pan"
	"path/filepath"
	"regexp"
	"strings"
)

func moveInit() (pan.Driver, pan.Driver) {
	internal.InitConfig()
	c, err := client.GetClient(pan.Cloudreve)
	if err != nil {
		panic(err)
	}
	q, err := client.GetClient(pan.Quark)
	if err != nil {
		panic(err)
	}
	return c, q
}

func StartMove() {
	fileChan := make(chan string)
	doneChan := make(chan struct{})
	moveConfig := internal.Config.Move
	cloudreve, quark := moveInit()
	go func() {
		for {
			select {
			case file := <-fileChan:
				if file == "" {
					fmt.Println("接收到全部下载完成")
					close(doneChan)
					return
				}
				absolutePath, err := filepath.Abs(moveConfig.TmpPath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				relativePath, err := filepath.Rel(absolutePath, filepath.Dir(file))
				if err != nil {
					fmt.Println(err)
					continue
				}
				relativePath = strings.ReplaceAll(relativePath, "\\", "/")
				if relativePath == "." {
					relativePath = "/"
				}
				err = quark.UploadFile(pan.UploadFileReq{
					LocalFile:  file,
					RemotePath: relativePath,
					Resumable:  true,
					SuccessDel: true,
					RemoteTransfer: func(remotePath, remoteName string) (string, string) {
						newFileName := remoteName
						newRemotePath := remotePath
						for _, removeStr := range moveConfig.RemoveStr {
							newFileName = strings.ReplaceAll(newFileName, removeStr, "")
							newRemotePath = strings.ReplaceAll(newRemotePath, removeStr, "")
						}
						// 使用正则表达式替换字符串
						if moveConfig.RemoveReg != "" {
							re := regexp.MustCompile(moveConfig.RemoveReg)
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
					fmt.Println(err)
					continue
				}
			}
		}
	}()
	go func() {
		downloadDir := &pan.PanObj{
			Name: strings.TrimLeft(moveConfig.RemotePath, "/"),
			Path: "/",
			Type: "dir",
		}
		objs, err := quark.List(pan.ListReq{
			Reload: true,
			Dir:    downloadDir,
		})
		if err != nil {
			fmt.Println(err)
			close(doneChan)
		}
		ignoreFiles := make([]string, 0)
		ignorePaths := make([]string, 0)
		for _, obj := range objs {
			if obj.Type == "dir" {
				ignorePaths = append(ignorePaths, obj.Name)
			} else {
				ignoreFiles = append(ignoreFiles, obj.Name)
			}
		}
		err = cloudreve.DownloadPath(pan.DownloadPathReq{
			RemotePath:  downloadDir,
			SkipFileErr: true,
			LocalPath:   filepath.Join(moveConfig.TmpPath, moveConfig.RemotePath),
			Concurrency: moveConfig.DownloadThread,
			ChunkSize:   moveConfig.DownloadChunkSize,
			OverCover:   false,
			IgnorePaths: ignorePaths,
			IgnoreFiles: ignoreFiles,
			DownloadCallback: func(localPath, localFile string) {
				fileChan <- localFile
			},
		})
		if err != nil {
			fmt.Println(err)
			close(doneChan)
		} else {
			close(fileChan)
		}
	}()

	select {
	case <-doneChan:
		return
	}

}
